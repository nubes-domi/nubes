package rpc

import (
	context "context"
	"nubes/sum/db"
	"nubes/sum/services"
	"nubes/sum/services/sessions"
	"nubes/sum/utils"

	"github.com/golang/protobuf/ptypes/empty"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type SessionsServerImpl struct {
	UnimplementedSessionsServer
}

func (s *SessionsServerImpl) GetAuthenticationMethods(ctx context.Context, req *GetAuthenticationMethodsRequest) (*GetAuthenticationMethodsResponse, error) {
	identifier := req.Identifier

	user, methods, err := sessions.GetAuthenticationMethods(identifier)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid_credentials")
	} else {
		return &GetAuthenticationMethodsResponse{
			Username:              user.Username,
			AuthenticationMethods: methods,
		}, nil
	}
}

func (s *SessionsServerImpl) Create(ctx context.Context, req *CreateSessionRequest) (*Session, error) {
	identifier := req.Identifier
	password := req.Password

	session, err := sessions.Create(identifier, password, req.Session.UserAgent, req.Session.IpAddress)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid_credentials")
	} else {
		return &Session{
			AccessToken: session.SignedToken,
			UpdatedAt:   timestamppb.New(session.UpdatedAt),
			UserAgent:   session.UserAgent,
			IpAddress:   session.IPAddress,
		}, nil
	}
}

func (s *SessionsServerImpl) Delete(ctx context.Context, req *DeleteSessionRequest) (*empty.Empty, error) {
	sessions.Delete(currentUser(ctx), req.GetSessionId())

	return &empty.Empty{}, nil
}

func (s *SessionsServerImpl) Get(ctx context.Context, req *GetSessionRequest) (*Session, error) {
	token := req.AuthenticationToken
	session, err := sessions.Get(token)
	if err != nil {
		return nil, services.ToGrpcError(err)
	} else {
		return &Session{
			AccessToken: token,
			UpdatedAt:   timestamppb.New(session.UpdatedAt),
			UserAgent:   session.UserAgent,
			IpAddress:   session.IPAddress,
		}, nil
	}
}

func (s *SessionsServerImpl) List(ctx context.Context, req *ListSessionsRequest) (*ListSessionsResponse, error) {
	list, err := sessions.List(currentUser(ctx), req.OrderBy)
	if err != nil {
		return nil, services.ToGrpcError(err)
	}

	return &ListSessionsResponse{
		Sessions: utils.Collect(list, func(s *db.UserSession) *Session {
			return &Session{
				UpdatedAt: timestamppb.New(s.UpdatedAt),
				UserAgent: s.UserAgent,
				IpAddress: s.IPAddress,
			}
		}),
	}, nil
}

func (s *SessionsServerImpl) Update(ctx context.Context, req *UpdateSessionRequest) (*Session, error) {
	session, err := sessions.Update(currentUser(ctx), currentSession(ctx), req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid_credentials")
	} else {
		return &Session{
			AccessToken: session.SignedToken,
			UpdatedAt:   timestamppb.New(session.UpdatedAt),
			UserAgent:   session.UserAgent,
			IpAddress:   session.IPAddress,
		}, nil
	}
}
