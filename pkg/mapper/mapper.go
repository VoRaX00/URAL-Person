package mapper

import (
	"persons/internal/domain"
	"persons/internal/domain/models"
)

func UserToGet(person models.User) domain.GetUser {
	return domain.GetUser{
		Login: person.Login,
		Image: person.Image,
	}
}

func RegisterToUser(person domain.RegisterUser) models.User {
	return models.User{
		Login:        person.Login,
		Email:        person.Email,
		PasswordHash: person.Password,
	}
}
