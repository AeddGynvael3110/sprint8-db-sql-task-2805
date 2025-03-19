package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

var (
	randSource = rand.NewSource(time.Now().UnixNano())
	randRange  = rand.New(randSource)
)

func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

func setupDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", "file:test.db?mode=memory&cache=shared")
	require.NoError(t, err, "Ошибка при подключении к БД")

	_, err = db.Exec(`CREATE TABLE parcel (
		number INTEGER PRIMARY KEY AUTOINCREMENT,
		client INTEGER,
		status TEXT,
		address TEXT,
		created_at TEXT
	)`)
	require.NoError(t, err, "Ошибка при создании таблицы")

	return db
}

func TestAddGetDelete(t *testing.T) {
	db := setupDB(t)
	defer db.Close()
	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel)
	require.NoError(t, err, "Ошибка при добавлении посылки")
	require.NotZero(t, id, "Идентификатор посылки не должен быть нулевым")

	retrievedParcel, err := store.Get(id)
	require.NoError(t, err, "Ошибка при получении посылки")
	require.Equal(t, parcel.Client, retrievedParcel.Client, "Клиент не совпадает")
	require.Equal(t, parcel.Status, retrievedParcel.Status, "Статус не совпадает")
	require.Equal(t, parcel.Address, retrievedParcel.Address, "Адрес не совпадает")
	require.Equal(t, parcel.CreatedAt, retrievedParcel.CreatedAt, "Дата создания не совпадает")

	err = store.Delete(id)
	require.NoError(t, err, "Ошибка при удалении посылки")

	_, err = store.Get(id)
	require.Error(t, err, "Ожидалась ошибка при получении удаленной посылки")
}

func TestSetAddress(t *testing.T) {
	db := setupDB(t)
	defer db.Close()
	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel)
	require.NoError(t, err, "Ошибка при добавлении посылки")

	newAddress := "new test address"
	err = store.SetAddress(id, newAddress)
	require.NoError(t, err, "Ошибка при обновлении адреса")

	retrievedParcel, err := store.Get(id)
	require.NoError(t, err, "Ошибка при получении посылки")
	require.Equal(t, newAddress, retrievedParcel.Address, "Адрес не обновился")
}

func TestSetStatus(t *testing.T) {
	db := setupDB(t)
	defer db.Close()
	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel)
	require.NoError(t, err, "Ошибка при добавлении посылки")

	err = store.SetStatus(id, ParcelStatusSent)
	require.NoError(t, err, "Ошибка при обновлении статуса")

	retrievedParcel, err := store.Get(id)
	require.NoError(t, err, "Ошибка при получении посылки")
	require.Equal(t, ParcelStatusSent, retrievedParcel.Status, "Статус не обновился")
}

func TestGetByClient(t *testing.T) {
	db := setupDB(t)
	defer db.Close()
	store := NewParcelStore(db)

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}

	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	for i := 0; i < len(parcels); i++ {
		id, err := store.Add(parcels[i])
		require.NoError(t, err, "Ошибка при добавлении посылки")
		require.NotZero(t, id, "Идентификатор посылки не должен быть нулевым")

		parcels[i].Number = id
		parcelMap[id] = parcels[i]
	}

	storedParcels, err := store.GetByClient(client)
	require.NoError(t, err, "Ошибка при получении посылок по клиенту")
	require.Equal(t, len(parcels), len(storedParcels), "Количество посылок не совпадает")

	for _, parcel := range storedParcels {
		expectedParcel, exists := parcelMap[parcel.Number]
		require.True(t, exists, "Посылка не найдена в map")
		require.Equal(t, expectedParcel.Client, parcel.Client, "Клиент не совпадает")
		require.Equal(t, expectedParcel.Status, parcel.Status, "Статус не совпадает")
		require.Equal(t, expectedParcel.Address, parcel.Address, "Адрес не совпадает")
		require.Equal(t, expectedParcel.CreatedAt, parcel.CreatedAt, "Дата создания не совпадает")
	}
}
