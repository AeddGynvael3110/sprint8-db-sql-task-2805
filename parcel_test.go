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
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer db.Close()

	_, err = db.Exec("CREATE TABLE parcel (number INTEGER PRIMARY KEY AUTOINCREMENT, client INTEGER, status TEXT, address TEXT, created_at TEXT)")
	require.NoError(t, err)

	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.Greater(t, id, 0)

	retrievedParcel, err := store.Get(id)
	assert.NoError(t, err)
	assert.Equal(t, parcel.Client, retrievedParcel.Client)
	assert.Equal(t, parcel.Status, retrievedParcel.Status)
	assert.Equal(t, parcel.Address, retrievedParcel.Address)
	assert.Equal(t, parcel.CreatedAt, retrievedParcel.CreatedAt)
	assert.Equal(t, id, retrievedParcel.Number)

	err = store.Delete(id)
	require.NoError(t, err)

	_, err = store.Get(id)
	require.Error(t, err)
}

func TestSetAddress(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer db.Close()

	_, err = db.Exec("CREATE TABLE parcel (number INTEGER PRIMARY KEY AUTOINCREMENT, client INTEGER, status TEXT, address TEXT, created_at TEXT)")
	require.NoError(t, err)

	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.Greater(t, id, 0)

	newAddress := "new test address"
	err = store.SetAddress(id, newAddress)
	require.NoError(t, err)

	retrievedParcel, err := store.Get(id)
	assert.NoError(t, err)
	assert.Equal(t, newAddress, retrievedParcel.Address)
}

func TestSetStatus(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer db.Close()

	_, err = db.Exec("CREATE TABLE parcel (number INTEGER PRIMARY KEY AUTOINCREMENT, client INTEGER, status TEXT, address TEXT, created_at TEXT)")
	require.NoError(t, err)

	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.Greater(t, id, 0)

	newStatus := ParcelStatusSent // используйте статус из константы
	err = store.SetStatus(id, newStatus)
	require.NoError(t, err)

	retrievedParcel, err := store.Get(id)
	assert.NoError(t, err)
	assert.Equal(t, newStatus, retrievedParcel.Status)
}

func TestGetByClient(t *testing.T) {
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

	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	for i := range parcels {
		id, err := store.Add(parcels[i])
		require.NoError(t, err)
		parcels[i].Number = id
		parcelMap[id] = parcels[i]
	}

	storedParcels, err := store.GetByClient(client)
	assert.NoError(t, err)
	assert.Len(t, storedParcels, len(parcels))

	for _, parcel := range storedParcels {
		assert.Contains(t, parcelMap, parcel.Number)
		assert.Equal(t, parcelMap[parcel.Number], parcel)
		assert.Equal(t, parcelMap[parcel.Number].Status, parcel.Status)
		assert.Equal(t, parcelMap[parcel.Number].Address, parcel.Address)
		assert.Equal(t, parcelMap[parcel.Number].CreatedAt, parcel.CreatedAt)
	}
}
