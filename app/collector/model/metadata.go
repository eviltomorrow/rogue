package model

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	databaseName   = "rogue"
	collectionName = "metadata"
)

func InsertMetadataMany(db *mongo.Client, source string, metadata []*Metadata, timeout time.Duration) (int64, error) {
	if len(metadata) == 0 {
		return 0, nil
	}
	var collection = db.Database(databaseName).Collection(collectionName)
	var data = make([]interface{}, 0, len(metadata))
	for _, md := range metadata {
		data = append(data, bson.M{
			"source":           source,
			"code":             md.Code,
			"name":             md.Name,
			"open":             md.Open,
			"yesterday_closed": md.YesterdayClosed,
			"high":             md.High,
			"low":              md.Low,
			"latest":           md.Latest,
			"volume":           md.Volume,
			"account":          md.Account,
			"date":             md.Date,
			"time":             md.Time,
			"suspend":          md.Suspend,
			"create_timestamp": time.Now().Unix(),
			"modify_timestamp": 0,
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	result, err := collection.InsertMany(ctx, data)
	if err != nil {
		return 0, err
	}
	if result == nil {
		return 0, fmt.Errorf("panic: InsertMany result is nil")
	}
	return int64(len(result.InsertedIDs)), nil
}

func DeleteMetadataByDate(db *mongo.Client, source string, code, date string, timeout time.Duration) (int64, error) {
	if date == "" {
		return 0, nil
	}
	var collection = db.Database(databaseName).Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	result, err := collection.DeleteMany(ctx, bson.M{
		"source": source,
		"code":   code,
		"date":   date,
	})
	if err != nil {
		return 0, err
	}
	if result == nil {
		return 0, fmt.Errorf("panic: DeleteMany result is nil")
	}
	return result.DeletedCount, nil
}

func SelectMetadataRange(db *mongo.Client, offset, limit int64, source string, date string, lastID string, timeout time.Duration) ([]*Metadata, error) {
	if date == "" {
		return nil, fmt.Errorf("invalid date")
	}
	if limit <= 0 {
		return nil, fmt.Errorf("invalid limit")
	}

	ctx, cannel := context.WithTimeout(context.Background(), timeout)
	defer cannel()

	var collection = db.Database(databaseName).Collection(collectionName)
	var opt = &options.FindOptions{}
	opt.SetLimit(limit)

	var filter = bson.M{
		"source": source,
		"date":   date,
	}
	if lastID != "" {
		objectID, err := primitive.ObjectIDFromHex(lastID)
		if err != nil {
			return nil, err
		}
		filter = bson.M{"_id": bson.M{"$gt": objectID}, "date": date}
	} else {
		opt.SetSkip(offset)
	}

	cur, err := collection.Find(ctx, filter, opt)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	var data = make([]*Metadata, 0, limit)
	for cur.Next(context.Background()) {
		var m = &Metadata{}
		if err := cur.Decode(m); err != nil {
			return nil, err
		}
		data = append(data, m)
	}

	if err = cur.Err(); err != nil {
		return nil, err
	}
	return data, nil
}

// Metadata trade data
type Metadata struct {
	ObjectID        string  `json:"_id" bson:"_id"`
	Source          string  `json:"source" bson:"source"`
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
