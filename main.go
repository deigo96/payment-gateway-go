package main

import (
	"net/http"
	"os"
	"os/signal"
	"topup-service/config"
	"topup-service/modules"
	"topup-service/router"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	con := config.GetConfig()
	dbCon := config.NewDatabaseConnection(con)
	controllers := modules.RegisterModules(dbCon, con)

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAccessControlAllowHeaders},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete, http.MethodOptions},
	}))
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	router.Router(e, &controllers)

	server := config.GetServer()

	e.Logger.Fatal(e.Start(server.Host + server.Port))
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
}
