package password

import (
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/sicet7/go-compose-it/src/utils"
	"golang.org/x/crypto/argon2"
	"strings"
)

// calibrated with: "docker run -it --rm --entrypoint kratos oryd/kratos:v0.5 hashers argon2 calibrate 1s"
var (
	defaultParams = NewArgon2Params(
		4194304,
		1,
		64,
		16,
		32,
	)
	argon2iType = HashType{
		name: "argon2i",
	}
	argon2idType = HashType{
		name: "argon2id",
	}
)

type Argon2Params struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

type Argon2Hash struct {
	hashType *HashType
	hash     []byte
	salt     []byte
	params   Argon2Params
}

func NewArgon2Params(
	memory uint32,
	iterations uint32,
	parallelism uint8,
	saltLength uint32,
	keyLength uint32,
) Argon2Params {
	return Argon2Params{
		memory:      memory,
		iterations:  iterations,
		parallelism: parallelism,
		saltLength:  saltLength,
		keyLength:   keyLength,
	}
}

func DefaultArgon2Params() Argon2Params {
	return defaultParams
}

func Argon2FromHash(hashString string) (Argon2Hash, error) {
	vals := strings.Split(hashString, "$")
	if len(vals) != 6 {
		return Argon2Hash{}, ErrInvalidHash
	}

	supportedTypes := map[string]*HashType{
		argon2iType.Name():  &argon2iType,
		argon2idType.Name(): &argon2idType,
	}

	hashType, exists := supportedTypes[vals[1]]

	if !exists {
		return Argon2Hash{}, ErrUnsupportedType
	}

	var version int
	_, err := fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return Argon2Hash{}, err
	}
	if version != argon2.Version {
		return Argon2Hash{}, ErrIncompatibleVersion
	}
	p := Argon2Params{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &p.memory, &p.iterations, &p.parallelism)
	if err != nil {
		return Argon2Hash{}, err
	}

	var salt []byte
	salt, err = base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return Argon2Hash{}, err
	}
	p.saltLength = uint32(len(salt))

	var hash []byte
	hash, err = base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return Argon2Hash{}, err
	}
	p.keyLength = uint32(len(hash))

	return Argon2Hash{
		hash:     hash,
		salt:     salt,
		hashType: hashType,
		params:   p,
	}, nil
}

func Argon2iFromPassword(password string, params Argon2Params) (Argon2Hash, error) {
	return createArgonPassword(&argon2iType, password, params)
}

func Argon2idFromPassword(password string, params Argon2Params) (Argon2Hash, error) {
	return createArgonPassword(&argon2idType, password, params)
}

func (a Argon2Params) String() string {
	return fmt.Sprintf(
		"m=%d,t=%d,p=%d",
		a.memory,
		a.iterations,
		a.parallelism,
	)
}

func (h *Argon2Hash) Type() *HashType {
	return h.hashType
}

func (h *Argon2Hash) Params() HashParams {
	return h.params
}

func (h *Argon2Hash) String() string {
	b64Salt := base64.RawStdEncoding.EncodeToString(h.salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(h.hash)
	return fmt.Sprintf(
		"$%s$v=%d$%s$%s$%s",
		h.hashType.Name(),
		argon2.Version,
		h.params.String(),
		b64Salt,
		b64Hash,
	)
}

func (h *Argon2Hash) VerifyPassword(password string) bool {
	passwordHash, err := createArgonHash(h.hashType, password, h.salt, h.params)
	if err != nil {
		return false
	}
	return subtle.ConstantTimeCompare(h.hash, passwordHash.hash) == 1
}

func (h *Argon2Hash) UnmarshalJSON(data []byte) error {
	var v string
	var hash Argon2Hash
	var err error
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}

	hash, err = Argon2FromHash(v)
	if err != nil {
		return err
	}

	h.hashType = hash.hashType
	h.params = hash.params
	h.hash = hash.hash
	h.salt = hash.salt
	return nil
}

func (h *Argon2Hash) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.String())
}

func createArgonPassword(hashType *HashType, password string, params Argon2Params) (Argon2Hash, error) {
	salt, err := utils.GenerateRandomBytes(int(params.saltLength))
	if err != nil {
		return Argon2Hash{}, err
	}
	return createArgonHash(hashType, password, salt, params)
}

func createArgonHash(
	hashType *HashType,
	password string,
	salt []byte,
	params Argon2Params,
) (Argon2Hash, error) {
	var hash []byte

	if hashType.Name() == argon2iType.Name() {
		hash = argon2.Key(
			[]byte(password),
			salt,
			params.iterations,
			params.memory,
			params.parallelism,
			params.keyLength,
		)
	} else if hashType.Name() == argon2idType.Name() {
		hash = argon2.IDKey(
			[]byte(password),
			salt,
			params.iterations,
			params.memory,
			params.parallelism,
			params.keyLength,
		)
	} else {
		return Argon2Hash{}, ErrUnsupportedType
	}

	return Argon2Hash{
		hash:     hash,
		salt:     salt,
		hashType: hashType,
		params:   params,
	}, nil
}
