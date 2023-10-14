package xrealip

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type XRealIP struct {
	IPs string
}

func (xri XRealIP) SetXRealIPInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	updatedAuthCtx := metadata.AppendToOutgoingContext(ctx, "X-Real-IP", xri.IPs)
	return invoker(updatedAuthCtx, method, req, reply, cc, opts...)
}
