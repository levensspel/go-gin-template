package server

import (
	"github.com/gin-gonic/gin"
	userHandler "github.com/levensspel/go-gin-template/handler/user"
	"github.com/levensspel/go-gin-template/logger"
	"github.com/levensspel/go-gin-template/middleware"
	pb "github.com/levensspel/go-gin-template/proto"
	dbTrxRepository "github.com/levensspel/go-gin-template/repository/db_trx"
	userRepository "github.com/levensspel/go-gin-template/repository/user"
	userService "github.com/levensspel/go-gin-template/service/user"
	"google.golang.org/grpc"
	"gorm.io/gorm"

	_ "github.com/levensspel/go-gin-template/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(r *gin.Engine, db *gorm.DB, g *grpc.Server) {
	logger := logger.NewlogHandler()

	// api := r.Group("/v1")
	// {
	// 	// untuk memanfaatkan api versioning, uncomment dan pakai ini
	// }

	dbTrxRepo := dbTrxRepository.NewDBTrxRepository(db)

	userRepo := userRepository.NewUserRepository(db)
	userSrv := userService.NewUserService(userRepo, dbTrxRepo, logger)
	userHdlr := userHandler.NewUserHandler(userSrv)

	grpcUserHandler := userHandler.NewUserGrpcHandler(userSrv)
	pb.RegisterUserServiceServer(g, grpcUserHandler)

	swaggerRoute := r.Group("/")
	{
		//Route untuk Swagger
		swaggerRoute.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	controllers := r.Group("/api")
	{
		user := controllers.Group("/user")
		{
			user.POST("/register", userHdlr.Register)
			user.POST("/login", userHdlr.Login)
			user.PUT("", middleware.Authorization, userHdlr.Update)
			user.DELETE("", middleware.Authorization, userHdlr.Delete)
		}
		// tambah route lainnya disini
	}
}
