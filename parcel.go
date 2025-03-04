package main

import (
	"database/sql"
	"errors"
	"fmt"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	res, err := s.db.Exec(
		"INSERT INTO parsel (client, status, address, created_at) VALUES (?, ?, ?, ?)",
		p.Client, p.Status, p.Address, p.CreatedAt,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to add parcel: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID: %w", err)
	}

	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	p := Parcel{}

	row := s.db.QueryRow("SELECT number, client, status, address, created_at FROM parsel WHERE number = ?", number)
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return Parcel{}, err
	}
	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	rows, err := s.db.Query("SELECT number, client, status, address, created_at FROM parsel WHERE client = ?", client)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var parcels []Parcel
	for rows.Next() {
		var p Parcel
		err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
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
	query := "UPDATE parsel SET status = ? WHERE number = ?"
	_, err := s.db.Exec(query, status, number)
	if err != nil {
		return err
	}
	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// Проверяем статус посылки
	var status string
	err := s.db.QueryRow("SELECT status FROM parsel WHERE number = ?", number).Scan(&status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("посылка с номером %d не найдена", number)
		}
		return fmt.Errorf("ошибка при получении статуса посылки: %w", err)
	}

	// Проверяем, что статус равен ParcelStatusRegistered
	if status != ParcelStatusRegistered {
		return fmt.Errorf("нельзя изменить адрес: статус посылки должен быть %s, текущий статус: %s", ParcelStatusRegistered, status)
	}

	// Обновляем адрес
	query := "UPDATE parsel SET address = ? WHERE number = ?"
	_, err = s.db.Exec(query, address, number)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении адреса: %w", err)
	}

	return nil
}

func (s ParcelStore) Delete(number int) error {
	var status string
	err := s.db.QueryRow(`
		DELETE FROM parsel 
		WHERE number = ? AND status = ? 
		RETURNING status
	`, number, ParcelStatusRegistered).Scan(&status)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("посылка с номером %d не найдена или её нельзя удалить", number)
		}
		return err
	}

	return nil
}
