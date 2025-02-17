package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)

	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	num, err := store.Add(parcel)
	assert.NoError(t, err)
	require.Greater(t, num, 0)
	// get
	secondParcel, err := store.Get(num)
	assert.NoError(t, err)
	assert.Equal(t, secondParcel.Number, num)
	assert.Equal(t, secondParcel.Client, parcel.Client)
	assert.Equal(t, secondParcel.Status, parcel.Status)
	assert.Equal(t, secondParcel.Address, parcel.Address)
	assert.Equal(t, secondParcel.CreatedAt, parcel.CreatedAt)
	// delete
	err = store.Delete(num)
	assert.NoError(t, err)
	_, err = store.Get(num)
	assert.Error(t, err)
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	// add
	store := NewParcelStore(db)
	parcel := getTestParcel()

	num, err := store.Add(parcel)
	require.NoError(t, err)
	// set address
	newAddress := "new test address"

	err = store.SetAddress(num, newAddress)
	assert.NoError(t, err)

	// check
	secondParcel, err := store.Get(num)
	assert.NoError(t, err)
	assert.Equal(t, secondParcel.Address, newAddress)

	// delete
	err = store.Delete(num)
	assert.NoError(t, err)
	_, err = store.Get(num)
	assert.Error(t, err)
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()
	// add
	store := NewParcelStore(db)
	parcel := getTestParcel()
	num, err := store.Add(parcel)
	require.NoError(t, err)

	// set status
	// обновите статус, убедитесь в отсутствии ошибки
	err = store.SetStatus(num, ParcelStatusSent)
	assert.NoError(t, err)

	// check
	secondParcel, err := store.Get(num)
	assert.NoError(t, err)
	assert.Equal(t, secondParcel.Status, ParcelStatusSent)
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

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
	for i := 0; i < len(parcels); i++ {
		id, err := store.Add(parcels[i])
		require.NoError(t, err)

		parcels[i].Number = id

		// сохраняем добавленную посылку в структуру map, чтобы её можно было легко достать по идентификатору посылки
		parcelMap[id] = parcels[i]
	}

	// get by client
	storedParcels, err := store.GetByClient(client)
	require.NoError(t, err)
	assert.Equal(t, len(storedParcels), len(parcels))

	// check
	for i := 0; i < len(storedParcels); i++ {
		assert.Equal(t, storedParcels[i], parcels[i])
	}
	for _, storedParcel := range storedParcels {
		expectedParcel, exists := parcelMap[storedParcel.Number]
		require.True(t, exists)
		// сравнтваем два parcel
		assert.Equal(t, expectedParcel, storedParcel)
	}
}
