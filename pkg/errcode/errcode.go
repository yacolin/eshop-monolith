package errcode

import "fmt"

// BizError 业务错误。
//
// Code 是稳定标识，一旦发布即永久保留，不得修改、删除或复用。
// 后续业务迭代如需废弃某错误码，仅标记为 Deprecated，保留定义但不使用。
type BizError struct {
	Code    int
	Message string
}

func (e *BizError) Error() string {
	return fmt.Sprintf("code: %d, message: %s", e.Code, e.Message)
}

// ============================================================================
// 错误码规范
//
// 1. 数字硬编码，对外稳定
//    错误码是 API 契约的一部分，客户端会依赖具体数字做逻辑判断。
//    一旦发布，该数字永久归属该错误，不得修改、删除或复用于其他错误。
//
// 2. 每个域预留一段区间
//    通用域 1001-1999，权限域 2001-2999。后续新增域在 3001+ 继续分配。
//
// 3. 新增错误码
//    在所属域的 var 块末尾追加一行，选区间内下一个未使用的数字。
//    运行 go test ./pkg/errcode/ 验证无重复。
//
// 4. 废弃错误码
//    保留定义，注释标记 Deprecated。禁止删除或修改 Code，
//    否则已集成该 code 的客户端会静默地误解响应。
//
//    例：
//    // Deprecated: 改用 ErrXXX
//    ErrOld = &BizError{Code: 10XX, Message: "..."}
//
// 5. 测试保障
//    TestNoDuplicateCodes 检测所有已定义错误码是否有重复，确保编码质量。
// ============================================================================

// ==================== 域：通用（1001-1999） ====================
var (
	ErrProductNotFound = &BizError{Code: 1001, Message: "product not found"}
	ErrInvalidParams   = &BizError{Code: 1002, Message: "invalid parameters"}
	ErrPaginationQuery = &BizError{Code: 1003, Message: "invalid pagination query"}
	ErrUnauthorized    = &BizError{Code: 1004, Message: "unauthorized"}
	ErrUserNotFound    = &BizError{Code: 1005, Message: "user not found"}
	ErrOrderNotFound   = &BizError{Code: 1006, Message: "order not found"}
	ErrDuplicateOrder  = &BizError{Code: 1007, Message: "duplicate order"}
	ErrPaymentFailed   = &BizError{Code: 1008, Message: "payment failed"}

	ErrInvalidCredentials = &BizError{Code: 1009, Message: "invalid credentials"}
	ErrNotFound           = &BizError{Code: 1010, Message: "resource not found"}

	ErrAccountDisabled           = &BizError{Code: 1011, Message: "account disabled"}
	ErrWechatClientNotConfigured = &BizError{Code: 1012, Message: "wechat client not configured"}
	ErrUsernameAlreadyExists     = &BizError{Code: 1013, Message: "username already exists"}
	ErrUnsupportedProvider       = &BizError{Code: 1014, Message: "unsupported provider"}
	ErrIdentityAlreadyBound      = &BizError{Code: 1015, Message: "identity already bound"}

	ErrInvalidToken = &BizError{Code: 1016, Message: "invalid token"}
	ErrTokenRevoked = &BizError{Code: 1017, Message: "token revoked"}

	ErrGenerateAccessToken     = &BizError{Code: 1018, Message: "generate access token failed"}
	ErrGenerateRefreshToken    = &BizError{Code: 1019, Message: "generate refresh token failed"}
	ErrSaveRefreshToken        = &BizError{Code: 1020, Message: "save refresh token failed"}
	ErrUnexpectedSigningMethod = &BizError{Code: 1021, Message: "unexpected signing method"}
	ErrParseToken              = &BizError{Code: 1022, Message: "parse token failed"}

	ErrDuplicateSKU           = &BizError{Code: 1023, Message: "duplicate sku"}
	ErrInsufficientInventory = &BizError{Code: 1024, Message: "insufficient inventory"}
)

// ==================== 域：权限（2001-2999） ====================
var (
	ErrPermissionNotFound      = &BizError{Code: 2001, Message: "permission not found"}
	ErrInsufficientPermissions = &BizError{Code: 2002, Message: "insufficient permissions"}
	ErrCannotModifySystemRole  = &BizError{Code: 2003, Message: "cannot modify system role"}
	ErrCannotDeleteSystemRole  = &BizError{Code: 2004, Message: "cannot delete system role"}
)
