package providers

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"testing/synctest"
	"time"

	"github.com/shoenig/test/must"

	fencev1 "github.com/binarymatt/fence/gen/fence/v1"
)

func TestNewCachedProvider(t *testing.T) {
	mockedProvider := NewMockFenceProvider(t)
	mockedProvider.EXPECT().Refresh(context.Background()).Return(nil)
	cp, err := NewCachedProvider(mockedProvider, 1*time.Second)
	must.NoError(t, err)
	must.NotNil(t, cp)
}
func TestNewCachedProvider_RefreshError(t *testing.T) {
	returnedErr := errors.New("oops")
	mockedProvider := NewMockFenceProvider(t)
	mockedProvider.EXPECT().Refresh(context.Background()).Return(returnedErr)
	state, err := NewCachedProvider(mockedProvider, 1*time.Second)
	must.ErrorIs(t, err, returnedErr)
	must.Nil(t, state)
}

func TestRefresh(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		mockedProvider := NewMockFenceProvider(t)
		ctx, cancel := context.WithCancel(t.Context())
		defer cancel()
		mockedProvider.EXPECT().Refresh(ctx).Return(nil).Twice()
		cp := &CachedProvider{ip: mockedProvider, refreshDuration: 10 * time.Second}
		go cp.Refresh(ctx)
		fmt.Println("sleeping")
		time.Sleep(20 * time.Second)
		synctest.Wait()
	})
}

func TestCachedIsAllowed(t *testing.T) {

	bob := &fencev1.UID{
		Type: "User",
		Id:   "bob",
	}
	resource := &fencev1.UID{
		Type: "Photo",
		Id:   "VacationPhoto94.jpg",
	}
	action := &fencev1.UID{
		Type: "Action",
		Id:   "view",
	}
	mockedProvider := NewMockFenceProvider(t)
	mockedProvider.EXPECT().Refresh(context.Background()).Return(nil).Once()
	state, err := NewCachedProvider(mockedProvider, time.Second)
	must.NoError(t, err)
	mockedProvider.EXPECT().IsAllowed(t.Context(), bob, action, resource).Return(&fencev1.IsAllowedResponse{Decision: true}, nil)
	resp, err := state.IsAllowed(t.Context(), bob, action, resource)
	must.NoError(t, err)
	must.True(t, resp.Decision)
}
