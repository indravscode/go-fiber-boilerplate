package controller

import (
	"app/src/config"
	"app/src/model"
	"app/src/response"
	"app/src/service"
	"app/src/validation"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AuthController struct {
	authService  *service.AuthService
	userService  *service.UserService
	tokenService *service.TokenService
	emailService *service.EmailService
}

func NewAuthController(
	authService *service.AuthService, userService *service.UserService,
	tokenService *service.TokenService, emailService *service.EmailService,
) *AuthController {
	return &AuthController{
		authService,
		userService,
		tokenService,
		emailService,
	}
}

// @Tags         Auth
// @Summary      Register as user
// @Accept       json
// @Produce      json
// @Param        request  body  validation.Register  true  "Request body"
// @Router       /auth/register [post]
// @Success      201  {object}  example.RegisterResponse
// @Failure      409  {object}  example.DuplicateEmail  "Email already taken"
func (a *AuthController) Register(c *fiber.Ctx) error {
	req := new(validation.Register)

	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	user, err := a.authService.Register(c, req)
	if err != nil {
		return err
	}

	tokens, err := a.tokenService.GenerateAuthTokens(c, user)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).
		JSON(response.SuccessWithTokens{
			Code:    fiber.StatusCreated,
			Status:  "success",
			Message: "Register successfully",
			User:    *user,
			Tokens:  *tokens,
		})
}

// @Tags         Auth
// @Summary      Login
// @Accept       json
// @Produce      json
// @Param        request  body  validation.Login  true  "Request body"
// @Router       /auth/login [post]
// @Success      200  {object}  example.LoginResponse
// @Failure      401  {object}  example.FailedLogin  "Invalid email or password"
func (a *AuthController) Login(c *fiber.Ctx) error {
	req := new(validation.Login)

	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	user, err := a.authService.Login(c, req)
	if err != nil {
		return err
	}

	tokens, err := a.tokenService.GenerateAuthTokens(c, user)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).
		JSON(response.SuccessWithTokens{
			Code:    fiber.StatusOK,
			Status:  "success",
			Message: "Login successfully",
			User:    *user,
			Tokens:  *tokens,
		})
}

// @Tags         Auth
// @Summary      Logout
// @Accept       json
// @Produce      json
// @Param        request  body  example.RefreshToken  true  "Request body"
// @Router       /auth/logout [post]
// @Success      200  {object}  example.LogoutResponse
// @Failure      404  {object}  example.NotFound  "Not found"
func (a *AuthController) Logout(c *fiber.Ctx) error {
	req := new(validation.Logout)

	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if err := a.authService.Logout(c, req); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).
		JSON(response.Common{
			Code:    fiber.StatusOK,
			Status:  "success",
			Message: "Logout successfully",
		})
}

// @Tags         Auth
// @Summary      Refresh auth tokens
// @Accept       json
// @Produce      json
// @Param        request  body  example.RefreshToken  true  "Request body"
// @Router       /auth/refresh-tokens [post]
// @Success      200  {object}  example.RefreshTokenResponse
// @Failure      401  {object}  example.Unauthorized  "Unauthorized"
func (a *AuthController) RefreshTokens(c *fiber.Ctx) error {
	req := new(validation.RefreshToken)

	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	tokens, err := a.authService.RefreshAuth(c, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).
		JSON(response.RefreshToken{
			Code:   fiber.StatusOK,
			Status: "success",
			Tokens: *tokens,
		})
}

// @Tags         Auth
// @Summary      Forgot password
// @Description  An email will be sent to reset password.
// @Accept       json
// @Produce      json
// @Param        request  body  validation.ForgotPassword  true  "Request body"
// @Router       /auth/forgot-password [post]
// @Success      200  {object}  example.ForgotPasswordResponse
// @Failure      404  {object}  example.NotFound  "Not found"
func (a *AuthController) ForgotPassword(c *fiber.Ctx) error {
	req := new(validation.ForgotPassword)

	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	resetPasswordToken, err := a.tokenService.GenerateResetPasswordToken(c, req)
	if err != nil {
		return err
	}

	if errEmail := a.emailService.SendResetPasswordEmail(req.Email, resetPasswordToken); errEmail != nil {
		return errEmail
	}

	return c.Status(fiber.StatusOK).
		JSON(response.Common{
			Code:    fiber.StatusOK,
			Status:  "success",
			Message: "A password reset link has been sent to your email address.",
		})
}

