package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
    "github.com/stretchr/testify/assert"
)

var (
	randSource = rand.NewSource(time.Now().UnixNano())
	randRange = rand.New(randSource)
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
    defer db.Close()

    store := NewParcelStore(db)

    parcel := getTestParcel()

    number, err := store.Add(parcel)
    if err != nil {
        require.NoError(t, err)
    }
    assert.NotEmpty(t, number)

    addedParcel, err := store.Get(number)
    if err != nil {
        require.NoError(t, err)
    }

    assert.EqualValues(t, parcel.Client, addedParcel.Client)
    assert.EqualValues(t, parcel.Status, addedParcel.Status)
    assert.EqualValues(t, parcel.Address, addedParcel.Address)
    assert.EqualValues(t, parcel.CreatedAt, addedParcel.CreatedAt)

    err = store.Delete(number)
    require.NoError(t, err)

    _, err = store.Get(number)
    require.Equal(t, err, sql.ErrNoRows)
}

func TestSetAddress(t *testing.T) {
    db, err := sql.Open("sqlite", "tracker.db")
    require.NoError(t, err)
    defer db.Close()

    store := NewParcelStore(db)

    parcel := getTestParcel()

    number, err := store.Add(parcel)
    require.NoError(t, err)
    assert.NotEmpty(t, number)

    newAddress := "new test address"
    err = store.SetAddress(number, newAddress)
    require.NoError(t, err)

    updatedParcel, err := store.Get(number)
    require.NoError(t, err)
    assert.Equal(t, newAddress, updatedParcel.Address)
}

func TestSetStatus(t *testing.T) {
    db, err := sql.Open("sqlite", "tracker.db")
    require.NoError(t, err)
    defer db.Close()

    store := NewParcelStore(db)

    parcel := getTestParcel()

    number, err := store.Add(parcel)
    require.NoError(t, err)
    require.NotEmpty(t, number)

    store.SetStatus(number, ParcelStatusDelivered)
    require.NoError(t, err)

    updatedParcel, err := store.Get(number)
    require.NoError(t, err)
    assert.Equal(t, ParcelStatusDelivered, updatedParcel.Status)
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
	require.Len(t, storedParcels, len(storedParcels))

	for _, parcel := range storedParcels {
		assert.NotEmpty(t, parcelMap[parcel.Number])
		assert.Equal(t, parcel, parcelMap[parcel.Number])
	}
}
