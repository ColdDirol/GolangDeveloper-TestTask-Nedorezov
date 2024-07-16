package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"test/model"

	"github.com/labstack/echo/v4"
)

var accounts = make(map[int]*model.Account)
var logChannel = make(chan string)

func init() { // запуск после импорта пакета
	go func() {
		for logMsg := range logChannel {
			slog.Info(logMsg)
		}
	}()
}

func PostAccount() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := len(accounts) + 1
		accounts[id] = model.NewAccount(id)
		logChannel <- fmt.Sprintf("Account [ID %d] created", id)
		return c.NoContent(http.StatusCreated)
	}
}

func PostDeposit() echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		account, ok := accounts[id]
		if !ok {
			return c.JSON(http.StatusNotFound, "Account not found")
		}
		var amount float64
		err = json.NewDecoder(c.Request().Body).Decode(&amount)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		go func() {
			err := account.Deposit(amount)
			if err != nil {
				logChannel <- fmt.Sprintf("[ID %d] Deposit failed: %s", id, err.Error())
			} else {
				logChannel <- fmt.Sprintf("[ID %d] Deposit: +%f, Balance: %f", id, amount, account.GetBalance())
			}
		}()
		return c.NoContent(http.StatusAccepted)
	}
}

func PostWithdraw() echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		account, ok := accounts[id]
		if !ok {
			return c.JSON(http.StatusNotFound, "Account not found")
		}
		var amount float64
		err = json.NewDecoder(c.Request().Body).Decode(&amount)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		go func() {
			err := account.Withdraw(amount)
			if err != nil {
				logChannel <- fmt.Sprintf("[ID %d] Withdraw failed: %s", id, err.Error())
			} else {
				logChannel <- fmt.Sprintf("[ID %d] Withdraw: -%f, Balance: %f", id, amount, account.GetBalance())
			}
		}()
		return c.NoContent(http.StatusAccepted)
	}
}

func GetBalance() echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		account, ok := accounts[id]
		if !ok {
			return c.JSON(http.StatusNotFound, "Account not found")
		}

		balanceChannel := make(chan float64)
		go func() {
			balanceChannel <- account.GetBalance()
		}()

		balance := <-balanceChannel
		logChannel <- fmt.Sprintf("[ID %d] Balance: %f", id, balance)
		return c.JSON(http.StatusOK, balance)
	}
}

func CloseLogChannel() {
	close(logChannel)
}
