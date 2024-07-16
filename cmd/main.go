package main

import (
	"test/api"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	defer api.CloseLogChannel() // Закрытие канала при завершении

	e.POST("/accounts", api.PostAccount())
	e.POST("/accounts/:id/deposit", api.PostDeposit())
	e.POST("/accounts/:id/withdraw", api.PostWithdraw())
	e.GET("/accounts/:id/balance", api.GetBalance())

	e.Logger.Fatal(e.Start(":8080"))
}
