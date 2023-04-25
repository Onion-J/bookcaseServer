package model

import (
	"database/sql/driver"
	"fmt"
	"time"
)

const timeFormat = "2006-01-02 15:04:05"
const timezone = "Asia/Shanghai"

type Time time.Time // Time 是与time.Time一样的结构体

// MarshalJSON 与 time.Now().MarshalJSON() 相比
// [34 50 48 50 50 45 48 54 45 50 52 84 49 52 58 51 55 58 48 57 46 57 55 57 51 57 52 50 43 48 56 58 48 48 34] <nil>
// [34 50 48 50 50 45 48 54 45 50 52 32 49 52 58 51 55 58 49 48 34] <nil>
// 方法作用于值类型接收者，Time本身没有改变，修改了Time的副本
func (t Time) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(timeFormat)+2) // 创建一个长度为21的空的byte型切片
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, timeFormat)
	b = append(b, '"')
	return b, nil
}

// UnmarshalJSON 方法作用于指针类型接收者，Time本身改变
// time.Parse("2006-01-02 15:04:05", "2022-06-24 14:37:10")
// time.ParseInLocation("2006-01-02 15:04:05", "2022-06-24 14:37:10", time.Local)
// 第一和第二个参数都必须为字符串，且第一个参数必须为"2006-01-02 15:04:05"
// time.Parse() 默认时区UTC,可以使用-700 MST 转换，time.Parse("2006-01-02 15:04:05 -0700 MST","2022-06-24 14:37:10 +0800 CST")
// time.ParseInLocation() 第三个参数必须为 time.Local
// time.Parse()和time.ParseInLocation()的参数都可以只有日期"2006-01-02" "2022-06-24" 则返回的时间为 "00:00:00"
// 2022-06-24 14:37:10 +0800 CST
func (t *Time) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(timeFormat, string(data), time.Local)
	*t = Time(now)
	return
}

// 与 time.Now().String() 相比
// 2022-06-24 14:37:10.0106214 +0800 CST m=+0.038902301
// 2022-06-24 14:37:10
// 转换成JSON的时间格式
func (t Time) String() string {
	return time.Time(t).Format(timeFormat)
}

// Local 2022-06-24 14:17:04.9972051 +0800 CST
func (t Time) local() time.Time {
	loc, _ := time.LoadLocation(timezone)
	return time.Time(t).In(loc)
}

// Value 2022-06-24 14:37:10.0128096 +0800 CST m=+0.041090501 <nil>
func (t Time) Value() (driver.Value, error) {
	var zeroTime time.Time
	var ti = time.Time(t)
	if ti.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return ti, nil
}

// Scan 实现sql里的接口
func (t *Time) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = Time(value)
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}
