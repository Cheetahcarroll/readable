package readable

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func GetString(object interface{}, tag string) string {
	return getString(object, 0, tag)
}

func GetJSONModel(object interface{}) string {
	return "\n{\n" + getJSONModelFull(object, 0) + "\n}"
}

func GetUnitTest(object interface{}, value string) string {
	if err := json.Unmarshal([]byte(value), &object); err != nil {
		return "Unable to get unit test"
	}
	return reflect.TypeOf(object).Elem().Name() + "{\n" + getUnitTest(object, 1) + "}"
}

func getString(object interface{}, indents int, tag string) string {
	var indent, value = "", ""
	for i := 0; i < indents; i++ {
		indent += "\t"
	}

	t, v, ok := getReflections(object)
	if !ok {
		return value
	}

	if v.Kind() != reflect.Struct && (v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct) {
		return value
	}
	for i := 0; i < v.NumField(); i++ {
		test := t.Field(i)
		val := v.Field(i)
		label, ok := test.Tag.Lookup(tag)
		if !val.IsValid() {
			continue
		}
		if ok {
			value += indent + label + ": "
		}
		if val.Kind() == reflect.Struct || (val.Kind() == reflect.Ptr && val.Elem().Kind() == reflect.Struct) {
			indents++
			value += "\n" + getString(val.Interface(), indents, tag)
		} else {
			if val.Kind() == reflect.Slice || (val.Kind() == reflect.Ptr && val.Elem().Kind() == reflect.Slice) {
				value += "\n"
				for j := 0; j < val.Len(); j++ {
					sliceVal := val.Index(j)

					if sliceVal.Kind() == reflect.Struct || (sliceVal.Kind() == reflect.Ptr && sliceVal.Elem().Kind() == reflect.Struct) {
						indents++
						value += getString(sliceVal.Interface(), indents, tag)
						indents--
					} else {
						value += convertReadable(sliceVal.Interface(), false, false) + "\n"
					}
				}
			} else {
				value += convertReadable(val.Interface(), false, false) + "\n"
			}
		}
	}

	return value
}

func ToJSONString(object interface{}, indents int) string {
	var indent, value = "", "{\n"
	for i := 0; i < indents; i++ {
		indent += "\t"
	}

	t, v, ok := getReflections(object)
	if !ok {
		return value
	}

	for i := 0; i < v.NumField(); i++ {
		test := t.Field(i)
		val := v.Field(i)
		label, ok := test.Tag.Lookup("json")
		if !val.IsValid() {
			continue
		}
		label = strings.ReplaceAll(label, ",omitempty", "")
		if label == "-" {
			ok = false
		}
		if i > 0 {
			value += ",\n"
		}
		if ok {
			value += indent + "\t\"" + label + "\": "
		}
		if test.Type.Kind() == reflect.Struct || (test.Type.Kind() == reflect.Ptr && test.Type.Elem().Kind() == reflect.Struct) {
			indents++
			if !val.IsZero() {
				value += ToJSONString(val.Interface(), indents)
			} else {
				value += "{}"
			}
		} else {
			if val.Kind() == reflect.Slice || (val.Kind() == reflect.Ptr && val.Elem().Kind() == reflect.Slice) {
				for j := 0; j < val.Len(); j++ {
					sliceVal := val.Index(j)

					if sliceVal.Kind() == reflect.Struct || (sliceVal.Kind() == reflect.Ptr && sliceVal.Elem().Kind() == reflect.Struct) {
						indents++
						value += ToJSONString(sliceVal.Interface(), indents)
						indents--
					} else {
						value += convertReadable(sliceVal.Interface(), false, true)
					}
				}
			} else {
				value += convertReadable(val.Interface(), false, true)
			}
		}
	}

	return value + "\n" + indent + "}"
}

