/* =================================================================================
 * NCLOUD SDK LAYER FOR TERRAFORM CODEGEN - DO NOT EDIT
 * ================================================================================= */

package ncloudsdk

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"database/sql/driver"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

type NClient struct {
	BaseURL    string
	HTTPClient *http.Client
	AccessKey  string
	SecretKey  string
}

var (
	ErrInvalidCopyDestination        = errors.New("copy destination must be non-nil and addressable")
	ErrInvalidCopyFrom               = errors.New("copy from must be non-nil and addressable")
	ErrMapKeyNotMatch                = errors.New("map's key type doesn't match")
	ErrNotSupported                  = errors.New("not supported")
	ErrFieldNameTagStartNotUpperCase = errors.New("copier field name tag must be start upper case")
	ErrInvalidEndpoint               = errors.New("invalid endpoint URL")
	ErrAPIRequest                    = errors.New("API request failed")
)

func NewClient(baseURL, accessKey, secretKey string) *NClient {
	return &NClient{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{},
		AccessKey:  accessKey,
		SecretKey:  secretKey,
	}
}

// MakeRequestWithContext() - Streamlined core logic of abstracted api call
//
// Manufacture main request call
func (n *NClient) MakeRequestWithContext(ctx context.Context, method, endpoint, reqBody string, query map[string]string) (map[string]interface{}, error) {
	url, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidEndpoint, err)
	}

	req, err := n.SetRequest(url, query, reqBody, method)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	// Make signature & set headers
	n.SetHeader(req, url, method)

	// Execute api call
	resp, err := n.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check if resp NoContent
	if resp.StatusCode == http.StatusNoContent {
		return map[string]interface{}{}, nil
	}

	// Parse response into map[string]interface{}
	var respBody map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return nil, err
	}

	return respBody, nil
}

func (n *NClient) SetRequest(url *url.URL, queryParams map[string]string, reqBody, method string) (*http.Request, error) {
	q := url.Query()
	for key, value := range queryParams {
		q.Add(key, value)
	}

	url.RawQuery = q.Encode()

	b := bytes.NewBuffer([]byte(reqBody))

	req, err := http.NewRequest(method, url.String(), b)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (n *NClient) SetHeader(req *http.Request, url *url.URL, method string) {

	// Check if query string exists.
	// If then, do not even add "?".
	queryString := ""
	if len(url.RawQuery) > 0 {
		queryString = "?" + url.RawQuery
	}

	timestamp := fmt.Sprintf("%d", time.Now().UnixMilli())
	signature := makeSignature(method, url.Path+queryString, timestamp, n.AccessKey, n.SecretKey)

	headers := map[string]string{
		"Content-Type":             "application/json",
		"x-ncp-apigw-timestamp":    timestamp,
		"x-ncp-iam-access-key":     n.AccessKey,
		"x-ncp-apigw-signature-v2": signature,
		"Cache-Control":            "no-cache",
		"Pragma":                   "no-cache",
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}
}

// For curl request
func makeSignature(method, url, timestamp, accessKey, secretKey string) string {
	message := fmt.Sprintf("%s %s\n%s\n%s",
		method,
		url,
		timestamp,
		accessKey,
	)

	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(message))

	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func diagOff[V, T interface{}](input func(ctx context.Context, elementType T, elements any) (V, diag.Diagnostics), ctx context.Context, elementType T, elements any) V {
	var emptyReturn V

	v, diags := input(ctx, elementType, elements)

	if diags.HasError() {
		diags.AddError("REFRESHING ERROR", "invalid diagOff operation")
		return emptyReturn
	}

	return v
}

// convertKeys recursively converts all keys in a map from camelCase to snake_case
func convertKeys(input interface{}) interface{} {
	switch v := input.(type) {
	case map[string]interface{}:
		newMap := make(map[string]interface{})
		for key, value := range v {
			// Convert the key to snake_case
			newKey := camelToSnake(key)
			// Recursively convert nested values
			newMap[newKey] = convertKeys(value)
		}
		return newMap
	case []interface{}:
		newSlice := make([]interface{}, len(v))
		for i, value := range v {
			newSlice[i] = convertKeys(value)
		}
		return newSlice
	default:
		return v
	}
}

func camelToSnake(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && unicode.IsUpper(r) {
			result.WriteRune('_')
		}
		result.WriteRune(unicode.ToLower(r))
	}
	return result.String()
}

