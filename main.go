package main

import (
	"preview/config"
	"preview/handler"
	"preview/helper"
	"preview/logger"
	"preview/models"

	"github.com/labstack/echo/v4"
)

func main() {
	db := config.Connect()
	handler := &handler.Repo{DB: db}
	db.AutoMigrate(&models.User{})
	e := echo.New()
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			logger.Logging(c).Info("Endpoint called")
			return next(c)
		}
	})

	e.POST("v1/ms-paylater/login", handler.Login)
	e.POST("v1/ms-paylater/register", handler.Register)
	service := e.Group("/v1/ms-paylater")
	service.Use(helper.Auth)
	{
		service.POST("/loan", handler.Loan)
		service.GET("/limit", handler.Limit)
		service.POST("/tarik-saldo", handler.DrawBalance)
		service.POST("/pay", handler.Pay)
	}

	e.Logger.Fatal(e.Start(":8080"))

}
