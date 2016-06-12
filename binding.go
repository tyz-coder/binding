package binding

import (
	"errors"
	"reflect"
	"fmt"
	"strconv"
)

const (
	k_MODEL_TAG                 = "model"
	k_MODEL_CLEANED_FUNC_PREFIX = "Cleaned"
	k_MODEL_NO_TAG              = "-"
	k_MODEL_CLEANED_DATA        = "CleanedData"
	k_MODEL_DEFAULT_FUNC_PREFIX = "Default"
)

func Bind(source map[string]interface{}, result interface{}) (err error) {
	return BindWithTag(source, result, k_MODEL_TAG)
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

	var cleanDataValue = objValue.FieldByName(k_MODEL_CLEANED_DATA)
	if cleanDataValue.IsValid() && cleanDataValue.IsNil() {
		cleanDataValue.Set(reflect.MakeMap(cleanDataValue.Type()))
	}
	return bindWithMap(objType, objValue, cleanDataValue, source, tag)
}

func bindWithMap(objType reflect.Type, objValue, cleanDataValue reflect.Value, source map[string]interface{}, tagName string) error {
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
				if err := bindWithMap(fieldValue.Addr().Type().Elem(), fieldValue, cleanDataValue, source, tagName); err != nil {
					return err
				}
				continue
			}

		} else if tag == k_MODEL_NO_TAG {
			continue
		}

		var value, exists = source[tag]
		if !exists {
			setDefaultValue(objValue, fieldValue, cleanDataValue, fieldStruct, tag)
			continue
		}

		//fieldValue.Set(reflect.ValueOf(value))
		if err := setValue(objValue, fieldValue, fieldStruct, value); err != nil {
			return err
		}

		if cleanDataValue.IsValid() {
			cleanDataValue.SetMapIndex(reflect.ValueOf(tag), fieldValue)
		}
	}
	return nil
}

func setDefaultValue(objValue, fieldValue, cleanDataValue reflect.Value, fieldStruct reflect.StructField, tag string) {
	var mName = k_MODEL_DEFAULT_FUNC_PREFIX + fieldStruct.Name
	var mValue = objValue.MethodByName(mName)
	if mValue.IsValid() == false {
		if objValue.CanAddr() {
			mValue = objValue.Addr().MethodByName(mName)
		}
	}

	if mValue.IsValid() {
		var rList = mValue.Call(nil)
		if cleanDataValue.IsValid() {
			cleanDataValue.SetMapIndex(reflect.ValueOf(tag), rList[0])
		}
		fieldValue.Set(rList[0])
	}
}

func setValue(objValue, fieldValue reflect.Value, fieldStruct reflect.StructField, value interface{}) error {
	var vValue = reflect.ValueOf(value)

	var mName = k_MODEL_CLEANED_FUNC_PREFIX + fieldStruct.Name
	var mValue = objValue.MethodByName(mName)
	if mValue.IsValid() == false {
		if objValue.CanAddr() {
			mValue = objValue.Addr().MethodByName(mName)
		}
	}

	if mValue.IsValid() {
		var rList = mValue.Call([]reflect.Value{vValue})
		if len(rList) > 1 {
			var rValue1 = rList[1]
			if rValue1.IsNil() == false {
				return rValue1.Interface().(error)
			}
		}
		fieldValue.Set(rList[0])
	} else {
		var valueKind = vValue.Kind()
		if valueKind == fieldValue.Kind() {
			return _setValueWithSameKind(fieldValue, fieldStruct, valueKind, vValue)
		} else {
			return _setValueWithDiffKind(fieldValue, fieldStruct, valueKind, vValue)
		}
	}
	return nil
}

func _setValueWithSameKind(fieldValue reflect.Value, fieldStruct reflect.StructField, valueKind reflect.Kind, value reflect.Value) error {
	switch valueKind {
	case reflect.String:
		fieldValue.SetString(value.String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fieldValue.SetInt(value.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
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

func _setValueWithDiffKind(fieldValue reflect.Value, fieldStruct reflect.StructField, valueKind reflect.Kind, value reflect.Value) (error) {
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
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		fieldValue.SetUint(uint64(f))
	case reflect.Float32, reflect.Float64:
		fieldValue.SetFloat(f)
	case reflect.Bool:
		if f >= 1.0000 {
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
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
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