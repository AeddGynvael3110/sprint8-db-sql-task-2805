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

	result, err := s.db.Exec("INSERT INTO parcel (client, status, address, created_at) VALUES (?, ?, ?, ?)",
		p.Client, p.Status, p.Address, p.CreatedAt)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	// верните идентификатор последней добавленной записи
	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// реализуйте чтение строки по заданному number
	// здесь из таблицы должна вернуться только одна строка

	var p Parcel
	err := s.db.QueryRow("SELECT number, client, status, address, created_at FROM parcel WHERE number = ?", number).
		Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return Parcel{}, fmt.Errorf("no parcel found with number %d", number)
		}
		return Parcel{}, err
	}
	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк

	rows, err := s.db.Query("SELECT number, client, status, address, created_at FROM parcel WHERE client = ?", client)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var parcels []Parcel
	for rows.Next() {
		p := Parcel{}
		if err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt); err != nil {
			return nil, err
		}
		parcels = append(parcels, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return parcels, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel

	_, err := s.db.Exec("UPDATE parcel SET status = ? WHERE number = ?", status, number)
	return err
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
	// Обновляем адрес, если статус registered
	result, err := s.db.Exec("UPDATE parcel SET address = ? WHERE number = ? AND status = ?", address, number, ParcelStatusRegistered)
	if err != nil {
		return err
	}

	// Проверяем, было ли изменено хотя бы одно значение
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("address can only be changed for registered parcels or parcel not found")
	}
	return nil
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered
	result, err := s.db.Exec("DELETE FROM parcel WHERE number = ? AND status = ?", number, ParcelStatusRegistered)
	if err != nil {
		return err
	}
	// Проверка, сколько строк было затронуто
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// Если строка не была удалена, значит, статус не 'registered'
	if rowsAffected == 0 {
		return fmt.Errorf("parcel with number %d not found or status is not registered", number)
	}

	return nil
}