func ClearDoubleQuote(s string) string {
	return strings.Replace(strings.Replace(strings.Replace(s, "\\", "", -1), "\"", "", -1), `"`, "", -1)
}

// Copier.go code from https://github.com/jinzhu/copier/blob/master/copier.go

// These flags define options for tag handling
const (
	// Denotes that a destination field must be copied to. If copying fails then a panic will ensue.
	tagMust uint8 = 1 << iota

	// Denotes that the program should not panic when the must flag is on and
	// value is not copied. The program will return an error instead.
	tagNoPanic

	// Ignore a destination field from being copied to.
	tagIgnore

	// Denotes the fact that the field should be overridden, no matter if the IgnoreEmpty is set
	tagOverride

	// Denotes that the value as been copied
	hasCopied

	// Some default converter types for a nicer syntax
	String  string  = ""
	Bool    bool    = false
	Int     int     = 0
	Float32 float32 = 0
	Float64 float64 = 0
)

// Option sets copy options
type Option struct {
	// setting this value to true will ignore copying zero values of all the fields, including bools, as well as a
	// struct having all it's fields set to their zero values respectively (see IsZero() in reflect/value.go)
	IgnoreEmpty   bool
	CaseSensitive bool
	DeepCopy      bool
	Converters    []TypeConverter
	// Custom field name mappings to copy values with different names in `fromValue` and `toValue` types.
	// Examples can be found in `copier_field_name_mapping_test.go`.
	FieldNameMapping []FieldNameMapping
}

func (opt Option) converters() map[converterPair]TypeConverter {
	var converters = map[converterPair]TypeConverter{}

	// save converters into map for faster lookup
	for i := range opt.Converters {
		pair := converterPair{
			SrcType: reflect.TypeOf(opt.Converters[i].SrcType),
			DstType: reflect.TypeOf(opt.Converters[i].DstType),
		}

		converters[pair] = opt.Converters[i]
	}

	return converters
}

type TypeConverter struct {
	SrcType interface{}
	DstType interface{}
	Fn      func(src interface{}) (dst interface{}, err error)
}

type converterPair struct {
	SrcType reflect.Type
	DstType reflect.Type
}

func (opt Option) fieldNameMapping() map[converterPair]FieldNameMapping {
	var mapping = map[converterPair]FieldNameMapping{}

	for i := range opt.FieldNameMapping {
		pair := converterPair{
			SrcType: reflect.TypeOf(opt.FieldNameMapping[i].SrcType),
			DstType: reflect.TypeOf(opt.FieldNameMapping[i].DstType),
		}

		mapping[pair] = opt.FieldNameMapping[i]
	}

	return mapping
}

type FieldNameMapping struct {
	SrcType interface{}
	DstType interface{}
	Mapping map[string]string
}

// Tag Flags
type flags struct {
	BitFlags  map[string]uint8
	SrcNames  tagNameMapping
	DestNames tagNameMapping
}

// Field Tag name mapping
type tagNameMapping struct {
	FieldNameToTag map[string]string
	TagToFieldName map[string]string
}

// Copy copy things
func Copy(toValue interface{}, fromValue interface{}) (err error) {
	return copier(toValue, fromValue, Option{})
}

// CopyWithOption copy with option
func CopyWithOption(toValue interface{}, fromValue interface{}, opt Option) (err error) {
	return copier(toValue, fromValue, opt)
}

