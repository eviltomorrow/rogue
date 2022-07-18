package grpcclient

import (
	"context"
	"fmt"
	"time"

	"github.com/eviltomorrow/rogue/app/email/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	ConnectTimeout = 5 * time.Second
	ServiceName    = "rogue-email"
)

func NewEmail() (pb.EmailClient, func() error, error) {
	ctx, cancel := context.WithTimeout(context.Background(), ConnectTimeout)
	defer cancel()

	var target = "etcd:///grpclb/" + ServiceName
	conn, err := grpc.DialContext(
		ctx,
		target,
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, nil, err
	}
	return pb.NewEmailClient(conn), func() error { return conn.Close() }, nil
}
