package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

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
	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
	id, err := store.Add(parcel)
	require.NoError(t, err, "Ошибка при добавлении посылки в базу данных")
	require.NotZero(t, id, "Идентификатор добавленной посылки должен быть больше нуля")

	parcel.Number = id

	// get
	// получите только что добавленную посылку, убедитесь в отсутствии ошибки
	// проверьте, что значения всех полей в полученном объекте совпадают со значениями полей в переменной parcel
	retrievedParcel, err := store.Get(parcel.Number)
	require.NoError(t, err, "Ошибка при получении посылки из базы данных")
	require.Equal(t, parcel, retrievedParcel, "Полученные данные не совпадают с исходными")

	// delete
	// удалите добавленную посылку, убедитесь в отсутствии ошибки
	err = store.Delete(parcel.Number) // удаляем добавленную посылку
	require.NoError(t, err, "Ошибка при удалении посылки из базы данных")
	// проверьте, что посылку больше нельзя получить из БД
	_, err = store.Get(parcel.Number)
	require.Error(t, err, "Посылка должна быть удалена, но она всё ещё существует")
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err, "Ошибка при подключении к базе данных")
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
	id, err := store.Add(parcel) // добавляем посылку в базу данных
	require.NoError(t, err, "Ошибка при добавлении посылки в базу данных")
	require.NotZero(t, id, "Идентификатор добавленной посылки должен быть больше нуля")

	parcel.Number = id

	// set address
	// обновите адрес, убедитесь в отсутствии ошибки
	newAddress := "new test address"
	err = store.SetAddress(parcel.Number, newAddress)
	require.NoError(t, err, "Ошибка при обновлении адреса посылки")

	// check
	// получите добавленную посылку и убедитесь, что адрес обновился
	updatedParcel, err := store.Get(parcel.Number)
	require.NoError(t, err, "Ошибка при получении обновленной посылки из базы данных")
	require.Equal(t, newAddress, updatedParcel.Address, "Адрес не был обновлен")

}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err, "Ошибка при подключении к базе данных")
	defer db.Close()

	store := NewParcelStore(db) // создаем экземпляр ParcelStore
	parcel := getTestParcel()

	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
	id, err := store.Add(parcel) // добавляем посылку в базу данных
	require.NoError(t, err, "Ошибка при добавлении посылки в базу данных")
	require.NotZero(t, id, "Идентификатор добавленной посылки должен быть больше нуля")

	parcel.Number = id

	// set status
	// обновите статус, убедитесь в отсутствии ошибки
	newStatus := ParcelStatusSent
	err = store.SetStatus(parcel.Number, newStatus)
	require.NoError(t, err, "Ошибка при обновлении статуса посылки")

	// check
	// получите добавленную посылку и убедитесь, что статус обновился
	updatedParcel, err := store.Get(parcel.Number)
	require.NoError(t, err, "Ошибка при получении обновленной посылки из базы данных")
	require.Equal(t, newStatus, updatedParcel.Status, "Статус не был обновлен")

}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")                     // настройка подключения к базе данных
	require.NoError(t, err, "Ошибка при подключении к базе данных") // проверка на ошибку при подключении
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
		require.NoError(t, err, "Ошибка при добавлении посылки в базу данных")
		require.NotZero(t, id, "Идентификатор добавленной посылки должен быть больше нуля")

		// обновляем идентификатор добавленной у посылки
		parcels[i].Number = id

		// сохраняем добавленную посылку в структуру map, чтобы её можно было легко достать по идентификатору посылки
		parcelMap[id] = parcels[i]
	}

	// get by client
	storedParcels, err := store.GetByClient(client)
	require.NoError(t, err, "Ошибка при получении посылок по идентификатору клиента")
	require.Len(t, storedParcels, len(parcels), "Количество полученных посылок не совпадает с количеством добавленных")

	// убедитесь в отсутствии ошибки
	// убедитесь, что количество полученных посылок совпадает с количеством добавленных

	// check
	for _, parcel := range storedParcels {
		expectedParcel, exists := parcelMap[parcel.Number]
		require.True(t, exists, "Посылка не найдена в карте по идентификатору")

		// убедитесь, что все поля у посылки совпадают
		require.Equal(t, expectedParcel.Client, parcel.Client, "Идентификатор клиента не совпадает")
		require.Equal(t, expectedParcel.Status, parcel.Status, "Статус не совпадает")
		require.Equal(t, expectedParcel.Address, parcel.Address, "Адрес не совпадает")
		require.Equal(t, expectedParcel.CreatedAt, parcel.CreatedAt, "Дата создания не совпадает")
	}
}