func copier(toValue interface{}, fromValue interface{}, opt Option) (err error) {
	var (
		isSlice    bool
		amount     = 1
		from       = indirect(reflect.ValueOf(fromValue))
		to         = indirect(reflect.ValueOf(toValue))
		converters = opt.converters()
		mappings   = opt.fieldNameMapping()
	)

	if !to.CanAddr() {
		return ErrInvalidCopyDestination
	}

	// Return is from value is invalid
	if !from.IsValid() {
		return ErrInvalidCopyFrom
	}

	fromType, isPtrFrom := indirectType(from.Type())
	toType, _ := indirectType(to.Type())

	if fromType.Kind() == reflect.Interface {
		fromType = reflect.TypeOf(from.Interface())
	}

	if toType.Kind() == reflect.Interface {
		toType, _ = indirectType(reflect.TypeOf(to.Interface()))
		oldTo := to
		to = reflect.New(reflect.TypeOf(to.Interface())).Elem()
		defer func() {
			oldTo.Set(to)
		}()
	}

	// Just set it if possible to assign for normal types
	if from.Kind() != reflect.Slice && from.Kind() != reflect.Struct && from.Kind() != reflect.Map && (from.Type().AssignableTo(to.Type()) || from.Type().ConvertibleTo(to.Type())) {
		if !isPtrFrom || !opt.DeepCopy {
			to.Set(from.Convert(to.Type()))
		} else {
			fromCopy := reflect.New(from.Type())
			fromCopy.Set(from.Elem())
			to.Set(fromCopy.Convert(to.Type()))
		}
		return
	}

	if from.Kind() != reflect.Slice && fromType.Kind() == reflect.Map && toType.Kind() == reflect.Map {
		if !fromType.Key().ConvertibleTo(toType.Key()) {
			return ErrMapKeyNotMatch
		}

		if to.IsNil() {
			to.Set(reflect.MakeMapWithSize(toType, from.Len()))
		}

		for _, k := range from.MapKeys() {
			toKey := indirect(reflect.New(toType.Key()))
			isSet, err := set(toKey, k, opt.DeepCopy, converters)
			if err != nil {
				return err
			}
			if !isSet {
				return fmt.Errorf("%w map, old key: %v, new key: %v", ErrNotSupported, k.Type(), toType.Key())
			}

			elemType := toType.Elem()
			if elemType.Kind() != reflect.Slice {
				elemType, _ = indirectType(elemType)
			}
			toValue := indirect(reflect.New(elemType))
			isSet, err = set(toValue, from.MapIndex(k), opt.DeepCopy, converters)
			if err != nil {
				return err
			}
			if !isSet {
				if err = copier(toValue.Addr().Interface(), from.MapIndex(k).Interface(), opt); err != nil {
					return err
				}
			}

			for {
				if elemType == toType.Elem() {
					to.SetMapIndex(toKey, toValue)
					break
				}
				elemType = reflect.PointerTo(elemType)
				toValue = toValue.Addr()
			}
		}
		return
	}

	if from.Kind() == reflect.Slice && to.Kind() == reflect.Slice {
		// Return directly if both slices are nil
		if from.IsNil() && to.IsNil() {
			return
		}
		if to.IsNil() {
			slice := reflect.MakeSlice(reflect.SliceOf(to.Type().Elem()), from.Len(), from.Cap())
			to.Set(slice)
		}
		if fromType.ConvertibleTo(toType) {
			for i := 0; i < from.Len(); i++ {
				if to.Len() < i+1 {
					to.Set(reflect.Append(to, reflect.New(to.Type().Elem()).Elem()))
				}
				isSet, err := set(to.Index(i), from.Index(i), opt.DeepCopy, converters)
				if err != nil {
					return err
				}
				if !isSet {
					// ignore error while copy slice element
					err = copier(to.Index(i).Addr().Interface(), from.Index(i).Interface(), opt)
					if err != nil {
						continue
					}
				}
			}

			if to.Len() > from.Len() {
				to.SetLen(from.Len())
			}

			return
		}
	}

	if fromType.Kind() != reflect.Struct || toType.Kind() != reflect.Struct {
		// skip not supported type
		return
	}

	if len(converters) > 0 {
		if ok, e := set(to, from, opt.DeepCopy, converters); e == nil && ok {
			// converter supported
			return
		}
	}

	if from.Kind() == reflect.Slice || to.Kind() == reflect.Slice {
		isSlice = true
		if from.Kind() == reflect.Slice {
			amount = from.Len()
		}
	}

	for i := 0; i < amount; i++ {
		var dest, source reflect.Value

		if isSlice {
			// source
			if from.Kind() == reflect.Slice {
				source = indirect(from.Index(i))
			} else {
				source = indirect(from)
			}
			// dest
			dest = indirect(reflect.New(toType).Elem())
		} else {
			source = indirect(from)
			dest = indirect(to)
		}

		if len(converters) > 0 {
			if ok, e := set(dest, source, opt.DeepCopy, converters); e == nil && ok {
				if isSlice {
					// FIXME: maybe should check the other types?
					if to.Type().Elem().Kind() == reflect.Ptr {
						to.Index(i).Set(dest.Addr())
					} else {
						if to.Len() < i+1 {
							reflect.Append(to, dest)
						} else {
							to.Index(i).Set(dest)
						}
					}
				} else {
					to.Set(dest)
				}

				continue
			}
		}

		destKind := dest.Kind()
		initDest := false
		if destKind == reflect.Interface {
			initDest = true
			dest = indirect(reflect.New(toType))
		}

		// Get tag options
		flgs, err := getFlags(dest, source, toType, fromType)
		if err != nil {
			return err
		}

		// check source
		if source.IsValid() {
			copyUnexportedStructFields(dest, source)

			// Copy from source field to dest field or method
			fromTypeFields := deepFields(fromType)
			for _, field := range fromTypeFields {
				name := field.Name

				// Get bit flags for field
				fieldFlags := flgs.BitFlags[name]

				// Check if we should ignore copying
				if (fieldFlags & tagIgnore) != 0 {
					continue
				}

				fieldNamesMapping := getFieldNamesMapping(mappings, fromType, toType)

				srcFieldName, destFieldName := getFieldName(name, flgs, fieldNamesMapping)

				if fromField := fieldByNameOrZeroValue(source, srcFieldName); fromField.IsValid() && !shouldIgnore(fromField, fieldFlags, opt.IgnoreEmpty) {
					// process for nested anonymous field
					destFieldNotSet := false
					if f, ok := dest.Type().FieldByName(destFieldName); ok {
						// only initialize parent embedded struct pointer in the path
						for idx := range f.Index[:len(f.Index)-1] {
							destField := dest.FieldByIndex(f.Index[:idx+1])

							if destField.Kind() != reflect.Ptr {
								continue
							}

							if !destField.IsNil() {
								continue
							}
							if !destField.CanSet() {
								destFieldNotSet = true
								break
							}

							// destField is a nil pointer that can be set
							newValue := reflect.New(destField.Type().Elem())
							destField.Set(newValue)
						}
					}

					if destFieldNotSet {
						break
					}

					toField := fieldByName(dest, destFieldName, opt.CaseSensitive)
					if toField.IsValid() {
						if toField.CanSet() {
							isSet, err := set(toField, fromField, opt.DeepCopy, converters)
							if err != nil {
								return err
							}
							if !isSet {
								if err := copier(toField.Addr().Interface(), fromField.Interface(), opt); err != nil {
									return err
								}
							}
							if fieldFlags != 0 {
								// Note that a copy was made
								flgs.BitFlags[name] = fieldFlags | hasCopied
							}
						}
					} else {
						// try to set to method
						var toMethod reflect.Value
						if dest.CanAddr() {
							toMethod = dest.Addr().MethodByName(destFieldName)
						} else {
							toMethod = dest.MethodByName(destFieldName)
						}

						if toMethod.IsValid() && toMethod.Type().NumIn() == 1 && fromField.Type().AssignableTo(toMethod.Type().In(0)) {
							toMethod.Call([]reflect.Value{fromField})
						}
					}
				}
			}

			// Copy from from method to dest field
			for _, field := range deepFields(toType) {
				name := field.Name
				srcFieldName, destFieldName := getFieldName(name, flgs, getFieldNamesMapping(mappings, fromType, toType))

				var fromMethod reflect.Value
				if source.CanAddr() {
					fromMethod = source.Addr().MethodByName(srcFieldName)
				} else {
					fromMethod = source.MethodByName(srcFieldName)
				}

				if fromMethod.IsValid() && fromMethod.Type().NumIn() == 0 && fromMethod.Type().NumOut() == 1 && !shouldIgnore(fromMethod, flgs.BitFlags[name], opt.IgnoreEmpty) {
					if toField := fieldByName(dest, destFieldName, opt.CaseSensitive); toField.IsValid() && toField.CanSet() {
						values := fromMethod.Call([]reflect.Value{})
						if len(values) >= 1 {
							set(toField, values[0], opt.DeepCopy, converters)
						}
					}
				}
			}
		}

		if isSlice && to.Kind() == reflect.Slice {
			if dest.Addr().Type().AssignableTo(to.Type().Elem()) {
				if to.Len() < i+1 {
					to.Set(reflect.Append(to, dest.Addr()))
				} else {
					isSet, err := set(to.Index(i), dest.Addr(), opt.DeepCopy, converters)
					if err != nil {
						return err
					}
					if !isSet {
						// ignore error while copy slice element
						err = copier(to.Index(i).Addr().Interface(), dest.Addr().Interface(), opt)
						if err != nil {
							continue
						}
					}
				}
			} else if dest.Type().AssignableTo(to.Type().Elem()) {
				if to.Len() < i+1 {
					to.Set(reflect.Append(to, dest))
				} else {
					isSet, err := set(to.Index(i), dest, opt.DeepCopy, converters)
					if err != nil {
						return err
					}
					if !isSet {
						// ignore error while copy slice element
						err = copier(to.Index(i).Addr().Interface(), dest.Interface(), opt)
						if err != nil {
							continue
						}
					}
				}
			}
		} else if initDest {
			to.Set(dest)
		}

	}

	return
}

