package userHandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/levensspel/go-gin-template/dto"
	"github.com/levensspel/go-gin-template/helper"
	service "github.com/levensspel/go-gin-template/service/user"
)

type UserHandler interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)

	RegisterWithGrpc(ctx *gin.Context) // Sampel API yang memanggil servis external via gRPC
}

type handler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) UserHandler {
	return &handler{service: service}
}

// Register new user
// @Tags users
// @Summary Register new user
// @Description Register the new user
// @Accept json
// @Produce json
// @Param data body dto.RequestRegister true "data"
// @Success 201 {object} helper.Response{data=helper.Response} "CREATED"
// @Failure 400 {object} helper.Response{errors=helper.ErrorResponse} "Bad Request"
// @Failure 409 {object} helper.Response{errors=helper.ErrorResponse} "data conflict, like email already exist"
// @Router /users/register [POST]
func (h handler) Register(ctx *gin.Context) {
	input := new(dto.RequestRegister)

	err := ctx.ShouldBindJSON(&input)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, helper.NewResponse(nil, err))
		return
	}

	response, err := h.service.RegisterUser(*input)

	if err != nil {
		ctx.JSON(
			helper.GetErrorStatusCode(err),
			helper.NewResponse(
				helper.GetErrorStatusCode(err),
				nil,
			),
		)
		return
	}

	ctx.JSON(http.StatusCreated, helper.NewResponse(response, err))
}

// Login user
// @Tags users
// @Summary Login user
// @Description Login user
// @Accept  json
// @Produce  json
// @Param data body dto.RequestLogin true "data"
// @Success 200 {object} helper.Response{data=dto.ResponseLogin} "OK"
// @Failure 400 {object} helper.Response{errors=helper.ErrorResponse} "Bad Request"
// @Failure 404 {object} helper.Response{errors=helper.ErrorResponse} "Record not found"
// @Router /users/login [POST]
func (h handler) Login(ctx *gin.Context) {
	input := new(dto.RequestLogin)
	err := ctx.ShouldBindJSON(&input)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, helper.NewResponse(nil, err))
		return
	}
	response, err := h.service.Login(*input)

	if err != nil {
		ctx.JSON(helper.GetErrorStatusCode(err), helper.NewResponse(nil, err))
		return
	}
	ctx.JSON(http.StatusOK, helper.Response{Data: response, Error: err})
}

// Update user
// @Tags users
// @Summary Update user
// @Description Update user
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer + user token"
// @Param data body dto.UserRequestUpdate true "data"
// @Success 200 {object} helper.Response{data=helper.Response} "OK"
// @Failure 400 {object} helper.Response{errors=helper.ErrorResponse} "Bad Request"
// @Failure 401 {object} helper.Response{errors=helper.ErrorResponse} "Unauthorization"
// @Router /users [PUT]
func (h handler) Update(ctx *gin.Context) {
	input := new(dto.RequestRegister)

	if err := ctx.ShouldBind(input); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, helper.NewResponse(nil, err))
		return
	}
	id := ctx.MustGet("user_id")
	input.Id = id.(string)
	response, err := h.service.Update(*input)

	if err != nil {
		ctx.JSON(helper.GetErrorStatusCode(err), helper.NewResponse(nil, err))
		return
	}
	ctx.JSON(http.StatusOK, helper.NewResponse(response, nil))
}

// DeleteByID user
// @Tags users
// @Summary Delete user
// @Description Delete user
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer + user token"
// @Success 200 {object} helper.Response{data=helper.Response} "OK"
// @Failure 404 {object} helper.Response{errors=helper.ErrorResponse} "Not Found"
// @Failure 400 {object} helper.Response{errors=helper.ErrorResponse} "Bad Request"
// @Failure 401 {object} helper.Response{errors=helper.ErrorResponse} "Unauthorization"
// @Router /users [DELETE]
func (h handler) Delete(ctx *gin.Context) {
	id := ctx.MustGet("user_id")

	err := h.service.DeleteByID(id.(string))

	if err != nil {
		ctx.JSON(helper.GetErrorStatusCode(err), helper.NewResponse(nil, err))
		return
	}

	message := map[string]interface{}{"message": "your account has been successfully deleted"}
	ctx.JSON(http.StatusOK, helper.NewResponse(message, nil))
}

// Register new user with a call to an external service via gRPC
// @Tags users
// @Summary Register new user
// @Description Register the new user
// @Accept json
// @Produce json
// @Param data body dto.RequestRegister true "data"
// @Success 201 {object} helper.Response{data=helper.Response} "CREATED"
// @Failure 400 {object} helper.Response{errors=helper.ErrorResponse} "Bad Request"
// @Failure 409 {object} helper.Response{errors=helper.ErrorResponse} "data conflict, like email already exist"
// @Router /users/register [POST]
func (h handler) RegisterWithGrpc(ctx *gin.Context) {
	input := new(dto.RequestRegister)

	err := ctx.ShouldBindJSON(&input)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, helper.NewResponse(nil, err))
		return
	}

	// Sampel handler yang memiliki alur memanggil servis eksternal via gRPC
	response, err := h.service.RegisterUserWithGrpc(*input)

	if err != nil {
		ctx.JSON(
			helper.GetErrorStatusCode(err),
			helper.NewResponse(
				helper.GetErrorStatusCode(err),
				err,
			),
		)
		return
	}

	ctx.JSON(http.StatusCreated, helper.NewResponse(response, err))
}
