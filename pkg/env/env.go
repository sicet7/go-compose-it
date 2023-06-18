package env

import (
	"errors"
	"golang.org/x/exp/constraints"
	"os"
	"strconv"
	"strings"
	"unsafe"
)

type validTypes interface {
	constraints.Integer | constraints.Float
}

func RequireNumber[T validTypes](name string) (result T, errResult error) {
	size := unsafe.Sizeof(result)
	stringVal, err := RequireString(name)
	if err != nil {
		return result, err
	}
	switch any(result).(type) {
	case float64:
	case float32:
		floatVal, floatErr := strconv.ParseFloat(stringVal, int(size))
		if floatErr != nil {
			return result, floatErr
		}
		return T(floatVal), nil
	}
	// truncating is fine here :-)
	intVal, intErr := strconv.ParseInt(stringVal, 10, int(size))
	if intErr != nil {
		return result, errors.New("invalid environment variable: \"" + name + "\" is not a integer value")
	}
	return T(intVal), nil
}

func RequireString(name string) (string, error) {
	value, ok := os.LookupEnv(name)
	if !ok || len(strings.TrimSpace(value)) == 0 {
		return "", errors.New("missing environment variable: \"" + name + "\"")
	}
	return value, nil
}

func RequireBool(name string) (bool, error) {
	stringVal, err := RequireString(name)
	if err != nil {
		return false, err
	}
	return parseBoolValue(stringVal)
}

func ReadNumber[T validTypes](name string, defaultValue T) (result T) {
	val, err := RequireNumber[T](name)
	if err != nil {
		return defaultValue
	}
	return val
}

func ReadString(name string, defaultValue string) string {
	val, err := RequireString(name)
	if err != nil {
		return defaultValue
	}
	return val
}

func ReadBool(name string, defaultValue bool) bool {
	val, err := RequireBool(name)
	if err != nil {
		return defaultValue
	}
	return val
}

func parseBoolValue(value string) (bool, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "true":
	case "1":
	case "yes":
	case "y":
		return true, nil
	case "false":
	case "0":
	case "no":
	case "n":
		return false, nil
	}
	return false, errors.New("failed to parse environment variable value")
}