func getFieldNamesMapping(mappings map[converterPair]FieldNameMapping, fromType reflect.Type, toType reflect.Type) map[string]string {
	var fieldNamesMapping map[string]string

	if len(mappings) > 0 {
		pair := converterPair{
			SrcType: fromType,
			DstType: toType,
		}
		if v, ok := mappings[pair]; ok {
			fieldNamesMapping = v.Mapping
		}
	}
	return fieldNamesMapping
}

func fieldByNameOrZeroValue(source reflect.Value, fieldName string) (value reflect.Value) {
	defer func() {
		if err := recover(); err != nil {
			value = reflect.Value{}
		}
	}()

	return source.FieldByName(fieldName)
}

func copyUnexportedStructFields(to, from reflect.Value) {
	if from.Kind() != reflect.Struct || to.Kind() != reflect.Struct || !from.Type().AssignableTo(to.Type()) {
		return
	}

	// create a shallow copy of 'to' to get all fields
	tmp := indirect(reflect.New(to.Type()))
	tmp.Set(from)

	// revert exported fields
	for i := 0; i < to.NumField(); i++ {
		if tmp.Field(i).CanSet() {
			tmp.Field(i).Set(to.Field(i))
		}
	}
	to.Set(tmp)
}

func shouldIgnore(v reflect.Value, bitFlags uint8, ignoreEmpty bool) bool {
	return ignoreEmpty && bitFlags&tagOverride == 0 && v.IsZero()
}

