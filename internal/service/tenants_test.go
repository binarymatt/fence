package service

import (
	"testing"

	"buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	"buf.build/go/protovalidate"
	"github.com/shoenig/test/must"
	"google.golang.org/protobuf/testing/protocmp"

	fencev1 "github.com/binarymatt/fence/gen/fence/v1"
)

func TestCreateTenant_Validation(t *testing.T) {
	cases := []struct {
		name string
		req  *fencev1.CreateTenantRequest
		err  error
	}{
		{
			name: "Valid",
			req: &fencev1.CreateTenantRequest{
				Tenant: &fencev1.Entity{
					Uid: &fencev1.UID{Type: "Tenant", Id: "test"},
				},
			},
			err: nil,
		},
		{
			name: "non tenant type",
			req: &fencev1.CreateTenantRequest{
				Tenant: &fencev1.Entity{
					Uid: &fencev1.UID{Type: "user", Id: "test"},
				},
			},
			err: &protovalidate.ValidationError{
				Violations: []*protovalidate.Violation{
					{
						Proto: &validate.Violation{
							Message: new("must be tenant entity"),
							RuleId:  new("tenant.type"),
						},
					},
				},
			},
		},
		{
			name: "no parents",
			req: &fencev1.CreateTenantRequest{
				Tenant: &fencev1.Entity{
					Uid: &fencev1.UID{Type: "Tenant", Id: "test"},
					Parents: []*fencev1.UID{
						{Type: "Tenant", Id: "test"},
					},
				},
			},
			err: &protovalidate.ValidationError{
				Violations: []*protovalidate.Violation{
					{
						Proto: &validate.Violation{
							Message: new("must have no parents"),
							RuleId:  new("tenant.no.parents"),
						},
					},
				},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := protovalidate.Validate(tc.req)
			must.Eq(t, tc.err, err, must.Cmp(protocmp.Transform()))
		})
	}
}
