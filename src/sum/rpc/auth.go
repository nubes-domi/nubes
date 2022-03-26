package rpc

import (
	context "context"
	"nubes/sum/db"
	"nubes/sum/utils"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	status "google.golang.org/grpc/status"
)

func getAuthorizationWhitelistedEndpoints() []string {
	return []string{"/sum.Sessions/Create", "/sum.Sessions/Get", "/sum.Sessions/GetAuthenticationMethods"}
}

func ServerInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	newCtx, err := authorize(ctx, info)
	if err != nil {
		return nil, err
	}

	return handler(newCtx, req)
}

func authorize(ctx context.Context, info *grpc.UnaryServerInfo) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx, status.Errorf(codes.InvalidArgument, "Could not retrieve request metadata")
	}

	if auth, ok := md["authorization"]; ok {
		if token, err := utils.JwtVerify(auth[0]); err == nil {
			if session, err := db.DB.UserSessions().FindById(token.JwtID()); err == nil {
				sub := token.Subject()
				if user, err := db.DB.Users().FindById(sub); err == nil {
					ctx = context.WithValue(ctx, "CurrentSession", session)
					ctx = context.WithValue(ctx, "CurrentUser", user)
					return ctx, nil
				}
			}
		}
	}

	if utils.Contains(getAuthorizationWhitelistedEndpoints(), info.FullMethod) {
		return ctx, nil
	} else {
		return ctx, status.Errorf(codes.Unauthenticated, "requires_authentication")
	}
}

func currentUser(ctx context.Context) *db.User {
	user, ok := ctx.Value("CurrentUser").(*db.User)
	if !ok {
		panic("Could not retrieve CurrentUser in gRPC")
	}

	return user
}

func currentSession(ctx context.Context) *db.UserSession {
	session, ok := ctx.Value("CurrentSession").(*db.UserSession)
	if !ok {
		panic("Could not retrieve CurrentSession in gRPC")
	}

	return session
}
