package state

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
	mockedState := NewMockFenceState(t)
	mockedState.EXPECT().refresh().Return(nil)
	state, err := NewCachedState(mockedState, 1*time.Second)
	must.NoError(t, err)
	must.NotNil(t, state)

}
func TestNewCachedState_RefreshError(t *testing.T) {
	returnedErr := errors.New("oops")
	mockedState := NewMockFenceState(t)
	mockedState.EXPECT().refresh().Return(returnedErr)
	state, err := NewCachedState(mockedState, 1*time.Second)
	must.ErrorIs(t, err, returnedErr)
	must.Nil(t, state)
}

func TestRefresh(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		mockedState := NewMockFenceState(t)
		mockedState.EXPECT().refresh().Return(nil).Times(2)
		state, err := NewCachedState(mockedState, time.Second)
		must.NoError(t, err)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		time.Sleep(4 * time.Second)
		err = state.Refresh(ctx)
		must.ErrorIs(t, err, context.DeadlineExceeded)
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
	mockedState := NewMockFenceState(t)
	mockedState.EXPECT().refresh().Return(nil).Once()
	state, err := NewCachedState(mockedState, time.Second)
	must.NoError(t, err)
	mockedState.EXPECT().IsAllowed(t.Context(), bob, action, resource).Return(nil)
	err = state.IsAllowed(t.Context(), bob, action, resource)
	must.NoError(t, err)
}
