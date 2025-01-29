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
	res, err := s.db.Exec("INSERT INTO parcel (number, client, status, address, created_at) VALUES (:number, :client, :status, :address, :created_at)",
        sql.Named("number", p.Number),
        sql.Named("client", p.Client),
        sql.Named("status", p.Status),
        sql.Named("address", p.Address),
        sql.Named("created_at", p.CreatedAt))
    if err != nil {
        return 0, err
    }
    last, _ := res.LastInsertId()
    return int(last), err
}

func (s ParcelStore) Get(number int) (Parcel, error) {
    rows, err := s.db.Query(fmt.Sprintf("SELECT * FROM parcel WHERE number = %s", number))
	p := Parcel{}
	if err != nil {
        return p, err
    }
    for rows.Next() {
        err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
        if err != nil {
            return p, err
        }
    }
	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк

	// заполните срез Parcel данными из таблицы
	var res []Parcel
    rows, err := s.db.Query(fmt.Sprintf("SELECT * FROM parcel WHERE client = %s", client))
    if err != nil {
        return res, err
    }
    for rows.Next() {
        p := Parcel{}
        err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
        if err != nil {
            return res, err
        }
        res = append(res, p)
    }
	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel

	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered

	return nil
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered

	return nil
}
