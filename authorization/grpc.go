package authorization

import (
	"context"
	"goclean"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func GRPCAuthorizer() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.PermissionDenied, goclean.ErrUnauthorized.Error())
		}

		var auth []string
		if auth = md.Get("Authorization"); len(auth) < 1 {
			return nil, status.Error(codes.PermissionDenied, goclean.ErrUnauthorized.Error())
		}

		if auth[0] != "password" {
			return nil, status.Error(codes.PermissionDenied, goclean.ErrUnauthorized.Error())
		}

		ctx = context.WithValue(ctx, goclean.ContextKeyPermissions, goclean.Greet)
		return handler(ctx, req)
	}
}
