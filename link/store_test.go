package link

import (
	"context"
	"errors"
	"linkd/bite"
	"linkd/sqlx/sqlxtest"
	"testing"
)

func TestStoreCreate(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	link := Link{Key: "go", URL: "https://go.dev/"}

	t.Run("ok", func(t *testing.T) {
		t.Parallel()

		store := NewStore(sqlxtest.Dial(ctx, t))
		if err := store.Create(ctx, link); err != nil {
			t.Errorf("Create(%q) err = %v, want nil", link.Key, err)
		}
	})

	t.Run("err_exists", func(t *testing.T) {
		t.Parallel()
		store := NewStore(sqlxtest.Dial(ctx, t))
		if err := store.Create(ctx, link); err != nil {
			t.Errorf("Create(%q) err = %v, want nil", link.Key, err)
		}
		if err := store.Create(ctx, link); !errors.Is(err, bite.ErrExists) {
			t.Errorf("Create(%q) err = %v, want %v", link.Key, err, bite.ErrExists)
		}
	})

}

func TestStoreRetrieve(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	link := Link{Key: "go", URL: "https://go.dev"}
	t.Run("ok", func(t *testing.T) {
		t.Parallel()
		store := NewStore(sqlxtest.Dial(ctx, t))
		if err := store.Create(ctx, link); err != nil {
			t.Errorf("Create(%q) err = %v, want nil", link.Key, err)
		}
		got, err := store.Retrieve(ctx, link.Key)
		if err != nil {
			t.Errorf("Retrieve(%q) err = %v, want nil", link.Key, err)
		}
		if got != link {
			t.Errorf("Retrieve(%q) = %#v, want %#v", link.Key, got, link)
		}
	})
	t.Run("err_not_exist", func(t *testing.T) {
		t.Parallel()
		store := NewStore(sqlxtest.Dial(ctx, t))
		_, err := store.Retrieve(ctx, "void")
		if !errors.Is(err, bite.ErrNotExists) {
			t.Errorf(`Retrieve(void) err = %v, want %v`, err, bite.ErrNotExists)
		}
	})
}
