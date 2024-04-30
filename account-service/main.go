package main

import (
	"account-service/account"
	"account-service/database"
	account_handler "account-service/handler"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {

	app := fiber.New()
	db := database.ConnectDatabase()
	account_hanlder := account_handler.NewHanlder(db)
	db.AutoMigrate(account.Account{})
	app.Get("/get-account-balance", account_hanlder.GetAccountBalance)
	log.Fatal(app.Listen(":8000"))
}
