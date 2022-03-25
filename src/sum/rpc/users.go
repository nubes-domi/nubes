package rpc

import (
	context "context"
	"nubes/sum/db"
	"nubes/sum/services"
	"nubes/sum/services/users"
	"nubes/sum/utils"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UsersServerImpl struct {
	UnimplementedUsersServer
}

func buildUserFromDb(u *db.User) *User {
	return &User{
		Admin:       u.Admin,
		Username:    u.Username,
		Name:        u.Name,
		Email:       u.Email,
		PhoneNumber: u.PhoneNumber,
		Picture:     u.Picture,
		Locale:      u.Locale,
		Zoneinfo:    u.Zoneinfo,
		UpdatedAt:   timestamppb.New(u.UpdatedAt),
	}
}

func (s *UsersServerImpl) Create(ctx context.Context, req *CreateUserRequest) (*User, error) {
	user, err := users.Create(currentUser(ctx), &db.User{
		Admin:       req.User.Admin,
		Username:    req.User.Username,
		Password:    req.User.Password,
		Name:        req.User.Name,
		Email:       req.User.Email,
		PhoneNumber: req.User.PhoneNumber,
		Picture:     req.User.Picture,
		// Birthdate:   req.User.Birthdate,
		Locale:   req.User.Locale,
		Zoneinfo: req.User.Zoneinfo,
	})
	if err != nil {
		return nil, services.ToGrpcError(err)
	}

	return buildUserFromDb(user), nil
}

func (s *UsersServerImpl) Delete(ctx context.Context, req *DeleteUserRequest) (*empty.Empty, error) {
	err := users.Delete(currentUser(ctx), req.UserId)
	return &empty.Empty{}, services.ToGrpcError(err)
}

func (s *UsersServerImpl) Get(ctx context.Context, req *GetUserRequest) (*User, error) {
	user, err := users.Get(currentUser(ctx), req.UserId)
	if err != nil {
		return nil, services.ToGrpcError(err)
	}

	return buildUserFromDb(user), nil
}

func (s *UsersServerImpl) List(ctx context.Context, req *ListUsersRequest) (*ListUsersResponse, error) {
	list, err := users.List(currentUser(ctx), req.OrderBy)
	if err != nil {
		return nil, services.ToGrpcError(err)
	}

	return &ListUsersResponse{
		Users: utils.Collect(list, buildUserFromDb),
	}, nil
}

func (s *UsersServerImpl) Update(ctx context.Context, req *UpdateUserRequest) (*User, error) {
	user, err := users.Update(currentUser(ctx), &db.User{
		Model:       db.Model{ID: req.User.Id},
		Admin:       req.User.Admin,
		Username:    req.User.Username,
		Password:    req.User.Password,
		Name:        req.User.Name,
		Email:       req.User.Email,
		PhoneNumber: req.User.PhoneNumber,
		Picture:     req.User.Picture,
		// Birthdate:   req.User.Birthdate,
		Locale:   req.User.Locale,
		Zoneinfo: req.User.Zoneinfo,
	})
	if err != nil {
		return nil, services.ToGrpcError(err)
	}

	return buildUserFromDb(user), nil
}
