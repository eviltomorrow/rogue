package repo

import (
	"strconv"

	"github.com/eviltomorrow/rogue/lib/zlog"
	"go.uber.org/zap"
)

func atof64(name string, loc int, val string) float64 {
	f64, err := strconv.ParseFloat(val, 64)
	if err != nil {
		zlog.Error("ParseFloat64 failure", zap.String("name", name), zap.Int("loc", loc), zap.String("val", val))
		return 0
	}
	return f64
}

func atou64(name string, loc int, val string) uint64 {
	u64, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		zlog.Error("ParseUint64 failure", zap.String("name", name), zap.Int("loc", loc), zap.String("val", val))
		return 0
	}
	return u64
}
