package service

import (
	"app/src/model"
	"app/src/utils"
	"app/src/validation"
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserService struct {
	Log      *logrus.Logger
	DB       *gorm.DB
	Validate *validator.Validate
}

func NewUserService(db *gorm.DB, validate *validator.Validate) *UserService {
	return &UserService{
		Log:      utils.Log,
		DB:       db,
		Validate: validate,
	}
}

func (s *UserService) GetUsers(c *fiber.Ctx, params *validation.QueryUser) ([]model.User, int64, error) {
	var users []model.User
	var totalResults int64

	if err := s.Validate.Struct(params); err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.Limit
	query := s.DB.WithContext(c.Context()).Order("created_at asc")

	if search := params.Search; search != "" {
		query = query.Where("name LIKE ? OR email LIKE ? OR role LIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	result := query.Find(&users).Count(&totalResults)
	if result.Error != nil {
		s.Log.Errorf("Failed to search users: %+v", result.Error)
		return nil, 0, result.Error
	}

	result = query.Limit(params.Limit).Offset(offset).Find(&users)
	if result.Error != nil {
		s.Log.Errorf("Failed to get all users: %+v", result.Error)
		return nil, 0, result.Error
	}

	return users, totalResults, result.Error
}

func (s *UserService) GetUserByID(c *fiber.Ctx, id string) (*model.User, error) {
	user := new(model.User)

	result := s.DB.WithContext(c.Context()).First(user, "id = ?", id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	if result.Error != nil {
		s.Log.Errorf("Failed get user by id: %+v", result.Error)
	}

	return user, result.Error
}

func (s *UserService) GetUserByEmail(c *fiber.Ctx, email string) (*model.User, error) {
	user := new(model.User)

	result := s.DB.WithContext(c.Context()).Where("email = ?", email).First(user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	if result.Error != nil {
		s.Log.Errorf("Failed get user by email: %+v", result.Error)
	}

	return user, result.Error
}

func (s *UserService) CreateUser(c *fiber.Ctx, req *validation.CreateUser) (*model.User, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		s.Log.Errorf("Failed hash password: %+v", err)
		return nil, err
	}

	user := &model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     req.Role,
	}

	result := s.DB.WithContext(c.Context()).Create(user)

	if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
		return nil, fiber.NewError(fiber.StatusConflict, "Email is already in use")
	}

	if result.Error != nil {
		s.Log.Errorf("Failed to create user: %+v", result.Error)
	}

	return user, result.Error
}

func (s *UserService) UpdateUser(c *fiber.Ctx, req *validation.UpdateUser, id string) (*model.User, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	if req.Email == "" && req.Name == "" && req.Password == "" {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid Request")
	}

	if req.Password != "" {
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			return nil, err
		}
		req.Password = hashedPassword
	}

	updateBody := &model.User{
		Name:     req.Name,
		Password: req.Password,
		Email:    req.Email,
	}

	result := s.DB.WithContext(c.Context()).Where("id = ?", id).Updates(updateBody)

	if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
		return nil, fiber.NewError(fiber.StatusConflict, "Email is already in use")
	}

	if result.RowsAffected == 0 {
		return nil, fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	if result.Error != nil {
		s.Log.Errorf("Failed to update user: %+v", result.Error)
	}

	user, err := s.GetUserByID(c, id)
	if err != nil {
		return nil, err
	}

	return user, result.Error
}

func (s *UserService) UpdatePassOrVerify(c *fiber.Ctx, req *validation.UpdatePassOrVerify, id string) error {
	if err := s.Validate.Struct(req); err != nil {
		return err
	}

	if req.Password == "" && !req.VerifiedEmail {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid Request")
	}

	if req.Password != "" {
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			return err
		}
		req.Password = hashedPassword
	}

	updateBody := &model.User{
		Password:      req.Password,
		VerifiedEmail: req.VerifiedEmail,
	}

	result := s.DB.WithContext(c.Context()).Where("id = ?", id).Updates(updateBody)

	if result.RowsAffected == 0 {
		return fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	if result.Error != nil {
		s.Log.Errorf("Failed to update user password or verifiedEmail: %+v", result.Error)
	}

	return result.Error
}

func (s *UserService) DeleteUser(c *fiber.Ctx, id string) error {
	user := new(model.User)

	result := s.DB.WithContext(c.Context()).Delete(user, "id = ?", id)

	if result.RowsAffected == 0 {
		return fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	if result.Error != nil {
		s.Log.Errorf("Failed to delete user: %+v", result.Error)
	}

	return result.Error
}

func (s *UserService) CreateGoogleUser(c *fiber.Ctx, req *validation.GoogleLogin) (*model.User, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	userFromDB, err := s.GetUserByEmail(c, req.Email)
	if err != nil {
		if err.Error() == "User not found" {
			user := &model.User{
				Name:          req.Name,
				Email:         req.Email,
				VerifiedEmail: req.VerifiedEmail,
			}

			if createErr := s.DB.WithContext(c.Context()).Create(user).Error; createErr != nil {
				s.Log.Errorf("Failed to create user: %+v", createErr)
				return nil, createErr
			}

			return user, nil
		}

		return nil, err
	}

	userFromDB.VerifiedEmail = req.VerifiedEmail
	if updateErr := s.DB.WithContext(c.Context()).Save(userFromDB).Error; updateErr != nil {
		s.Log.Errorf("Failed to update user: %+v", updateErr)
		return nil, updateErr
	}

	return userFromDB, nil
}
