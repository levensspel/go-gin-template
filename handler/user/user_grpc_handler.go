package userHandler

import (
	"context"
	"net/http"

	"github.com/levensspel/go-gin-template/dto"
	"github.com/levensspel/go-gin-template/helper"
	pb "github.com/levensspel/go-gin-template/proto"
	service "github.com/levensspel/go-gin-template/service/user"
)

type grpcHandler struct {
	service service.UserService
	pb.UnimplementedUserServiceServer
}

func NewUserGrpcHandler(service service.UserService) *grpcHandler {
	return &grpcHandler{
		service:                        service,
		UnimplementedUserServiceServer: pb.UnimplementedUserServiceServer{},
	}
}

func (h *grpcHandler) RegisterUser(_ context.Context, in *pb.RequestRegister) (*pb.ResponseRegister, error) {
	input := &dto.RequestRegister{
		Id:       in.GetId(),
		Username: in.GetUsername(),
		Email:    in.GetEmail(),
		Password: in.GetPassword(),
	}

	response, err := h.service.RegisterUser(*input)

	if err != nil {
		return &pb.ResponseRegister{
			StatusCode: http.StatusBadRequest,
			Message:    helper.ErrBadRequest.Error(),
		}, err
	}

	return &pb.ResponseRegister{
		UserId:     response.Id,
		StatusCode: http.StatusCreated,
	}, nil
}
