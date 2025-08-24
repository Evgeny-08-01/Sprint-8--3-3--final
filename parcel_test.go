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
 // настройте подключение к БД
	 db, err := sql.Open("sqlite", "tracker.db")
    if err != nil {
        fmt.Println(err)
		require.NotNil(t,err,"ошибка при подключении к  db, err := sql.Open(`sqlite`, `tracker.db`)")
        return
    }
    defer db.Close()
	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
newLine,err:=db.Exec("INSERT INTO parcel (client,status,address,created_at) VALUES (:Client,:Status,:Address,:CreatedAt)", 
		sql.Named("Client", parcel.Client),
        sql.Named("Status", parcel.Status),
		sql.Named("Address", parcel.Address),
        sql.Named("CreatedAt", parcel.CreatedAt))
		assert.Nil(t,err,"ошибка при добавлении посылки")
	lastID,err:=newLine.LastInsertId()
	if assert.Nil(t,err,"ошибка при добавлении посылки"){
	assert.NotEmpty(t,lastID,"нет идентификатора посылки")}
	
	// get
	// получите только что добавленную посылку, убедитесь в отсутствии ошибки

	parcelReal, err := store.Get(int(lastID))
 assert.Nil(t,err,"в parcelReal-добавленная посылка, ошибка")
// проверьте, что значения всех полей в полученном объекте совпадают со значениями полей в переменной parcel
assert.Equal(t,parcel.Address,parcelReal.Address,"в parcelReal-адрес не совпадает")
assert.Equal(t,parcel.Client,parcelReal.Client,"в parcelReal-клиент не совпадает")
assert.Equal(t,parcel.CreatedAt,parcelReal.CreatedAt,"в parcelReal-время создания не совпадает")
//assert.Equal(t,parcel.Number,parcelReal.Number,"в parcelReal-номер не совпадает")
assert.Equal(t,parcel.Status,parcelReal.Status,"в parcelReal-статус не совпадает")
	// delete
	// удалите добавленную посылку, убедитесь в отсутствии ошибки
 err = store.Delete(int(lastID))
    assert.Nil(t,err,"ошибка при удалении Delete")
	// проверьте, что посылку больше нельзя получить из БД
parcelReal, err = store.Get(int(lastID))
if assert.NotNil(t,err,"Get не выдает ошибку. ОШИБКА"){
	assert.Empty(t,parcelReal,"после Get есть инф о посылке. ОШИБКА")}
	}
// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	// prepare
	// настройте подключение к БД
 db, err := sql.Open("sqlite", "tracker.db")
    if err != nil {
        fmt.Println(err)
		require.NotNil(t,err,"ошибка при подключении к  db, err := sql.Open(`sqlite`, `tracker.db`)")
        return
    }
    defer db.Close()
	store := NewParcelStore(db)
	parcel := getTestParcel()
	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
newLine,err:=db.Exec("INSERT INTO parcel (client,status,address,created_at) VALUES (:Client,:Status,:Address,:CreatedAt)", 
		sql.Named("Client", parcel.Client),
        sql.Named("Status", parcel.Status),
		sql.Named("Address", parcel.Address),
        sql.Named("CreatedAt", parcel.CreatedAt))
		assert.Nil(t,err,"ошибка при добавлении посылки")
	lastID,err:=newLine.LastInsertId()
	if assert.Nil(t,err,"ошибка при добавлении посылки"){
	assert.NotEmpty(t,lastID,"идентификатора посылки нет")}
	// set address
	// обновите адрес, убедитесь в отсутствии ошибки
	newAddress := "new test address"
	err =store.SetAddress(int(lastID),newAddress)
    assert.Nil(t,err,"новый адрес не добавлен. ОШИБКА")
	// check
	// получите добавленную посылку и убедитесь, что адрес обновился
 parcelReal, err := store.Get(int(lastID))
 assert.Nil(t,err,"ошибка при получении инф о добавленной посылке")
 assert.Equal(t,parcelReal.Address,newAddress,"адреса не совпали после SetAddress")
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	// prepare
	// настройте подключение к БД
 db, err := sql.Open("sqlite", "tracker.db")
    if err != nil {
        fmt.Println(err)
		require.NotNil(t,err,"ошибка при подключении к  db, err := sql.Open(`sqlite`, `tracker.db`)")
        return
    }
    defer db.Close()
	store := NewParcelStore(db)
	parcel := getTestParcel()
	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
