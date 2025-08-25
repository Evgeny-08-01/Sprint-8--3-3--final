package main

import (
	"database/sql"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	// реализуйте добавление строки в таблицу parcel, используйте данные из переменной p
	res, err := s.db.Exec("INSERT INTO parcel (client,status,address,created_at) VALUES (:Client,:Status,:Address,:CreatedAt)",
		sql.Named("Client", p.Client),
		sql.Named("Status", p.Status),
		sql.Named("Address", p.Address),
		sql.Named("CreatedAt", p.CreatedAt))
	if err != nil {
		return 0, err
	}
	// верните идентификатор последней добавленной записи
	lastID, err := res.LastInsertId()
	return int(lastID), err
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// реализуйте чтение строки по заданному number
	// здесь из таблицы должна вернуться только одна строка
	row := s.db.QueryRow("SELECT*  FROM parcel WHERE number = :number", sql.Named("number", number))
	// заполните объект Parcel данными из таблицы
	p := Parcel{}
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	return p, err
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк
	rows, err := s.db.Query("SELECT* FROM parcel WHERE client=:client", sql.Named("client", client))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// заполните срез Parcel данными из таблицы
	var res []Parcel

	for rows.Next() {
		//r-промежуточная структура для заполнения среза
		r := Parcel{}
		err = rows.Scan(&r.Number, &r.Client, &r.Status, &r.Address, &r.CreatedAt)
		res = append(res, r)
	}
	return res, err
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel
	_, err := s.db.Exec("UPDATE parcel SET status = :Status WHERE number = :number", sql.Named("Status",
		status), sql.Named("number", number))
	return err
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
	stat, err := s.Get(number)
	if stat.Status == ParcelStatusRegistered {
		_, err := s.db.Exec("UPDATE parcel SET address = :address WHERE number = :number", sql.Named("address",
			address), sql.Named("number", number))
		return err
	}
	return err
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered
	stat, err := s.Get(number)
	if err != nil {
		return err
	}
	if stat.Status == ParcelStatusRegistered {
		_, err = s.db.Exec("DELETE FROM parcel WHERE number = :number", sql.Named("number", number))
	}
	return err
}
