package goToolsUtils

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"
)

func TickToStr(intTime int64)string  {
	length := len(strconv.FormatInt(intTime,10))
	part:=int64(math.Pow10(length-10))
	if part ==0{
		fmt.Printf("intTime:%v\n", intTime)
		fmt.Printf("length-10%v\n",length-10)

	}
	ret :=time.Unix(intTime/part, intTime%part).Format("2006-01-02 15:04:05")
	return ret
}

func StrToTick(strTime string)(int64,error)  {
	var err error
	var timeTemplate string
	if len(strTime) == 19{
		timeTemplate = "2006-01-02 15:04:05"
	}else if len(strTime) == 16{
		timeTemplate = "2006-01-02 15:04"
	}else if len(strTime) == 10{
		timeTemplate = "2006-01-02"
	}else {
		err = errors.New("invalid strTime:" + strTime)
	}
	stamp, err := time.ParseInLocation(timeTemplate, strTime, time.Local) //使用parseInLocation将字符串格式化返回本地时区时间
	if err != nil{
		err = errors.New("time.ParseInLocation failed :"+ err.Error())
	}

	//log.Println(stamp.Unix())  //输出：1546926630
	return stamp.Unix(),err
}

func GetFromTick(intTime int64, targetType string)(int,error)  {
	length := len(strconv.FormatInt(intTime,10))
	part:=int64(math.Pow10(length-10))
	ret :=time.Unix(intTime/part, intTime%part)
	var output int
	var err error
	switch targetType {
	case "year":
		output = ret.Year()
	case "month":
		output = int(ret.Month())
	case "day":
		output = ret.Day()
	case "hour":
		output = ret.Hour()
	case "minute":
		output = ret.Minute()
	case "second":
		output = ret.Second()
	default:
		err = errors.New(fmt.Sprintf("param targetType of GetFromTick is out of enum:%v", targetType))
	}
	return output,err
}

//将时间进一位到 ceilTo: option["minute","hour","day","month","year"]
func TimeCeil(tick int64, ceilTo string)int64  {
	if len(strconv.FormatInt(tick,10))!=10{
		panic("TimeCeil(tick int64 should be 10")
	}
	var output = tick
	TimeStr := TickToStr(tick)
	switch ceilTo {
	case "minute":
		if TimeStr[17:19] != "00"{
			TimeStr = TimeStr[:15] +":00"
			output,_ = StrToTick(TimeStr)
			output +=  1*60
		}
	case "hour":
		if !(TimeStr[14:16] == "00" && TimeStr[17:19] == "00"){
			TimeStr = TimeStr[:13] +":00:00"
			output,_ = StrToTick(TimeStr)
			output +=  60*60
		}
	default:
		panic("unfinished about day month week year")
	}
	return output
}

//将时间退一位到 ceilTo: option["minute","hour","day","month","year"]
func TimeFloor(tick int64, ceilTo string)int64  {
	if len(strconv.FormatInt(tick,10))!=10{
		panic("TimeCeil(tick int64 should be 10")
	}
	var output int64 = tick
	TimeStr := TickToStr(tick)
	switch ceilTo {
	case "minute":
		if TimeStr[17:19] != "00"{
			TimeStr = TimeStr[:15] +":00"
			output,_ = StrToTick(TimeStr)
		}
	case "hour":
		if TimeStr[14:16] != "00" || TimeStr[17:19] != "00"{
			TimeStr = TimeStr[:13] +":00:00"
			output,_ = StrToTick(TimeStr)
		}
	default:
		panic("unfinished about day month week year")
	}
	return output
}