var deepFieldsLock sync.RWMutex
var deepFieldsMap = make(map[reflect.Type][]reflect.StructField)

func deepFields(reflectType reflect.Type) []reflect.StructField {
	deepFieldsLock.RLock()
	cache, ok := deepFieldsMap[reflectType]
	deepFieldsLock.RUnlock()
	if ok {
		return cache
	}
	var res []reflect.StructField
	if reflectType, _ = indirectType(reflectType); reflectType.Kind() == reflect.Struct {
		fields := make([]reflect.StructField, 0, reflectType.NumField())

		for i := 0; i < reflectType.NumField(); i++ {
			v := reflectType.Field(i)
			// PkgPath is the package path that qualifies a lower case (unexported)
			// field name. It is empty for upper case (exported) field names.
			// See https://golang.org/ref/spec#Uniqueness_of_identifiers
			if v.PkgPath == "" {
				fields = append(fields, v)
				if v.Anonymous {
					// also consider fields of anonymous fields as fields of the root
					fields = append(fields, deepFields(v.Type)...)
				}
			}
		}
		res = fields
	}

	deepFieldsLock.Lock()
	deepFieldsMap[reflectType] = res
	deepFieldsLock.Unlock()
	return res
}

func indirect(reflectValue reflect.Value) reflect.Value {
	for reflectValue.Kind() == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}
	return reflectValue
}

