package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/Yandex-Practicum/go-db-sql-final/internal/database"
	"github.com/Yandex-Practicum/go-db-sql-final/internal/models"
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
func getTestParcel() models.Parcel {
	return models.Parcel{
		Client:    1000,
		Status:    models.ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// TestAddGetDelete проверяет добавление, получение и удаление посылки
func TestAddGetDelete(t *testing.T) {
	// prepare
	db, err := database.Open()
	require.NoError(t, err)
	store := models.NewParcelStore(db)
	parcel := getTestParcel()

	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	// get
	// получите только что добавленную посылку, убедитесь в отсутствии ошибки
	// проверьте, что значения всех полей в полученном объекте совпадают со значениями полей в переменной parcel
	parc, err := store.Get(id)
	require.NoError(t, err)
	require.Equal(t, parcel.Client, parc.Client)
	require.Equal(t, parcel.Status, parc.Status)
	require.Equal(t, parcel.Address, parc.Address)
	require.Equal(t, parcel.CreatedAt, parc.CreatedAt)

	// delete
	// удалите добавленную посылку, убедитесь в отсутствии ошибки
	err = store.Delete(id)
	require.NoError(t, err)

	// проверьте, что посылку больше нельзя получить из БД
	_, err = store.Get(id)
	require.Equal(t, sql.ErrNoRows, err)

}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	// prepare
	db, err := database.Open()
	require.NoError(t, err)
	store := models.NewParcelStore(db)
	parcel := getTestParcel()

	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	// set address
	// обновите адрес, убедитесь в отсутствии ошибки
	newAddress := "new test address"
	err = store.SetAddress(id, newAddress)
	require.NoError(t, err)

	// check
	// получите добавленную посылку и убедитесь, что адрес обновился
	parc, err := store.Get(id)
	require.NoError(t, err)
	require.Equal(t, newAddress, parc.Address)
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	// prepare
	db, err := database.Open()
	require.NoError(t, err)
	store := models.NewParcelStore(db)
	parcel := getTestParcel()

	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	// set status
	// обновите статус, убедитесь в отсутствии ошибки
	err = store.SetStatus(id, models.ParcelStatusSent)
	require.NoError(t, err)

	// check
	// получите добавленную посылку и убедитесь, что статус обновился
	parc, err := store.Get(id)
	require.NoError(t, err)
	require.Equal(t, models.ParcelStatusSent, parc.Status)
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	// prepare
	db, err := database.Open()
	require.NoError(t, err)
	store := models.NewParcelStore(db)

	parcels := []models.Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]models.Parcel{}

	// задаём всем посылкам один и тот же идентификатор клиента
	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	// add
	for i := 0; i < len(parcels); i++ {
		// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
		id, err := store.Add(parcels[i])
		require.NoError(t, err)
		require.NotEmpty(t, id)

		// обновляем идентификатор добавленной у посылки
		parcels[i].Number = id

		// сохраняем добавленную посылку в структуру map, чтобы её можно было легко достать по идентификатору посылки
		parcelMap[id] = parcels[i]
	}

	// get by client
	// получите список посылок по идентификатору клиента, сохранённого в переменной client
	// убедитесь в отсутствии ошибки
	// убедитесь, что количество полученных посылок совпадает с количеством добавленных
	storedParcels, err := store.GetByClient(client)
	require.NoError(t, err)
	require.Equal(t, len(parcels), len(storedParcels))

	// check
	for _, parcel := range storedParcels {
		// в parcelMap лежат добавленные посылки, ключ - идентификатор посылки, значение - сама посылка
		// убедитесь, что все посылки из storedParcels есть в parcelMap
		// убедитесь, что значения полей полученных посылок заполнены верно
		_, ok := parcelMap[parcel.Number]
		require.Equal(t, true, ok)
		require.Equal(t, parcelMap[parcel.Number].Client, parcel.Client)
		require.Equal(t, parcelMap[parcel.Number].Status, parcel.Status)
		require.Equal(t, parcelMap[parcel.Number].Address, parcel.Address)
		require.Equal(t, parcelMap[parcel.Number].CreatedAt, parcel.CreatedAt)
	}
}
