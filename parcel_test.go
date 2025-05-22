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

func TestAddGetDelete(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	store := NewParcelStore(db)
	parcel := getTestParcel()
	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotZero(t, id)
	parcel.Number = id
	stored, err := store.Get(id)
	require.NoError(t, err)
	require.Equal(t, parcel, stored)
	require.NoError(t, store.Delete(id))
	_, err = store.Get(id)
	require.ErrorIs(t, err, sql.ErrNoRows)
}

func TestSetAddress(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	store := NewParcelStore(db)
	p := getTestParcel()
	id, err := store.Add(p)
	require.NoError(t, err)
	newAddress := "new test address"
	require.NoError(t, store.SetAddress(id, newAddress))
	stored, err := store.Get(id)
	require.NoError(t, err)
	require.Equal(t, newAddress, stored.Address)
}

func TestSetStatus(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	store := NewParcelStore(db)
	p := getTestParcel()
	id, err := store.Add(p)
	require.NoError(t, err)
	require.NoError(t, store.SetStatus(id, ParcelStatusSent))
	stored, err := store.Get(id)
	require.NoError(t, err)
	require.Equal(t, ParcelStatusSent, stored.Status)
}

func TestGetByClient(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	store := NewParcelStore(db)
	parcels := []Parcel{getTestParcel(), getTestParcel(), getTestParcel()}
	client := randRange.Intn(10_000_000)
	for i := range parcels {
		parcels[i].Client = client
		id, err := store.Add(parcels[i])
		require.NoError(t, err)
		parcels[i].Number = id
	}
	stored, err := store.GetByClient(client)
	require.NoError(t, err)
	require.Len(t, stored, len(parcels))
	m := map[int]Parcel{}
	for _, p := range parcels {
		m[p.Number] = p
	}
	for _, p := range stored {
		exp, ok := m[p.Number]
		require.True(t, ok)
		require.Equal(t, exp, p)
	}
}
