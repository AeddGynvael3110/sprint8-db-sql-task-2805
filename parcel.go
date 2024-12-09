package main

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	// реализуйте добавление строки в таблицу parcel, используйте данные из переменной p
	res, err := s.db.Exec("INSERT INTO parsel (client, status, address, createdAt) VALUES (:client, :status, :address, :createdAt)",
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("createdAt", p.CreatedAt))
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	number, err := res.LastInsertId()
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	// верните идентификатор последней добавленной записи
	return int(number), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// реализуйте чтение строки по заданному number
	// здесь из таблицы должна вернуться только одна строка
	row := s.db.QueryRow("SELECT id, client, status, address, createdAt FROM parsel WHERE id = :id",
		sql.Named("id", number))

	var (
		ID        int
		client    int
		status    string
		address   string
		createdAt string
	)

	err := row.Scan(&number, &client, &status, &address, &createdAt)
	if err != nil {
		fmt.Println(err)
		return Parcel{}, err
	}
	// заполните объект Parcel данными из таблицы
	p := Parcel{
		Number:    ID,
		Client:    client,
		Status:    status,
		Address:   address,
		CreatedAt: createdAt,
	}
	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк
	rows, err := s.db.Query("SELECT number, client, status, address, createdAt FROM parser WHERE client = :client", sql.Named("client", client))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()
	// заполните срез Parcel данными из таблицы
	var res []Parcel
	for rows.Next() {
		p := Parcel{}

		err := rows.Scan(&p.Number, &p.Client, p.Status, p.Address, p.CreatedAt)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		res = append(res, p)
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel
	_, err := s.db.Exec("UPDATE parser SET status = :status WHERE id = :id",
		sql.Named("status", status),
		sql.Named("id", number))
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
	row := s.db.QueryRow("SELECT status, address FROM parsel WHERE id = :id",
		sql.Named("id", number))

	var status string

	err := row.Scan(&status)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if status == "registered" {
		_, err := s.db.Exec("UPDATE parser SET address = :address WHERE id = :id",
			sql.Named("address", address),
			sql.Named("id", number))
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	return nil
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered
	row := s.db.QueryRow("SELECT status, address FROM parsel WHERE id = :id",
		sql.Named("id", number))

	var status string

	err := row.Scan(&status)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if status == "registered" {
		_, err := s.db.Exec("DELETE FROM parser WHERE id = :id", sql.Named("id", number))
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	return nil
}
