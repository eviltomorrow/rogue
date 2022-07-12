package etcd

import (
	"context"
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

var (
	Endpoints      []string
	ConnectTimeout = 5 * time.Second
	ExecuteTimeout = 20 * time.Second
)

func NewClient() (*clientv3.Client, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   Endpoints,
		DialTimeout: ConnectTimeout,
		LogConfig: &zap.Config{
			Level:            zap.NewAtomicLevelAt(zap.ErrorLevel),
			Development:      false,
			Encoding:         "json",
			EncoderConfig:    zap.NewProductionEncoderConfig(),
			OutputPaths:      []string{"stderr"},
			ErrorOutputPaths: []string{"stderr"},
		},
	})
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), ConnectTimeout)
	defer cancel()

	for i, endpoint := range Endpoints {
		_, err = client.Status(ctx, endpoint)
		if err != nil {
			if i == len(Endpoints)-1 {
				return nil, fmt.Errorf("connect to etcd service failure, nest error: no valid endpoint, endpoints: %v", Endpoints)
			}
		} else {
			break
		}
	}
	return client, nil
}
