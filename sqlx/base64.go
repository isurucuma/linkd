package sqlx

import (
	"database/sql/driver"
	"encoding/base64"
	"fmt"
)

type Base64String string

func (s Base64String) Value() (driver.Value, error) {
	dst := []byte(s)
	return base64.StdEncoding.EncodeToString(dst), nil
}

func (s *Base64String) Scan(src any) error {
	ss, ok := src.(string)
	if !ok {
		return fmt.Errorf("%q is %T not string", ss, src)
	}

	dst, err := base64.StdEncoding.DecodeString(ss)
	if err != nil {
		return fmt.Errorf("decoding %q: %w", ss, err)
	}
	*s = Base64String(dst)
	return nil
}

func (s Base64String) String() string {
	return string(s)
}
