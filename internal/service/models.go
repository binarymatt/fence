package service

import (
	"math"
	"time"

	"github.com/cedar-policy/cedar-go"
	"github.com/uptrace/bun"
	"google.golang.org/protobuf/types/known/structpb"

	fencev1 "github.com/binarymatt/fence/gen/fence/v1"
)

type Base struct {
	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
	// DeletedAt bun.NullTime `bun:"deleted_at,soft_delete"`
}

type Policy struct {
	bun.BaseModel `bun:"table:policies,alias:p"`
	ID            string `bun:"id,pk"`
	Content       string `bun:"content"`
	Base
}

func (p *Policy) ToProto() *fencev1.Policy {
	return &fencev1.Policy{
		Id:         p.ID,
		Definition: p.Content,
	}
}

type UID struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type Entity struct {
	bun.BaseModel `bun:"table:entities,alias:e"`
	ID            string       `bun:",pk"`
	Type          string       `bun:"type,pk"`
	Parents       []UID        `bun:"parents,type:json"`
	Attributes    cedar.Record `bun:"attributes,type:json"`
	Tags          cedar.Record `bun:"tags,type:json"`
	Base
}

func (e *Entity) ToProto() *fencev1.Entity {
	parents := make([]*fencev1.UID, len(e.Parents))
	for i, p := range e.Parents {
		parents[i] = &fencev1.UID{Id: p.ID, Type: p.Type}
	}
	return &fencev1.Entity{
		Uid:     &fencev1.UID{Id: e.ID, Type: e.Type},
		Parents: parents,
	}
}

func isWhole(x float64) bool {
	return math.Ceil(x) == x
}
func structToCedarValue(v *structpb.Value) cedar.Value {
	switch v.GetKind().(type) {
	case *structpb.Value_StringValue:
		return cedar.String(v.GetStringValue())
	case *structpb.Value_BoolValue:
		return cedar.Boolean(v.GetBoolValue())
	case *structpb.Value_NumberValue:
		f := v.GetNumberValue()
		if isWhole(f) {
			return cedar.Long(v.GetNumberValue())
		} else {
			d, err := cedar.NewDecimalFromFloat(f)
			if err != nil {
				return nil
			}
			return d
		}
	}
	return nil
}
func fenceToRecord(values map[string]*structpb.Value) cedar.Record {
	m := cedar.RecordMap{}
	for k, v := range values {
		m[cedar.String(k)] = structToCedarValue(v)
	}
	return cedar.NewRecord(m)
}
func fenceToDBUID(ui *fencev1.UID) UID {
	return UID{
		ID:   ui.GetId(),
		Type: ui.GetType(),
	}
}
