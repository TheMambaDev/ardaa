package routes

import (
	"ardaa/pkg/user"
	"ardaa/web/handlers"
	middleware "ardaa/web/middlewares"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type router struct {
	db  *gorm.DB
	app *fiber.App
}

func NewRouter(db *gorm.DB, app *fiber.App) *router {
	return &router{
		db:  db,
		app: app,
	}
}

func (r *router) Setup() fiber.Router {
	api := r.app.Group("/")

	// Initiate User Resources and Inject Dependencies
	userApi := api.Group("user")
	userRepo := user.NewUserRepo(r.db)
	userService := user.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	userApi.Post("/", userHandler.Register)
	userApi.Post("/login", userHandler.Login)

	// after this middleware, all user api routes will be authenticated
	userApi.Use(middleware.Authenticate)
	userApi.Get("/", userHandler.Me)
	userApi.Post("/logout", userHandler.Logout)
	userApi.Put("/", userHandler.Update)

	return api
}
