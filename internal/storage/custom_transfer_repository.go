package storage

type CustomTransferRepository struct {
	storage *Storage
}

func NewCustomTransferRepository(storage *Storage) *CustomTransferRepository {
	return &CustomTransferRepository{
		storage: storage,
	}
}
