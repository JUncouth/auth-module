package auth

import (
	authpb "auth/src/proto"
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCService struct {
	authpb.UnimplementedAuthorizationServiceServer

	service AuthService
}

func NewGRPCService(service AuthService) *GRPCService {
	return &GRPCService{
		service: service,
	}
}

func (s *GRPCService) Register(
	ctx context.Context,
	req *authpb.RegisterRequest,
) (*authpb.AuthResponse, error) {
	result, err := s.service.Register(
		ctx,
		req.Email,
		req.FirstName,
		req.LastName,
		req.PhoneNumber,
		req.Password,
	)
	if err != nil {
		return nil, mapError(err)
	}

	return &authpb.AuthResponse{
		UserId: result.UserID.String(),
		Token:  result.Token,
	}, nil
}

func (s *GRPCService) Authenticate(
	ctx context.Context,
	req *authpb.AuthenticateRequest,
) (*authpb.AuthResponse, error) {
	result, err := s.service.Authenticate(
		ctx,
		req.Email,
		req.Password,
	)
	if err != nil {
		return nil, mapError(err)
	}

	return &authpb.AuthResponse{
		UserId: result.UserID.String(),
		Token:  result.Token,
	}, nil
}

func (s *GRPCService) ChangePassword(
	ctx context.Context,
	req *authpb.ChangePasswordRequest,
) (*authpb.OperationResponse, error) {
	if err := s.service.ChangePassword(
		ctx,
		req.Email,
		req.OldPassword,
		req.NewPassword,
	); err != nil {
		return nil, mapError(err)
	}

	return &authpb.OperationResponse{
		Success: true,
		Message: "password changed",
	}, nil
}

func (s *GRPCService) ResetPassword(
	ctx context.Context,
	req *authpb.ResetPasswordRequest,
) (*authpb.OperationResponse, error) {
	if err := s.service.ResetPassword(
		ctx,
		req.Email,
	); err != nil {
		return nil, mapError(err)
	}

	return &authpb.OperationResponse{
		Success: true,
		Message: "password reset",
	}, nil
}

func (s *GRPCService) VerifyToken(
	ctx context.Context,
	req *authpb.VerifyTokenRequest,
) (*authpb.VerifyTokenResponse, error) {
	info, err := s.service.VerifyToken(
		ctx,
		req.Token,
	)
	if err != nil {
		return nil, mapError(err)
	}

	return &authpb.VerifyTokenResponse{
		UserId: info.UserID.String(),
		Email:  info.Email,
	}, nil
}

func mapError(err error) error {
	switch {
	case errors.Is(err, ErrInvalidInput):
		return status.Error(
			codes.InvalidArgument,
			err.Error(),
		)

	case errors.Is(err, ErrEmailAlreadyExists):
		return status.Error(
			codes.AlreadyExists,
			err.Error(),
		)

	case errors.Is(err, ErrInvalidCredentials):
		return status.Error(
			codes.Unauthenticated,
			err.Error(),
		)

	case errors.Is(err, ErrInvalidToken):
		return status.Error(
			codes.Unauthenticated,
			err.Error(),
		)

	case errors.Is(err, ErrUserNotFound):
		return status.Error(
			codes.NotFound,
			err.Error(),
		)

	default:
		return status.Error(
			codes.Internal,
			"internal authorization error",
		)
	}
}
