package generic_repository

type GenericRepositoryInterface interface {
	Insert(item GenericItemInterface) error
	Update(item GenericItemInterface) error
	FindById(receiver *GenericItemInterface, id string) (error)
	FindAll(receiverType GenericItemInterface) ([]*GenericItemInterface ,error)
	Delete(item GenericItemInterface) error
}
