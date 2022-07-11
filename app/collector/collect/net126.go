package collect

import (
	"fmt"
	"strings"
	"time"

	"github.com/eviltomorrow/rogue/lib/httpclient"
	"github.com/eviltomorrow/rogue/lib/model"
	"github.com/eviltomorrow/rogue/lib/zlog"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
)

var Net126Header = map[string]string{
	"Connection":                "keep-alive",
	"Cache-Control":             "max-age=0",
	"Upgrade-Insecure-Requests": "1",
	"User-Agent":                "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.99 Safari/537.36",
	"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
	"Accept-Encoding":           "gzip, deflate",
	"Accept-Language":           "zh-CN,zh;q=0.9,en;q=0.8,da;q=0.7,pt;q=0.6,ja;q=0.5",
}

func FetchMetadataFromNet126(codes []string) ([]*model.Metadata, error) {
	var c = make([]string, 0, len(codes))
	for _, code := range codes {
		if strings.HasPrefix(code, "sh") {
			c = append(c, strings.ReplaceAll(code, "sh", "0"))
		}
		if strings.HasPrefix(code, "sz") {
			c = append(c, strings.ReplaceAll(code, "sz", "1"))
		}
	}
	if len(c) == 0 {
		return nil, nil
	}

	var url = fmt.Sprintf("https://api.money.126.net/data/feed/%s,money.api", strings.Join(c, ","))
	data, err := httpclient.Get(url, 20*time.Second, Net126Header)
	if err != nil {
		return nil, fmt.Errorf("url: %v, nest error: %v", url, err)
	}

	data = strings.TrimPrefix(data, "_ntes_quote_callback(")
	data = strings.TrimSuffix(data, ");")
	data = strings.TrimSpace(data)

	var result = make([]*model.Metadata, 0, len(codes))
	kv, err := parseNet126DataToMap(data)
	if err != nil {
		zlog.Error("parseNet126DataToMap failure", zap.String("data", data), zap.Error(err))
	}
	for key, val := range kv {
		metadata, err := parseNet126LineToMetadata(key, val)
		if err != nil {
			zlog.Error("parseNet126LineToMetadata failure", zap.String("val", val.String()), zap.Error(err))
		}
		if metadata != nil {
			result = append(result, metadata)
		}
	}
	return result, nil
}

func parseNet126DataToMap(data string) (map[string]quote, error) {
	var result = make(map[string]quote, 32)
	if err := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal([]byte(data), &result); err != nil {
		return nil, err
	}
	return result, nil
}

func parseNet126LineToMetadata(code string, quote quote) (*model.Metadata, error) {
	if quote.Volume == 0 {
		return nil, nil
	}

	t, err := time.ParseInLocation("2006/01/02 15:04:05", quote.Time, time.Local)
	if err != nil {
		return nil, err
	}
	var result = &model.Metadata{
		Code:            fmt.Sprintf("%s%s", strings.ToLower(quote.Type), quote.Symbol),
		Name:            quote.Name,
		Open:            quote.Open,
		YesterdayClosed: quote.YestClose,
		Latest:          quote.Price,
		High:            quote.High,
		Low:             quote.Low,
		Volume:          quote.Volume,
		Account:         quote.Turnover,
		Date:            t.Format("2006-01-02"),
		Time:            t.Format("15:04:05"),
		Suspend:         suspendNormal,
	}
	return result, nil
}

type quote struct {
	Code      string  `json:"code"`
	Name      string  `json:"name"`
	Open      float64 `json:"open"`
	YestClose float64 `json:"yestclose"`
	Price     float64 `json:"price"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
	Volume    uint64  `json:"volume"`
	Turnover  float64 `json:"turnover"`
	Time      string  `json:"time"`
	Type      string  `json:"type"`
	Symbol    string  `json:"symbol"`
}

func (q quote) String() string {
	buf, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(&q)
	return string(buf)
}