newLine,err:=db.Exec("INSERT INTO parcel (client,status,address,created_at) VALUES (:Client,:Status,:Address,:CreatedAt)", 
		sql.Named("Client", parcel.Client),
        sql.Named("Status", parcel.Status),
		sql.Named("Address", parcel.Address),
        sql.Named("CreatedAt", parcel.CreatedAt))
		assert.Nil(t,err,"ошибкак при добавлении посылки")
	lastID,err:=newLine.LastInsertId()
	if assert.Nil(t,err,"ошибка при добавлении посылки"){
	assert.NotEmpty(t,lastID,"идентификатора посылки нет")}
	// set status
	// обновите статус, убедитесь в отсутствии ошибки
	newStatus:=ParcelStatusSent
	err =store.SetStatus(int(lastID),newStatus)
    assert.Nil(t,err,"новый статус не присвоен")
	// check
	// получите добавленную посылку и убедитесь, что статус обновился
	 parcelReal, err := store.Get(int(lastID))
 assert.Nil(t,err,"ошибка при получении инф о добавленной посылке")
 assert.Equal(t,parcelReal.Address,parcel.Address,"статус не совпал после SetStatus")
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	// prepare
// настройте подключение к БД
 db, err := sql.Open("sqlite", "tracker.db")
    if err != nil {
        fmt.Println(err)
		require.NotNil(t,err,"ошибка при подключении к  db, err := sql.Open(`sqlite`, `tracker.db`)")
        return
    }
    defer db.Close()
    store := NewParcelStore(db)
	//parcel := getTestParcel()
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
newLine,err:=db.Exec("INSERT INTO parcel (client,status,address,created_at) VALUES (:Client,:Status,:Address,:CreatedAt)", 
		sql.Named("Client", parcels[i].Client),
        sql.Named("Status", parcels[i].Status),
		sql.Named("Address", parcels[i].Address),
        sql.Named("CreatedAt", parcels[i].CreatedAt))
		assert.Nil(t,err,"ошибка при добавлении посылки")
	lastID,err:=newLine.LastInsertId()
	if assert.Nil(t,err,"ошибка при добавлении посылки"){
	assert.NotEmpty(t,lastID,"идентификатора посылки нет")}
	// обновляем идентификатор добавленной у посылки
	 parcels[i].Number=int(lastID)

	// сохраняем добавленную посылку в структуру map, чтобы её можно было легко достать по идентификатору посылки
		parcelMap[i] = parcels[i]
	}

	// get by client
	 // получите список посылок по идентификатору клиента, сохранённого в переменной client
    storedParcels, err := store.GetByClient(client)

	// убедитесь в отсутствии ошибки
	// убедитесь, что количество полученных посылок совпадает с количеством добавленных
 assert.Nil(t,err,"ошибка при получении списка посылок")
 assert.Equal(t,len(parcels),len(storedParcels),"количество полученных и добавленных посылок нет совпало")

	// check
	for i, parcel := range storedParcels {
 assert.Equal(t,parcel,parcelMap[i],"посылки не совпадают: посылка №%d(%d)я итерация",i+1,i)



		// в parcelMap лежат добавленные посылки, ключ - идентификатор посылки, значение - сама посылка
		// убедитесь, что все посылки из storedParcels есть в parcelMap
		// убедитесь, что значения полей полученных посылок заполнены верно
	}
}
