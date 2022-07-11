package grpcclient

import (
	"context"
	"testing"

	"github.com/eviltomorrow/rogue/lib/etcd"
	"github.com/eviltomorrow/rogue/lib/grpclb"
	"github.com/eviltomorrow/rogue/lib/self"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/resolver"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestPing(t *testing.T) {
	_assert := assert.New(t)

	self.ServiceName = "rogue-email"
	etcd.Endpoints = []string{"127.0.0.1:2379"}
	client, err := etcd.NewClient()
	_assert.Nil(err)

	defer client.Close()

	resolver.Register(grpclb.NewBuilder(client))
	stub, err := New()
	_assert.Nil(err)

	resp, err := stub.Ping(context.Background(), &emptypb.Empty{})
	_assert.Nil(err)
	t.Log(resp.Value)
}