func ToFlatJSONString(object interface{}) string {
	value := "{"
	t, v, ok := getReflections(object)
	if !ok {
		return value
	}

	for i := 0; i < v.NumField(); i++ {
		test := t.Field(i)
		val := v.Field(i)
		label, ok := test.Tag.Lookup("json")
		if !val.IsValid() {
			continue
		}
		label = strings.ReplaceAll(label, ",omitempty", "")
		if label == "-" {
			label = test.Name
		}
		if i > 0 {
			value += ","
		}
		if ok {
			value += "\"" + label + "\":"
		}
		if test.Type.Kind() == reflect.Struct || (test.Type.Kind() == reflect.Ptr && test.Type.Elem().Kind() == reflect.Struct) {
			if !val.IsZero() && val.CanInterface() {
				value += ToFlatJSONString(val.Interface())
			} else {
				value += "{}"
			}
		} else if val.CanInterface() {
			if val.Kind() == reflect.Slice || (val.Kind() == reflect.Ptr && val.Elem().Kind() == reflect.Slice) {
				for j := 0; j < val.Len(); j++ {
					sliceVal := val.Index(j)

					if sliceVal.Kind() == reflect.Struct || (sliceVal.Kind() == reflect.Ptr && sliceVal.Elem().Kind() == reflect.Struct) {
						value += ToFlatJSONString(sliceVal.Interface())
					} else {
						value += convertReadable(sliceVal.Interface(), false, true)
					}
				}
			} else {
				value += convertReadable(val.Interface(), false, true)
			}
		}
	}

	return value + "}"
}

func getJSONModelFull(object interface{}, indents int) string {
	var (
		indent, value = "", ""
		appended      = false
	)
	for i := 0; i < indents; i++ {
		indent += "\t"
	}

	_, v, ok := getReflections(object)
	if !ok {
		return value
	}

	for i := 0; i < v.NumField(); i++ {
		if !appended && value != "" {
			appended = true
		}
		val := v.Field(i)
		if val.Kind() == reflect.Ptr {
			if val.CanInterface() && val.Elem().Kind() == reflect.Struct {
				val = val.Elem()
			} else if val.CanInterface() {
				val = reflect.ValueOf(val.Interface())
			}
		}

		label, ok := v.Type().Field(i).Tag.Lookup("json")
		label = strings.ReplaceAll(label, ",omitempty", "")
		if label == "-" {
			continue
		}

		if ok {
			if appended {
				value += ",\n"
			}
			appended = true
			value += indent + "\"" + label + "\": "

			switch val.Kind() {
			case reflect.Struct:
				switch t := val.Interface().(type) {
				case time.Time:
					value += convertReadable(t, false, true)
				default:
					indents++
					value += "{\n" + getJSONModelFull(val.Interface(), indents) + "\n" + indent + "}"
					indents--
				}
			case reflect.Ptr:
				fallthrough
			default:
				value += convertReadable(val.Interface(), false, true)
			}
		} else if val.CanInterface() {
			if val.Kind() == reflect.Struct {
				if appended {
					value += ",\n"
				}
				indents++
				append := getJSONModelFull(val.Interface(), indents)
				indents--
				appended = true
				value += append
			}
		}
	}

	return value
}

func getUnitTest(object interface{}, indents int) string {
	var indent, value = "", ""
	for i := 0; i < indents; i++ {
		indent += "\t"
	}

	t, v, ok := getReflections(object)
	if !ok {
		return value
	}

	for i := 0; i < v.NumField(); i++ {
		test := t.Field(i)
		val := v.Field(i)
		if !val.IsValid() || val.IsZero() {
			continue
		}
		value += indent + test.Name + ": "
		switch {
		case test.Type.Kind() == reflect.Struct:
			indents++
			value += "{\n" + getUnitTest(val.Interface(), indents) + indent + "},\n"
		case test.Type.Kind() == reflect.Ptr && test.Type.Elem().Kind() == reflect.Struct:
			indents++
			value += "&" + test.Type.Elem().Name() + "{\n" + getUnitTest(val.Interface(), indents) + indent + "},\n"
		case val.Kind() == reflect.Slice || (val.Kind() == reflect.Ptr && val.Elem().Kind() == reflect.Slice):
			indents++
			if test.Type.Elem().Kind() == reflect.Ptr {
				value += "[]*" + test.Type.Elem().Elem().Name() + "{\n"
			} else {
				value += "[]" + test.Type.Elem().Name() + "{\n"
			}
			for j := 0; j < val.Len(); j++ {
				sliceVal := val.Index(j)

				switch {
				case sliceVal.Kind() == reflect.Struct:
					indents++
					value += indent + "\t" + "{\n" + getUnitTest(sliceVal.Interface(), indents) + indent + "\t},\n"
					indents--
				case sliceVal.Kind() == reflect.Ptr && sliceVal.Elem().Kind() == reflect.Struct:
					indents++
					value += indent + "\t" + "{\n" + getUnitTest(sliceVal.Interface(), indents) + indent + "\t},\n"
					indents--
				default:
					value += indent + "\t" + convertReadable(sliceVal.Interface(), false, true) + ",\n"
				}
			}
			value += "\t},\n"
		default:
			value += convertReadable(val.Interface(), false, true) + ",\n"
		}
	}

	return value
}

