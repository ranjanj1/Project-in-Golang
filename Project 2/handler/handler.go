package handler

import (
	"Project2/mapper"
	tokenv1 "Project2/proto"
	"Project2/service"
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TokenHandler struct {
	tokenv1.UnimplementedTokenServiceServer
	tickService *service.TokenService
	mapper      mapper.Mapper
}

func ProvideTokenHandler(tokenSvc *service.TokenService, mapper mapper.Mapper) *TokenHandler {
	return &TokenHandler{
		tickService: tokenSvc,
		mapper:      mapper,
	}
}

func (th TokenHandler) CreateToken(_ context.Context, request *tokenv1.CreateTokenRequest) (*tokenv1.CreateTokenResponse, error) {
	th.tickService.Create(request.Id)
	return &tokenv1.CreateTokenResponse{Status: fmt.Sprintf("token with id %v created", request.Id)}, nil
}

func (th TokenHandler) WriteToken(_ context.Context, request *tokenv1.WriteTokenRequest) (*tokenv1.WriteTokenResponse, error) {
	err := th.tickService.Write(th.mapper.Token(request.Token))
	if err != nil {
		return &tokenv1.WriteTokenResponse{Status: "token not found"}, status.Error(codes.NotFound, err.Error())
	}
	return &tokenv1.WriteTokenResponse{Status: fmt.Sprintf("token with id %v saved", request.Token.Id)}, nil
}

func (th TokenHandler) ReadToken(_ context.Context, request *tokenv1.ReadTokenRequest) (*tokenv1.ReadTokenResponse, error) {
	token, err := th.tickService.Read(request.Id)
	if err != nil {
		return &tokenv1.ReadTokenResponse{Token: nil, Status: "token not found"}, status.Error(codes.NotFound, err.Error())
	}
	return &tokenv1.ReadTokenResponse{Token: th.mapper.TokenPb(token), Status: "token found"}, nil
}

func (th TokenHandler) DropToken(_ context.Context, request *tokenv1.DropTokenRequest) (*tokenv1.DropTokenResponse, error) {
	err := th.tickService.Drop(request.Id)
	if err != nil {
		return &tokenv1.DropTokenResponse{Status: "token not found"}, status.Error(codes.NotFound, err.Error())
	}
	return &tokenv1.DropTokenResponse{Status: fmt.Sprintf("token with id %v dropped", request.Id)}, nil
}
