package user

import (
	"time"

	"github.com/bisrimusthofa/acesport/helper"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(input RegisterUserInput) (User, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository: repository}
}

func (s *service) Register(input RegisterUserInput) (User, error) {
	user := User{
		Id:        uuid.New().String(),
		Name:      input.Name,
		Email:     input.Email,
		Phone:     input.Phone,
		Role:      "Reseller",
		CreatedAt: time.Now(),
	}

	password, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	helper.PanicIfError(err)
	user.Password = string(password)

	dataUser, err := s.repository.Save(user)
	if err != nil {
		return user, err
	}

	return dataUser, nil
}