func indirectType(reflectType reflect.Type) (_ reflect.Type, isPtr bool) {
	for reflectType.Kind() == reflect.Ptr || reflectType.Kind() == reflect.Slice {
		reflectType = reflectType.Elem()
		isPtr = true
	}
	return reflectType, isPtr
}

func set(to, from reflect.Value, deepCopy bool, converters map[converterPair]TypeConverter) (bool, error) {
	if !from.IsValid() {
		return true, nil
	}
	if ok, err := lookupAndCopyWithConverter(to, from, converters); err != nil {
		return false, err
	} else if ok {
		return true, nil
	}

	if to.Kind() == reflect.Ptr {
		// set `to` to nil if from is nil
		if from.Kind() == reflect.Ptr && from.IsNil() {
			to.Set(reflect.Zero(to.Type()))
			return true, nil
		} else if to.IsNil() {
			// `from`         -> `to`
			// sql.NullString -> *string
			if fromValuer, ok := driverValuer(from); ok {
				v, err := fromValuer.Value()
				if err != nil {
					return true, nil
				}
				// if `from` is not valid do nothing with `to`
				if v == nil {
					return true, nil
				}
			}
			// allocate new `to` variable with default value (eg. *string -> new(string))
			to.Set(reflect.New(to.Type().Elem()))
		} else if from.Kind() != reflect.Ptr && from.IsZero() {
			to.Set(reflect.Zero(to.Type()))
			return true, nil
		}
		// depointer `to`
		to = to.Elem()
	}

	if deepCopy {
		toKind := to.Kind()
		if toKind == reflect.Interface && to.IsNil() {
			if reflect.TypeOf(from.Interface()) != nil {
				to.Set(reflect.New(reflect.TypeOf(from.Interface())).Elem())
				toKind = reflect.TypeOf(to.Interface()).Kind()
			}
		}
		if from.Kind() == reflect.Ptr && from.IsNil() {
			to.Set(reflect.Zero(to.Type()))
			return true, nil
		}
		if _, ok := to.Addr().Interface().(sql.Scanner); !ok && (toKind == reflect.Struct || toKind == reflect.Map || toKind == reflect.Slice) {
			return false, nil
		}
	}

	// try convert directly
	if from.Type().ConvertibleTo(to.Type()) {
		to.Set(from.Convert(to.Type()))
		return true, nil
	}

	// try Scanner
	if toScanner, ok := to.Addr().Interface().(sql.Scanner); ok {
		// `from`  -> `to`
		// *string -> sql.NullString
		if from.Kind() == reflect.Ptr {
			// if `from` is nil do nothing with `to`
			if from.IsNil() {
				return true, nil
			}
			// depointer `from`
			from = indirect(from)
		}
		// `from` -> `to`
		// string -> sql.NullString
		// set `to` by invoking method Scan(`from`)
		err := toScanner.Scan(from.Interface())
		if err == nil {
			return true, nil
		}
	}

	// try Valuer
	if fromValuer, ok := driverValuer(from); ok {
		// `from`         -> `to`
		// sql.NullString -> string
		v, err := fromValuer.Value()
		if err != nil {
			return false, nil
		}
		// if `from` is not valid do nothing with `to`
		if v == nil {
			return true, nil
		}
		rv := reflect.ValueOf(v)
		if rv.Type().AssignableTo(to.Type()) {
			to.Set(rv)
			return true, nil
		}
		if to.CanSet() && rv.Type().ConvertibleTo(to.Type()) {
			to.Set(rv.Convert(to.Type()))
			return true, nil
		}
		return false, nil
	}

	// from is ptr
	if from.Kind() == reflect.Ptr {
		return set(to, from.Elem(), deepCopy, converters)
	}

	return false, nil
}

