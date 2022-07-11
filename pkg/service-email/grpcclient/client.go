package grpcclient

import (
	"context"
	"fmt"
	"time"

	"github.com/eviltomorrow/rogue/app/email/pb"
	"github.com/eviltomorrow/rogue/lib/self"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	ConnectTimeout = 5 * time.Second
)

func New() (pb.EmailClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), ConnectTimeout)
	defer cancel()

	var target = "etcd:///grpclb/" + self.ServiceName
	conn, err := grpc.DialContext(
		ctx,
		target,
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, err
	}
	return pb.NewEmailClient(conn), nil
}
