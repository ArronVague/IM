package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Bind 函数从 HTTP 请求中获取 Content-Type 头部的值。这个值表示了请求 Body 的媒体类型。
//
// 然后，函数检查 Content-Type 是否包含 "application/json"。如果包含，表示 Body 中的数据是 JSON 格式，然后调用 BindJson 函数将 Body 中的 JSON 数据解析并绑定到 obj 对象中。
func Bind(req *http.Request, obj interface{}) error {
	contentType := req.Header.Get("Content-Type")
	//如果是简单的json
	if strings.Contains(strings.ToLower(contentType), "application/json") {
		return BindJson(req, obj)
	}
	if strings.Contains(strings.ToLower(contentType), "application/x-www-form-urlencoded") {
		return BindForm(req, obj)
	}
	return errors.New("当前方法暂不支持")
}

// BindJson 函数接收同样的参数：一个 HTTP 请求 req 和一个空接口 obj。函数的目的是将 HTTP 请求的 Body 中的 JSON 数据解析并绑定到 obj 对象中。
func BindJson(req *http.Request, obj interface{}) error {
	s, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(s, obj)
	return err
}

// BindForm 将 HTTP 请求中的 Form 数据解析并绑定到 ptr 对象中。
// 首先，函数调用 req.ParseForm() 来解析 HTTP 请求中的 Form 数据。这个方法会解析 URL 中的查询字符串以及请求体中的 application/x-www-form-urlencoded 数据，并将解析后的数据存储在 req.Form 中。
// 然后，函数使用 fmt.Println(req.Form.Encode()) 打印出编码后的 Form 数据。req.Form.Encode() 会将 Form 数据编码为 URL 编码的字符串，每个键值对以 = 连接，键值对之间以 & 分隔。
// 接着，函数调用 mapForm(ptr, req.Form) 来将 Form 数据映射到 ptr 对象中。
func BindForm(req *http.Request, ptr interface{}) error {
	err := req.ParseForm()
	if err != nil {
		return err
	}
	fmt.Println(req.Form.Encode())
	err = mapForm(ptr, req.Form)
	return err
}

// 自动绑定方法
// 借鉴了gin
// 改良了时间绑定
func mapForm(ptr interface{}, form map[string][]string) error {
	//reflect.TypeOf(ptr) 返回 ptr 的类型，Elem() 方法返回这个类型的元素类型。例如，如果 ptr 是一个 *Person 类型（一个指向 Person 类型的指针），那么 reflect.TypeOf(ptr) 返回 *Person，reflect.TypeOf(ptr).Elem() 返回 Person。
	//
	//reflect.ValueOf(ptr) 返回 ptr 的值，Elem() 方法返回这个值的元素值。例如，如果 ptr 是一个指向 Person 实例的指针，那么 reflect.ValueOf(ptr) 返回这个指针，reflect.ValueOf(ptr).Elem() 返回这个指针指向的 Person 实例。
	typ := reflect.TypeOf(ptr).Elem()
	val := reflect.ValueOf(ptr).Elem()
	//使用反射来遍历 ptr 指向的结构体的所有字段，并获取每个字段的类型和值
	for i := 0; i < typ.NumField(); i++ {
		typeField := typ.Field(i)
		structField := val.Field(i)
		//如果字段是不导出的（即字段名的首字母是小写的），那么这个字段就不能设置，structField.CanSet() 会返回 false。
		if !structField.CanSet() {
			continue
		}

		structFieldKind := structField.Kind()
		inputFieldName := typeField.Tag.Get("form")
		if inputFieldName == "" {
			inputFieldName = typeField.Name

			//函数检查字段的类型。如果字段是一个结构体，那么函数递归调用 mapForm 来映射这个结构体的字段。
			// if "form" tag is nil, we inspect if the field is a struct.
			// this would not make sense for JSON parsing, but it does for a form
			// since data is flattened
			if structFieldKind == reflect.Struct {
				err := mapForm(structField.Addr().Interface(), form)
				if err != nil {
					return err
				}
				continue
			}
		}
		inputValue, exists := form[inputFieldName]
		if !exists {
			continue
		}

		numElems := len(inputValue)
		if structFieldKind == reflect.Slice && numElems > 0 {
			sliceOf := structField.Type().Elem().Kind()
			slice := reflect.MakeSlice(structField.Type(), numElems, numElems)
			for i := 0; i < numElems; i++ {
				if err := setWithProperType(sliceOf, inputValue[i], slice.Index(i)); err != nil {
					return err
				}
			}
			val.Field(i).Set(slice)
		} else {
			if _, isTime := structField.Interface().(time.Time); isTime {
				if err := setTimeField(inputValue[0], typeField, structField); err != nil {
					return err
				}
				continue
			}
			if err := setWithProperType(typeField.Type.Kind(), inputValue[0], structField); err != nil {
				return err
			}
		}
	}
	return nil
}

