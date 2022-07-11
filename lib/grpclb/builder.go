package grpclb

import (
	"context"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/status"
)

type builder struct {
	c *clientv3.Client
}

func NewBuilder(client *clientv3.Client) resolver.Builder {
	return &builder{c: client}
}

func (b builder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &Resolver{
		c:      b.c,
		target: target.URL.Path,
		cc:     cc,
	}
	r.ctx, r.cancel = context.WithCancel(context.Background())

	em, err := endpoints.NewManager(r.c, r.target)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "resolver: failed to new endpoint manager: %s", err)
	}
	r.wch, err = em.NewWatchChannel(r.ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "resolver: failed to new watch channer: %s", err)
	}

	r.wg.Add(1)
	go r.watch()
	return r, nil
}

func (b *builder) Scheme() string {
	return "etcd"
}

func (b *builder) Close() error {
	return nil
}
