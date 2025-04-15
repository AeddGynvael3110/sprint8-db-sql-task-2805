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
	res, err := s.db.Exec("INSERT INTO parcel (client, status, address, created_at) VALUES (:client, :status, :address, :created_at)",
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt))

	if err != nil {
		return 0, errors.New("can not add values")
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, errors.New("there is no last insert id")
	}
	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	row := s.db.QueryRow("SELECT number, client, status, address, created_at FROM parcel WHERE number = ?", number)

	p := Parcel{}
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return p, errors.New("can not scan values")
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	row, err := s.db.Query("SELECT number, client, status, address, created_at FROM parcel WHERE client = ?", client)
	if err != nil {
		return nil, errors.New("can not get data")
	}

	defer row.Close()

	var res []Parcel
	for row.Next() {
		p := Parcel{}
		err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			return res, errors.New("can not scan values")
		}

		res = append(res, p)
	}

	if err := row.Err(); err != nil {
		return res, errors.New("can not read data")
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	_, err := s.Get(number)
	if err != nil {
		return errors.New("can not get data")
	}

	_, err = s.db.Exec("UPDATE parcel SET status = ? WHERE number = ?", status, number)
	if err != nil {
		return errors.New("can not update data")
	}
	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	p, err := s.Get(number)
	if err != nil {
		return errors.New("can not get data")
	}
	if p.Status != ParcelStatusRegistered {
		return errors.New("cannot update address: status is not 'registered'")
	}

	_, err = s.db.Exec("UPDATE parcel SET address = ? WHERE number = ?", address, number)
	if err != nil {
		return errors.New("can not update data")
	}
	return nil
}

func (s ParcelStore) Delete(number int) error {
	p, err := s.Get(number)
	if err != nil {
		return errors.New("can not get data")
	}

	if p.Status != ParcelStatusRegistered {
		return errors.New("cannot delete parcel: status is not 'registered'")
	}

	_, err = s.db.Exec("DELETE FROM parcel WHERE number = ?", number)
	if err != nil {
		return errors.New("can not update data")
	}

	return nil
}
