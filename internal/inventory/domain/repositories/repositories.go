package repositories

type IinventoryRepository interface{}

func NewInventoryRepository(db interface{}) IinventoryRepository {
	return &inventoryRepository{}
}

type inventoryRepository struct{}