// lookupAndCopyWithConverter looks up the type pair, on success the TypeConverter Fn func is called to copy src to dst field.
func lookupAndCopyWithConverter(to, from reflect.Value, converters map[converterPair]TypeConverter) (copied bool, err error) {
	pair := converterPair{
		SrcType: from.Type(),
		DstType: to.Type(),
	}

	if cnv, ok := converters[pair]; ok {
		result, err := cnv.Fn(from.Interface())
		if err != nil {
			return false, err
		}

		if result != nil {
			to.Set(reflect.ValueOf(result))
		} else {
			// in case we've got a nil value to copy
			to.Set(reflect.Zero(to.Type()))
		}

		return true, nil
	}

	return false, nil
}

// parseTags Parses struct tags and returns uint8 bit flags.
func parseTags(tag string) (flg uint8, name string, err error) {
	for _, t := range strings.Split(tag, ",") {
		switch t {
		case "-":
			flg = tagIgnore
			return
		case "must":
			flg = flg | tagMust
		case "nopanic":
			flg = flg | tagNoPanic
		case "override":
			flg = flg | tagOverride
		default:
			if unicode.IsUpper([]rune(t)[0]) {
				name = strings.TrimSpace(t)
			} else {
				err = ErrFieldNameTagStartNotUpperCase
			}
		}
	}
	return
}

// getTagFlags Parses struct tags for bit flags, field name.
func getFlags(dest, src reflect.Value, toType, fromType reflect.Type) (flags, error) {
	flgs := flags{
		BitFlags: map[string]uint8{},
		SrcNames: tagNameMapping{
			FieldNameToTag: map[string]string{},
			TagToFieldName: map[string]string{},
		},
		DestNames: tagNameMapping{
			FieldNameToTag: map[string]string{},
			TagToFieldName: map[string]string{},
		},
	}

	var toTypeFields, fromTypeFields []reflect.StructField
	if dest.IsValid() {
		toTypeFields = deepFields(toType)
	}
	if src.IsValid() {
		fromTypeFields = deepFields(fromType)
	}

	// Get a list dest of tags
	for _, field := range toTypeFields {
		tags := field.Tag.Get("copier")
		if tags != "" {
			var name string
			var err error
			if flgs.BitFlags[field.Name], name, err = parseTags(tags); err != nil {
				return flags{}, err
			} else if name != "" {
				flgs.DestNames.FieldNameToTag[field.Name] = name
				flgs.DestNames.TagToFieldName[name] = field.Name
			}
		}
	}

	// Get a list source of tags
	for _, field := range fromTypeFields {
		tags := field.Tag.Get("copier")
		if tags != "" {
			var name string
			var err error

			if _, name, err = parseTags(tags); err != nil {
				return flags{}, err
			} else if name != "" {
				flgs.SrcNames.FieldNameToTag[field.Name] = name
				flgs.SrcNames.TagToFieldName[name] = field.Name
			}
		}
	}

	return flgs, nil
}

