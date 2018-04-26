package generic_repository

type GenericItemInterface interface {
	GetId() string
	SetId(id string)
	GetSerializedItem() ([]byte, error)
	SetItem(serializedItem []byte) error
	Clone() *GenericItemInterface
}