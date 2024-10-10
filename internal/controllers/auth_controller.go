package controllers

import (
	"fmt"
	"mintyplex-api/internal/models"
	"mintyplex-api/internal/services"
	"mintyplex-api/internal/utils"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthController struct {
	authService services.AuthService
	userService services.UserService
}

func NewAuthController(authService services.AuthService, userService services.UserService) AuthController {
	return AuthController{authService, userService}
}

func (ac *AuthController) SignUpUser(c *fiber.Ctx) error {
	var user models.SignUpUser

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	if user.Password != user.PasswordConfirm {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Passwords do not match"})
	}

	newUser, err := ac.authService.SignUpUser(&user)
	if err != nil {
		if strings.Contains(err.Error(), "email already exist") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"status": "error", "message": err.Error()})
		}
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	dbResponse := &models.DBResponse{
		ID:        newUser.ID,
		Name: newUser.Name,
		Email:     newUser.Email,
		CreatedAt: newUser.CreatedAt,
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"user": models.FilteredResponse(dbResponse),
		},
	})
}

func (ac *AuthController) SignInUser(c *fiber.Ctx) error {
	var credentials models.SignInInput

	if err := c.BodyParser(&credentials); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	user, err := ac.userService.GetUserByEmail(credentials.Email)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Invalid email or password"})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	// Generate Tokens
	expirationDuration, err := time.ParseDuration(os.Getenv("ACCESS_TOKEN_EXPIRED_IN"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail", "message": "Invalid token expiration format"})
	}

	accessToken, err := utils.CreateToken(expirationDuration, user.ID.Hex(), os.Getenv("ACCESS_TOKEN_PRIVATE_KEY"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	refreshTokenPrivateKey := os.Getenv("REFRESH_TOKEN_PRIVATE_KEY")
	refreshTokenExpiredIn, err := time.ParseDuration(os.Getenv("REFRESH_TOKEN_EXPIRED_IN"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail", "message": "Invalid refresh token expiration format"})
	}

	refreshToken, err := utils.CreateToken(refreshTokenExpiredIn, user.ID.Hex(), refreshTokenPrivateKey)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	accessTokenMaxAge, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_MAXAGE"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail", "message": "Invalid AccessTokenMaxAge value"})
	}

	refreshTokenMaxAge, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_MAXAGE"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail", "message": "Invalid RefreshTokenMaxAge value"})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		MaxAge:   accessTokenMaxAge * 60,
		Path:     "/",
		Domain:   "localhost",
		Secure:   false,
		HTTPOnly: true,
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		MaxAge:   refreshTokenMaxAge * 60,
		Path:     "/",
		Domain:   "localhost",
		Secure:   false,
		HTTPOnly: true,
	})

	c.Cookie(&fiber.Cookie{
		Name:     "logged_in",
		Value:    "true",
		MaxAge:   accessTokenMaxAge * 60,
		Path:     "/",
		Domain:   "localhost",
		Secure:   false,
		HTTPOnly: false,
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "access_token": accessToken})
}

func (ac *AuthController) RefreshAccessToken(c *fiber.Ctx) error {
	message := "Could not refresh access token"

	cookie := c.Cookies("refresh_token")
	if cookie == "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "message": message})
	}

	refreshTokenPublicKey := os.Getenv("REFRESH_TOKEN_PUBLIC_KEY")
	sub, err := utils.ValidateToken(cookie, refreshTokenPublicKey)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	user, err := ac.userService.GetUserById(fmt.Sprint(sub))
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "message": "The user belonging to this token no longer exists"})
	}

	accessTokenExpiresIn, err := time.ParseDuration(os.Getenv("ACCESS_TOKEN_EXPIRED_IN"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail", "message": "Invalid token expiration format"})
	}

	accessTokenPrivateKey := os.Getenv("ACCESS_TOKEN_PRIVATE_KEY")
	accessToken, err := utils.CreateToken(accessTokenExpiresIn, user.ID.Hex(), accessTokenPrivateKey)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	accessTokenMaxAge, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_MAXAGE"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail", "message": "Invalid AccessTokenMaxAge value"})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		MaxAge:   accessTokenMaxAge * 60,
		Path:     "/",
		Domain:   "localhost",
		Secure:   false,
		HTTPOnly: true,
	})

	c.Cookie(&fiber.Cookie{
		Name:     "logged_in",
		Value:    "true",
		MaxAge:   accessTokenMaxAge * 60,
		Path:     "/",
		Domain:   "localhost",
		Secure:   false,
		HTTPOnly: false,
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "access_token": accessToken})
}

func (ac *AuthController) GoogleAuth(c *fiber.Ctx) error {
	code := c.Query("code")
	pathUrl := "/"

	if c.Query("state") != "" {
		pathUrl = c.Query("state")
	}

	if code == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": "authorization code not provided ",
		})
	}

	tokenRes, err := utils.GetGoogleAuthToken(code)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"status":  "fail",
			"message": "line 170 [auth_controller.go]" + err.Error(),
		})
	}

	user, err := utils.GetGoogleUser(tokenRes.Access_token, tokenRes.ID_token)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"status":  "fail",
			"message": "line 182 [auth_controller.go]" + err.Error(),
		})
	}

	createdAt := time.Now()
	resBody := *&models.UpdateDBUser{
		Email:     user.Email,
		Name:      user.Name,
		Provider:  "google",
		Verified:  true,
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
	}

	updatedUser, err := ac.userService.UpsertUser(user.Email, &resBody)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"status":  "fail",
			"message": "line 196 [auth_controller.go]" + err.Error(),
		})
	}

	accessTokenExpiresInInv, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_EXPIRED_IN"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid access token expiration time" + err.Error(),
		})
	}
	accessTokenExpiresIn := time.Duration(accessTokenExpiresInInv) * time.Second

	accessTokenPrivateKey := os.Getenv("ACCESS_TOKEN_PRIVATE_KEY")

	accessToken, err := utils.CreateToken(accessTokenExpiresIn, updatedUser.ID.Hex(), accessTokenPrivateKey)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	refreshTokenExpiresInInv, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_EXPIRED_IN"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid refresh token expiration time" + err.Error(),
		})
	}
	refreshTokenExpiresIn := time.Duration(refreshTokenExpiresInInv) * time.Second

	refreshTokenPrivateKey := os.Getenv("REFRESH_TOKEN_PRIVATE_KEY")

	refreshToken, err := utils.CreateToken(refreshTokenExpiresIn, updatedUser.ID.Hex(), refreshTokenPrivateKey)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"status":  "fail",
			"message": "<>" + err.Error(),
		})
	}

	accessTokenMaxAge, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_MAXAGE"))

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		MaxAge:   accessTokenMaxAge * 60,
		Path:     "/",
		Domain:   "localhost",
		Secure:   false,
		HTTPOnly: true,
	})

	refreshTokenMaxAge, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_MAXAGE"))

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		MaxAge:   refreshTokenMaxAge * 60,
		Path:     "/",
		Domain:   "localhost",
		Secure:   false,
		HTTPOnly: true,
	})

	c.Cookie(&fiber.Cookie{
		Name:     "logged_in",
		Value:    "true",
		MaxAge:   accessTokenMaxAge * 60,
		Path:     "/",
		Domain:   "localhost",
		Secure:   false,
		HTTPOnly: false,
	})

	clientOrigin := os.Getenv("CLIENT_ORIGIN")

	return c.Redirect(fmt.Sprint(clientOrigin, pathUrl), fiber.StatusTemporaryRedirect)

}
