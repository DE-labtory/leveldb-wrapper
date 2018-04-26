package generic_repository

import (
	"errors"

	"github.com/it-chain/leveldb-wrapper"
)

var NoDataError string = "NoDataError"
var AlreadyDataError string = "NoDataError"
type GenericLevelDBRepository struct {
	leveldb *leveldbwrapper.DB
}

func NewGenericLevelDBRepository(path string) *GenericLevelDBRepository {
	db := leveldbwrapper.CreateNewDB(path)
	db.Open()
	return &GenericLevelDBRepository{
		leveldb: db,
	}
}

func (gr GenericLevelDBRepository) Insert(item GenericItemInterface) error{
	if item.GetId() == "" {
		return errors.New("item ID is empty")
	}

	b, err := gr.leveldb.Get([]byte(item.GetId()))
	if len(b) != 0 {
		return errors.New(AlreadyDataError)
	}

	b, err = item.GetSerializedItem()
	if err != nil {
		return err
	}
	if err = gr.leveldb.Put([]byte(item.GetId()), b, true); err != nil {
		return err
	}

	return nil
}

func (gr GenericLevelDBRepository) Update(item GenericItemInterface) error {
	oldItem := item.Clone()
	err := gr.FindById(oldItem, item.GetId())
	if err != nil {
		return err
	}
	b, err := item.GetSerializedItem()
	if err != nil {
		return err
	}
	err = gr.leveldb.Put([]byte(item.GetId()), b, true)
	if err != nil {
		return err
	}
	return nil
}

func (gr GenericLevelDBRepository) FindById(receiver *GenericItemInterface, id string) error {
	b, err := gr.leveldb.Get([]byte(id))
	if err != nil {
		return err
	}
	if len(b) == 0 {
		return errors.New(NoDataError)
	}
	err = (*receiver).SetItem(b)
	if err != nil {
		return err
	}

	return nil
}

func (gr GenericLevelDBRepository) FindAll(receiverType GenericItemInterface) ([]*GenericItemInterface, error) {
	iter := gr.leveldb.GetIteratorWithPrefix([]byte(""))
	genericItems := []*GenericItemInterface{}
	for iter.Next() {
		if len(iter.Value()) == 0{
			continue
		}
		cloneData := receiverType.Clone()
		(*cloneData).SetId(string(iter.Key()))
		err := (*cloneData).SetItem(iter.Value())

		if err != nil {
			return nil, err
		}

		genericItems = append(genericItems, cloneData)
	}

	return genericItems, nil
}

func (gr GenericLevelDBRepository) Delete(item GenericItemInterface) error {
	return gr.leveldb.Delete([]byte(item.GetId()), true)
}

func (gr GenericLevelDBRepository) Close ()  {
	gr.leveldb.Close()
}