package services

import (
	"errors"
	"nubes/sum/db"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type ForbiddenError struct{}

func (r *ForbiddenError) Error() string {
	return "forbidden"
}

type ConflictError struct {
	Detail string
}

func (r *ConflictError) Error() string {
	return r.Detail
}

func ToGrpcError(e error) error {
	var forbidden *ForbiddenError
	var validation *db.ValidationError
	var conflict *ConflictError

	if errors.Is(e, gorm.ErrRecordNotFound) {
		return status.Errorf(codes.NotFound, "not_found")
	}

	if errors.As(e, &forbidden) {
		return status.Errorf(codes.PermissionDenied, forbidden.Error())
	}

	if errors.As(e, &validation) {
		return status.Errorf(codes.InvalidArgument, validation.Error())
	}

	if errors.As(e, &conflict) {
		return status.Errorf(codes.Aborted, conflict.Error())
	}

	return e
}
