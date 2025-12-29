package providers

import (
	"context"
	"errors"
	"testing"
	"testing/synctest"
	"time"

	"github.com/shoenig/test/must"

	fencev1 "github.com/binarymatt/fence/gen/fence/v1"
)

func TestNewCachedState(t *testing.T) {
	mockedState := NewMockFenceProvider(t)
	mockedState.EXPECT().Refresh(context.Background()).Return(nil)
	state, err := NewCachedProvider(mockedState, 1*time.Second)
	must.NoError(t, err)
	must.NotNil(t, state)

}
func TestNewCachedState_RefreshError(t *testing.T) {
	returnedErr := errors.New("oops")
	mockedState := NewMockFenceProvider(t)
	mockedState.EXPECT().Refresh(context.Background()).Return(returnedErr)
	state, err := NewCachedProvider(mockedState, 1*time.Second)
	must.ErrorIs(t, err, returnedErr)
	must.Nil(t, state)
}

func TestRefresh(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		mockedState := NewMockFenceProvider(t)
		mockedState.EXPECT().Refresh(context.Background()).Return(nil)
		state, err := NewCachedProvider(mockedState, time.Second)
		must.NoError(t, err)
		must.NotNil(t, state)
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
	mockedState := NewMockFenceProvider(t)
	mockedState.EXPECT().Refresh(context.Background()).Return(nil).Once()
	state, err := NewCachedProvider(mockedState, time.Second)
	must.NoError(t, err)
	mockedState.EXPECT().IsAllowed(t.Context(), bob, action, resource).Return(nil)
	err = state.IsAllowed(t.Context(), bob, action, resource)
	must.NoError(t, err)
}
