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

func getTestParcel() Parcel {

	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

func SetupDb() (*sql.DB, error) {

	db, err := sql.Open("sqlite", "tracker.db")

	if err != nil {
		return nil, err
	}

	return db, nil
}

func TestAddGetDelete(t *testing.T) {

	db, err := SetupDb()
	require.NoError(t, err)

	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	parcel.Number, err = store.Add(parcel)
	require.NoError(t, err)
	require.Greater(t, parcel.Number, 0)

	storedParcel, err := store.Get(parcel.Number)
	require.NoError(t, err)

	assert.Equal(t, parcel, storedParcel)
	err = store.Delete(parcel.Number)
	require.NoError(t, err)

	parcel, err = store.Get(parcel.Number)
	require.Error(t, err)
}

func TestSetAddress(t *testing.T) {

	db, err := SetupDb()
	require.NoError(t, err)

	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	parcelID, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotNil(t, parcelID)

	newAddress := "new test address"
	err = store.SetAddress(parcelID, newAddress)
	require.NoError(t, err)

	updatedParcel, err := store.Get(parcelID)
	require.NoError(t, err)
	require.Equal(t, newAddress, updatedParcel.Address)
}

func TestSetStatus(t *testing.T) {

	db, err := SetupDb()
	require.NoError(t, err)

	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	parcelID, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotNil(t, parcelID)

	err = store.SetStatus(parcelID, ParcelStatusSent)
	require.NoError(t, err)

	updatedParcel, err := store.Get(parcelID)
	require.NoError(t, err)
	require.Equal(t, ParcelStatusSent, updatedParcel.Status)
}

func TestGetByClient(t *testing.T) {

	db, err := SetupDb()
	require.NoError(t, err)

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
		require.NotEmpty(t, id)

		parcels[i].Number = id
		parcelMap[id] = parcels[i]
	}

	storedParcels, err := store.GetByClient(client)
	require.NoError(t, err)
	assert.Len(t, storedParcels, len(parcelMap))

	for _, parcel := range storedParcels {
		_, ok := parcelMap[parcel.Number]
		assert.True(t, ok)
		assert.Equal(t, parcel, parcelMap[parcel.Number])
	}
}
