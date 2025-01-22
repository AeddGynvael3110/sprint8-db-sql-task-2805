package main

import (
	"fmt"

	"github.com/Yandex-Practicum/go-db-sql-final/internal/database"
	"github.com/Yandex-Practicum/go-db-sql-final/internal/models"
	"github.com/Yandex-Practicum/go-db-sql-final/internal/services"
	log "github.com/sirupsen/logrus"
)

func main() {
	// настройте подключение к БД
	db, openDbErr := database.Open()
	defer db.Close()

	if openDbErr != nil {
		log.WithFields(log.Fields{
			"error_message": fmt.Sprintf("scan results error: %v", openDbErr),
		}).Error()
		return
	}

	// создайте объект ParcelStore функцией NewParcelStore
	store := models.NewParcelStore(db)
	service := services.NewParcelService(store)

	// регистрация посылки
	client := 1
	address := "Псков, д. Пушкина, ул. Колотушкина, д. 5"
	p, err := service.Register(client, address)
	if err != nil {
		fmt.Println(err)
		return
	}

	// изменение адреса
	newAddress := "Саратов, д. Верхние Зори, ул. Козлова, д. 25"
	err = service.ChangeAddress(p.Number, newAddress)
	if err != nil {
		fmt.Println(err)
		return
	}

	// изменение статуса
	err = service.NextStatus(p.Number)
	if err != nil {
		fmt.Println(err)
		return
	}

	// вывод посылок клиента
	err = service.PrintClientParcels(client)
	if err != nil {
		fmt.Println(err)
		return
	}

	// попытка удаления отправленной посылки
	err = service.Delete(p.Number)
	if err != nil {
		fmt.Println(err)
		return
	}

	// вывод посылок клиента
	// предыдущая посылка не должна удалиться, т.к. её статус НЕ «зарегистрирована»
	err = service.PrintClientParcels(client)
	if err != nil {
		fmt.Println(err)
		return
	}

	// регистрация новой посылки
	p, err = service.Register(client, address)
	if err != nil {
		fmt.Println(err)
		return
	}

	// удаление новой посылки
	err = service.Delete(p.Number)
	if err != nil {
		fmt.Println(err)
		return
	}

	// вывод посылок клиента
	// здесь не должно быть последней посылки, т.к. она должна была успешно удалиться
	err = service.PrintClientParcels(client)
	if err != nil {
		fmt.Println(err)
		return
	}
}
