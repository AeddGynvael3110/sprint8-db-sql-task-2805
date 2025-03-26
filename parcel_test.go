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
	db, err := sql.open("sqlite", "tracker.db")// настройте подключение к БД
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
number, err := store.Add(parcel)
require.NoError(t, err)
require.NotNil(t, number)
	// get
	// получите только что добавленную посылку, убедитесь в отсутствии ошибки
	// проверьте, что значения всех полей в полученном объекте совпадают со значениями полей в переменной parcel
getParcel, err := store.Get(number)
require.NoError(t,err)
parcel.Number = number
assert.Equal(t, parcel, getParcel)

	// delete
	// удалите добавленную посылку, убедитесь в отсутствии ошибки
	// проверьте, что посылку больше нельзя получить из БД
	err = store.Delete(number)
	require.NoError(t, err)
	_, err = store.Get(number)
	require.Error(t, err)
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")// настройте подключение к БД
reqire.NoError(t, err)
defer db.Close()
	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
store := NewParcelStore(db)
parcel := getTestParcel()
number, err := store.Add(parcel)
reqire.NoError(t, err)
reqire.NotNill(t, number)

	// set address
	// обновите адрес, убедитесь в отсутствии ошибки
	newAddress := "new test address"
	err = store.SetAddress(number, newAddress)
	require.NoError(t, err)

	// check
	// получите добавленную посылку и убедитесь, что адрес обновился
	getParcel, err := store.Get(number)
	reqire.NoError(t, err)
	assert.Equal(t, "new test address", getParcel.Address)
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")// настройте подключение к БД
	reqire.NoError(t, err)
	defer db.Close
	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
store := NewParcelStore(db)
parcel := getTestParcel()
number, err := store.Add(parcel)
require.NoError(t, err)
require.NotNil(t, number)
	// set status
	// обновите статус, убедитесь в отсутствии ошибки
err = store.SetStatus(number, ParcelStatusSent)
reqire.NoError(t, err)
	// check
	// получите добавленную посылку и убедитесь, что статус обновился
	getParcel, err := store.Get(number)
	reqire.NoError(t, err)
	assert.Equal(t, ParcelStatusSent, getParcel.Status)
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	// prepare
	db, err := sql.open("sqlite", "tracker.db")// настройте подключение к БД
require.NoError(t, err)
defer db. Close ()
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
		id, err := store.Add(parcels[i]) // добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
requre.NoError(t, err)

		// обновляем идентификатор добавленной у посылки
		parcels[i].Number = id

		// сохраняем добавленную посылку в структуру map, чтобы её можно было легко достать по идентификатору посылки
		parcelMap[id] = parcels[i]
	}

	// get by client
	storedParcels, err := store.GetByClient(client)// получите список посылок по идентификатору клиента, сохранённого в переменной client
	requre.NoError(t, err)// убедитесь в отсутствии ошибки
	// убедитесь, что количество полученных посылок совпадает с количеством добавленных
assert.Len(t, storeParcels, 3)

	// check
	for _, parcel := range storedParcels {
		// в parcelMap лежат добавленные посылки, ключ - идентификатор посылки, значение - сама посылка
		// убедитесь, что все посылки из storedParcels есть в parcelMap
		// убедитесь, что значения полей полученных посылок заполнены верно
	assert.Equal(t, parcelMap[parcel.Number], parcel)
	}
}
