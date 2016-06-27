package binding

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

const (
	k_BINDING_TAG                 = "binding"
	k_BINDING_CLEANED_FUNC_PREFIX = "Cleaned"
	k_BINDING_NO_TAG              = "-"
	k_BINDING_CLEANED_DATA        = "CleanedData"
	k_BINDING_DEFAULT_FUNC_PREFIX = "Default"
)

func Bind(source map[string]interface{}, result interface{}) (err error) {
	return BindWithTag(source, result, k_BINDING_TAG)
}

func BindWithTag(source map[string]interface{}, result interface{}, tag string) (err error) {
	var objType = reflect.TypeOf(result)
	var objValue = reflect.ValueOf(result)
	var objValueKind = objValue.Kind()

	if objValueKind == reflect.Struct {
		return errors.New("obj is struct")
	}

	if objValue.IsNil() {
		return errors.New("obj is nil")
	}

	for {
		if objValueKind == reflect.Ptr && objValue.IsNil() {
			objValue.Set(reflect.New(objType.Elem()))
		}

		if objValueKind == reflect.Ptr {
			objValue = objValue.Elem()
			objType = objType.Elem()
			objValueKind = objValue.Kind()
			continue
		}
		break
	}

	var cleanDataValue = objValue.FieldByName(k_BINDING_CLEANED_DATA)
	if cleanDataValue.IsValid() && cleanDataValue.IsNil() {
		cleanDataValue.Set(reflect.MakeMap(cleanDataValue.Type()))
	}
	return bindWithMap(objType, objValue, objValue, cleanDataValue, source, tag)
}

func bindWithMap(objType reflect.Type, currentObjValue, objValue, cleanDataValue reflect.Value, source map[string]interface{}, tagName string) error {
	var numField = objType.NumField()
	for i := 0; i < numField; i++ {
		var fieldStruct = objType.Field(i)
		var fieldValue = objValue.Field(i)

		if !fieldValue.CanSet() {
			continue
		}

		var tag = fieldStruct.Tag.Get(tagName)

		if tag == "" {
			tag = fieldStruct.Name

			if fieldValue.Kind() == reflect.Ptr {
				if fieldValue.IsNil() {
					fieldValue.Set(reflect.New(fieldValue.Type().Elem()))
				}
				fieldValue = fieldValue.Elem()
			}

			if fieldValue.Kind() == reflect.Struct {
				if err := bindWithMap(fieldValue.Addr().Type().Elem(), currentObjValue, fieldValue, cleanDataValue, source, tagName); err != nil {
					return err
				}
				continue
			}

		} else if tag == k_BINDING_NO_TAG {
			continue
		}

		var value, exists = source[tag]
		if !exists {
			setDefaultValue(currentObjValue, objValue, fieldValue, cleanDataValue, fieldStruct, tag)
			continue
		}

		if err := setValue(currentObjValue, objValue, fieldValue, fieldStruct, value); err != nil {
			return err
		}

		if cleanDataValue.IsValid() {
			cleanDataValue.SetMapIndex(reflect.ValueOf(tag), fieldValue)
		}
	}
	return nil
}

func getFuncWithName(funcName string, currentObjValue, objValue reflect.Value) reflect.Value {
	var funcValue = currentObjValue.MethodByName(funcName)
	if funcValue.IsValid() == false {
		if currentObjValue.CanAddr() {
			funcValue = currentObjValue.Addr().MethodByName(funcName)
		}
	}
	if funcValue.IsValid() == false && currentObjValue != objValue {
		return getFuncWithName(funcName, objValue, objValue)
	}
	return funcValue
}

func setDefaultValue(currentObjValue, objValue, fieldValue, cleanDataValue reflect.Value, fieldStruct reflect.StructField, tag string) {
	var funcValue = getFuncWithName(k_BINDING_DEFAULT_FUNC_PREFIX + fieldStruct.Name, currentObjValue, objValue)
	if funcValue.IsValid() {
		var rList = funcValue.Call(nil)
		if cleanDataValue.IsValid() {
			cleanDataValue.SetMapIndex(reflect.ValueOf(tag), rList[0])
		}
		fieldValue.Set(rList[0])
	}
}

