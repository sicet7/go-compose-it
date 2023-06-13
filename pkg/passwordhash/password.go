package passwordhash

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
)

var (
	ErrInvalidHash         = errors.New("the encoded hash is not in the correct format")
	ErrUnsupportedType     = errors.New("unsupported hash type")
	ErrIncompatibleVersion = errors.New("incompatible version")
)

type HashType struct {
	name string
}

func (ht HashType) Name() string {
	return ht.name
}

type HashParams interface {
	fmt.Stringer
}

type Hash interface {
	fmt.Stringer
	json.Unmarshaler
	Type() *HashType
	Params() HashParams
	VerifyPassword(string) bool
}

func Parse(hashString string) (Hash, error) {
	matched, err := regexp.MatchString("^\\$2[a-zA-Z]\\$.+$", hashString)

	if err == nil && matched {
		b, bErr := BcryptFromHash(hashString)
		return &b, bErr
	}

	a, aErr := Argon2FromHash(hashString)
	return &a, aErr
}
