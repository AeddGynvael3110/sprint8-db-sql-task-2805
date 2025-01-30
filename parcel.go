package main

import (
	"database/sql"
	"fmt"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	// реализуйте добавление строки в таблицу parcel, используйте данные из переменной p
	queryString := "insert into parcel (client, status, address, created_at) values ($1, $2, $3, $4)"
	req, err := s.db.Exec(queryString, p.Client, p.Status, p.Address, p.CreatedAt)
	if err != nil {
		return 0, fmt.Errorf("возникла ошибка при добавлении посылки: %w", err)
	}

	id, err := req.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("возникла ошибка при получении ID запроса: %w", err)
	}

	// верните идентификатор последней добавленной записи
	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// реализуйте чтение строки по заданному number
	// здесь из таблицы должна вернуться только одна строка
	queryString := "SELECT number, client, status, address, created_at FROM parcel WHERE number = $1"
	respRow := s.db.QueryRow(queryString, number)

	// заполните объект Parcel данными из таблицы
	p := Parcel{}
	err := respRow.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return p, fmt.Errorf("возникла ошибка при получении клиента по номеру: %w", err)
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк

	// заполните срез Parcel данными из таблицы
	var res []Parcel
	queryString := "SELECT number, client, status, address, created_at FROM parcel WHERE client = $1"
	respRows, err := s.db.Query(queryString, client)
	if err != nil {
		return res, fmt.Errorf("возникла ошибка при получении строк по заданному клиенту: %w", err)
	}
	defer respRows.Close()

	for respRows.Next() {
		p := Parcel{}
		err := respRows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			return res, fmt.Errorf("возникла ошибка при получении строки: %w", err)
		}
		res = append(res, p)
	}

	if err := respRows.Err(); err != nil {
		return res, fmt.Errorf("возникла ошибка при итерации по строкам: %w", err)
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel

	queryString := "UPDATE parcel SET status = $1 WHERE number = $2"
	_, err := s.db.Exec(queryString, status, number)
	if err != nil {
		return fmt.Errorf("возникла ошибка при обновлении статуса посылки: %w", err)
	}

	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered

	queryString := "UPDATE parcel SET address = $1 WHERE number = $2 AND status = $3"
	_, err := s.db.Exec(queryString, address, number, ParcelStatusRegistered)
	if err != nil {
		return fmt.Errorf("возникла ошибка при обновлении адреса: %w", err)
	}

	return nil
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered

	queryString := "delete FROM parcel WHERE number = $1 AND status = $2"
	_, err := s.db.Exec(queryString, number, ParcelStatusRegistered)
	if err != nil {
		return fmt.Errorf("возникла ошибка при удалении строки: %w", err)
	}

	return nil
}
