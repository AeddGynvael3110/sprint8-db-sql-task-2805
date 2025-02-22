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
	query := `INSERT INTO parcel (client, address, status, created_at) VALUES(?, ?, ?, ?)`
	result, err := s.db.Exec(query, p.Client, p.Address, p.Status, p.CreatedAt)
	if err != nil {
		return 0, fmt.Errorf("ошибка добавления %w", err)
	}
	// реализуйте добавление строки в таблицу parcel, используйте данные из переменной p
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("ошибка получения последнего id %w", err)
	}
	// верните идентификатор последней добавленной записи
	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// реализуйте чтение строки по заданному number
	query := `SELECT number, client, status, address FROM parcel WHERE number = ?`
	row := s.db.QueryRow(query, number)

	// здесь из таблицы должна вернуться только одна строка

	// заполните объект Parcel данными из таблицы
	p := Parcel{}
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address)
	if err != nil {
		if err == sql.ErrNoRows {
			return Parcel{}, fmt.Errorf("посылка с номером %d не найдена", number)
		}
		return Parcel{}, fmt.Errorf("ошибка получения посылки %w", err)
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк
	query := `SELECT number, client, status, address FROM parcel WHERE client = ?`
	rows, err := s.db.Query(query, client)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения посылок %w", err)
	}
	defer rows.Close()

	// заполните срез Parcel данными из таблицы
	var res []Parcel
	for rows.Next() {
		p := Parcel{}
		err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address)
		if err != nil {
			return nil, fmt.Errorf("ошибка скана %w", err)
		}
		res = append(res, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка обработки %w", err)
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel
	query := `UPDATE parcel SET status = ? WHERE number = ?`
	_, err := s.db.Exec(query, status, number)
	if err != nil {
		return fmt.Errorf("ошибка обновления статуса %w", err)
	}
	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
	query := `UPDATE parcel SET address = ? WHERE number = ? AND status = 'registered'`
	_, err := s.db.Exec(query, address, number)
	if err != nil {
		return fmt.Errorf("ошибка обновления адреса %w", err)
	}

	return nil
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered
	query := `DELETE FROM parcel WHERE number = ? AND status = 'registered'`
	_, err := s.db.Exec(query, number)
	if err != nil {
		return fmt.Errorf("ошибка удаления %w", err)
	}

	return nil
}
