package util

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"strconv"
	"strings"
	"unsafe"
)

// byte to string.
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// strings to bytes
func StringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}

func InArray(needle string, haystack []string) bool {
	if len(haystack) == 0 {
		return false
	}

	for _, v := range haystack {
		if strings.Compare(needle, v) == 0 {
			return true
		}
	}

	return false
}

// 判断数组是否存在
func Contains(array interface{}, val interface{}) (index int) {
	index = -1
	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(array)
		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) {
				index = i
				return
			}
		}
	}
	return
}

// 获取随机数
func RandNumber(number int) string {
	var r strings.Builder
	r.Grow(number)
	for i := 0; i < number; i++ {
		f, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			continue
		}

		r.WriteString(f.String())
	}

	return r.String()
}

func IntJoin(a []int, sep string) string {
	switch len(a) {
	case 0:
		return ""
	case 1:
		return strconv.Itoa(a[0])
	}
	n := len(sep) * (len(a) - 1)
	for i := 0; i < len(a); i++ {
		n += len(strconv.Itoa(a[i]))
	}

	var b strings.Builder
	b.Grow(n)
	b.WriteString(strconv.Itoa(a[0]))
	for _, s := range a[1:] {
		b.WriteString(sep)
		b.WriteString(strconv.Itoa(s))
	}

	return b.String()
}

func SplitToInt(s, sep string) []int64 {
	if len(s) == 0 {
		return []int64{}
	}
	arr := strings.Split(s, sep)
	target := make([]int64, 0, len(arr))
	for _, v := range arr {
		parse, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			continue
		}

		target = append(target, parse)
	}

	return target
}

func SnakeString(str string) string {
	data := make([]byte, 0, len(str)*2)
	j := false
	num := len(str)
	for i := 0; i < num; i++ {
		d := str[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}

		if d != '_' {
			j = true
		}

		data = append(data, d)
	}

	return strings.ToLower(string(data[:]))
}

func SliceToMap(param interface{}, key string, mapParam interface{}) error {
	paramV := reflect.ValueOf(param)
	mapArgV := reflect.ValueOf(mapParam)
	mapArgT := reflect.TypeOf(mapParam)
	if mapArgT.Kind() != reflect.Map {
		return errors.New("mapParam只接收map类型参数")
	}
	if paramV.Kind() != reflect.Slice {
		return errors.New("param只接收slice类型参数")
	}
	if mapArgV.IsNil() {
		return errors.New("mapParam请初始化后传入")
	}
	for p := 0; p < paramV.Len(); p++ {
		sliceV := paramV.Index(p)
		if sliceV.Kind() != reflect.Ptr {
			continue
		}
		slicePV := reflect.Indirect(sliceV)

		sliceChildPV := slicePV.FieldByName(key)
		if !sliceChildPV.IsValid() {
			return errors.New("未找到指定key参数")
		}
		//获取map的key类型
		structPT, _ := reflect.TypeOf(slicePV.Interface()).FieldByName(key)
		var (
			sKey    reflect.Value
			intData int64
		)
		switch structPT.Type.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			intData = sliceChildPV.Int()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			intData = int64(sliceChildPV.Uint())
		case reflect.Float32, reflect.Float64:
			intData = int64(sliceChildPV.Float())
		case reflect.String:
			intData, _ = strconv.ParseInt(sliceChildPV.String(), 10, 64)
		default:
			return errors.New("指定key类型无法转化map")
		}

		//判断mapParam的key类型
		switch mapArgT.Key().Kind() {
		case reflect.Int:
			sKey = reflect.ValueOf(int(intData))
		case reflect.Int8:
			sKey = reflect.ValueOf(int8(intData))
		case reflect.Int16:
			sKey = reflect.ValueOf(int16(intData))
		case reflect.Int32:
			sKey = reflect.ValueOf(int32(intData))
		case reflect.Int64:
			sKey = reflect.ValueOf(intData)
		case reflect.Uint:
			sKey = reflect.ValueOf(uint(intData))
		case reflect.Uint8:
			sKey = reflect.ValueOf(uint8(intData))
		case reflect.Uint16:
			sKey = reflect.ValueOf(uint16(intData))
		case reflect.Uint32:
			sKey = reflect.ValueOf(uint32(intData))
		case reflect.Uint64:
			sKey = reflect.ValueOf(uint64(intData))
		case reflect.String:
			sKey = reflect.ValueOf(strconv.FormatInt(intData, 10))
		case reflect.Float32:
			sKey = reflect.ValueOf(float32(intData))
		case reflect.Float64:
			sKey = reflect.ValueOf(float64(intData))
		}
		mapArgV.SetMapIndex(sKey, sliceV)
	}
	return nil
}

