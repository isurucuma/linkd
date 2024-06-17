package sqlxtest

import (
	"context"
	"fmt"
	"testing"

	"linkd/sqlx"
)

func Dial(ctx context.Context, tb testing.TB) *sqlx.DB { // TB helps to visible this Dial in both functionanl tests and benchmark tests
	tb.Helper() // this will inform go test as this is a helper function. Therefore if there is something faild in this function then
	// that line number and this file name will not gets printed, instead that particular test function's line where the Dial
	// gets called will be pointed.
	dsn := fmt.Sprintf(
		"file:%s?mode=memory&cache=shared",
		tb.Name(), // tb.Name() will printout the Name of the test function
	)
	db, err := sqlx.Dial(ctx, sqlx.DefaultDriver, dsn)
	if err != nil {
		tb.Fatalf("dialing test db: %v", err)
	}
	tb.Cleanup(func() { // this cleanup is crucial, kind of like similar to defer but rather than running the underline closure
		// when Dial exists this cleanup runs the closure when the actual test function exists
		if err := db.Close(); err != nil {
			tb.Log("closing test db:", err)
		}
	})
	return db
}
