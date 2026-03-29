package utils

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"time"
)

type Timestamp time.Time

// 实现 driver.Valuer 接口，用于将 Timestamp 转换为数据库可存储的值
func (t Timestamp) Value() (driver.Value, error) {
	// 转换为 time.Time 并格式化为数据库支持的时间格式
	return time.Time(t), nil
}

// 实现 sql.Scanner 接口，用于从数据库读取值并转换为 Timestamp
func (t *Timestamp) Scan(value interface{}) error {
	if value == nil {
		*t = Timestamp(time.Time{})
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		*t = Timestamp(v)
		return nil
	default:
		return fmt.Errorf("cannot scan type %T into Timestamp", value)
	}
}

// MarshalJSON：将 Timestamp 序列化为毫秒级时间戳
func (t Timestamp) MarshalJSON() ([]byte, error) {
	// 如果是零值，返回 null
	if time.Time(t).IsZero() {
		return []byte("null"), nil
	}

	// 这里使用 UnixNano() 方法获取纳秒级别时间戳，然后除以 1e6 转换为毫秒级别
	ts := time.Time(t).UnixNano() / 1e6 // 毫秒
	return []byte(strconv.FormatInt(ts, 10)), nil
}

// UnmarshalJSON：反序列化前端传的 13 位毫秒时间戳
func (t *Timestamp) UnmarshalJSON(b []byte) error {
	s := string(b)
	if s == "null" || s == "" {
		*t = Timestamp(time.Time{})
		return nil
	}
	// 解析 13 位时间戳
	ts, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return err
	}
	*t = Timestamp(time.Unix(0, ts*1e6))
	return nil
}

// ToTime: 辅助方法，取出标准的 time.Time
func (t Timestamp) ToTime() time.Time {
	return time.Time(t)
}

// GenerateID 生成唯一ID
func GenerateID() string {
	return fmt.Sprintf("%d%d", time.Now().UnixNano(), time.Now().Nanosecond())
}
