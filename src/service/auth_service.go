package service

import (
	"app/src/config"
	"app/src/model"
	"app/src/response"
	"app/src/utils"
	"app/src/validation"
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AuthService struct {
	log          *logrus.Logger
	db           *gorm.DB
	validate     *validator.Validate
	userService  *UserService
	tokenService *TokenService
}

func NewAuthService(
	db *gorm.DB, validate *validator.Validate, userService *UserService, tokenService *TokenService,
) *AuthService {
	return &AuthService{
		log:          utils.Log,
		db:           db,
		validate:     validate,
		userService:  userService,
		tokenService: tokenService,
	}
}

func (s *AuthService) Register(c *fiber.Ctx, req *validation.Register) (*model.User, error) {
	if err := s.validate.Struct(req); err != nil {
		return nil, err
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		s.log.Errorf("Failed hash password: %+v", err)
		return nil, err
	}

	user := &model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
	}

	result := s.db.WithContext(c.Context()).Create(user)
	if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
		return nil, fiber.NewError(fiber.StatusConflict, "Email already taken")
	}

	if result.Error != nil {
		s.log.Errorf("Failed create user: %+v", result.Error)
	}

	return user, result.Error
}

func (s *AuthService) Login(c *fiber.Ctx, req *validation.Login) (*model.User, error) {
	if err := s.validate.Struct(req); err != nil {
		return nil, err
	}

	user, err := s.userService.GetUserByEmail(c, req.Email)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid email or password")
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid email or password")
	}

	return user, nil
}

func (s *AuthService) Logout(c *fiber.Ctx, req *validation.Logout) error {
	if err := s.validate.Struct(req); err != nil {
		return err
	}

	token, err := s.tokenService.GetTokenByUserID(c, req.RefreshToken)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Token not found")
	}

	err = s.tokenService.DeleteToken(c, config.TokenTypeRefresh, token.UserID.String())

	return err
}

func (s *AuthService) RefreshAuth(c *fiber.Ctx, req *validation.RefreshToken) (*response.Tokens, error) {
	if err := s.validate.Struct(req); err != nil {
		return nil, err
	}

	token, err := s.tokenService.GetTokenByUserID(c, req.RefreshToken)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Please authenticate")
	}

	user, err := s.userService.GetUserByID(c, token.UserID.String())
	if err != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Please authenticate")
	}

	newTokens, err := s.tokenService.GenerateAuthTokens(c, user)
	if err != nil {
		return nil, fiber.ErrInternalServerError
	}

	return newTokens, err
}

func (s *AuthService) ResetPassword(c *fiber.Ctx, query *validation.Token, req *validation.UpdatePassOrVerify) error {
	if err := s.validate.Struct(query); err != nil {
		return err
	}

	userID, err := utils.VerifyToken(query.Token, config.JWTSecret, config.TokenTypeResetPassword)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid Token")
	}

	user, err := s.userService.GetUserByID(c, userID)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Password reset failed")
	}

	if errUpdate := s.userService.UpdatePassOrVerify(c, req, user.ID.String()); errUpdate != nil {
		return errUpdate
	}

	if errToken := s.tokenService.DeleteToken(c, config.TokenTypeResetPassword, user.ID.String()); errToken != nil {
		return errToken
	}

	return nil
}

func (s *AuthService) VerifyEmail(c *fiber.Ctx, query *validation.Token) error {
	if err := s.validate.Struct(query); err != nil {
		return err
	}

	userID, err := utils.VerifyToken(query.Token, config.JWTSecret, config.TokenTypeVerifyEmail)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid Token")
	}

	user, err := s.userService.GetUserByID(c, userID)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Verify email failed")
	}

	if errToken := s.tokenService.DeleteToken(c, config.TokenTypeVerifyEmail, user.ID.String()); errToken != nil {
		return errToken
	}

	updateBody := &validation.UpdatePassOrVerify{
		VerifiedEmail: true,
	}

	if errUpdate := s.userService.UpdatePassOrVerify(c, updateBody, user.ID.String()); errUpdate != nil {
		return errUpdate
	}

	return nil
}
