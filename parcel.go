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
	query := fmt.Sprintf("INSERT INTO parcel (client, status, address, created_at) VALUES ('%d', '%s', '%s', '%s')",
		p.Client, p.Status, p.Address, p.CreatedAt)
	result, err := s.db.Exec(query)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	query := fmt.Sprintf("SELECT number, client, status, address, created_at FROM parcel WHERE number = '%d'", number)
	row := s.db.QueryRow(query)
	p := Parcel{}
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return p, err
	}
	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	query := fmt.Sprintf("SELECT number, client, status, address, created_at FROM parcel WHERE client = '%d'", client)
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []Parcel
	for rows.Next() {
		p := Parcel{}
		err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	query := fmt.Sprintf("UPDATE parcel SET status = '%s' WHERE number = '%d'", status, number)
	_, err := s.db.Exec(query)
	return err
}

func (s ParcelStore) SetAddress(number int, address string) error {
	query := fmt.Sprintf("UPDATE parcel SET address = '%s' WHERE number = '%d' AND status = '%s'",
		address, number, ParcelStatusRegistered)
	_, err := s.db.Exec(query)
	return err
}

func (s ParcelStore) Delete(number int) error {
	query := fmt.Sprintf("DELETE FROM parcel WHERE number = '%d' AND status = '%s'", number, ParcelStatusRegistered)
	_, err := s.db.Exec(query)
	return err
}