func ToString(intData interface{}) string {
	v := reflect.ValueOf(intData)
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(v.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%.6f", v.Float())
	case reflect.String:
		return v.String()
	default:
		return ""
	}
}

func UniqueSlice(data interface{}) error {
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Ptr {
		return errors.New("需传入地址参数")
	}
	sl := reflect.Indirect(v)
	if sl.Kind() != reflect.Slice {
		return errors.New("需传入切片类型参数")
	}
	if sl.Len() == 0 {
		return nil
	}
	slItem := sl.Index(0)
	var (
		uniqueMap reflect.Value
		mapType   reflect.Type
		newType   reflect.Type
		newValue  reflect.Value
	)
	newValue = reflect.ValueOf(true)
	switch slItem.Kind() {
	case reflect.Int:
		newType = reflect.TypeOf(make([]int, 0))
		mapType = reflect.TypeOf(make(map[int]bool))
	case reflect.Int8:
		newType = reflect.TypeOf(make([]int8, 0))
		mapType = reflect.TypeOf(make(map[int8]bool))
	case reflect.Int16:
		newType = reflect.TypeOf(make([]int16, 0))
		mapType = reflect.TypeOf(make(map[int16]bool))
	case reflect.Int32:
		newType = reflect.TypeOf(make([]int32, 0))
		mapType = reflect.TypeOf(make(map[int32]bool))
	case reflect.Int64:
		newType = reflect.TypeOf(make([]int64, 0))
		mapType = reflect.TypeOf(make(map[int64]bool))
	case reflect.Uint:
		newType = reflect.TypeOf(make([]uint, 0))
		mapType = reflect.TypeOf(make(map[uint]bool))
	case reflect.Uint8:
		newType = reflect.TypeOf(make([]uint8, 0))
		mapType = reflect.TypeOf(make(map[uint8]bool))
	case reflect.Uint16:
		newType = reflect.TypeOf(make([]uint16, 0))
		mapType = reflect.TypeOf(make(map[uint16]bool))
	case reflect.Uint32:
		newType = reflect.TypeOf(make([]uint32, 0))
		mapType = reflect.TypeOf(make(map[uint32]bool))
	case reflect.Uint64:
		newType = reflect.TypeOf(make([]uint64, 0))
		mapType = reflect.TypeOf(make(map[uint64]bool))
	case reflect.Float32:
		newType = reflect.TypeOf(make([]float32, 0))
		mapType = reflect.TypeOf(make(map[float32]bool))
	case reflect.Float64:
		newType = reflect.TypeOf(make([]float64, 0))
		mapType = reflect.TypeOf(make(map[float64]bool))
	case reflect.String:
		newType = reflect.TypeOf(make([]string, 0))
		mapType = reflect.TypeOf(make(map[string]bool))
	}
	uniqueMap = reflect.MakeMap(mapType)
	unique := reflect.MakeSlice(newType, 0, 0)
	for i := 0; i < sl.Len(); i++ {
		var isHas bool
		for _, x := range uniqueMap.MapKeys() {
			if x.Interface() == sl.Index(i).Interface() {
				isHas = true
			}
		}
		if !isHas {
			unique = reflect.Append(unique, sl.Index(i))
		}
		uniqueMap.SetMapIndex(sl.Index(i), newValue)
	}
	sl.Set(unique)
	return nil
}
