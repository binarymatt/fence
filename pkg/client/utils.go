package client

import fencev1 "github.com/binarymatt/fence/gen/fence/v1"

func NewUID(uidType, id string) *fencev1.UID {
	return &fencev1.UID{
		Type: uidType,
		Id:   id,
	}
}
