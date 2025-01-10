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
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	// верните идентификатор последней добавленной записи
	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// реализуйте чтение строки по заданному number
	// здесь из таблицы должна вернуться только одна строка
	sqlRow := s.db.QueryRow("SELECT number, client, status, address, created_at FROM parcel WHERE number = :number",
		sql.Named("number", number))

	// заполните объект Parcel данными из таблицы
	p := Parcel{}
	err := sqlRow.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return Parcel{}, err
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк
	sqlRows, err := s.db.Query("SELECT number, client, status, address, created_at FROM parcel where client = :client",
		sql.Named("client", client))
	if err != nil {
		return nil, err
	}
	defer sqlRows.Close()

	// заполните срез Parcel данными из таблицы
	var res []Parcel

	for sqlRows.Next() {
		parcel := Parcel{}
		err := sqlRows.Scan(&parcel.Number, &parcel.Client, &parcel.Status, &parcel.Address, &parcel.CreatedAt)
		if err != nil {
			return nil, err
		}
		res = append(res, parcel)
	}

	if err := sqlRows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel
	_, err := s.db.Exec("UPDATE parcel SET status = :status WHERE number = :number",
		sql.Named("status", status),
		sql.Named("number", number))
	if err != nil {
		return err
	}

	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
	res, err := s.db.Exec("UPDATE parcel SET address = :address WHERE number AND status = 'registered'",
		sql.Named("addresss", address),
		sql.Named("number", number))
	if err != nil {
		return err
	}

	rowsAff, err := res.RowsAffected()

	if rowsAff == 0 {
		return fmt.Errorf("no such row with number = %d and status = %s for new address", number, "registered")
	}

	if err != nil {
		return err
	}

	return nil
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered
	res, err := s.db.Exec("DELETE FROM parcel, WHERE number = :number AND status = 'registered'",
		sql.Named("number", number))
	if err != nil {
		return err
	}

	rowsAff, err := res.RowsAffected()
	if rowsAff == 0 {
		return fmt.Errorf("no such row with number = %d and status = %s for delete", number, "registered")
	}
	if err != nil {
		return err
	}

	return nil
}
