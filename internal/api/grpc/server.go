package grpc_api

import (
	"context"
	"fmt"

	"github.com/Sugyk/auth_service/internal/api/grpc/pb"
	"github.com/Sugyk/auth_service/internal/models"
	"github.com/Sugyk/auth_service/pkg/logger"
)

type Service interface {
	Register(ctx context.Context, login string, password string) error
	Login(ctx context.Context, login string, password string) (string, error)
}

type Server struct {
	pb.UnimplementedAuthServiceServer
	service Service
	logger  logger.Logger
}

func NewServer(service Service, logger logger.Logger) *Server {
	return &Server{
		service: service,
		logger:  logger,
	}
}

func (s *Server) handleError(ctx context.Context, err error) error {
	appErr := appErrorFrom(err)
	s.logger.Error(ctx, "request error", "error", appErr.Error(), "cause", appErr.Cause())
	return toGRPCError(appErr)
}

func (s *Server) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	reqModel := models.RegisterRequest{
		Login:    req.GetLogin(),
		Password: req.GetPassword(),
	}
	if err := reqModel.Validate(); err != nil {
		return nil, s.handleError(ctx, err)
	}

	if err := s.service.Register(ctx, reqModel.Login, reqModel.Password); err != nil {
		return nil, s.handleError(ctx, err)
	}

	return &pb.RegisterResponse{
		Message: fmt.Sprintf("user with login '%s' created", reqModel.Login),
	}, nil
}

func (s *Server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	reqModel := models.LoginRequest{
		Login:    req.GetLogin(),
		Password: req.GetPassword(),
	}
	if err := reqModel.Validate(); err != nil {
		return nil, s.handleError(ctx, err)
	}

	accessToken, err := s.service.Login(ctx, reqModel.Login, reqModel.Password)
	if err != nil {
		return nil, s.handleError(ctx, err)
	}

	return &pb.LoginResponse{
		AccessToken: accessToken,
	}, nil
}
