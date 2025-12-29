package service

import (
	"context"
	"testing"

	"github.com/shoenig/test/must"
	"google.golang.org/protobuf/types/known/structpb"

	fencev1 "github.com/binarymatt/fence/gen/fence/v1"
)

func TestCreateEntity(t *testing.T) {
	cases := []struct {
		name      string
		entity    *fencev1.Entity
		setupMock func(context.Context, *MockDataStore)
		err       error
	}{
		{
			name: "basic happy path",
			entity: &fencev1.Entity{
				Uid:     &fencev1.UID{Type: "User", Id: "bob"},
				Parents: []*fencev1.UID{{Type: "Group", Id: "admins"}},
				Attributes: map[string]*structpb.Value{
					"age": structpb.NewNumberValue(40),
				},
				Tags: map[string]*structpb.Value{
					"a": structpb.NewStringValue("b"),
				},
			},
			setupMock: func(ctx context.Context, ds *MockDataStore) {
				ds.EXPECT().addEntity(ctx, &fencev1.Entity{
					Uid:     &fencev1.UID{Type: "User", Id: "bob"},
					Parents: []*fencev1.UID{{Type: "Group", Id: "admins"}},
					Attributes: map[string]*structpb.Value{
						"age": structpb.NewNumberValue(40),
					},
					Tags: map[string]*structpb.Value{
						"a": structpb.NewStringValue("b"),
					},
				}).Return(nil)
			},
		},
		{
			name: "existing entity",
			entity: &fencev1.Entity{
				Uid:     &fencev1.UID{Type: "User", Id: "bob"},
				Parents: []*fencev1.UID{{Type: "Group", Id: "admins"}},
				Attributes: map[string]*structpb.Value{
					"age": structpb.NewNumberValue(40),
				},
				Tags: map[string]*structpb.Value{
					"a": structpb.NewStringValue("b"),
				},
			},
			setupMock: func(ctx context.Context, ds *MockDataStore) {
				expected := &fencev1.Entity{
					Uid:     &fencev1.UID{Type: "User", Id: "bob"},
					Parents: []*fencev1.UID{{Type: "Group", Id: "admins"}},
					Attributes: map[string]*structpb.Value{
						"age": structpb.NewNumberValue(40),
					},
					Tags: map[string]*structpb.Value{
						"a": structpb.NewStringValue("b"),
					},
				}
				ds.EXPECT().addEntity(ctx, expected).Return(ErrEntityAlreadyExists)
			},
			err: ErrEntityAlreadyExists,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s, ds := setupTest(t)
			tc.setupMock(t.Context(), ds)

			req := &fencev1.CreateEntitiesRequest{
				Entities: []*fencev1.Entity{tc.entity},
			}
			_, err := s.CreateEntities(t.Context(), req)
			must.ErrorIs(t, err, tc.err)
		})
	}
}

func TestDeleteEntity(t *testing.T) {
	cases := []struct {
		name      string
		id        string
		typ       string
		setupMock func(context.Context, *MockDataStore)
		err       error
	}{
		{
			name: "happy path",
			id:   "bob",
			typ:  "User",
			setupMock: func(ctx context.Context, ds *MockDataStore) {
				ds.EXPECT().deleteEntity(ctx, &fencev1.UID{
					Id:   "bob",
					Type: "User",
				}).Return(nil)
			},
		},
		{
			name: "does not exist",
			id:   "jane",
			typ:  "User",
			setupMock: func(ctx context.Context, ds *MockDataStore) {
				uid := &fencev1.UID{
					Id:   "jane",
					Type: "User",
				}
				ds.EXPECT().deleteEntity(ctx, uid).Return(ErrEntityNotFound)
			},
			err: ErrEntityNotFound,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s, ds := setupTest(t)
			tc.setupMock(t.Context(), ds)

			req := &fencev1.DeleteEntityRequest{
				Uid: &fencev1.UID{
					Id:   tc.id,
					Type: tc.typ,
				},
			}
			_, err := s.DeleteEntity(t.Context(), req)
			must.ErrorIs(t, err, tc.err)
		})
	}

}
