package grpcclient

import (
	"context"
	"testing"

	"github.com/eviltomorrow/rogue/app/email/pb"
	"github.com/eviltomorrow/rogue/lib/etcd"
	"github.com/eviltomorrow/rogue/lib/grpclb"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/resolver"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestPing(t *testing.T) {
	_assert := assert.New(t)

	etcd.Endpoints = []string{"127.0.0.1:2379"}
	client, err := etcd.NewClient()
	_assert.Nil(err)

	defer client.Close()

	resolver.Register(grpclb.NewBuilder(client))
	stub, closeFunc, err := NewEmail()
	_assert.Nil(err)
	defer closeFunc()

	resp, err := stub.Ping(context.Background(), &emptypb.Empty{})
	_assert.Nil(err)
	t.Log(resp.Value)
}

func TestSend(t *testing.T) {
	_assert := assert.New(t)

	etcd.Endpoints = []string{"127.0.0.1:2379"}
	client, err := etcd.NewClient()
	_assert.Nil(err)

	defer client.Close()

	resolver.Register(grpclb.NewBuilder(client))
	stub, closeFunc, err := NewEmail()
	_assert.Nil(err)
	defer closeFunc()

	_, err = stub.Send(context.Background(), &pb.Mail{
		To: []*pb.Contact{
			{Name: "eviltomorrow", Address: "eviltomorrow@163.com"},
		},
		Subject:     "This is one test",
		Body:        "Hello world",
		ContentType: pb.Mail_TEXT_PLAIN,
	})
	_assert.Nil(err)

}
