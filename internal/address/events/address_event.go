package events

type AddressCreatedEvent struct {
	AddressID int64 `json:"address_id"`
	UserID    int64 `json:"user_id"`
}

type AddressUpdatedEvent struct {
	AddressID int64 `json:"address_id"`
	UserID    int64 `json:"user_id"`
}

type AddressDeletedEvent struct {
	AddressID int64 `json:"address_id"`
	UserID    int64 `json:"user_id"`
}
