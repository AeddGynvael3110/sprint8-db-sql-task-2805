package main

import (
	"database/sql"
	"fmt"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	_, _ = db.Exec(`CREATE TABLE IF NOT EXISTS parcel(
		number INTEGER PRIMARY KEY AUTOINCREMENT,
		client INTEGER NOT NULL,
		status TEXT NOT NULL,
		address TEXT NOT NULL,
		created_at TEXT NOT NULL
	)`)
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	res, err := s.db.Exec(`INSERT INTO parcel(client,status,address,created_at) VALUES(?,?,?,?)`,
		p.Client, p.Status, p.Address, p.CreatedAt)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	row := s.db.QueryRow(`SELECT number,client,status,address,created_at FROM parcel WHERE number=?`, number)
	var p Parcel
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	return p, err
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	rows, err := s.db.Query(`SELECT number,client,status,address,created_at FROM parcel WHERE client=?`, client)
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
	return res, rows.Err()
}

func (s ParcelStore) SetStatus(number int, status string) error {
	res, err := s.db.Exec(`UPDATE parcel SET status=? WHERE number=?`, status, number)
	if err != nil {
		return err
	}
	aff, _ := res.RowsAffected()
	if aff == 0 {
		return fmt.Errorf("parcel %d not found", number)
	}
	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	res, err := s.db.Exec(`UPDATE parcel SET address=? WHERE number=? AND status=?`,
		address, number, ParcelStatusRegistered)
	if err != nil {
		return err
	}
	aff, _ := res.RowsAffected()
	if aff == 0 {
		return fmt.Errorf("cannot change address")
	}
	return nil
}

func (s ParcelStore) Delete(number int) error {
	res, err := s.db.Exec(`DELETE FROM parcel WHERE number=? AND status=?`,
		number, ParcelStatusRegistered)
	if err != nil {
		return err
	}
	aff, _ := res.RowsAffected()
	if aff == 0 {
		return fmt.Errorf("cannot delete parcel")
	}
	return nil
}
