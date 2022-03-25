package rpc

import (
	context "context"
	"fmt"
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

func (s *SessionsServerImpl) Create(ctx context.Context, req *CreateSessionRequest) (*Session, error) {
	fmt.Printf("GOT REQ\n")

	username := req.Username
	password := req.Password

	session, err := sessions.Create(username, password, req.Session.UserAgent, req.Session.IpAddress)
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
