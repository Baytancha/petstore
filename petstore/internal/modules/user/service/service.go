package service

import (
	"errors"
	"test/internal/infrastructure/filters"
	"test/internal/infrastructure/helpers"
	"test/internal/infrastructure/validator"
	"test/internal/modules/user/repository"

	"test/internal/models"
)

var (
	ErrValidation      = errors.New("validation error")
	ErrRecordNotFound  = errors.New("record not found")
	ErrEditConflict    = errors.New("edit conflict")
	ErrDuplicateEmail  = errors.New("duplicate email")
	ErrNoUser          = errors.New("user doesn't exist")
	ErrWrongPassword   = errors.New("wrong password")
	ErrInvalidPassword = errors.New("invalid password")
)

// r.Post("/user", ctrl.UserHandler.CreateUser)//ctrl.Auth.Register
// 	r.Get("/user/login", ctrl.Auth.Login)
// 	r.Get("/user/logout", ctrl.Auth.Logout)
// 	r.Get("/user/{username}", ctrl.UserHandler.GetUserByEmail)
// 	r.Put("/user/{username}", ctrl.UserHandler.UpdateUser)
// 	r.Delete("/user/{username}", ctrl.UserHandler.DeleteUser)
// 	r.Post("/user/CreateWithList", ctrl.UserHandler.ListUsers)
// 	r.Post("/user/CreateWithArray", ctrl.UserHandler.CreateArray)

type IUserService interface {
	Login(username, password string) (*models.User, *string, error)
	GetUserByName(email string) (*models.User, error)
	GetUserById(id int64) (*models.User, error)
	CreateUser(password string, user *models.User) error
	UpdateUser(dto any, name string) error
	DeleteUser(name string) error
	ListUsers(filters filters.Filters) ([]*models.User, filters.Metadata, error)
}

type UserService struct {
	storage repository.IUserStorage
}

func NewUserService(repo repository.IUserStorage) *UserService {
	return &UserService{storage: repo}
}

func (s *UserService) Login(username, password string) (*models.User, *string, error) {
	user, err := s.storage.GetByName(username)
	if err != nil {
		switch {
		case errors.Is(err, ErrRecordNotFound):
			return nil, nil, ErrRecordNotFound
		default:
			return nil, nil, err
		}
	}

	ok, err := user.Password.Matches(password)
	if !ok && err != nil {
		return nil, nil, err
	} else if !ok && err == nil {
		return nil, nil, ErrWrongPassword
	}
	token := helpers.GenerateToken(username)
	return user, &token, nil
}

func (s *UserService) CreateUser(password string, user *models.User) error {
	err := user.Password.Set(password)
	if err != nil {
		return ErrInvalidPassword
	}
	v := validator.New()

	if ValidateUser(v, user); !v.Valid() {
		return ErrValidation
	}
	err = s.storage.Insert(user)

	if err != nil {
		switch {
		case errors.Is(err, ErrDuplicateEmail):
			return ErrDuplicateEmail
		default:
			return err
		}
	}
	return nil
}

func (s *UserService) GetUserByName(username string) (*models.User, error) {
	user, err := s.storage.GetByName(username)
	if err != nil {
		switch {
		case errors.Is(err, ErrRecordNotFound):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return user, nil
}

func (s *UserService) UpdateUser(dto any, name string) error {
	user, err := s.storage.GetByName(name)
	if err != nil {
		switch {
		case errors.Is(err, ErrRecordNotFound):
			return ErrNoUser
		default:
			return err
		}
	}

	if dto.(struct {
		Name  *string `json:"name"`
		Email *string `json:"email"`
	}).Name != nil {
		user.Name = *dto.(struct {
			Name  *string `json:"name"`
			Email *string `json:"email"`
		}).Name
	}
	if dto.(struct {
		Name  *string `json:"name"`
		Email *string `json:"email"`
	}).Email != nil {
		user.Email = *dto.(struct {
			Name  *string `json:"name"`
			Email *string `json:"email"`
		}).Email
	}

	v := validator.New()

	if ValidateUser(v, user); !v.Valid() {
		return ErrValidation
	}

	err = s.storage.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, ErrEditConflict):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (s *UserService) GetUserById(id int64) (*models.User, error) {
	user, err := s.storage.Get(id)
	if err != nil {
		return nil, ErrRecordNotFound
	}
	return user, nil
}

func (s *UserService) DeleteUser(name string) error {
	deleted, err := s.storage.GetByName(name)
	if err != nil {
		switch {
		case errors.Is(err, ErrRecordNotFound):
			return ErrRecordNotFound
		default:
			return err
		}
	}
	err = s.storage.Delete(deleted.ID)
	if err != nil {
		switch {
		case errors.Is(err, ErrRecordNotFound):
			return ErrRecordNotFound
		default:
			return err
		}
	}
	return nil
}

func (s *UserService) ListUsers(filters filters.Filters) ([]*models.User, filters.Metadata, error) {
	models, meta, err := s.storage.GetAll(filters)
	if err != nil {
		return nil, meta, err
	}
	return models, meta, nil
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}

func ValidateUser(v *validator.Validator, user *models.User) {
	v.Check(user.Name != "", "name", "must be provided")
	v.Check(len(user.Name) <= 500, "name", "must not be more than 500 bytes long")

	ValidateEmail(v, user.Email)

	if user.Password.Plaintext != nil {
		ValidatePasswordPlaintext(v, *user.Password.Plaintext)
	}

	if user.Password.Hash == nil {
		panic("missing password hash for user")
	}
}
