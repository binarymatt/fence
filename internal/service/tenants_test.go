package service

import (
	"context"
	"testing"

	"connectrpc.com/connect/v2"
	"connectrpc.com/connect/v2/connectinprocess"
	cvalidate "connectrpc.com/validate"
	"github.com/shoenig/test/must"
	"github.com/stretchr/testify/mock"

	fencev1 "github.com/binarymatt/fence/gen/fence/v1"
	"github.com/binarymatt/fence/gen/fence/v1/fencev1connect"
)

func TestCreateTenant(t *testing.T) {
	ctx := context.Background()
	svc, mockDS := setupTest(t)
	server := connect.NewServer(cvalidate.NewServerInterceptor())
	fencev1connect.RegisterFenceAdminServiceHandler(server, svc)
	client := fencev1connect.NewFenceAdminServiceClient(
		connect.NewClient(connectinprocess.New(server)),
	)

	expectedTenant := &fencev1.Tenant{
		Id: &fencev1.UID{
			Type: "Tenant",
			Id:   "123",
		},
		Name: "CreateTest",
	}
	mockDS.EXPECT().AddTenant(mock.AnythingOfType("*context.valueCtx"), expectedTenant).Return(nil).Once()

	_, err := client.CreateTenant(ctx, &fencev1.CreateTenantRequest{
		Tenant: &fencev1.Tenant{
			Id: &fencev1.UID{
				Type: "Tenant",
				Id:   "123",
			},
			Name: "CreateTest",
		},
	})
	must.NoError(t, err)

	_, err = client.CreateTenant(ctx, &fencev1.CreateTenantRequest{
		Tenant: &fencev1.Tenant{
			Id: &fencev1.UID{
				Type: "Test",
				Id:   "123",
			},
		},
	})
	must.Eq(t, "invalid_argument: validation error: tenant.id: must be tenant entity", err.Error())

}
