package main

import (
	"database/sql"
	"fmt"
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
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
	num, err := store.Add(parcel)
	assert.NoError(t, err)
	require.Greater(t, num, 0)
	// get
	// получите только что добавленную посылку, убедитесь в отсутствии ошибки
	// проверьте, что значения всех полей в полученном объекте совпадают со значениями полей в переменной parcel
	secondParcel, err := store.Get(num)
	assert.NoError(t, err)
	assert.Equal(t, num, secondParcel.Number)
	assert.Equal(t, parcel.Client, secondParcel.Client)
	assert.Equal(t, parcel.Status, secondParcel.Status)
	assert.Equal(t, parcel.Address, secondParcel.Address)
	assert.Equal(t, parcel.CreatedAt, secondParcel.CreatedAt)
	// delete
	// удалите добавленную посылку, убедитесь в отсутствии ошибки
	// проверьте, что посылку больше нельзя получить из БД
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
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
	store := NewParcelStore(db)
	parcel := getTestParcel()

	num, err := store.Add(parcel)
	require.NoError(t, err)
	// set address
	// обновите адрес, убедитесь в отсутствии ошибки
	newAddress := "new test address"

	err = store.SetAddress(num, newAddress)
	assert.NoError(t, err)

	// check
	// получите добавленную посылку и убедитесь, что адрес обновился
	secondParcel, err := store.Get(num)
	assert.NoError(t, err)
	assert.Equal(t, newAddress, secondParcel.Address)

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
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
	store := NewParcelStore(db)
	parcel := getTestParcel()
	num, err := store.Add(parcel)
	require.NoError(t, err)

	// set status
	// обновите статус, убедитесь в отсутствии ошибки
	err = store.SetStatus(num, ParcelStatusSent)
	assert.NoError(t, err)

	// check
	// получите добавленную посылку и убедитесь, что статус обновился
	secondParcel, err := store.Get(num)
	assert.NoError(t, err)
	assert.Equal(t, ParcelStatusSent, secondParcel.Status)
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
		// обновляем идентификатор добавленной у посылки
		parcels[i].Number = id

		// сохраняем добавленную посылку в структуру map, чтобы её можно было легко достать по идентификатору посылки
		parcelMap[id] = parcels[i]
	}

	// get by client
	storedParcels, err := store.GetByClient(client) // получите список посылок по идентификатору клиента, сохранённого в переменной client
	// убедитесь в отсутствии ошибки
	require.NoError(t, err)
	assert.Equal(t, len(parcels), len(storedParcels))
	// убедитесь, что количество полученных посылок совпадает с количеством добавленных

	// check
	for _, storedParcel := range storedParcels {
		expectedParcel, exists := parcelMap[storedParcel.Number]
		require.True(t, exists)

		assert.Equal(t, expectedParcel.Client, storedParcel.Client)
		assert.Equal(t, expectedParcel.Status, storedParcel.Status)
		assert.Equal(t, expectedParcel.Address, storedParcel.Address)
		assert.Equal(t, expectedParcel.CreatedAt, storedParcel.CreatedAt)
	}
}
