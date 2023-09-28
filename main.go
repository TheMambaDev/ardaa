package main

import (
	"ardaa/config"
	"ardaa/pkg/database"
	"ardaa/web/routes"
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func init() {
	config.Setup()
}

func main() {
	defer config.LogFile.Close()

	app := fiber.New()

	// setting the default logger of fiber
	fiber_log, err := os.OpenFile("logs/fiber.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic("could not open the fiber log file")
	}
	defer fiber_log.Close()

	// use middlewares here
	app.Use(logger.New(logger.Config{
		Format:     "{\"time\":\"${time}\",\"status\":\"${status}\",\"method\":\"${method}\",\"path\":\"${path}\",\"ip\":\"${ip}\",\"latency\":\"${latency}\"}\n",
		TimeFormat: "2006-01-02 15:04:05",
		Output:     fiber_log,
	}))

	// let's connect the database now
	db, err := database.NewDatabase(os.Getenv("MYSQL_DSN")).ConnectMysql() // can't connect sqlite for now because of DB Enum issues
	if err != nil {
		panic(err)
	}

	_ = database.Automigrate(db)

	// routes on the stove
	router := routes.NewRouter(db, app)
	router.Setup()

	slog.Error("Main: ", "Server Error: ", app.Listen(":3000"))
}
