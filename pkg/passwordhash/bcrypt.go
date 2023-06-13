package passwordhash

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
)

var bcryptType = HashType{
	name: "bcrypt",
}

type BcryptParams struct {
	cost int
}

type BcryptHash struct {
	params   BcryptParams
	hashType *HashType
	hash     []byte
}

func BcryptFromHash(hash string) (BcryptHash, error) {
	bytes := []byte(hash)
	cost, err := bcrypt.Cost(bytes)
	if err != nil {
		return BcryptHash{}, err
	}
	params := BcryptParams{
		cost: cost,
	}
	return BcryptHash{
		hash:     bytes,
		params:   params,
		hashType: &bcryptType,
	}, nil
}

func BcryptFromPassword(password string, cost int) (BcryptHash, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return BcryptHash{}, err
	}
	return BcryptHash{
		hash: bytes,
		params: BcryptParams{
			cost: cost,
		},
		hashType: &bcryptType,
	}, nil
}

func (BcryptParams) String() string {
	return ""
}

func (h *BcryptHash) Type() *HashType {
	return h.hashType
}

func (h *BcryptHash) String() string {
	return string(h.hash)
}

func (h *BcryptHash) Params() HashParams {
	return h.params
}

func (h *BcryptHash) VerifyPassword(password string) bool {
	return bcrypt.CompareHashAndPassword(h.hash, []byte(password)) == nil
}

func (h *BcryptHash) UnmarshalJSON(data []byte) error {
	var v string
	var hash BcryptHash
	var err error
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}

	hash, err = BcryptFromHash(v)
	if err != nil {
		return err
	}

	h.hashType = hash.hashType
	h.params = hash.params
	h.hash = hash.hash
	return nil
}

func (h *BcryptHash) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.String())
}