// checkBitFlags Checks flags for error or panic conditions.
func checkBitFlags(flagsList map[string]uint8) (err error) {
	// Check flag conditions were met
	for name, flgs := range flagsList {
		if flgs&hasCopied == 0 {
			switch {
			case flgs&tagMust != 0 && flgs&tagNoPanic != 0:
				err = fmt.Errorf("field %s has must tag but was not copied", name)
				return
			case flgs&(tagMust) != 0:
				panic(fmt.Sprintf("Field %s has must tag but was not copied", name))
			}
		}
	}
	return
}

func getFieldName(fieldName string, flgs flags, fieldNameMapping map[string]string) (srcFieldName string, destFieldName string) {
	// get dest field name
	if name, ok := fieldNameMapping[fieldName]; ok {
		srcFieldName = fieldName
		destFieldName = name
		return
	}

	if srcTagName, ok := flgs.SrcNames.FieldNameToTag[fieldName]; ok {
		destFieldName = srcTagName
		if destTagName, ok := flgs.DestNames.TagToFieldName[srcTagName]; ok {
			destFieldName = destTagName
		}
	} else {
		if destTagName, ok := flgs.DestNames.TagToFieldName[fieldName]; ok {
			destFieldName = destTagName
		}
	}
	if destFieldName == "" {
		destFieldName = fieldName
	}

	// get source field name
	if destTagName, ok := flgs.DestNames.FieldNameToTag[fieldName]; ok {
		srcFieldName = destTagName
		if srcField, ok := flgs.SrcNames.TagToFieldName[destTagName]; ok {
			srcFieldName = srcField
		}
	} else {
		if srcField, ok := flgs.SrcNames.TagToFieldName[fieldName]; ok {
			srcFieldName = srcField
		}
	}

	if srcFieldName == "" {
		srcFieldName = fieldName
	}
	return
}

func driverValuer(v reflect.Value) (i driver.Valuer, ok bool) {
	if !v.CanAddr() {
		i, ok = v.Interface().(driver.Valuer)
		return
	}

	i, ok = v.Addr().Interface().(driver.Valuer)
	return
}

func fieldByName(v reflect.Value, name string, caseSensitive bool) reflect.Value {
	if caseSensitive {
		return v.FieldByName(name)
	}

	return v.FieldByNameFunc(func(n string) bool { return strings.EqualFold(n, name) })
}

func ConvertStructToStringMap(input interface{}) (interface{}, error) {
	val := reflect.ValueOf(input)
	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("input must be a struct")
	}

	result := make(map[string]string)
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldName := val.Type().Field(i).Tag.Get("json")
		if fieldName == "" {
			fieldName = val.Type().Field(i).Name
		}

		switch field.Kind() {
		case reflect.Slice:
			var sliceResult []string
			for j := 0; j < field.Len(); j++ {
				elem := field.Index(j).Interface()
				jsonValue, err := json.Marshal(elem)
				if err != nil {
					return nil, err
				}
				sliceResult = append(sliceResult, string(jsonValue))
			}
			jsonSlice, err := json.Marshal(sliceResult)
			if err != nil {
				return nil, err
			}
			result[fieldName] = string(jsonSlice)
		default:
			jsonValue, err := json.Marshal(field.Interface())
			if err != nil {
				return nil, err
			}
			result[fieldName] = string(jsonValue)
		}
	}

	return result, nil
}
