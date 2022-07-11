package model

import (
	"encoding/json"
	"fmt"
)

// Metadata trade data
type Metadata struct {
	ObjectID        string  `json:"_id" bson:"_id"`
	Code            string  `json:"code" bson:"code"`
	Name            string  `json:"name" bson:"name"`                         // 0 股票简称
	Open            float64 `json:"open" bson:"open"`                         // 1 今日开盘价格
	YesterdayClosed float64 `json:"yesterday_closed" bson:"yesterday_closed"` // 2 昨日收盘价格
	Latest          float64 `json:"latest" bson:"latest"`                     // 3 最近成交价格
	High            float64 `json:"high" bson:"high"`                         // 4 最高成交价
	Low             float64 `json:"low" bson:"low"`                           // 5 最低成交价
	Volume          uint64  `json:"volume" bson:"volume"`                     // 8 成交数量（股）
	Account         float64 `json:"account" bson:"account"`                   // 9 成交金额（元）
	Date            string  `json:"date" bson:"date"`                         // 30 日期
	Time            string  `json:"time" bson:"time"`                         // 31 时间
	Suspend         string  `json:"suspend" bson:"suspend"`                   // 32 停牌状态
}

func (m *Metadata) String() string {
	buf, err := json.Marshal(m)
	if err != nil {
		return fmt.Sprintf("Metadata marshal json failure, nest error: %v", err)
	}
	return string(buf)
}
