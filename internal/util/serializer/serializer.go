package serializer

// import (
// 	"reflect"
// 	"slices"
// 	"strings"

// 	string_utils "github.com/ebobola-dev/socially-app-go-server/internal/util/strings"
// )

// const (
// 	safeTag  = "safe"
// 	shortTag = "short"
// )

// var tags = []string{safeTag, shortTag}

// type SerializeOptions struct {
// 	Safe  bool
// 	Short bool
// }

// func Struct(v any, options SerializeOptions) map[string]interface{} {
// 	val := reflect.ValueOf(v)
// 	if val.Kind() == reflect.Pointer {
// 		val = val.Elem()
// 	}
// 	typ := val.Type()

// 	out := map[string]interface{}{}

// 	for i := 0; i < typ.NumField(); i++ {
// 		field := typ.Field(i)
// 		tag := field.Tag.Get("serializer")
// 		if tag == "" {
// 			continue
// 		}
// 		flags := strings.Split(tag, ",")
// 		if len(flags) == 0 || flags[0] == "" {
// 			continue
// 		}
// 		if slices.Contains(flags, "safe") && !options.Safe {
// 			continue
// 		}
// 		if !slices.Contains(flags, "short") && options.Short {
// 			continue
// 		}

// 		fieldName := string_utils.ToSnakeCase(field.Name)
// 		if !slices.Contains(tags, flags[0]) {
// 			fieldName = flags[0]
// 		}
// 		out[fieldName] = val.Field(i).Interface()
// 	}

// 	return out
// }
