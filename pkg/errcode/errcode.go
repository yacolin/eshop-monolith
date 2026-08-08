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
//    通用域 1001-1999，权限域 2001-2999，评论评分域 3001-3999。
//    后续新增域在 4001+ 继续分配。
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

	ErrDuplicateSKU          = &BizError{Code: 1023, Message: "duplicate sku"}
	ErrInsufficientInventory = &BizError{Code: 1024, Message: "insufficient inventory"}
)

// ==================== 域：权限（2001-2999） ====================
var (
	ErrPermissionNotFound      = &BizError{Code: 2001, Message: "permission not found"}
	ErrInsufficientPermissions = &BizError{Code: 2002, Message: "insufficient permissions"}
	ErrCannotModifySystemRole  = &BizError{Code: 2003, Message: "cannot modify system role"}
	ErrCannotDeleteSystemRole  = &BizError{Code: 2004, Message: "cannot delete system role"}
	ErrCannotOperateSelf       = &BizError{Code: 2005, Message: "cannot delete or disable self or builtin admin"}
)

// ==================== 域：品牌/类目（4001-4099） ====================
var (
	ErrBrandNotFound          = &BizError{Code: 4001, Message: "brand not found"}
	ErrBrandNameExists        = &BizError{Code: 4002, Message: "brand name already exists"}
	ErrCategoryNotFound       = &BizError{Code: 4010, Message: "category not found"}
	ErrCategoryNameExists     = &BizError{Code: 4011, Message: "category name already exists"}
	ErrCategoryHasChildren    = &BizError{Code: 4012, Message: "category has children, cannot delete"}
	ErrCategoryParentNotFound = &BizError{Code: 4013, Message: "parent category not found"}
	ErrCategoryLevelExceed    = &BizError{Code: 4014, Message: "category level exceeds maximum (3)"}
	ErrCategoryLevelInvalid   = &BizError{Code: 4015, Message: "category level must be between 1 and 3"}
	ErrAttributeNotFound      = &BizError{Code: 4020, Message: "attribute not found"}
	ErrSPUNotFound            = &BizError{Code: 4030, Message: "product not found"}
	ErrSKUCodeExists          = &BizError{Code: 4031, Message: "sku code already exists"}
	ErrSPUInvalidStatus       = &BizError{Code: 4032, Message: "invalid product status transition"}
	ErrProductAttrDuplicate   = &BizError{Code: 4033, Message: "duplicate product attribute"}
	ErrSPUHasNoSKU            = &BizError{Code: 4034, Message: "product must have at least one sku"}

	// ==================== 域：库存（5001-5099） ====================
	ErrInventoryNotFound  = &BizError{Code: 5001, Message: "inventory not found"}
	ErrInsufficientStock  = &BizError{Code: 5002, Message: "insufficient stock"}
	ErrInvalidStockChange = &BizError{Code: 5003, Message: "invalid stock change"}

	// ==================== 域：交易（6001-6099） ====================
	ErrPaymentNotFound = &BizError{Code: 6001, Message: "payment not found"}
	ErrRefundNotFound  = &BizError{Code: 6010, Message: "refund not found"}
	ErrRefundFailed    = &BizError{Code: 6011, Message: "refund failed"}

	// ==================== 域：用户/地址（9001-9099） ====================
	ErrAddressLimit    = &BizError{Code: 9001, Message: "address limit reached"}
	ErrAddressNotFound = &BizError{Code: 9002, Message: "address not found"}

	// ==================== 域：订单（7001-7099） ====================
	ErrInvalidOrderStatus = &BizError{Code: 7002, Message: "invalid order status transition"}

	// ==================== 域：营销（8001-8099） ====================
	ErrPromotionNotFound        = &BizError{Code: 8001, Message: "promotion not found"}
	ErrPromotionRuleInvalid     = &BizError{Code: 8002, Message: "invalid promotion rule"}
	ErrPromotionProductConflict = &BizError{Code: 8003, Message: "promotion product already exists"}
	ErrCouponExpired            = &BizError{Code: 8010, Message: "coupon expired"}
	ErrCouponSoldOut            = &BizError{Code: 8011, Message: "coupon sold out"}
	ErrCouponAlreadyClaimed     = &BizError{Code: 8012, Message: "coupon already claimed"}
	ErrOrderItemNotFound        = &BizError{Code: 7003, Message: "order item not found"}
)

// ==================== 域：评论评分（3001-3999） ====================
var (
	ErrReviewNotFound          = &BizError{Code: 3001, Message: "review not found"}
	ErrReviewDuplicate         = &BizError{Code: 3002, Message: "review already exists for this order item"}
	ErrReviewNotPurchased      = &BizError{Code: 3003, Message: "only purchased products can be reviewed"}
	ErrReviewInvalidRating     = &BizError{Code: 3004, Message: "rating must be between 1 and 5"}
	ErrReviewMediaLimitExceed  = &BizError{Code: 3005, Message: "media count exceeds the limit"}
	ErrReviewNotOwner          = &BizError{Code: 3006, Message: "not the owner of the review"}
	ErrReviewPendingModeration = &BizError{Code: 3007, Message: "review is pending moderation"}
)
