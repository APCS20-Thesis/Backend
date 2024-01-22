package service

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
)

type AuthInterceptor struct {
	jwtManager      *JWTManager
	accessibleRoles map[string][]string
}

func NewAuthInterceptor(jwtManager *JWTManager, accessibleRoles map[string][]string) *AuthInterceptor {
	return &AuthInterceptor{jwtManager, accessibleRoles}
}

func (interceptor *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		log.Println("Full method in interceptor: ", info.FullMethod)
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "Metadata is not provided")
		}

		claims, err := interceptor.authorize(ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}
		if claims != nil {
			md.Append("account_uuid", claims.UUID)
			ctx = metadata.NewIncomingContext(ctx, md)
		}

		return handler(ctx, req)
	}
}

func (interceptor *AuthInterceptor) authorize(ctx context.Context, method string) (*UserClaims, error) {
	accessibleRoles, ok := interceptor.accessibleRoles[method]
	if !ok {
		// everyone can access
		return nil, nil
	}

	accessToken, err := GetMetadata(ctx, "authorization")
	if err != nil {
		log.Fatalln("Cannot get accessToken from context", err)
		return nil, err
	}

	claims, err := interceptor.jwtManager.Verify(accessToken)
	if err != nil {
		return nil, err
	}

	for _, role := range accessibleRoles {
		if role == claims.Role {
			return claims, nil
		}
	}

	return nil, status.Error(codes.PermissionDenied, "No permission to access")
}

func AddMetadata(ctx context.Context, md metadata.MD, keys []string) {
	//TODO: Try to implement it
}

func GetMetadata(ctx context.Context, key string) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "Metadata is not provided")
	}
	values := md[key]
	if len(values) == 0 {
		return "", status.Errorf(codes.Unauthenticated, "Key is not provided")
	}

	return values[0], nil
}

//func unpackFullMethod(fullMethod string) (string, string) {
//	methodInfo := strings.Split(fullMethod, "/")
//	if methodName, ok := constant.MapServiceName[methodInfo[1]]; ok {
//		return methodName, methodInfo[2]
//	}
//	return methodInfo[1], methodInfo[2]
//}