func Compare(object1, object2 interface{}, tag string) string {
	value := ""
	t, v1, ok := getReflections(object1)
	if !ok {
		v1 = reflect.New(t)
	}
	_, v2, ok := getReflections(object2)
	if !ok {
		v2 = reflect.New(t)
	}

	if v1.Kind() != reflect.Struct && (v1.Kind() != reflect.Ptr || v1.Elem().Kind() != reflect.Struct) {
		return value
	}
	for i := 0; i < v1.NumField(); i++ {
		test := t.Field(i)
		val1 := v1.Field(i)
		val2 := v2.Field(i)
		label, _ := test.Tag.Lookup(tag)
		if !val1.IsValid() {
			continue
		}
		if !val2.IsValid() {
			val2 = reflect.New(test.Type)
		}
		if val1.Kind() == reflect.Struct || (val1.Kind() == reflect.Ptr && val1.Elem().Kind() == reflect.Struct) {
			var newVal, newVal2 interface{}
			if val1.CanInterface() {
				newVal = val1.Interface()
			}
			if val2.CanInterface() {
				newVal2 = val2.Interface()
			}
			value += Compare(newVal, newVal2, tag)
		} else {
			if val1.Kind() == reflect.Slice || (val1.Kind() == reflect.Ptr && val1.Elem().Kind() == reflect.Slice) {
				for j := 0; j < val1.Len(); j++ {
					sliceVal1 := val1.Index(j)
					var (
						sliceVal2       reflect.Value
						newVal, newVal2 interface{}
					)
					if val2.Len() > j {
						sliceVal2 = val2.Index(j)
					}
					if sliceVal2.CanInterface() {
						newVal = sliceVal2.Interface()
					} else {
						newVal = sliceVal2
					}
					newVal2 = sliceVal1.Interface()
					if sliceVal1.Kind() == reflect.Struct || (sliceVal1.Kind() == reflect.Ptr && sliceVal1.Elem().Kind() == reflect.Struct) {
						value += Compare(newVal2, newVal, tag)
					} else if label != "" && newVal != newVal2 {
						value += fmt.Sprintf(label, convertReadable(newVal2, true, false), convertReadable(newVal, true, false)) + "\n"
					}
				}
			} else if label != "" {
				var newVal, newVal2 interface{}
				if val2.CanInterface() {
					newVal = val2.Interface()
				}
				newVal2 = val1.Interface()
				if newVal != newVal2 {
					value += fmt.Sprintf(label, convertReadable(newVal2, true, false), convertReadable(newVal, true, false)) + "\n"
				}
			}
		}
	}

	return value
}

func DeepCopy(fromObj, toObj interface{}) bool {
	t, v, _ := getReflections(fromObj)
	t2, _, _ := getReflections(toObj)
	if t.Name() != t2.Name() {
		return false
	}

	for i := 0; i < v.NumField(); i++ {
		test := t.Field(i)
		val := v.Field(i)
		if test.Type.Kind() == reflect.Struct || (test.Type.Kind() == reflect.Ptr && test.Type.Elem().Kind() == reflect.Struct) {
			if !val.IsZero() && val.CanInterface() {

			} else {

			}
		} else if val.CanInterface() {

		}
	}
	return true
}

