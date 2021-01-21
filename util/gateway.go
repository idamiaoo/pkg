package util

import (
	"github.com/gin-gonic/gin"
	"reflect"
	"strconv"
)

// 获取参数 从网关获取的user_id 和school_id的参数
func BindJson(c *gin.Context, param interface{}) error {
	var err error
	switch c.Request.Method {
	case "POST", "PUT", "PATCH", "DELETE":
		err = c.ShouldBindJSON(param)
	case "GET":
		err = c.ShouldBindQuery(param)
	}

	if err != nil {
		return err
	}

	return bindGatewayParam(c, param)
}

func bindGatewayParam(c *gin.Context, param interface{}) error {
	objV := reflect.ValueOf(param)
	objT := reflect.TypeOf(param)
	if objV.Kind() == reflect.Ptr {
		objV = objV.Elem()
		objT = objT.Elem()
	}

	if objV.Kind() != reflect.Struct {
		return nil
	}

	for i := 0; i < objV.NumField(); i++ {
		fieldV := objV.Field(i)
		fieldT := objT.Field(i)
		if fieldV.Kind() == reflect.Struct {
			for j := 0; j < fieldV.NumField(); j++ {
				setValue(c, fieldV.Field(j), reflect.TypeOf(fieldV.Interface()).Field(j))
			}
		} else if fieldV.Kind() == reflect.Ptr {
			if fieldV.IsNil() {
				fieldV.Set(reflect.New(fieldV.Type().Elem()))
			}
			fieldV = reflect.Indirect(fieldV)
			if fieldV.Kind() != reflect.Struct {
				continue
			}
			for j := 0; j < fieldV.NumField(); j++ {
				setValue(c, fieldV.Field(j), reflect.TypeOf(fieldV.Interface()).Field(j))
			}
		} else {
			setValue(c, fieldV, fieldT)
		}
	}

	return nil
}

func GetUserId(c *gin.Context) int64 {
	id, ok := c.Get("id")
	if !ok {
		return 0
	}

	return reflectId(id)
}

func GetSchoolId(c *gin.Context) int64 {
	id, ok := c.Get("school_id")
	if !ok {
		return 0
	}

	return reflectId(id)
}

func reflectId(id interface{}) int64 {
	v := reflect.ValueOf(id)
	switch v.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int64(v.Uint())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int()
	case reflect.String:
		i, _ := strconv.ParseInt(v.String(), 10, 64)
		return i
	case reflect.Float32, reflect.Float64:
		return int64(v.Float())
	default:
		return 0
	}
}

func setValue(c *gin.Context, value reflect.Value, t reflect.StructField) {
	var data int64
	key := t.Tag.Get("json")
	switch key {
	case "school_id":
		data = GetSchoolId(c)
	case "user_id":
		data = GetUserId(c)
	default:
		return
	}
	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value.SetInt(data)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		value.SetUint(uint64(data))
	case reflect.String:
		value.SetString(strconv.Itoa(int(data)))
	}
}
