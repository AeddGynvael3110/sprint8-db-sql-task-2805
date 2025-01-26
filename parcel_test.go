package main

import (
	"database/sql"
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
	if err != nil {
		require.NoError(t, err)
	}
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	number, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotEmpty(t, number)

	stored, err := store.Get(number)
	require.NoError(t, err)
	require.Equal(t, stored.Client, parcel.Client)
	require.Equal(t, stored.Status, parcel.Status)
	require.Equal(t, stored.Address, parcel.Address)
	require.Equal(t, stored.CreatedAt, parcel.CreatedAt)

	err = store.Delete(number)
	require.NoError(t, err)

	_, err = store.Get(number)
	require.Equal(t, sql.ErrNoRows, err)
}

func TestSetAddress(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		require.NoError(t, err)
	}
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	number, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotEmpty(t, number)

	newAddress := "new test address"
	err = store.SetAddress(number, newAddress)
	require.NoError(t, err)

	stored, err := store.Get(number)
	require.NoError(t, err)
	require.Equal(t, stored.Address, newAddress)
}

func TestSetStatus(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		require.NoError(t, err)
	}
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	number, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotEmpty(t, number)

	err = store.SetStatus(number, ParcelStatusSent)
	require.NoError(t, err)

	stored, err := store.Get(number)
	require.NoError(t, err)
	require.Equal(t, stored.Status, ParcelStatusSent)
}

func TestGetByClient(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		require.NoError(t, err)
	}
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
	require.Equal(t, len(storedParcels), len(parcels))

	for _, parcel := range storedParcels {
		require.Equal(t, parcelMap[parcel.Number].Address, parcel.Address)
		require.Equal(t, parcelMap[parcel.Number].CreatedAt, parcel.CreatedAt)
		require.Equal(t, parcelMap[parcel.Number].Status, parcel.Status)
		require.Equal(t, parcelMap[parcel.Number].Client, parcel.Client)
	}
}