func getReflections(object interface{}) (reflect.Type, reflect.Value, bool) {
	t := reflect.TypeOf(object)
	v := reflect.ValueOf(object)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
		if !v.IsValid() {
			return nil, reflect.Value{}, false
		}
	}
	return t, v, true
}

func convertReadable(elem interface{}, nilPtr, byType bool) string {
	switch t := elem.(type) {
	case bool:
		return strconv.FormatBool(t)
	case int:
		return strconv.Itoa(int(t))
	case int8:
		return strconv.Itoa(int(t))
	case int16:
		return strconv.Itoa(int(t))
	case int32:
		return strconv.Itoa(int(t))
	case int64:
		return strconv.Itoa(int(t))
	case uint:
		return strconv.FormatUint(uint64(t), 10)
	case uint8:
		return strconv.FormatUint(uint64(t), 10)
	case uint16:
		return strconv.FormatUint(uint64(t), 10)
	case uint32:
		return strconv.FormatUint(uint64(t), 10)
	case uint64:
		return strconv.FormatUint(uint64(t), 10)
	case float32:
		return strconv.FormatFloat(float64(t), 'G', -1, 64)
	case float64:
		return strconv.FormatFloat(float64(t), 'G', -1, 64)
	case string:
		return "\"" + t + "\""
	case time.Time:
		return "\"" + t.Format(time.RFC3339) + "\""
	case *int:
		if t != nil {
			return strconv.Itoa(*t)
		}
		if nilPtr {
			return "null"
		}
		return "0"
	case *int8:
		if t != nil {
			return strconv.Itoa(int(*t))
		}
		if nilPtr {
			return "null"
		}
		return "0"
	case *int16:
		if t != nil {
			return strconv.Itoa(int(*t))
		}
		if nilPtr {
			return "null"
		}
		return "0"
	case *int32:
		if t != nil {
			return strconv.Itoa(int(*t))
		}
		if nilPtr {
			return "null"
		}
		return "0"
	case *int64:
		if t != nil {
			return strconv.Itoa(int(*t))
		}
		if nilPtr {
			return "null"
		}
		return "0"
	case *uint:
		if t != nil {
			return strconv.FormatUint(uint64(*t), 10)
		}
		if nilPtr {
			return "null"
		}
		return "0"
	case *uint8:
		if t != nil {
			return strconv.FormatUint(uint64(*t), 10)
		}
		if nilPtr {
			return "null"
		}
		return "0"
	case *uint16:
		if t != nil {
			return strconv.FormatUint(uint64(*t), 10)
		}
		if nilPtr {
			return "null"
		}
		return "0"
	case *uint32:
		if t != nil {
			return strconv.FormatUint(uint64(*t), 10)
		}
		if nilPtr {
			return "null"
		}
		return "0"
	case *uint64:
		if t != nil {
			return strconv.FormatUint(*t, 10)
		}
		if nilPtr {
			return "null"
		}
		return "0"
	case *float32:
		if t != nil {
			return strconv.FormatFloat(float64(*t), 'G', -1, 64)
		}
		if nilPtr {
			return "null"
		}
		return "0.00"
	case *float64:
		if t != nil {
			return strconv.FormatFloat(float64(*t), 'G', -1, 64)
		}
		if nilPtr {
			return "null"
		}
		return "0.00"
	case *string:
		if t != nil {
			return "\"" + *t + "\""
		}
		if nilPtr {
			return "null"
		}
		return "\"\""
	case *time.Time:
		if t != nil {
			return "\"" + t.Format(time.RFC3339) + "\""
		}
		if nilPtr {
			return "null"
		}
		return "\"" + time.Now().Format(time.RFC3339) + "\""
	case reflect.Value:
		if nilPtr {
			return "null"
		}
		return t.String()
	default:
		if nilPtr {
			return "null"
		}
		if byType {
			return reflect.TypeOf(t).String()
		}
		return ""
	}
}
