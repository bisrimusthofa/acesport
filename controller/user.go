package controller

import (
	"net/http"

	"github.com/bisrimusthofa/acesport/helper"
	"github.com/bisrimusthofa/acesport/user"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService user.Service
}

func NewUserController(userService user.Service) *UserController {
	return &UserController{userService: userService}
}

func (controller *UserController) Register(c *gin.Context) {
	var input user.RegisterUserInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.ErrorValidationResponse(err)
		errorMessages := gin.H{"errors": errors}

		response := helper.APIResponse(
			http.StatusUnprocessableEntity,
			"error",
			"The data given was invalid",
			errorMessages,
		)
		c.JSON(http.StatusBadRequest, response)
		panic(err)
	}

	newUser, err := controller.userService.Register(input)
	if err != nil {
		response := helper.APIResponse(
			http.StatusBadRequest,
			"error",
			"Failed register account",
			nil,
		)
		c.JSON(http.StatusBadRequest, response)
		panic(err)
	}

	dataFormated := user.FormatUser(newUser)

	// token
	dataFormated.Token = "blablabla"

	response := helper.APIResponse(
		http.StatusCreated,
		"success",
		"Account has been created",
		dataFormated,
	)
	c.JSON(http.StatusOK, response)
}
