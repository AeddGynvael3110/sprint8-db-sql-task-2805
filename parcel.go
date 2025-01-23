package main

import (
	"database/sql"
	"errors"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	query := `INSERT INTO parcel (client, status, address, created_at) VALUES (?, ?, ?, ?)`

	res, err := s.db.Exec(query, p.Client, p.Status, p.Address, p.CreatedAt)

	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()

	return int(id), err
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	query := `SELECT number, client, status, address, created_at FROM parcel WHERE number = ?`
	row := s.db.QueryRow(query, number)

	p := Parcel{}
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)

	if err == sql.ErrNoRows {
		return Parcel{}, errors.New("посылка не найдена")
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	query := `SELECT number, client, status, address, created_at FROM parcel WHERE client = ?`

	rows, err := s.db.Query(query, client)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var res []Parcel

	for rows.Next() {
		var p Parcel

		if err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt); err != nil {
			return nil, err
		}

		res = append(res, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	query := `UPDATE parcel SET status = ? WHERE number = ?`

	_, err := s.db.Exec(query, status, number)

	return err
}

func (s ParcelStore) SetAddress(number int, address string) error {
	query := `UPDATE parcel SET address = ? WHERE number = ? AND status = ?`

	res, err := s.db.Exec(query, address, number, ParcelStatusRegistered)

	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil || rowsAffected == 0 {
		return errors.New("нельзя изменить адрес, посылка не имеет статус 'зарегистрирована'")
	}

	return nil
}

func (s ParcelStore) Delete(number int) error {
	query := `DELETE FROM parcel WHERE number = ? AND status = ?`

	res, err := s.db.Exec(query, number, ParcelStatusRegistered)

	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil || rowsAffected == 0 {
		return errors.New("нельзя удалить посылку, она не имеет статус 'зарегистрирована'")
	}

	return nil
}
