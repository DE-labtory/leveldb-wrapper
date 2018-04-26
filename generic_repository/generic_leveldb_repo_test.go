package generic_repository

import (
	"testing"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
)

//setting data for testing
type FooItem struct {
	Id string
	Data1 string
	Data2 int
	Data3 float32
}

func (f FooItem) GetId() string {
	return f.Id
}

func (f FooItem) SetId(id string) {
	f.Id = id
}

func (f FooItem) GetSerializedItem() ([]byte, error) {
	data, err := json.Marshal(f)
	if err != nil {
		panic(fmt.Sprintf("Error encoding : %s", err))
	}
	return data, nil
}

func (f *FooItem) SetItem(serializedItem []byte) error {
	err := json.Unmarshal(serializedItem, f)
	if err != nil {
		panic(fmt.Sprintf("Error decoding : %s", err))
	}
	return nil
}

func (f FooItem) Clone() *GenericItemInterface {
	var cloneInterface GenericItemInterface = &f
	return &cloneInterface
}

func Deserialize(serializedBytes []byte, object interface{}) error {
	if len(serializedBytes) == 0 {
		return nil
	}
	err := json.Unmarshal(serializedBytes, object)
	if err != nil {
		panic(fmt.Sprintf("Error decoding : %s", err))
	}
	return err
}
func Serialize(object interface{}) ([]byte, error) {
	data, err := json.Marshal(object)
	if err != nil {
		panic(fmt.Sprintf("Error encoding : %s", err))
	}
	return data, nil
}

func newRepo() *GenericLevelDBRepository {
	path := "test_db_path"
	repo := NewGenericLevelDBRepository(path)
	return repo
}


//testing
func TestNewGenericLevelDBRepository(t *testing.T) {
	path := "test_db_path"
	defer os.RemoveAll(path)
	repo := NewGenericLevelDBRepository(path)
	defer repo.Close()
	assert.DirExists(t,path)
}

func TestGenericLevelDBRepository_Insert(t *testing.T) {
	path := "test_db_path"
	r := newRepo()
	defer os.RemoveAll(path)
	defer r.Close()
	inputItem := &FooItem{
		Id:    "first",
		Data1: "data1",
		Data2: 2,
		Data3: 3,
	}

	//case 1 not exist data
	err := r.Insert(inputItem)
	assert.NoError(t,err)

	data, err := r.leveldb.Get([]byte(inputItem.Id))
	readItem := &FooItem{}
	err = Deserialize(data,readItem)
	assert.NoError(t,err)

	assert.Equal(t,inputItem.Id,readItem.Id)
	assert.Equal(t,inputItem.Data1,readItem.Data1)
	assert.Equal(t,inputItem.Data2,readItem.Data2)
	assert.Equal(t,inputItem.Data3,readItem.Data3)

	//case 2 already exist data
	err = r.Insert(inputItem)
	assert.EqualError(t,err,AlreadyDataError)

	os.RemoveAll(path)
}

func TestGenericLevelDBRepository_Update(t *testing.T) {
	path := "test_db_path"
	r := newRepo()
	defer os.RemoveAll(path)
	defer r.Close()

	//case 1 data is exist
	inputItem := &FooItem{
		Id:    "first",
		Data1: "data1",
		Data2: 2,
		Data3: 3,
	}
	b,err := Serialize(inputItem)
	assert.NoError(t,err)
	err = r.leveldb.Put([]byte(inputItem.Id),b,true)
	assert.NoError(t,err)
	inputItem.Data1="updated"
	err = r.Update(inputItem)
	assert.NoError(t,err)
	data, err := r.leveldb.Get([]byte(inputItem.Id))
	readItem := &FooItem{}
	err = Deserialize(data,readItem)
	assert.NoError(t,err)
	assert.Equal(t,"updated",inputItem.Data1)

	//case 2 data is not exist
	inputItem.Id = "notExistId"
	err = r.Update(inputItem)
	assert.EqualError(t,err,NoDataError)

}

func TestGenericLevelDBRepository_FindAll(t *testing.T) {

	path := "test_db_path"
	r := newRepo()
	defer os.RemoveAll(path)
	defer r.Close()

	inputItem1 := &FooItem{
		Id:    "first",
		Data1: "data1",
		Data2: 2,
		Data3: 3,
	}
	inputItem2 := &FooItem{
		Id:    "second",
		Data1: "data1",
		Data2: 2,
		Data3: 3,
	}
	inputItem3 := &FooItem{
		Id:    "third",
		Data1: "data1",
		Data2: 2,
		Data3: 3,
	}

	err := r.Insert(inputItem1)
	assert.NoError(t, err)

	err = r.Insert(inputItem2)
	assert.NoError(t, err)

	err = r.Insert(inputItem3)
	assert.NoError(t, err)


	readDatas, err := r.FindAll(&FooItem{})
	assert.NoError(t, err)


	assert.Equal(t,(*readDatas[0]).GetId(),inputItem1.GetId())
	get , err := (*readDatas[0]).GetSerializedItem()
	assert.NoError(t,err)
	inputSerial, err := inputItem1.GetSerializedItem()
	assert.NoError(t,err)
	assert.Equal(t,get,inputSerial)

	assert.Equal(t,(*readDatas[1]).GetId(),inputItem2.GetId())
	get , err = (*readDatas[1]).GetSerializedItem()
	assert.NoError(t,err)
	inputSerial, err = inputItem2.GetSerializedItem()
	assert.NoError(t,err)
	assert.Equal(t,get,inputSerial)

	assert.Equal(t,(*readDatas[2]).GetId(),inputItem3.GetId())
	get , err = (*readDatas[2]).GetSerializedItem()
	assert.NoError(t,err)
	inputSerial, err = inputItem3.GetSerializedItem()
	assert.NoError(t,err)
	assert.Equal(t,get,inputSerial)

}

func TestGenericLevelDBRepository_Delete(t *testing.T) {
	path := "test_db_path"
	r := newRepo()
	defer os.RemoveAll(path)
	defer r.Close()


	inputItem := &FooItem{
		Id:    "first",
		Data1: "data1",
		Data2: 2,
		Data3: 3,
	}
	b,err := Serialize(inputItem)
	assert.NoError(t,err)
	err = r.leveldb.Put([]byte(inputItem.Id),b,true)
	assert.NoError(t,err)


	err = r.Delete(inputItem)
	assert.NoError(t,err)

	data,err := r.leveldb.Get([]byte(inputItem.Id))
	assert.NoError(t,err)
	assert.Nil(t,data)
}













