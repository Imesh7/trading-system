package account_handler

import (
	"account-service/account"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type accountHandler struct {
	db *gorm.DB
}

func NewHanlder(db *gorm.DB) accountHandler {
	return accountHandler{
		db: db,
	}
}

func (a accountHandler) GetAccountBalance(c *fiber.Ctx) error {
	accountData := account.Account{
		UserId:    2,
		Amount:    1000.00,
		CreatedAt: time.Now(),
	}

	return c.Status(200).JSON(accountData)
}
