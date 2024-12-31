package helper

import (
	"errors"
)

type FunctionCaller string

const (
	UserRepoCreate FunctionCaller = "userRepo.Create"
	DbTrxRepoBegin FunctionCaller = "dbTrxRepo.Begin"
)

var ErrorBadRequest = errors.New("invalid request format")
var ErrorNotFound = errors.New("data not found")
var ErrorInternalServerError = errors.New("internal server error")

var ErrorEmailRegistered = errors.New("email is already registered")
var ErrorUsernameRegistered = errors.New("username is already registered")
var ErrorInvalidLogin = errors.New("invalid email or password")

var WORK_DIR string