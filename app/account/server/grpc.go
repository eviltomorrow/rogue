package server

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/eviltomorrow/rogue/app/account/pb"
	"github.com/eviltomorrow/rogue/lib/etcd"
	"github.com/eviltomorrow/rogue/lib/grpcmiddleware"
	"github.com/eviltomorrow/rogue/lib/self"
	"github.com/eviltomorrow/rogue/lib/util"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type GRPC struct {
	Client     *clientv3.Client
	ctx        context.Context
	cancel     func()
	server     *grpc.Server
	revokeFunc func() error

	pb.UnimplementedAccountServer
}

// Add(context.Context, *User) (*wrapperspb.StringValue, error)
// Del(context.Context, *wrapperspb.StringValue) (*emptypb.Empty, error)
// List(*emptypb.Empty, Account_ListServer) error
// Find(context.Context, *wrapperspb.StringValue) (*User, error)

func (g *GRPC) Add(ctx context.Context, req *pb.User) (*wrapperspb.StringValue, error) {
	return nil, nil
}

func (g *GRPC) Del(ctx context.Context, req *wrapperspb.StringValue) (*emptypb.Empty, error) {
	return nil, nil
}
func (g *GRPC) List(_ *emptypb.Empty, ls pb.Account_ListServer) error {
	return nil
}

func (g *GRPC) Find(ctx context.Context, req *wrapperspb.StringValue) (*pb.User, error) {
	return nil, nil
}

func (g *GRPC) StartupGRPC() error {
	port, err := util.GetAvailablePort()
	if err != nil {
		return err
	}
	localIP, err := util.GetLocalIP2()
	if err != nil {
		return err
	}

	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", localIP, port))
	if err != nil {
		return err
	}

	g.server = grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpcmiddleware.UnaryServerRecoveryInterceptor,
			grpcmiddleware.UnaryServerLogInterceptor,
		),
		grpc.ChainStreamInterceptor(
			grpcmiddleware.StreamServerRecoveryInterceptor,
			grpcmiddleware.StreamServerLogInterceptor,
		),
	)

	g.ctx, g.cancel = context.WithCancel(context.Background())

	g.revokeFunc, err = etcd.RegisterService(g.ctx, self.ServiceName, localIP, port, 10, g.Client)
	if err != nil {
		return err
	}

	reflection.Register(g.server)
	pb.RegisterAccountServer(g.server, g)
	go func() {
		if err := g.server.Serve(listen); err != nil {
			log.Fatalf("[F] GRPC Server startup failure, nest error: %v", err)
		}
	}()
	return nil
}

func (g *GRPC) ShutdownGRPC() error {
	g.cancel()

	if g.revokeFunc != nil {
		g.revokeFunc()
	}

	if g.server != nil {
		g.server.Stop()
	}
	return nil
}
