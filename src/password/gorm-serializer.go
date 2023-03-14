package password

import (
	"context"
	"encoding/json"
	"fmt"
	"gorm.io/gorm/schema"
	"reflect"
)

type PasswordSerializer struct {
}

func init() {
	schema.RegisterSerializer("password", PasswordSerializer{})
}

func (PasswordSerializer) Scan(
	ctx context.Context,
	field *schema.Field,
	dst reflect.Value,
	dbValue interface{},
) error {
	switch value := dbValue.(type) {
	case string:
		val, err := Parse(value)
		if err != nil {
			return fmt.Errorf("load of password hash failed %#v", err)
		}
		return field.Set(ctx, dst, val)
	default:
		return fmt.Errorf("unsupported data %#v", dbValue)
	}
}

func (PasswordSerializer) Value(
	ctx context.Context,
	field *schema.Field,
	dst reflect.Value,
	fieldValue interface{},
) (interface{}, error) {
	val, err := json.Marshal(fieldValue)
	return string(val), err
}
