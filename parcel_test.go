package main

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
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

func TestAddGetDelete(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer func() {
		_ = db.Close()
	}()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	number, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotNil(t, number)

	getParcel, err := store.Get(number)
	require.NoError(t, err)
	parcel.Number = number
	assert.Equal(t, parcel, getParcel)

	err = store.Delete(number)
	require.NoError(t, err)
	_, err = store.Get(number)
	require.Error(t, err)
}

func TestSetAddress(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer func() {
		_ = db.Close()
	}()

	store := NewParcelStore(db)
	parcel := getTestParcel()
	number, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotNil(t, number)

	newAddress := "new test address"
	err = store.SetAddress(number, newAddress)
	require.NoError(t, err)

	getParcel, err := store.Get(number)
	require.NoError(t, err)
	assert.Equal(t, "new test address", getParcel.Address)
}

func TestSetStatus(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer func() {
		_ = db.Close()
	}()

	store := NewParcelStore(db)
	parcel := getTestParcel()
	number, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotNil(t, number)

	err = store.SetStatus(number, ParcelStatusSent)
	require.NoError(t, err)

	getParcel, err := store.Get(number)
	require.NoError(t, err)
	assert.Equal(t, ParcelStatusSent, getParcel.Status)
}

func TestGetByClient(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db") // настройте подключение к БД
	require.NoError(t, err)
	defer func() {
		_ = db.Close()
	}()
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
	assert.Len(t, storedParcels, 3)

	for _, parcel := range storedParcels {
		assert.Equal(t, parcelMap[parcel.Number], parcel)
	}
}
