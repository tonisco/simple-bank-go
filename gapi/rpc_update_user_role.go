package gapi

import (
	"context"
	"errors"

	db "github.com/tonisco/simple-bank-go/db/sqlc"
	"github.com/tonisco/simple-bank-go/pb"
	"github.com/tonisco/simple-bank-go/util"
	"github.com/tonisco/simple-bank-go/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) UpdateUserRole(ctx context.Context, req *pb.UpdateUserRoleRequest) (*pb.UpdateUserRoleResponse, error) {
	_, err := server.authorizeUser(ctx, []string{util.BankerRole})

	if err != nil {
		return nil, unauthenticatedError(err)
	}

	violations := validateUpdateUserRoleRequest(req)

	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	args := db.UpdateUserRoleParams{
		Username: req.GetUsername(),
		Role:     req.GetRole(),
	}

	user, err := server.store.UpdateUserRole(ctx, args)

	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to create user: %s", err)
	}
	rsp := &pb.UpdateUserRoleResponse{
		User: convertUser(user),
	}
	return rsp, nil
}

func validateUpdateUserRoleRequest(req *pb.UpdateUserRoleRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if err := val.ValidateRole(req.GetRole()); err != nil {
		violations = append(violations, fieldViolation("role", err))

	}

	return violations
}