func setWithProperType(valueKind reflect.Kind, val string, structField reflect.Value) error {
	switch valueKind {
	case reflect.Int:
		return setIntField(val, 0, structField)
	case reflect.Int8:
		return setIntField(val, 8, structField)
	case reflect.Int16:
		return setIntField(val, 16, structField)
	case reflect.Int32:
		return setIntField(val, 32, structField)
	case reflect.Int64:
		return setIntField(val, 64, structField)
	case reflect.Uint:
		return setUintField(val, 0, structField)
	case reflect.Uint8:
		return setUintField(val, 8, structField)
	case reflect.Uint16:
		return setUintField(val, 16, structField)
	case reflect.Uint32:
		return setUintField(val, 32, structField)
	case reflect.Uint64:
		return setUintField(val, 64, structField)
	case reflect.Bool:
		return setBoolField(val, structField)
	case reflect.Float32:
		return setFloatField(val, 32, structField)
	case reflect.Float64:
		return setFloatField(val, 64, structField)
	case reflect.String:
		structField.SetString(val)
	default:
		return errors.New("unknown type")
	}
	return nil
}

func setIntField(val string, bitSize int, field reflect.Value) error {
	if val == "" {
		val = "0"
	}
	intVal, err := strconv.ParseInt(val, 10, bitSize)
	if err == nil {
		field.SetInt(intVal)
	}
	return err
}

func setUintField(val string, bitSize int, field reflect.Value) error {
	if val == "" {
		val = "0"
	}
	uintVal, err := strconv.ParseUint(val, 10, bitSize)
	if err == nil {
		field.SetUint(uintVal)
	}
	return err
}

func setBoolField(val string, field reflect.Value) error {
	if val == "" {
		val = "false"
	}
	boolVal, err := strconv.ParseBool(val)
	if err == nil {
		field.SetBool(boolVal)
	}
	return nil
}

func setFloatField(val string, bitSize int, field reflect.Value) error {
	if val == "" {
		val = "0.0"
	}
	floatVal, err := strconv.ParseFloat(val, bitSize)
	if err == nil {
		field.SetFloat(floatVal)
	}
	return err
}

func setTimeField(val string, structField reflect.StructField, value reflect.Value) error {
	timeFormat := structField.Tag.Get("time_format")
	//2018-01-02 01:02:03

	if timeFormat == "" {
		timeFormat = "2006-01-02 15:04:05"
		val = strings.Replace(val, "/", "-", -1)
		num := len(strings.Split(val, " "))
		if num == 1 {
			val = val + " 00:00:00"
		} else {
			//2018-01-02 00
			num = len(strings.Split(val, ":"))

			if num == 1 {
				val = val + ":00:00"
			} else if num == 2 {
				val = val + ":00"
			}
		}

	}

	if val == "" {
		value.Set(reflect.ValueOf(time.Time{}))
		return nil
	}

	l := time.Local
	if isUTC, _ := strconv.ParseBool(structField.Tag.Get("time_utc")); isUTC {
		l = time.UTC
	}

	if locTag := structField.Tag.Get("time_location"); locTag != "" {
		loc, err := time.LoadLocation(locTag)
		if err != nil {
			return err
		}
		l = loc
	}

	t, err := time.ParseInLocation(timeFormat, val, l)
	if err != nil {
		return err
	}

	value.Set(reflect.ValueOf(t))
	return nil
}
