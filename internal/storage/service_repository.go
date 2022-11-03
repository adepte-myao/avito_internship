package storage

type ServiceRepository struct {
	storage *Storage
}

func NewServiceRepository(storage *Storage) *ServiceRepository {
	return &ServiceRepository{
		storage: storage,
	}
}
