package mapper

import (
	"persons/internal/domain"
	"persons/internal/domain/models"
)

func PersonToGet(person models.Person) domain.GetPerson {
	return domain.GetPerson{
		Login: person.Login,
		Image: person.Image,
	}
}

func RegisterToPerson(person domain.RegisterPerson) models.Person {
	return models.Person{
		Login:        person.Login,
		Email:        person.Email,
		PasswordHash: person.Password,
	}
}
