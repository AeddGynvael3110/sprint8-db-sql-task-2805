package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3" // Используйте драйвер вашей БД
)

func main() {
	// Подключение к базе данных SQLite (или другой БД)
	db, err := sql.Open("sqlite3", "parcel.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Создание хранилища посылок
	store := NewParcelStore(db)

	// Добавление новой посылки
	parcel := Parcel{
		Client:    123,
		Status:    "registered",
		Address:   "ул. Пушкина, д. Колотушкина",
		CreatedAt: "2024-06-18",
	}

	id, err := store.Add(parcel) // Вызов метода Add из вашего кода
	if err != nil {
		fmt.Println("Ошибка при добавлении посылки:", err)
	} else {
		fmt.Println("Посылка добавлена с ID:", id)
	}
}
