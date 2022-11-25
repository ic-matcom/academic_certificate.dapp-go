package lib

import (
	"dapp/schema"
	"dapp/schema/dto"
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"

	"os"
	"reflect"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// DepObtainUserDid this tries to get the user DID store in the previously generated auth Bearer token.
func DepObtainUserDid(ctx iris.Context) dto.InjectedParam {
	tkData := jwt.Get(ctx).(*dto.AccessTokenData)

	// returning the DID and Identifier (Username)
	return tkData.Claims
}

func ParamsToStruct(ctx iris.Context, resStruct any) error {
	paramsMap := ctx.URLParams()
	paramsEncoded, err := json.Marshal(paramsMap)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(paramsEncoded, &resStruct); err != nil {
		return err
	}
	return nil
}

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
func MapToSliceOfValues[T any](dMap map[string]T) []T {
	var values []T
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

// ToMap struct to Map[string]interface{}
func ToMap(in interface{}, tagName string) (map[string]interface{}, error) {
	out := make(map[string]interface{})

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct { // Non-structural return error
		return nil, fmt.Errorf("ToMap only accepts struct or struct pointer; got %T", v)
	}

	t := v.Type()
	// Traversing structure fields
	// Specify the tagName value as the key in the map; the field value as the value in the map
	for i := 0; i < v.NumField(); i++ {
		fi := t.Field(i)
		if tagValue := fi.Tag.Get(tagName); tagValue != "" {
			out[tagValue] = v.Field(i).Interface()
		}
	}
	return out, nil
}
