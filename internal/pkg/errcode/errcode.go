package errcode

import "fmt"

type BizError struct {
	Code    int
	Message string
}

func (e *BizError) Error() string {
	return fmt.Sprintf("code: %d, message: %s", e.Code, e.Message)
}

// 常用业务错误
var (
	ErrProductNotFound = &BizError{Code: 1001, Message: "product not found"}
	ErrStockNotEnough  = &BizError{Code: 1002, Message: "stock not enough"}
	ErrInvalidParams   = &BizError{Code: 1003, Message: "invalid parameters"}
	ErrPaginationQuery = &BizError{Code: 1004, Message: "invalid pagination query"}
	ErrUnauthorized    = &BizError{Code: 1005, Message: "unauthorized"}
	ErrUserNotFound    = &BizError{Code: 1006, Message: "user not found"}
	ErrOrderNotFound   = &BizError{Code: 1007, Message: "order not found"}
	ErrDuplicateOrder  = &BizError{Code: 1008, Message: "duplicate order"}
	ErrPaymentFailed   = &BizError{Code: 1009, Message: "payment failed"}

	ErrInvalidCredentials     = &BizError{Code: 1010, Message: "invalid credentials"}
	ErrEmailAlreadyRegistered = &BizError{Code: 1011, Message: "email already registered"}
	ErrUserAlreadyRegistered  = &BizError{Code: 1012, Message: "user already registered"}

	ErrNotFound = &BizError{Code: 1013, Message: "resource not found"}

	ErrAccountDisabled           = &BizError{Code: 1014, Message: "account disabled"}
	ErrWechatClientNotConfigured = &BizError{Code: 1015, Message: "wechat client not configured"}
	ErrUsernameAlreadyExists     = &BizError{Code: 1016, Message: "username already exists"}
	ErrUnsupportedProvider       = &BizError{Code: 1017, Message: "unsupported provider"}
	ErrIdentityAlreadyBound      = &BizError{Code: 1018, Message: "identity already bound"}

	ErrInvalidToken = &BizError{Code: 1019, Message: "invalid token"}
	ErrTokenRevoked = &BizError{Code: 1020, Message: "token revoked"}

	ErrGenerateAccessToken     = &BizError{Code: 1021, Message: "generate access token failed"}
	ErrGenerateRefreshToken    = &BizError{Code: 1022, Message: "generate refresh token failed"}
	ErrSaveRefreshToken        = &BizError{Code: 1023, Message: "save refresh token failed"}
	ErrUnexpectedSigningMethod = &BizError{Code: 1024, Message: "unexpected signing method"}
	ErrParseToken              = &BizError{Code: 1025, Message: "parse token failed"}

	// 权限相关错误
	ErrPermissionNotFound      = &BizError{Code: 2001, Message: "permission not found"}
	ErrPermissionAlreadyExists = &BizError{Code: 2002, Message: "permission already exists"}
	ErrInvalidRoleName         = &BizError{Code: 2003, Message: "invalid role name"}
	ErrInsufficientPermissions = &BizError{Code: 2004, Message: "insufficient permissions"}
	ErrCannotModifySystemRole  = &BizError{Code: 2005, Message: "cannot modify system role"}
	ErrCannotDeleteSystemRole  = &BizError{Code: 2006, Message: "cannot delete system role"}
)
