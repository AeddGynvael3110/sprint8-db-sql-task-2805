package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

var (
	// randSource источник псевдо случайных чисел.
	// Для повышения уникальности в качестве seed
	// используется текущее время в unix формате (в виде числа)
	randSource = rand.NewSource(time.Now().UnixNano())
	// randRange использует randSource для генерации случайных чисел
	randRange = rand.New(randSource)
)

// getTestParcel возвращает тестовую посылку
func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// TestAddGetDelete проверяет добавление, получение и удаление посылки
func TestAddGetDelete(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", ":memory:") // используем временную базу данных для тестирования
	require.NoError(t, err)
	defer db.Close()

	// создаем таблицу для тестов
	_, err = db.Exec("CREATE TABLE parcel (number INTEGER PRIMARY KEY AUTOINCREMENT, client INTEGER, status TEXT, address TEXT, created_at TEXT)")
	require.NoError(t, err)

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.Greater(t, id, 0)

	// get
	retrievedParcel, err := store.Get(id)
	assert.NoError(t, err)
	assert.Equal(t, parcel.Client, retrievedParcel.Client)
	assert.Equal(t, parcel.Status, retrievedParcel.Status)
	assert.Equal(t, parcel.Address, retrievedParcel.Address)
	assert.Equal(t, parcel.CreatedAt, retrievedParcel.CreatedAt)
	assert.Equal(t, id, retrievedParcel.Number)

	// delete
	err = store.Delete(id)
	require.NoError(t, err)

	// check
	_, err = store.Get(id)
	require.Error(t, err)
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer db.Close()

	_, err = db.Exec("CREATE TABLE parcel (number INTEGER PRIMARY KEY AUTOINCREMENT, client INTEGER, status TEXT, address TEXT, created_at TEXT)")
	require.NoError(t, err)

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.Greater(t, id, 0)

	// set address
	newAddress := "new test address"
	err = store.SetAddress(id, newAddress)
	require.NoError(t, err)

	// check
	retrievedParcel, err := store.Get(id)
	assert.NoError(t, err)
	assert.Equal(t, newAddress, retrievedParcel.Address)
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer db.Close()

	_, err = db.Exec("CREATE TABLE parcel (number INTEGER PRIMARY KEY AUTOINCREMENT, client INTEGER, status TEXT, address TEXT, created_at TEXT)")
	require.NoError(t, err)

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.Greater(t, id, 0)

	// set status
	newStatus := "shipped" // любой статус, отличный от registered
	err = store.SetStatus(id, newStatus)
	require.NoError(t, err)

	// check
	retrievedParcel, err := store.Get(id)
	assert.NoError(t, err)
	assert.Equal(t, newStatus, retrievedParcel.Status)
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer db.Close()

	_, err = db.Exec("CREATE TABLE parcel (number INTEGER PRIMARY KEY AUTOINCREMENT, client INTEGER, status TEXT, address TEXT, created_at TEXT)")
	require.NoError(t, err)

	store := NewParcelStore(db)
	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}

	// задаём всем посылкам один и тот же идентификатор клиента
	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	// add
	for i := range parcels {
		id, err := store.Add(parcels[i])
		require.NoError(t, err)
		parcels[i].Number = id
		parcelMap[id] = parcels[i]
	}

	// get by client
	storedParcels, err := store.GetByClient(client)
	assert.NoError(t, err)
	assert.Len(t, storedParcels, len(parcels))

	// check
	for _, parcel := range storedParcels {
		assert.Contains(t, parcelMap, parcel.Number)
		assert.Equal(t, parcelMap[parcel.Number], parcel)
		assert.Equal(t, parcelMap[parcel.Number].Status, parcel.Status)
		assert.Equal(t, parcelMap[parcel.Number].Address, parcel.Address)
		assert.Equal(t, parcelMap[parcel.Number].CreatedAt, parcel.CreatedAt)
	}
}
