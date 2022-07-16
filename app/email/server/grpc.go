package server

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/eviltomorrow/rogue/app/email/conf"
	"github.com/eviltomorrow/rogue/app/email/pb"
	"github.com/eviltomorrow/rogue/lib/etcd"
	"github.com/eviltomorrow/rogue/lib/grpcmiddleware"
	"github.com/eviltomorrow/rogue/lib/self"
	"github.com/eviltomorrow/rogue/lib/smtp"
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
	SMTP       *conf.SMTP

	pb.UnimplementedEmailServer
}

// Send(context.Context, *Mail) (*emptypb.Empty, error)
// Ping(context.Context, *emptypb.Empty) (*wrapperspb.StringValue, error)

func (g *GRPC) Send(ctx context.Context, req *pb.Mail) (*emptypb.Empty, error) {
	if req == nil {
		return nil, fmt.Errorf("illegal parameter, nest error: mail is nil")
	}
	if len(req.To) == 0 {
		return nil, fmt.Errorf("illegal parameter, nest error: to is nil")
	}

	var contentType = smtp.TextHTML
	switch req.ContentType {
	case pb.Mail_TEXT_PLAIN:
		contentType = smtp.TextPlain
	default:
	}
	var message = &smtp.Message{
		From: smtp.Contact{
			Name:    g.SMTP.Alias,
			Address: g.SMTP.Username,
		},
		Subject:     req.Subject,
		Body:        req.Body,
		ContentType: contentType,
	}

	var to = make([]smtp.Contact, 0, len(req.To))
	for _, c := range req.To {
		if c != nil {
			to = append(to, smtp.Contact{Name: c.Name, Address: c.Address})
		}
	}
	message.To = to

	var cc = make([]smtp.Contact, 0, len(req.Cc))
	for _, c := range req.Cc {
		if c != nil {
			cc = append(cc, smtp.Contact{Name: c.Name, Address: c.Address})
		}
	}
	message.Cc = cc

	var bcc = make([]smtp.Contact, 0, len(req.Bcc))
	for _, c := range req.Bcc {
		if c != nil {
			bcc = append(bcc, smtp.Contact{Name: c.Name, Address: c.Address})
		}
	}
	message.Bcc = bcc

	if err := smtp.SendWithSSL(g.SMTP.Server, g.SMTP.Username, g.SMTP.Password, message); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (g *GRPC) Ping(context.Context, *emptypb.Empty) (*wrapperspb.StringValue, error) {
	return &wrapperspb.StringValue{Value: "pong"}, nil
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
	pb.RegisterEmailServer(g.server, g)
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
