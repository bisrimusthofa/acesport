package main

import (
	"github.com/bisrimusthofa/acesport/auth"
	"github.com/bisrimusthofa/acesport/controller"
	"github.com/bisrimusthofa/acesport/helper"
	"github.com/bisrimusthofa/acesport/middleware"
	"github.com/bisrimusthofa/acesport/user"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:@tcp(103.127.97.60:3306)/acesport?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	helper.PanicIfError(err)

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	authService := auth.NewJwtService()
	userController := controller.NewUserController(userService, authService)

	router := gin.Default()
	api := router.Group("/api/v1")

	//auth
	api.POST("/register", userController.Register)
	api.POST("/login", userController.Login)
	api.POST("/email-checker", userController.CheckEmailAvailable)

	// profile
	api.POST("/upload-avatar", middleware.AuthMiddleware(authService, userService), userController.UploadAvatar)
	api.PUT("/update-profile", middleware.AuthMiddleware(authService, userService), userController.UpdateProfile)

	err = router.Run(":3000")
	helper.PanicIfError(err)
}
