package events

type AddressCreatedEvent struct {
	AddressID int64 `json:"address_id"`
	UserID    int64 `json:"user_id"`
}

func (e AddressCreatedEvent) RoutingKey() string { return "address.created" }

type AddressUpdatedEvent struct {
	AddressID int64 `json:"address_id"`
	UserID    int64 `json:"user_id"`
}

func (e AddressUpdatedEvent) RoutingKey() string { return "address.updated" }

type AddressDeletedEvent struct {
	AddressID int64 `json:"address_id"`
	UserID    int64 `json:"user_id"`
}

func (e AddressDeletedEvent) RoutingKey() string { return "address.deleted" }