func setValue(currentObjValue, objValue, fieldValue reflect.Value, fieldStruct reflect.StructField, value interface{}) error {
	var vValue = reflect.ValueOf(value)
	var fieldValueKind = fieldValue.Kind()

	var mValue = getFuncWithName(k_BINDING_CLEANED_FUNC_PREFIX + fieldStruct.Name, currentObjValue, objValue)
	if mValue.IsValid() {
		var rList = mValue.Call([]reflect.Value{vValue})
		if len(rList) > 1 {
			var rValue1 = rList[1]
			if rValue1.IsNil() == false {
				return rValue1.Interface().(error)
			}
		}
		fieldValue.Set(rList[0])
	} else if fieldValueKind == reflect.Slice && fieldValue.IsNil() == false {
		var valueLen int
		if vValue.Kind() == reflect.Slice {
			valueLen = vValue.Len()
			var s = reflect.MakeSlice(fieldValue.Type(), valueLen, valueLen)
			for i:=0; i<valueLen; i++ {
				if err := _setValue(s.Index(i), fieldStruct, vValue.Index(i)); err != nil {
					return err
				}
			}
			fieldValue.Set(s)
		} else {
			valueLen = 1
			var s = reflect.MakeSlice(fieldValue.Type(), valueLen, valueLen)
			if err := _setValue(s.Index(0), fieldStruct, vValue); err != nil {
				return err
			}
			fieldValue.Set(s)
		}
	} else {
		return _setValue(fieldValue, fieldStruct, vValue)
	}
	return nil
}

func _setValue(fieldValue reflect.Value, fieldStruct reflect.StructField, value reflect.Value) error {
	var valueKind = value.Kind()
	var fieldKind = fieldValue.Kind()
	if valueKind == fieldKind {
		return _setValueWithSameKind(fieldValue, fieldStruct, valueKind, value)
	}
	return _setValueWithDiffKind(fieldValue, fieldStruct, valueKind, value)
}

func _setValueWithSameKind(fieldValue reflect.Value, fieldStruct reflect.StructField, valueKind reflect.Kind, value reflect.Value) error {
	switch valueKind {
	case reflect.String:
		fieldValue.SetString(value.String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fieldValue.SetInt(value.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		fieldValue.SetUint(value.Uint())
	case reflect.Float32, reflect.Float64:
		fieldValue.SetFloat(value.Float())
	case reflect.Bool:
		fieldValue.SetBool(value.Bool())
		//	case reflect.Complex64, reflect.Complex128:
		//		fieldValue.SetComplex(value.Complex())
	default:
		return errors.New(fmt.Sprintf("Unknown type: %s", fieldStruct.Name))
	}
	return nil
}

func _setValueWithDiffKind(fieldValue reflect.Value, fieldStruct reflect.StructField, valueKind reflect.Kind, value reflect.Value) error {
	var f, err = floatValue(valueKind, value)
	if err != nil {
		return errors.New(fmt.Sprintln("[" + fieldStruct.Name + "]" + err.Error()))
	}

	var fieldValueKind = fieldValue.Kind()

	switch fieldValueKind {
	case reflect.String:
		fieldValue.SetString(fmt.Sprintf("%f", f))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fieldValue.SetInt(int64(f))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		fieldValue.SetUint(uint64(f))
	case reflect.Float32, reflect.Float64:
		fieldValue.SetFloat(f)
	case reflect.Bool:
		if f >= 0.990 {
			fieldValue.SetBool(true)
		} else {
			fieldValue.SetBool(false)
		}
		//	case reflect.Complex64, reflect.Complex128:
		//		fieldValue.SetComplex(value.Complex())
	default:
		return errors.New(fmt.Sprintf("Unknown type: %s", fieldStruct.Name))
	}
	return nil
}

func floatValue(valueKind reflect.Kind, value reflect.Value) (float64, error) {
	switch valueKind {
	case reflect.String:
		var v, e = strconv.ParseFloat(value.String(), 64)
		return v, e
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(value.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return float64(value.Uint()), nil
	case reflect.Float32, reflect.Float64:
		return value.Float(), nil
	case reflect.Bool:
		var b = value.Bool()
		if b {
			return 1.0, nil
		}
		return 0.0, nil
	}
	return 0.0, nil
}
