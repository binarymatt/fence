package service

import (
	"testing"

	"github.com/cedar-policy/cedar-go"
	"github.com/shoenig/test/must"

	"github.com/binarymatt/fence/internal/translation"
)

func TestMapToRecordRoundTrip(t *testing.T) {
	decimal, err := cedar.NewDecimalFromFloat(1.5)
	must.NoError(t, err)

	cases := map[string]cedar.Record{
		"empty": cedar.NewRecord(cedar.RecordMap{}),
		"primitives": cedar.NewRecord(cedar.RecordMap{
			"name":  cedar.String("alice"),
			"admin": cedar.True,
			"age":   cedar.Long(30),
		}),
		"nested": cedar.NewRecord(cedar.RecordMap{
			"tags":   cedar.NewSet(cedar.String("a"), cedar.String("b")),
			"nested": cedar.NewRecord(cedar.RecordMap{"level": cedar.Long(2)}),
		}),
		"entity and extension": cedar.NewRecord(cedar.RecordMap{
			"owner": cedar.NewEntityUID("User", "bob"),
			"ratio": decimal,
		}),
	}
	for name, record := range cases {
		t.Run(name, func(t *testing.T) {
			got := translation.MapToRecord(recordToMap(record))
			must.True(t, record.Equal(got), must.Sprintf("expected %s, got %s", record, got))
		})
	}
}
