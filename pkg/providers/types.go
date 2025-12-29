package state

import (
	"fmt"

	"github.com/cedar-policy/cedar-go"

	fencev1 "github.com/binarymatt/fence/gen/fence/v1"
)

func UIDToString(u *fencev1.UID) string {
	return fmt.Sprintf(`%s::"%s"`, u.Type, u.Id)
}

func uidToCedar(u *fencev1.UID) cedar.EntityUID {
	return cedar.NewEntityUID(cedar.EntityType(u.Type), cedar.String(u.Id))
}
