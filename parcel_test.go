package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite" // Для поддержки SQLite
)

func setupDB() *sql.DB {
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		panic(err)
	}
	return db
}

func getTestParcel() Parcel {
	return Parcel{
		Client:    rand.Intn(1000),
		Status:    ParcelStatusRegistered,
		Address:   "Test Address",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

func TestAddGetDelete(t *testing.T) {
	db := setupDB()
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// Add
	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotZero(t, id)

	// Get
	storedParcel, err := store.Get(id)
	require.NoError(t, err)
	require.Equal(t, parcel.Address, storedParcel.Address)

	// Delete
	err = store.Delete(id)
	require.NoError(t, err)

	// Check delete
	storedParcel, err = store.Get(id)
	require.NoError(t, err)
	require.Zero(t, storedParcel.Number)
}

func TestSetAddress(t *testing.T) {
	db := setupDB()
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// Add
	id, err := store.Add(parcel)
	require.NoError(t, err)

	// Set Address
	newAddress := "New Test Address"
	err = store.SetAddress(id, newAddress)
	require.NoError(t, err)

	// Check
	storedParcel, err := store.Get(id)
	require.NoError(t, err)
	require.Equal(t, newAddress, storedParcel.Address)
}

func TestSetStatus(t *testing.T) {
	db := setupDB()
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// Add
	id, err := store.Add(parcel)
	require.NoError(t, err)

	// Set Status
	err = store.SetStatus(id, ParcelStatusSent)
	require.NoError(t, err)

	// Check
	storedParcel, err := store.Get(id)
	require.NoError(t, err)
	require.Equal(t, ParcelStatusSent, storedParcel.Status)
}

func TestGetByClient(t *testing.T) {
	db := setupDB()
	defer db.Close()

	store := NewParcelStore(db)

	clientID := rand.Intn(1000)
	parcels := []Parcel{
		{Client: clientID, Status: ParcelStatusRegistered, Address: "Addr1", CreatedAt: time.Now().Format(time.RFC3339)},
		{Client: clientID, Status: ParcelStatusSent, Address: "Addr2", CreatedAt: time.Now().Format(time.RFC3339)},
	}

	for i := range parcels {
		id, err := store.Add(parcels[i])
		require.NoError(t, err)
		parcels[i].Number = id
	}

	// Get by Client
	storedParcels, err := store.GetByClient(clientID)
	require.NoError(t, err)
	require.Len(t, storedParcels, 2)
}
