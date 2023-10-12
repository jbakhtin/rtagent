package trustedsubnets

import (
	"context"
	"github.com/go-faster/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"strings"
)

type TrustedSubnet struct {
	Subnets string
}

func (ts TrustedSubnet) TrustedSubnetsInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	trustedSubnets := metadata.ValueFromIncomingContext(ctx, "x-real-ip")
	if trustedSubnets == nil {
		return nil, errors.New("subnet required, but not passed")
	} else {
		if !strings.Contains(trustedSubnets[0], ts.Subnets) {
			return nil, errors.New("subnet check failed")
		}
	}

	return handler(ctx, req)
}
