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
	insertQuery := "insert into parcel (client, status, address, created_at) values ($1, $2, $3, $4)"
	req, err := s.db.Exec(insertQuery, p.Client, p.Status, p.Address, p.CreatedAt)
	if err != nil {
		return 0, fmt.Errorf("ошибка при добавлении посылки = %w", err)

	}
	id, err := req.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("ошибка при получении добавленного id = %w", err)
	}
	// верните идентификатор последней добавленной записи
	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// реализуйте чтение строки по заданному number
	// здесь из таблицы должна вернуться только одна строка

	// заполните объект Parcel данными из таблицы
	p := Parcel{}
	selectQuery := "select number, client, status, address, created_at from parcel where number = $1"
	res := s.db.QueryRow(selectQuery, number)
	err := res.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return p, fmt.Errorf("ошибка при получении клиента по номеру: %w", err)

	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк

	// заполните срез Parcel данными из таблицы
	var res []Parcel
	selectQuery := ("select number, client, status, address, created_at from parcel where client = $1")
	rows, err := s.db.Query(selectQuery, client)
	if err != nil {
		return res, fmt.Errorf("ошибка при чтении по заданному client: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		p := Parcel{}
		err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			return res, fmt.Errorf("ошибка при заполнении среза: %w", err)
		}
		res = append(res, p)
	}
	if err := rows.Err(); err != nil {
		return res, fmt.Errorf("ошибка rows.Next(): %w", err)
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel
	updateQuery := "update parcel set status = $1 where number = $2"
	_, err := s.db.Exec(updateQuery, status, number)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении статуса: %w", err)
	}
	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
	updateQuery := "update parcel set address = $1 where number = $2 and status = 'registered'"
	_, err := s.db.Exec(updateQuery, address, number)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении статуса: %w", err)
	}
	return nil
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered
	deleteQuery := "delete from parcel where number = $1 and status = 'registered'"
	_, err := s.db.Exec(deleteQuery, number)
	if err != nil {
		return fmt.Errorf("ошибка при удалении строки: %w", err)
	}
	return nil
}
