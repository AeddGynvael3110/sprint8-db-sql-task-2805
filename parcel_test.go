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
	randSource = rand.NewSource(time.Now().UnixNano())
	randRange  = rand.New(randSource)
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
	db, err := sql.Open("sqlite", "tracker.db")
	require.Empty(t, err)
	defer db.Close()
	store := NewParcelStore(db)
	parcel := getTestParcel()

	res, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotEqual(t, res, 0)

	p, err := store.Get(res)
	require.NoError(t, err)
	assert.Equal(t, p.Address, parcel.Address)
	assert.Equal(t, p.Client, parcel.Client)
	assert.Equal(t, p.CreatedAt, parcel.CreatedAt)
	assert.Equal(t, p.Status, parcel.Status)
	assert.Equal(t, p.Number, res)

	err1 := store.Delete(res)
	require.NoError(t, err1)
	p, err2 := store.Get(res)
	require.NotEmpty(t, err2)
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	require.Empty(t, err)
	defer db.Close()
	store := NewParcelStore(db)
	parcel := getTestParcel()

	res, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotEmpty(t, res)

	newAddress := "new test address"
	err1 := store.SetAddress(res, newAddress)
	require.NoError(t, err1)

	res1, err := store.Get(res)
	require.NoError(t, err)
	require.Equal(t, res1.Address, newAddress)
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	require.Empty(t, err)
	defer db.Close()
	store := NewParcelStore(db)
	parcel := getTestParcel()

	res, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotEmpty(t, res)

	err1 := store.SetStatus(res, ParcelStatusRegistered)
	require.NoError(t, err1)

	res1, err := store.Get(res)
	require.NoError(t, err)
	require.Equal(t, res1.Status, ParcelStatusRegistered)
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	require.Empty(t, err)
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
		require.NoError(t, err)

		parcels[i].Number = id
		parcelMap[id] = parcels[i]
	}

	storedParcels, err := store.GetByClient(client)
	require.NoError(t, err)
	require.Equal(t, len(storedParcels), len(parcels))

	for _, storedParcel := range storedParcels {
		expectedParcel, exists := parcelMap[storedParcel.Number]
		require.True(t, exists, "Parcel with number %d not found in parcelMap", storedParcel.Number)
		assert.Equal(t, expectedParcel, storedParcel)
	}
}
