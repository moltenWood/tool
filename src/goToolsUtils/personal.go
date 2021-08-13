package goToolsUtils

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"io"
	"math"
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"
)

//用字符串从结构体中获得某属性的值
func GetValueFromStructByStr(v interface{}, name string) (interface{}, error) {
	if v == nil {
		panic("v cannot be nil")
	}
	if reflect.TypeOf(v).Kind() != reflect.Ptr {
		panic("v must be ptr")
	}

	value := reflect.ValueOf(v)

	ref := value.Elem()
	field := ref.FieldByName(name)
	if !field.IsValid() {
		return nil, errors.New(fmt.Sprintf("%s not exits", name))
	}
	return field.Interface(), nil
}

////用字符串给结构体赋值 必须要用 unsafe 还是算了
func SetValueToStructByStr(v interface{}, name string, ) {
	//if v == nil {
	//	panic("v cannot be nil")
	//}
	//if reflect.TypeOf(v).Kind() != reflect.Ptr {
	//	panic("v must be ptr")
	//}
	//
	//value := reflect.ValueOf(v)
	//
	//ref := value.Elem()
	//field := ref.FieldByName(name)
	//field.Interface() =
	//if !field.IsValid() {
	//	return nil, errors.New(fmt.Sprintf("%s not exits", name))
	//}
	//return field.Interface(), nil

	//参考这个 必须要用 unsafe 还是算了
	//if age.Kind()==reflect.Int{
	//	*(*int)(unsafe.Pointer(age.Addr().Pointer())) =24
	//}
}

//判断字符串是否为数字，支持小数点和负号，仅支持十进制数，数字内不能有空格
func IsNumber(target string) bool {
	target = strings.TrimSpace(target)
	if strings.HasPrefix(target, "-") {
		target = target[1:]
	}
	pointCount := 0
	isNumber := true
	for _, v := range target {
		if v == 46 {
			if pointCount != 0 {
				isNumber = false
				break
			}
			pointCount += 1
		} else if !unicode.IsDigit(v) {
			isNumber = false
			break
		}
	}
	return isNumber
}

//结构体转map[string]string
func StructConvertMapByTag(obj interface{}) map[string]string {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]string)
	for i := 0; i < t.NumField(); i++ {
		filed := t.Field(i)
		tag := filed.Tag
		key := ""
		if tag != "" {
			key = tag.Get("columnName")
			if key == "_" {
				continue
			}
		}
		if key == "" {
			key = filed.Name
		}
		head := strings.Split(filed.Name, "")[0]
		if strings.ToUpper(head) == head {
			value := v.Field(i).Interface()
			data[key] = StringAllType(value)

		}
	}
	return data
}

func StringAllType(value interface{}) string {
	var key string
	if value == nil {
		return key
	}

	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	case decimal.Decimal:
		key = value.(decimal.Decimal).String()
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}

	return key
}

func WriteFile(filename string, data []byte, perm os.FileMode) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	n, err := f.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}

// 获得后缀的数字
func GetSuffixInt(str string) (int, error) {
	ru := []rune(str)
	var position int
	for i := len(ru) - 1; i >= 0; i-- {
		if unicode.IsDigit(ru[i]) == false {
			position = i
			break
		}
	}
	output, err := strconv.Atoi(str[position+1:])
	return output, err
}

// 强制舍弃尾数
func FormatFloatFloor(num float64, decimal int) (float64, error) {
	// 默认乘1
	d := float64(1)
	if decimal > 0 {
		// 10的N次方
		d = math.Pow10(decimal)
	}
	// math.trunc作用就是返回浮点数的整数部分
	// 再除回去，小数点后无效的0也就不存在了
	res := strconv.FormatFloat(math.Floor(num*d)/d, 'f', -1, 64)
	return strconv.ParseFloat(res, 64)
}

func GetRandomString(length int) string {
	rand.Seed(time.Now().Unix())
	str := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func GetGCDFromMany(many []int) int {
	var output int = many[0]
	for i := 1; i < len(many); i++ {
		output = Gcd(many[i], output)
	}
	return output
}

//最大公约数
func Gcd(x, y int) int {
	var tmp int
	for {
		tmp = (x % y)
		if tmp > 0 {
			x = y
			y = tmp
		} else {
			return y
		}
	}
}

//最小共倍数
func Lcm(x, y int) int {
	return x * y / Gcd(x, y)
}

//获得float有几位有效小数位
func GetDecimalCount(num float64) int {
	splitRet := strings.Split(fmt.Sprintf("%v", num), ".")
	if len(splitRet) == 1 {
		return 0
	} else {
		return len(strings.Split(fmt.Sprintf("%v", num), ".")[1])
	}
}

// string add string 转成decimal后相加 最后 在转成string
func StringAddString(param0 string, param1 string) string {
	param0d, err := decimal.NewFromString(param0)
	if err != nil {
		panic("StringAddString decimal.NewFromString failed param0:" + param0)
	}
	param1d, err := decimal.NewFromString(param1)
	if err != nil {
		panic("StringAddString decimal.NewFromString failed param1d:" + param1)
	}
	return param0d.Add(param1d).String()
}
