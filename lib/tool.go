package lib

import (
	"dapp/schema"
	"encoding/json"
	"fmt"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"os"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

func Contains(arr []string, elem string) bool {
	for _, e := range arr {
		if elem == e {
			return true
		}
	}
	return false
}

func SliceToMap(slice []string, dMap map[string]string) {
	for _, data := range slice {
		if _, ok := dMap[data]; !ok {
			dMap[data] = ""
		}
	}
}

// MapToSliceOfKey Convert map to slice of keys.
func MapToSliceOfKey(dMap map[string]string) []string {
	slice := make([]string, 0)
	for key := range dMap {
		slice = append(slice, key)
	}
	return slice
}

// MapToSliceOfValues Convert map to slice of values.
func MapToSliceOfValues(dMap map[string]any) []any {
	var values []any
	for _, value := range dMap {
		values = append(values, value)
	}
	return values
}

func GetEnvOrDefault(env, defaultVal string) string {
	value, ok := os.LookupEnv(env)
	if !ok {
		value = defaultVal
	}
	return value
}

// GetEnvOrError Check environment variable. If the environment variable does not exist,
// the panic function raised
func GetEnvOrError(key string) string {
	value, found := os.LookupEnv(key)
	if found {
		return value
	}

	errorString := fmt.Sprintf("The environment variable '%s' is not set: %s \n%s", key, schema.ErrInvalidEnvVar, "-) Read the readme!!")
	panic(errorString)
}

// GetBoolOrDefault Note that the method returns default value if the string
// cannot be parsed!
func GetBoolOrDefault(value string, defaultVal bool) bool {
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return defaultVal
	}
	return parsed
}

func UpdateJSON(request interface{}, stateDB interface{}) ([]byte, error) {
	// JSON encoding
	requestMarshal, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	// sobreescribiendo el activo con los valores del request
	err = json.Unmarshal(requestMarshal, stateDB)
	if err != nil {
		return nil, err
	}
	// result JSON encoding
	resMarshal, err := json.Marshal(stateDB)
	if err != nil {
		return nil, err
	}

	return resMarshal, nil
}

// ConcatenateBytes is useful for combining multiple arrays of bytes, especially for
// signatures or digests over multiple fields
func ConcatenateBytes(data ...[]byte) []byte {
	finalLength := 0
	for _, slice := range data {
		finalLength += len(slice)
	}
	result := make([]byte, finalLength)
	last := 0
	for _, slice := range data {
		for i := range slice {
			result[i+last] = slice[i]
		}
		last += len(slice)
	}
	return result
}

func DeepCopy(v interface{}) (interface{}, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	vptr := reflect.New(reflect.TypeOf(v))
	err = json.Unmarshal(data, vptr.Interface())
	if err != nil {
		return nil, err
	}
	return vptr.Elem().Interface(), err
}

func NormalizeString(text string, upper bool) string {
	if upper {
		text = strings.ToUpper(text)
	}

	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), runes.Remove(runes.In(unicode.Space)), norm.NFC) // Mn: nonspacing marks
	result, _, err := transform.String(t, text)
	if err != nil {
		return text
	}

	return result
}

func Unique(input []interface{}) []interface{} {
	u := make([]interface{}, 0, len(input))
	m := make(map[string]bool)

	for _, val := range input {
		if _, ok := m[val.(string)]; !ok {
			m[val.(string)] = true
			u = append(u, val)
		}
	}

	return u
}
func UniqueStrings(input []string) []string {
	u := make([]string, 0, len(input))
	m := make(map[string]bool)

	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}

	return u
}

func TrimDoubleQuotes(text string) string {
	text = strings.ReplaceAll(text, " ", "")
	text = strings.ReplaceAll(text, "\"", "")

	return text
}
