package handlers

import (
	"ardaa/domain"
	"errors"
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

type userHandler struct {
	service domain.UserService
}

func NewUserHandler(service domain.UserService) *userHandler {
	return &userHandler{
		service: service,
	}
}

func (u *userHandler) Login(ctx *fiber.Ctx) error {
	// check if the user is already logged in
	if ctx.Get("Authorization") != "" {
		token := ctx.Get("Authorization")

		// check if the user is logged in with the same token
		if u.service.IsLoggedIn(token) {
			return ctx.JSON(fiber.Map{
				"message": "Already logged in",
			})
		}

		// maybe the token is expired and is being used
		if u.service.IsLoggedIn(ctx.IP()) {
			return ctx.JSON(fiber.Map{
				"message": "Already logged in",
			})
		}
	}

	var user domain.LoginUser
	err := ctx.BodyParser(&user)
	if err != nil {
		ctx.Status(400)
		return ctx.JSON(fiber.Map{
			"error": "Empty request body",
		})
	}

	user.Ip = []byte(ctx.IP())
	res, err := u.service.Login(&user)
	if err != nil {
		ctx.Status(400)
		return ctx.JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	ctx.Status(200)
	return ctx.JSON(res)
}

func (u *userHandler) Register(ctx *fiber.Ctx) error {
	var user domain.RegisterUser
	err := ctx.BodyParser(&user)
	if err != nil {

		ctx.Status(400)
		return ctx.JSON(fiber.Map{
			"error": "Empty request body",
		})
	}

	user.Ip = []byte(ctx.IP())
	registeredUser, err := u.service.Register(user)
	if err != nil {
		slog.Error("Registeration error", err)

		ctx.Status(400)
		return ctx.JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	ctx.Status(201)
	return ctx.JSON(registeredUser)
}

// Update user profile
func (u *userHandler) Update(ctx *fiber.Ctx) error {
	return errors.New("not implemented")
}

// User profile
func (u *userHandler) Me(ctx *fiber.Ctx) error {
	user_id := ctx.Locals("user").(string)

	res, err := u.service.Me(user_id)
	if err != nil {
		ctx.Status(400)
		return ctx.JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	ctx.Status(200)
	return ctx.JSON(res)
}

// logout user
func (u *userHandler) Logout(ctx *fiber.Ctx) error {
	user_id := ctx.Locals("user").(string)
	err := u.service.Logout(user_id)
	if err != nil {
		ctx.Status(400)
		return ctx.JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	ctx.Status(200)
	return ctx.JSON(fiber.Map{
		"message": "Logout successful",
	})
}
