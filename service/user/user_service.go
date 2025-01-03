package user_service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"github.com/levensspel/go-gin-template/auth"
	"github.com/levensspel/go-gin-template/dto"
	"github.com/levensspel/go-gin-template/entity"
	"github.com/levensspel/go-gin-template/helper"
	"github.com/levensspel/go-gin-template/logger"
	pb "github.com/levensspel/go-gin-template/proto"
	dbTrxRepository "github.com/levensspel/go-gin-template/repository/db_trx"
	userRepository "github.com/levensspel/go-gin-template/repository/user"
	"github.com/levensspel/go-gin-template/validation"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	RegisterUser(input dto.RequestRegister) (dto.ResponseRegister, error)
	Login(input dto.RequestLogin) (dto.ResponseLogin, error)
	Update(input dto.RequestRegister) (dto.Response, error)
	DeleteByID(input string) error

	RegisterUserWithGrpc(input dto.RequestRegister) (dto.ResponseRegister, error)
}

type service struct {
	userRepo        userRepository.UserRepository
	dbTrxRepo       dbTrxRepository.DBTrxRepository
	grpcUserService pb.UserServiceClient
	logger          logger.Logger
}

func NewUserService(
	userRepo userRepository.UserRepository,
	dbTrxRepo dbTrxRepository.DBTrxRepository,
	grpcUserService pb.UserServiceClient,
	logger logger.Logger,
) UserService {
	return &service{
		userRepo:        userRepo,
		dbTrxRepo:       dbTrxRepo,
		grpcUserService: grpcUserService,
		logger:          logger,
	}
}

func (s *service) RegisterUser(input dto.RequestRegister) (dto.ResponseRegister, error) {
	err := validation.ValidateUserCreate(input, s.userRepo)

	if err != nil {
		return dto.ResponseRegister{}, err
	}

	user := entity.User{}

	user.Id = uuid.New().String()
	user.Username = input.Username
	user.Email = input.Email
	user.CreatedAt = time.Now().Unix()
	user.UpdatedAt = time.Now().Unix()
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		s.logger.Error(err.Error(), helper.UserServiceRegister, err)
		return dto.ResponseRegister{}, err
	}

	user.Password = string(passwordHash)

	err = s.userRepo.Create(user, nil)
	if err != nil {
		s.logger.Error(err.Error(), helper.UserServiceRegister, err)
		return dto.ResponseRegister{}, err
	}
	response := dto.ResponseRegister{
		Id: user.Id,
	}

	return response, nil
}

func (s *service) Login(input dto.RequestLogin) (dto.ResponseLogin, error) {
	err := validation.ValidateUserLogin(input)
	if err != nil {
		return dto.ResponseLogin{}, err
	}
	email := input.Email
	password := input.Password
	user, err := s.userRepo.FindByEmail(email, nil)

	if err != nil {
		s.logger.Error(err.Error(), helper.UserServiceLogin, err)
		return dto.ResponseLogin{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		s.logger.Error(err.Error(), helper.UserServiceLogin, err)
		return dto.ResponseLogin{}, helper.ErrorInvalidLogin
	}

	jwtService := auth.NewJWTService()

	token, err := jwtService.GenerateToken(user.Id)

	if err != nil {
		s.logger.Error(err.Error(), helper.UserServiceLogin, err)
		return dto.ResponseLogin{}, err
	}

	response := dto.ResponseLogin{}
	response.Token = token
	return response, nil
}

func (s *service) Update(input dto.RequestRegister) (dto.Response, error) {
	user := entity.User{}
	user.Id = input.Id
	user.Username = input.Username
	user.Email = input.Email
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)

	if err != nil {
		s.logger.Error(err.Error(), helper.UserServiceUpdate, err)
		return dto.Response{}, err
	}

	user.Password = string(passwordHash)
	user.UpdatedAt = time.Now().Unix()
	updatedUser, err := s.userRepo.Update(user, nil)
	if err != nil {
		s.logger.Error(err.Error(), helper.UserServiceUpdate, err)
		return dto.Response{}, err
	}
	response := dto.Response{}
	copier.Copy(&response, &updatedUser)

	return response, nil
}

func (s *service) DeleteByID(id string) error {
	return s.userRepo.DeleteByID(id, nil)
}

// Sampel servis user yang memanggil servis gRPC eksternal
func (s *service) RegisterUserWithGrpc(input dto.RequestRegister) (dto.ResponseRegister, error) {
	err := validation.ValidateUserCreate(input, s.userRepo)
	if err != nil {
		return dto.ResponseRegister{}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// baris ini akan mengirim request gRPC ke servis eksternal (microservice lain)
	grpcResponse, err := s.grpcUserService.RegisterUser(ctx, &pb.RequestRegister{
		Id:       uuid.New().String(),
		Username: input.Username,
		Email:    input.Email,
		Password: input.Password,
	})
	if err != nil {
		s.logger.Error(err.Error(), helper.GrpcUserServiceRegister, err)
		return dto.ResponseRegister{}, errors.New(grpcResponse.GetMessage())
	}

	response := dto.ResponseRegister{
		Id: grpcResponse.UserId,
	}

	return response, nil
}
