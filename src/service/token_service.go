package service

import (
	"app/src/config"
	"app/src/model"
	res "app/src/response"
	"app/src/utils"
	"app/src/validation"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TokenService struct {
	Log         *logrus.Logger
	DB          *gorm.DB
	Validate    *validator.Validate
	UserService *UserService
}

func NewTokenService(db *gorm.DB, validate *validator.Validate, userService *UserService) *TokenService {
	return &TokenService{
		Log:         utils.Log,
		DB:          db,
		Validate:    validate,
		UserService: userService,
	}
}

func (s *TokenService) GenerateToken(userID string, expires time.Time, tokenType string) (string, error) {
	claims := jwt.MapClaims{
		"sub":  userID,
		"iat":  time.Now().Unix(),
		"exp":  expires.Unix(),
		"type": tokenType,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(config.JWTSecret))
}

func (s *TokenService) SaveToken(c *fiber.Ctx, token, userID, tokenType string, expires time.Time) error {
	if err := s.DeleteToken(c, tokenType, userID); err != nil {
		return err
	}

	tokenDoc := &model.Token{
		Token:   token,
		UserID:  uuid.MustParse(userID),
		Type:    tokenType,
		Expires: expires,
	}

	result := s.DB.WithContext(c.Context()).Create(tokenDoc)

	if result.Error != nil {
		s.Log.Errorf("Failed save token: %+v", result.Error)
	}

	return result.Error
}

func (s *TokenService) DeleteToken(c *fiber.Ctx, tokenType string, userID string) error {
	tokenDoc := new(model.Token)

	result := s.DB.WithContext(c.Context()).
		Where("type = ? AND user_id = ?", tokenType, userID).
		Delete(tokenDoc)

	if result.Error != nil {
		s.Log.Errorf("Failed to delete token: %+v", result.Error)
	}

	return result.Error
}

func (s *TokenService) DeleteAllToken(c *fiber.Ctx, userID string) error {
	tokenDoc := new(model.Token)

	result := s.DB.WithContext(c.Context()).Where("user_id = ?", userID).Delete(tokenDoc)

	if result.Error != nil {
		s.Log.Errorf("Failed to delete all token: %+v", result.Error)
	}

	return result.Error
}

func (s *TokenService) GetTokenByUserID(c *fiber.Ctx, tokenStr string) (*model.Token, error) {
	userID, err := utils.VerifyToken(tokenStr, config.JWTSecret, config.TokenTypeRefresh)
	if err != nil {
		return nil, err
	}

	tokenDoc := new(model.Token)

	result := s.DB.WithContext(c.Context()).
		Where("token = ? AND user_id = ?", tokenStr, userID).
		First(tokenDoc)

	if result.Error != nil {
		s.Log.Errorf("Failed get token by user id: %+v", err)
		return nil, result.Error
	}

	return tokenDoc, nil
}

func (s *TokenService) GenerateAuthTokens(c *fiber.Ctx, user *model.User) (*res.Tokens, error) {
	accessTokenExpires := time.Now().UTC().Add(time.Minute * time.Duration(config.JWTAccessExp))
	accessToken, err := s.GenerateToken(user.ID.String(), accessTokenExpires, config.TokenTypeAccess)
	if err != nil {
		s.Log.Errorf("Failed generate token: %+v", err)
		return nil, err
	}

	refreshTokenExpires := time.Now().UTC().Add(time.Hour * 24 * time.Duration(config.JWTRefreshExp))
	refreshToken, err := s.GenerateToken(user.ID.String(), refreshTokenExpires, config.TokenTypeRefresh)
	if err != nil {
		s.Log.Errorf("Failed generate token: %+v", err)
		return nil, err
	}

	if err = s.SaveToken(c, refreshToken, user.ID.String(), config.TokenTypeRefresh, refreshTokenExpires); err != nil {
		return nil, err
	}

	return &res.Tokens{
		Access: res.TokenExpires{
			Token:   accessToken,
			Expires: accessTokenExpires,
		},
		Refresh: res.TokenExpires{
			Token:   refreshToken,
			Expires: refreshTokenExpires,
		},
	}, nil
}

func (s *TokenService) GenerateResetPasswordToken(c *fiber.Ctx, req *validation.ForgotPassword) (string, error) {
	if err := s.Validate.Struct(req); err != nil {
		return "", err
	}

	user, err := s.UserService.GetUserByEmail(c, req.Email)
	if err != nil {
		return "", err
	}

	expires := time.Now().UTC().Add(time.Minute * time.Duration(config.JWTResetPasswordExp))
	resetPasswordToken, err := s.GenerateToken(user.ID.String(), expires, config.TokenTypeResetPassword)
	if err != nil {
		s.Log.Errorf("Failed generate token: %+v", err)
		return "", err
	}

	if err = s.SaveToken(c, resetPasswordToken, user.ID.String(), config.TokenTypeResetPassword, expires); err != nil {
		return "", err
	}

	return resetPasswordToken, nil
}

func (s *TokenService) GenerateVerifyEmailToken(c *fiber.Ctx, user *model.User) (*string, error) {
	expires := time.Now().UTC().Add(time.Minute * time.Duration(config.JWTVerifyEmailExp))
	verifyEmailToken, err := s.GenerateToken(user.ID.String(), expires, config.TokenTypeVerifyEmail)
	if err != nil {
		s.Log.Errorf("Failed generate token: %+v", err)
		return nil, err
	}

	if err = s.SaveToken(c, verifyEmailToken, user.ID.String(), config.TokenTypeVerifyEmail, expires); err != nil {
		return nil, err
	}

	return &verifyEmailToken, nil
}