// @Tags         Auth
// @Summary      Reset password
// @Accept       json
// @Produce      json
// @Param        token   query  string  true  "The reset password token"
// @Param        request  body  validation.UpdatePassOrVerify  true  "Request body"
// @Router       /auth/reset-password [post]
// @Success      200  {object}  example.ResetPasswordResponse
// @Failure      401  {object}  example.FailedResetPassword  "Password reset failed"
func (a *AuthController) ResetPassword(c *fiber.Ctx) error {
	req := new(validation.UpdatePassOrVerify)
	query := &validation.Token{
		Token: c.Query("token"),
	}

	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if err := a.authService.ResetPassword(c, query, req); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).
		JSON(response.Common{
			Code:    fiber.StatusOK,
			Status:  "success",
			Message: "Update password successfully",
		})
}

// @Tags         Auth
// @Summary      Send verification email
// @Description  An email will be sent to verify email.
// @Security BearerAuth
// @Produce      json
// @Router       /auth/send-verification-email [post]
// @Success      200  {object}  example.SendVerificationEmailResponse
// @Failure      401  {object}  example.Unauthorized  "Unauthorized"
func (a *AuthController) SendVerificationEmail(c *fiber.Ctx) error {
	user, _ := c.Locals("user").(*model.User)

	verifyEmailToken, err := a.tokenService.GenerateVerifyEmailToken(c, user)
	if err != nil {
		return err
	}

	if errEmail := a.emailService.SendVerificationEmail(user.Email, *verifyEmailToken); errEmail != nil {
		return errEmail
	}

	return c.Status(fiber.StatusOK).
		JSON(response.Common{
			Code:    fiber.StatusOK,
			Status:  "success",
			Message: "Please check your email for a link to verify your account",
		})
}

// @Tags         Auth
// @Summary      Verify email
// @Produce      json
// @Param        token   query  string  true  "The verify email token"
// @Router       /auth/verify-email [post]
// @Success      200  {object}  example.VerifyEmailResponse
// @Failure      401  {object}  example.FailedVerifyEmail  "Verify email failed"
func (a *AuthController) VerifyEmail(c *fiber.Ctx) error {
	query := &validation.Token{
		Token: c.Query("token"),
	}

	if err := a.authService.VerifyEmail(c, query); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).
		JSON(response.Common{
			Code:    fiber.StatusOK,
			Status:  "success",
			Message: "Verify email successfully",
		})
}

// @Tags         Auth
// @Summary      Login with google
// @Description  This route initiates the Google OAuth2 login flow. Please try this in your browser.
// @Router       /auth/google [get]
// @Success      200  {object}  example.GoogleLoginResponse
func (a *AuthController) GoogleLogin(c *fiber.Ctx) error {
	// Generate a random state
	state := uuid.New().String()

	c.Cookie(&fiber.Cookie{
		Name:   "oauth_state",
		Value:  state,
		MaxAge: 30,
	})

	url := config.AppConfig.GoogleLoginConfig.AuthCodeURL(state)

	return c.Status(fiber.StatusSeeOther).Redirect(url)
}

func (a *AuthController) GoogleCallback(c *fiber.Ctx) error {
	state := c.Query("state")
	storedState := c.Cookies("oauth_state")

	if state != storedState {
		return fiber.NewError(fiber.StatusUnauthorized, "States don't Match!")
	}

	code := c.Query("code")
	googlecon := config.GoogleConfig()

	token, err := googlecon.Exchange(context.Background(), code)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(
		c.Context(), http.MethodGet,
		"https://www.googleapis.com/oauth2/v2/userinfo?access_token="+token.AccessToken,
		nil,
	)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	userData, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	googleUser := new(validation.GoogleLogin)
	if errJSON := json.Unmarshal(userData, googleUser); errJSON != nil {
		return errJSON
	}

	user, err := a.userService.CreateGoogleUser(c, googleUser)
	if err != nil {
		return err
	}

	tokens, err := a.tokenService.GenerateAuthTokens(c, user)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).
		JSON(response.SuccessWithTokens{
			Code:    fiber.StatusOK,
			Status:  "success",
			Message: "Login successfully",
			User:    *user,
			Tokens:  *tokens,
		})

	// TODO: replace this url with the link to the oauth google success page of your front-end app
	// googleLoginURL := fmt.Sprintf("http://link-to-app/google/success?access_token=%s&refresh_token=%s",
	// 	tokens.Access.Token, tokens.Refresh.Token)

	// return c.Status(fiber.StatusSeeOther).Redirect(googleLoginURL)
}
