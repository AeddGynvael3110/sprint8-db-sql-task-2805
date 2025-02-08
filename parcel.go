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
	res, err := s.db.Exec(` 
		INSERT INTO parcel (client, status, address, created_at) VALUES (?, ?, ?, ?)`,
		p.Client, p.Status, p.Address, p.CreatedAt)

	if err != nil {
		return 0, fmt.Errorf("Ошибка Add(): Вставка данных в parcel: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("Ошибка Add(): Получение последнего ID: %w", err)
	}

	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// реализуйте чтение строки по заданному number
	// здесь из таблицы должна вернуться только одна строка
	res := s.db.QueryRow("SELECT number, client, status, address, created_at FROM parcel WHERE number = ?", number)

	// заполните объект Parcel данными из таблицы
	p := Parcel{}
	err := res.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return Parcel{}, fmt.Errorf("Ошибка Get(): Заполние объекта Parcel данными из таблицы: %w", err)
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк
	rows, err := s.db.Query("SELECT number, client, status, address, created_at FROM parcel WHERE client = ?", client)
	if err != nil {
		return nil, fmt.Errorf("Ошибка GetByClient(): Реализуйте чтение строк из таблицы parcel по заданному client: %w", err)
	}
	// заполните срез Parcel данными из таблицы
	var res []Parcel

	defer rows.Close()

	for rows.Next() {
		p := Parcel{}
		err = rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			return res, fmt.Errorf("Ошибка GetByClient(): Заполние среза Parcel данными из таблицы: %w", err)
		}
		res = append(res, p)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("Ошибка GetByClient(): ошибка при обработке строк: %w", err)
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel
	_, err := s.db.Exec("UPDATE parcel SET status = ? WHERE number = ?", status, number)
	if err != nil {
		return fmt.Errorf("Ошибка SetStatus(): реализуйте обновление статуса в таблице parcel: %w", err)
	}

	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
	_, err := s.db.Exec("UPDATE parcel SET address = ? WHERE number = ? AND status = ?",
		address, number, ParcelStatusRegistered)
	if err != nil {
		return fmt.Errorf("Ошибка SetAddress(): обновление адреса: %w", err)
	}

	return nil
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered

	_, err := s.db.Exec("DELETE FROM parcel WHERE number = ? AND status = ?", number, ParcelStatusRegistered)
	if err != nil {
		return fmt.Errorf("Ошибка Delete():Не удалось DELETE FROM parcel WHERE number = %d: %w", number, err)
	}

	return nil
}
