package sqlx

import (
	"context"
	"testing"
)

func TestDial(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // if we call this without defer then the test should be failed as the context passed to the Dial is already closed
	db, err := Dial(ctx, DefaultDriver, ":memory:")
	if err != nil {
		t.Errorf("got err %q, want nil", err)
	}

	if db == nil {
		t.Errorf("got nil, want non nil")
	}
}
