package _vscode

import (
	"database/sql"
	"github.com/Yandex-Practicum/go-db-sql-final"
	"log"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p main.Parcel) (int, error)
	// реализуйте добавление строки в таблицу parcel, используйте данные из переменной p
	res, err := s.db.Exec("insert into parcel (client, status, address, created_at)"+
		"values (:client, :status, :address, :created_at)",
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt))
	if err != nil {
		log.Println(err)
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Println(err)
		return 0, err
	}

	// верните идентификатор последней добавленной записи
	return int(id), nil
}
//fufkndlnvjd
func (s ParcelStore) Get(number int) (main.Parcel, error) {
	// реализуйте чтение строки по заданному number
	// здесь из таблицы должна вернуться только одна строка

	// заполните объект Parcel данными из таблицы
	p := main.Parcel{}
	row := s.db.QueryRow("select * from parcel where number = $1", number)
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err == sql.ErrNoRows {
		return p, err
	}

	return p, err
}

func (s ParcelStore) GetByClient(client int) ([]main.Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк

	// заполните срез Parcel данными из таблицы
	var res []main.Parcel

	rows, err := s.db.Query("select * from parcel where client = $1", client)
	if err != nil {
		log.Println(err)
	}

	defer rows.Close()
	for rows.Next() {
		p := main.Parcel{}
		err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			log.Println(err)
			return res, err
		}
		res = append(res, p)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		return res, err
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel
	_, err := s.db.Exec("update parcel set status = $1 where number = $2", status, number)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
	//p, err := s.Get(number)
	//
	//if p.Status == ParcelStatusRegistered {
	//	_, err := s.db.Exec("update parcel set address = $1 where number = $2", address, number)
	//	if err != nil {
	//		log.Println(err)
	//		return err
	//	}
	//
	//	return nil
	//} else {
	//	log.Printf("При статусе  %s, обновление запрещено\n", p.Status)
	//}
	//return err
	_, err := s.db.Exec("update parcel set address = $1 where number = $2 and status = $3", address, number, main.ParcelStatusRegistered)
	if err != nil {
		log.Println(err)
		return err

	}

	return nil
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered
	//p, err := s.Get(number)
	//if p.Status == ParcelStatusRegistered {
	//	_, err := s.db.Exec("delete from parcel where number = $1", number)
	//	if err != nil {
	//		log.Println(err)
	//		return err
	//	}

	if _, err := s.db.Exec("delete from parcel where number = $1 and status = $2", number, main.ParcelStatusRegistered); err != nil {
		log.Println(err)
		return err
	}

	return nil

}
