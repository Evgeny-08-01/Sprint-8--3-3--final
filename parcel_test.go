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
	// настройте подключение к БД
	db, err := sql.Open("sqlite", "tracker.db")
			require.NoError(t, err, "ошибка при подключении к  db, err := sql.Open(`sqlite`, `tracker.db`)")
	defer db.Close()
	store := NewParcelStore(db)
	parcel := getTestParcel()
	service := NewParcelService(store)

	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
	newLine, err := service.Register(parcel.Client, parcel.Address)
	require.NoError(t, err, "ошибка при добавлении посылки")
	assert.NotEmpty(t, newLine.Number, "нет идентификатора посылки")

	// get
	// получите только что добавленную посылку, убедитесь в отсутствии ошибки
	parcelReal, err := store.Get(newLine.Number)
	assert.NoError(t, err, "в parcelReal-добавленная посылка, ошибка")
	parcel.Number = newLine.Number
	assert.Equal(t, parcel, parcelReal, "полученная из БД добавленная посылка-не совпадает с добавленной тестом посылкой")
	
	// delete
	// удалите добавленную посылку, убедитесь в отсутствии ошибки
	err = store.Delete(newLine.Number)
	assert.NoError(t, err, "ошибка при удалении Delete")
	// проверьте, что посылку больше нельзя получить из БД
	_, err = store.Get(newLine.Number)
	require.ErrorIs(t, err, sql.ErrNoRows,"посылка не была удалена. ОШИБКА")
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	// prepare
	// настройте подключение к БД
	db, err := sql.Open("sqlite", "tracker.db")
		require.NoError(t, err, "ошибка при подключении к  db, err := sql.Open(`sqlite`, `tracker.db`)")
	defer db.Close()
	store := NewParcelStore(db)
	parcel := getTestParcel()
	service := NewParcelService(store)
	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
	newLine, err := service.Register(parcel.Client, parcel.Address)
	require.NoError(t, err, "ошибка при добавлении посылки")
	assert.NotEmpty(t, newLine.Number, "нет идентификатора посылки")

	// set address
	// обновите адрес, убедитесь в отсутствии ошибки
	newAddress := "new test address"
	err = store.SetAddress(newLine.Number, newAddress)
	require.NoError(t, err, "ошибка при вызове SetAddress")
	// check
	// получите добавленную посылку и убедитесь, что адрес обновился
	parcelReal, err := store.Get(newLine.Number)
	require.NoError(t, err, "ошибка при получении посылки")
	assert.Equal(t, parcelReal.Address, newAddress, "адреса не совпали после SetAddress")
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	// prepare
	// настройте подключение к БД
	db, err := sql.Open("sqlite", "tracker.db")
		require.NoError(t, err, "ошибка при подключении к  db, err := sql.Open(`sqlite`, `tracker.db`)")
	defer db.Close()
	store := NewParcelStore(db)
	parcel := getTestParcel()
	service := NewParcelService(store)
	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
	newLine, err := service.Register(parcel.Client, parcel.Address)
	require.NoError(t, err, "ошибка при добавлении посылки")
	require.NotEmpty(t, newLine.Number, "нет идентификатора посылки")
	// set status
	// обновите статус, убедитесь в отсутствии ошибки
	newStatus := ParcelStatusSent
	err = store.SetStatus(newLine.Number, newStatus)
	require.NoError(t, err, "ошибка при вызове SetStatus")
	// check
	// получите добавленную посылку и убедитесь, что статус обновился
	parcelReal, err := store.Get(newLine.Number)
	require.NoError(t, err, "ошибка при добавлении посылки, при вызове Get")
	assert.Equal(t, parcelReal.Address, parcel.Address, "статус не совпал после SetStatus")
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	// prepare
	// настройте подключение к БД
	db, err := sql.Open("sqlite", "tracker.db")
		require.NoError(t, err, "ошибка при подключении к  db, err := sql.Open(`sqlite`, `tracker.db`)")
	defer db.Close()
	store := NewParcelStore(db)
	service := NewParcelService(store)
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
		// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
		id, err := service.Register(parcels[i].Client, parcels[i].Address)
		if !assert.NoError(t, err, "ошибка при добавлении посылки при %d итерации, при вызове Register",i+1)||!assert.NotEmpty(t, 
			id.Number, "нет идентификатора посылки"){ 
					return}
		
		// обновляем идентификатор добавленной у посылки
		parcels[i].Number = id.Number

		// сохраняем добавленную посылку в структуру map, чтобы её можно было легко достать по идентификатору посылки
		parcelMap[id.Number] = parcels[i]
	}
	// get by client
	// получите список посылок по идентификатору клиента, сохранённого в переменной client
	// убедитесь в отсутствии ошибки
	// убедитесь, что количество полученных посылок совпадает с количеством добавленных
	storedParcels, err := store.GetByClient(client)
	require.NoError(t, err, "ошибка при получении списка посылок client")
	assert.Len(t, storedParcels, len(parcels), "количество добавленных и полученных посылок не совпадает")

	// check
	/*for _, parcel := range storedParcels {
		if _, ok := parcelMap[parcel.Number]; ok {
			assert.Equal(t, parcel, parcelMap[parcel.Number], "значения полей посылки №%d не совпадают", parcel.Number)
		} else {
			assert.True(t, ok, "посылка с ID%d не существует", parcel.Number)
		}
		// в parcelMap лежат добавленные посылки, ключ - идентификатор посылки, значение - сама посылка
		// убедитесь, что все посылки из storedParcels есть в parcelMap
		// убедитесь, что значения полей полученных посылок заполнены верно

	}
}*/

for _, parcel := range storedParcels {
		expectedParcel, ok := parcelMap[parcel.Number]
		require.True(t, ok)
		require.Equal(t, expectedParcel, parcel)}}