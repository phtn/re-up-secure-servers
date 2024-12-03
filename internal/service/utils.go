package service

import (
	"reflect"
	"strings"
)

// StructEmptyOptions allows specifying which fields to check for emptiness
type StructEmptyOptions struct {
	// Fields to check. If empty, checks all fields
	Fields []string
	// Ignore unexported fields
	IgnoreUnexported bool
}

// IsStructEmpty checks if specified fields (or all fields) are empty
func IsStructEmpty(s interface{}, opts ...StructEmptyOptions) bool {
	v := reflect.ValueOf(s)

	// Check if it's actually a struct
	if v.Kind() != reflect.Struct {
		L.Fail(r, "hcode params-validation", "input not a struct")
	}

	// Default options if not provided
	option := StructEmptyOptions{}
	if len(opts) > 0 {
		option = opts[0]
	}

	// Iterate through all fields
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := v.Type().Field(i)

		// Skip unexported fields if specified
		if option.IgnoreUnexported && !field.CanInterface() {
			continue
		}

		// Check if we should examine this field
		shouldCheck := len(option.Fields) == 0 ||
			containsCaseInsensitive(option.Fields, fieldType.Name)

		if shouldCheck && !isZeroValue(field) {
			return false
		}
	}

	return true
}

// containsCaseInsensitive checks if a slice contains a string case-insensitively
func containsCaseInsensitive(slice []string, str string) bool {
	for _, s := range slice {
		if strings.EqualFold(s, str) {
			return true
		}
	}
	return false
}

// isZeroValue checks if a reflect.Value is at its zero value
func isZeroValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Func, reflect.Map, reflect.Slice:
		return v.IsNil()
	case reflect.Array:
		for i := 0; i < v.Len(); i++ {
			if !isZeroValue(v.Index(i)) {
				return false
			}
		}
		return true
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			if !isZeroValue(v.Field(i)) {
				return false
			}
		}
		return true
	case reflect.Ptr:
		return v.IsNil() || isZeroValue(v.Elem())
	}

	// Compare with zero value
	zero := reflect.Zero(v.Type())
	return reflect.DeepEqual(v.Interface(), zero.Interface())
}

// Example structs to demonstrate
type User struct {
	ID       int
	Name     string
	Email    string
	Active   bool
	Metadata map[string]string
}

type ComplexStruct struct {
	SimpleField  string
	NestedStruct struct {
		Field1 int
		Field2 string
	}
	Pointer *int
	Slice   []string
}

// func main() {
// 	// User struct examples
// 	user := User{
// 		ID:    1,
// 		Name:  "John Doe",
// 		Email: "",
// 	}

// 	// Check all fields
// 	fmt.Println("All fields empty:", IsStructEmpty(user)) // false

// 	// Check only specific fields
// 	fmt.Println("Only Name empty:",
// 		IsStructEmpty(user, StructEmptyOptions{
// 			Fields: []string{"Name"},
// 		})) // false

// 	fmt.Println("Only Email empty:",
// 		IsStructEmpty(user, StructEmptyOptions{
// 			Fields: []string{"Email"},
// 		})) // true

// 	// Complex struct examples
// 	complexStruct := ComplexStruct{
// 		SimpleField: "",
// 		Pointer:     new(int),
// 	}

// 	fmt.Println("Complex struct - all fields:",
// 		IsStructEmpty(complexStruct)) // false

// 	fmt.Println("Complex struct - only SimpleField:",
// 		IsStructEmpty(complexStruct, StructEmptyOptions{
// 			Fields: []string{"SimpleField"},
// 		})) // true

// 	fmt.Println("Complex struct - ignore unexported:",
// 		IsStructEmpty(complexStruct, StructEmptyOptions{
// 			IgnoreUnexported: true,
// 		})) // false
// }
