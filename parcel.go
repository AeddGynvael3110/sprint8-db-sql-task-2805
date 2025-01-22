package main

import (
	"database/sql"
	"log"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	// Добавление объект Parcel в таблицу tracker.db.parcel
	res, errExec := s.db.Exec(`
INSERT INTO parcel(client,
                   status,
                   address,
                   created_at)
    VALUES(:client,
           :status,
           :addr,
           :date);`,
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("addr", p.Address),
		sql.Named("date", p.CreatedAt))

	if errExec != nil {
		return 0, errExec
	}

	id, errID := res.LastInsertId()
	if errID != nil {
		return 0, errID
	}

	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// Возвращаем объект из таблицы tracker.db.parcel
	// tracker.db.parcel.number равный number
	row := s.db.QueryRow(`
SELECT number,
       client,
       status,
       address,
       created_at
FROM parcel
WHERE number = :num;`,
		sql.Named("num", number))

	// заполните объект Parcel данными из таблицы
	p := Parcel{}
	err := row.Scan(
		&p.Number,
		&p.Client,
		&p.Status,
		&p.Address,
		&p.CreatedAt)

	return p, err
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	//  Чтение строк из таблицы parcel по заданному client
	rows, err := s.db.Query(`
SELECT number,
       client,
       status,
       address,
       created_at
FROM parcel
WHERE client=:client;`,
		sql.Named("client", client))

	if err != nil {
		return nil, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Printf("sql.Rows error - %v", err)
		}
	}()

	var res []Parcel

	for rows.Next() {
		p := Parcel{}

		errScan := rows.Scan(
			&p.Number,
			&p.Client,
			&p.Status,
			&p.Address,
			&p.CreatedAt,
		)

		if errScan != nil {
			return nil, errScan
		}

		res = append(res, p)
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// Oбновление статуса в таблице parcel
	_, err := s.db.Exec(`
UPDATE parcel
SET status = :status
WHERE number = :num;`,
		sql.Named("num", number),
		sql.Named("status", status))

	return err
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// Обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
	_, err := s.db.Exec(`
UPDATE parcel
SET address = :addr
WHERE number = :num AND 
      status = :reg;`,
		sql.Named("addr", address),
		sql.Named("num", number),
		sql.Named("reg", ParcelStatusRegistered))

	return err
}

func (s ParcelStore) Delete(number int) error {
	// Удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered
	_, err := s.db.Exec(`
DELETE FROM parcel
WHERE number = :num AND
      status = :reg;`,
		sql.Named("num", number),
		sql.Named("reg", ParcelStatusRegistered))

	return err
}
