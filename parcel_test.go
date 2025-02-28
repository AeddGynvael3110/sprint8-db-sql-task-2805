package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

var (
	randSource = rand.NewSource(time.Now().UnixNano())
	randRange  = rand.New(randSource)
)

// getTestParcel возвращает тестовую посылку
func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    "registered",
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// setupDB настраивает базу данных для тестов
func setupDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)

	_, err = db.Exec(`
        CREATE TABLE parsel (
            number INTEGER PRIMARY KEY AUTOINCREMENT,
            client INTEGER NOT NULL,
            status TEXT NOT NULL,
            address TEXT NOT NULL,
            created_at TEXT NOT NULL
        );
    `)
	require.NoError(t, err)

	return db
}

// TestAddGetDelete проверяет добавление, получение и удаление посылки
func TestAddGetDelete(t *testing.T) {
	// Подготовка
	db := setupDB(t)
	store := NewParcelStore(db)
	parcel := getTestParcel()

	// Добавление
	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotZero(t, id)

	// Получение
	retrievedParcel, err := store.Get(id)
	require.NoError(t, err)
	require.Equal(t, parcel.Client, retrievedParcel.Client)
	require.Equal(t, parcel.Status, retrievedParcel.Status)
	require.Equal(t, parcel.Address, retrievedParcel.Address)
	require.Equal(t, parcel.CreatedAt, retrievedParcel.CreatedAt)

	// Удаление
	err = store.Delete(id)
	require.NoError(t, err)

	// Проверка, что посылка удалена
	_, err = store.Get(id)
	require.Error(t, err)
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	// Подготовка
	db := setupDB(t)
	store := NewParcelStore(db)
	parcel := getTestParcel()

	// Добавление
	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotZero(t, id)

	// Обновление адреса
	newAddress := "new test address"
	err = store.SetAddress(id, newAddress)
	require.NoError(t, err)

	// Проверка
	retrievedParcel, err := store.Get(id)
	require.NoError(t, err)
	require.Equal(t, newAddress, retrievedParcel.Address)
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	// Подготовка
	db := setupDB(t)
	store := NewParcelStore(db)
	parcel := getTestParcel()

	// Добавление
	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotZero(t, id)

	// Обновление статуса
	newStatus := "delivered"
	err = store.SetStatus(id, newStatus)
	require.NoError(t, err)

	// Проверка
	retrievedParcel, err := store.Get(id)
	require.NoError(t, err)
	require.Equal(t, newStatus, retrievedParcel.Status)
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	// Подготовка
	db := setupDB(t)
	store := NewParcelStore(db)

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}

	// Задаём всем посылкам один и тот же идентификатор клиента
	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	// Добавление
	for i := 0; i < len(parcels); i++ {
		id, err := store.Add(parcels[i])
		require.NoError(t, err)
		require.NotZero(t, id)

		parcels[i].Number = id
		parcelMap[id] = parcels[i]
	}

	// Получение по клиенту
	storedParcels, err := store.GetByClient(fmt.Sprintf("%d", client))
	require.NoError(t, err)
	require.Equal(t, len(parcels), len(storedParcels))

	// Проверка
	for _, parcel := range storedParcels {
		expectedParcel, exists := parcelMap[parcel.Number]
		require.True(t, exists)
		require.Equal(t, expectedParcel.Client, parcel.Client)
		require.Equal(t, expectedParcel.Status, parcel.Status)
		require.Equal(t, expectedParcel.Address, parcel.Address)
		require.Equal(t, expectedParcel.CreatedAt, parcel.CreatedAt)
	}
}
