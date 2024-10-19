package evaluate

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/PaesslerAG/gval"
	log "github.com/sirupsen/logrus"
)

var varselector = gval.NewLanguage(
	gval.VariableSelector(func(path gval.Evaluables) gval.Evaluable {
		return func(c context.Context, v interface{}) (interface{}, error) {
			defer recoverGvalFunc("variable selector")
			keys, err := path.EvalStrings(c, v)
			if err != nil {
				return nil, err
			}
			for i, k := range keys {
				if strings.HasPrefix(k, "ENV_") {
					ctx := c.Value("ctx")
					if eCtx, ok := ctx.(Context); ok {
						varName := k[4:]
						v, ok = reflectSelect(varName, eCtx.EnvVars)
						if !ok {
							log.Debugf("unknown parameter %s", strings.Join(keys[:i+1], "."))
							return noIdentifierPresented{}, nil
						}
						continue
					}
				}

				switch o := v.(type) {
				case gval.Selector:
					v, err = o.SelectGVal(c, k)
					if err != nil {
						return nil, fmt.Errorf("failed to select '%s' on %T: %w", k, o, err)
					}
					continue
				case map[interface{}]interface{}:
					v = o[k]
					continue
				case map[string]interface{}:
					v = o[k]
					continue
				case []interface{}:
					if i, err := strconv.Atoi(k); err == nil && i >= 0 && len(o) > i {
						v = o[i]
						continue
					}
				default:
					var ok bool
					v, ok = reflectSelect(k, o)
					if !ok {
						log.Debugf("unknown parameter %s", strings.Join(keys[:i+1], "."))
						return noIdentifierPresented{}, nil
					}
				}
			}
			return v, nil
		}
	}),
)

func reflectSelect(key string, value interface{}) (selection interface{}, ok bool) {
	vv := reflect.ValueOf(value)
	vvElem := resolvePotentialPointer(vv)

	switch vvElem.Kind() {
	case reflect.Map:
		mapKey, ok := reflectConvertTo(vv.Type().Key().Kind(), key)
		if !ok {
			return nil, false
		}

		vvElem = vv.MapIndex(reflect.ValueOf(mapKey))
		vvElem = resolvePotentialPointer(vvElem)

		if vvElem.IsValid() {
			return vvElem.Interface(), true
		}
	case reflect.Slice:
		if i, err := strconv.Atoi(key); err == nil && i >= 0 && vv.Len() > i {
			vvElem = resolvePotentialPointer(vv.Index(i))
			return vvElem.Interface(), true
		}
	case reflect.Struct:
		field := vvElem.FieldByName(key)
		if field.IsValid() {
			return field.Interface(), true
		}

		method := vv.MethodByName(key)
		if method.IsValid() {
			return method.Interface(), true
		}
	}
	return nil, false
}

func resolvePotentialPointer(value reflect.Value) reflect.Value {
	if value.Kind() == reflect.Ptr {
		return value.Elem()
	}
	return value
}

func reflectConvertTo(k reflect.Kind, value string) (interface{}, bool) {
	switch k {
	case reflect.String:
		return value, true
	case reflect.Int:
		if i, err := strconv.Atoi(value); err == nil {
			return i, true
		}
	}
	return nil, false
}
