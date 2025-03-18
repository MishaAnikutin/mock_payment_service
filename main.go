package main

import (
	"example.com/m/src/common"
	"example.com/m/src/infrastructure"
	"example.com/m/src/presentation"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Инициализация зависимостей
	deps, err := common.Inject()
	if err != nil {
		panic(err)
	}

	// Выполнение миграций
	if err := infrastructure.UpgradeHead(deps.DB); err != nil {
		panic(err)
	}

	// Запуск HTTP сервера
	router := presentation.RegisterRouter(deps.TransferHandlers)
	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
