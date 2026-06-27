package events

// UserRegisteredEvent 用户注册事件
type UserRegisteredEvent struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func (e UserRegisteredEvent) RoutingKey() string { return "user.registered" }

// UserLoggedInEvent 用户登录事件
type UserLoggedInEvent struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	IP       string `json:"ip"`
}

func (e UserLoggedInEvent) RoutingKey() string { return "user.logged-in" }

// UserProfileUpdatedEvent 用户资料更新事件
type UserProfileUpdatedEvent struct {
	UserID   int64  `json:"user_id"`
	Nickname string `json:"nickname"`
}

func (e UserProfileUpdatedEvent) RoutingKey() string { return "user.profile-updated" }
