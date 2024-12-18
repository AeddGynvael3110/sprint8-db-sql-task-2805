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

	res, err := s.db.Exec("INSERT INTO parcel (client, status, address, created_at) VALUES (:client, :status, :address, :created_at)",
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt))
	if err != nil {
		fmt.Errorf("query row failed %w", err)
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		fmt.Errorf("get last index failed %w", err)
		return 0, err
	}
	//fmt.Println(res.LastInsertId())
	// верните идентификатор последней добавленной записи
	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// реализуйте чтение строки по заданному number
	// здесь из таблицы должна вернуться только одна строка

	// заполните объект Parcel данными из таблицы
	p := Parcel{}

	row := s.db.QueryRow("SELECT number, client, status, address, created_at FROM parcel WHERE number = :number", sql.Named("number", number))
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		fmt.Println(err)
		return p, err
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк

	// заполните срез Parcel данными из таблицы
	var res []Parcel
	rows, err := s.db.Query("SELECT number, client, status, address, created_at FROM parcel WHERE client = :client", sql.Named("client", client))
	if err != nil {
		fmt.Errorf("query row failed %w", err)
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		p := Parcel{}
		err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			fmt.Errorf("scan rows failed %w", err)
			return res, err
		}
		res = append(res, p)
	}
	err = rows.Err()
	if err == nil {
		fmt.Errorf("scan next rows failed %w", err)
		return res, err
	}
	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel
	_, err := s.db.Exec("UPDATE parcel SET status = :status WHERE number = :id",
		sql.Named("status", status),
		sql.Named("id", number))
	if err != nil {
		fmt.Errorf("update row failed %w", err)
		return err
	}
	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM parcel WHERE number = :id and status = :status",
		sql.Named("id", number),
		sql.Named("status", ParcelStatusRegistered)).Scan(&count)
	if err != nil || count == 0 {
		fmt.Errorf("query row failed %w", err)
		return err
	}
	_, err = s.db.Exec("UPDATE parcel SET address = :address WHERE number = :id",
		sql.Named("address", address),
		sql.Named("id", number))
	if err != nil {
		fmt.Errorf("update row failed %w", err)
		return err
	}
	return nil
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM parcel WHERE number = :id and status = :status",
		sql.Named("id", number),
		sql.Named("status", ParcelStatusRegistered)).Scan(&count)
	if err != nil || count == 0 {
		fmt.Errorf("query row failed %w", err)
		return err
	}
	_, err = s.db.Exec("DELETE FROM parcel WHERE number = :id", sql.Named("id", number))
	if err != nil {
		fmt.Errorf("delete row failed %w", err)
		return err
	}
	return nil
}
