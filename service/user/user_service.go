package user_service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/levensspel/go-gin-template/auth"
	"github.com/levensspel/go-gin-template/dto"
	"github.com/levensspel/go-gin-template/entity"
	"github.com/levensspel/go-gin-template/helper"
	"github.com/levensspel/go-gin-template/logger"
	repositories "github.com/levensspel/go-gin-template/repository/user"
	"github.com/levensspel/go-gin-template/validation"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	RegisterUser(input dto.RequestRegister) (dto.ResponseRegister, error)
	Login(input dto.RequestLogin) (dto.ResponseLogin, error)
	Update(input dto.RequestRegister) (dto.Response, error)
	DeleteByID(id string) error
}

type service struct {
	userRepo repositories.UserRepository
	logger   logger.Logger
}

func NewUserService(
	userRepo repositories.UserRepository,
	logger logger.Logger,
) UserService {
	return &service{
		userRepo: userRepo,
		logger:   logger,
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
		s.logger.Error(err.Error(), helper.GenerateFromPassword, passwordHash)
		return dto.ResponseRegister{}, err
	}
	user.Password = string(passwordHash)

	err = s.userRepo.Create(context.Background(), user)
	if err != nil {
		s.logger.Error(err.Error(), helper.UserServiceRegister, user)
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
	s.logger.Info("Finding user account by email", helper.FunctionCaller("UserService.Login"), input)
	user, err := s.userRepo.GetUserbyEmail(context.Background(), input.Email)
	if err != nil {
		s.logger.Error(err.Error(), helper.UserServiceLogin, input)
		return dto.ResponseLogin{}, err
	}
	s.logger.Info("Found available user", helper.FunctionCaller("UserService.Login"), user)

	if len(user) == 0 {
		return dto.ResponseLogin{}, helper.ErrorInvalidLogin
	}

	err = bcrypt.CompareHashAndPassword([]byte(user[0].Password), []byte(input.Password))
	if err != nil {
		s.logger.Error(err.Error(), helper.UserServiceLogin, err)
		return dto.ResponseLogin{}, helper.ErrorInvalidLogin
	}

	jwtService := auth.NewJWTService()

	token, err := jwtService.GenerateToken(user[0].Id)

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
	err = s.userRepo.Update(context.Background(), user)
	if err != nil {
		s.logger.Error(err.Error(), helper.UserServiceUpdate, err)
		return dto.Response{}, err
	}
	response := dto.Response{}
	response.Id = input.Id
	response.Email = input.Email
	response.Username = input.Username

	return response, nil
}

func (s *service) DeleteByID(id string) error {
	err := s.userRepo.Delete(context.Background(), id)
	if err != nil {
		s.logger.Error(err.Error(), helper.UserServiceUpdate, err)
		return err
	}
	return err
}
