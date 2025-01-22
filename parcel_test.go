package main

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"reflect"
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
	db, errDB := sql.Open("sqlite", "tracker.db")

	require.Nil(t, errDB)
	require.NotNil(t, db)

	defer func() {
		err := db.Close()
		assert.Nil(t, err)
	}()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	id, errAdd := store.Add(parcel)
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
	require.Nil(t, errAdd)
	require.NotEmpty(t, id)

	// get
	p, errGet := store.Get(id)
	// получите только что добавленную посылку, убедитесь в отсутствии ошибки
	require.Nil(t, errGet)

	// проверьте, что значения всех полей в полученном объекте совпадают со значениями полей в переменной parcel
	// это невозможно т к number получаем во время записи в базу данных
	wont := reflect.ValueOf(parcel)
	get := reflect.ValueOf(p)

	// protects us from out of range in cycle
	require.Equal(t, wont.Type(), get.Type())

	lenFields := wont.NumField()

	// skip number -> i:=1
	for i := 1; i < lenFields; i++ {
		assert.Equal(
			t,
			wont.Field(i).Interface(),
			get.Field(i).Interface(),
		)
	}

	// delete
	errDel := store.Delete(id)
	// удалите добавленную посылку, убедитесь в отсутствии ошибки
	require.Nil(t, errDel)

	// проверьте, что посылку больше нельзя получить из БД
	_, errCheck := store.Get(id)
	require.NotNil(t, errCheck)
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	db, errDB := sql.Open("sqlite", "tracker.db")

	require.Nil(t, errDB)
	require.NotNil(t, db)

	defer func() {
		err := db.Close()
		assert.Nil(t, err)
	}()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	id, errAdd := store.Add(parcel)
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
	require.Nil(t, errAdd)
	require.NotEmpty(t, id)

	// set address
	// обновите адрес, убедитесь в отсутствии ошибки
	newAddress := "new test address"
	errSet := store.SetAddress(id, newAddress)
	require.Nil(t, errSet)

	// check
	getParcel, errGet := store.Get(id)
	// получите добавленную посылку и убедитесь, что адрес обновился
	require.Nil(t, errGet)
	assert.Equal(t, newAddress, getParcel.Address)
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	db, errDB := sql.Open("sqlite", "tracker.db")

	require.Nil(t, errDB)
	require.NotNil(t, db)

	defer func() {
		err := db.Close()
		assert.Nil(t, err)
	}()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	id, errAdd := store.Add(parcel)
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
	require.Nil(t, errAdd)
	require.NotEmpty(t, id)

	// set status
	errSet := store.SetStatus(id, ParcelStatusSent)
	// обновите статус, убедитесь в отсутствии ошибки
	require.Nil(t, errSet)

	// check
	p, errGet := store.Get(id)
	// получите добавленную посылку и убедитесь, что статус обновился
	require.Nil(t, errGet)
	assert.Equal(t, ParcelStatusSent, p.Status)
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	db, errDB := sql.Open("sqlite", "tracker.db")

	require.Nil(t, errDB)
	require.NotNil(t, db)

	defer func() {
		err := db.Close()
		assert.Nil(t, err)
	}()

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
		// отправляем посылку в БД, убедиться в отсутствии ошибки и наличии идентификатора
		id, err := store.Add(parcels[i])
		require.Nil(t, err)
		require.NotEmpty(t, id)

		// обновляем идентификатор добавленной у посылки
		parcels[i].Number = id

		// сохраняем добавленную посылку в структуру map, чтобы её можно было легко достать по идентификатору посылки
		parcelMap[id] = parcels[i]
	}

	// get by client
	// получить список посылок по идентификатору клиента, сохранённого в переменной client
	storedParcels, errGet := store.GetByClient(client)
	// убедитесь в отсутствии ошибки
	require.Nil(t, errGet)
	// убедитесь, что количество полученных посылок совпадает с количеством добавленных
	assert.Len(t, storedParcels, len(parcels))

	// check
	for _, parcel := range storedParcels {
		// в parcelMap лежат добавленные посылки, ключ - идентификатор посылки, значение - сама посылка
		// убедитесь, что все посылки из storedParcels есть в parcelMap
		p, ex := parcelMap[parcel.Number]
		assert.NotEmpty(t, ex)

		wont := reflect.ValueOf(p)
		get := reflect.ValueOf(parcel)

		require.Equal(t, wont.Type(), get.Type())

		lenFields := wont.NumField()

		// убедитесь, что значения полей полученных посылок заполнены верно
		for i := 0; i < lenFields; i++ {
			assert.Equal(
				t,
				wont.Field(i).Interface(),
				get.Field(i).Interface(),
			)
		}
	}
}
