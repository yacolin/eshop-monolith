package shared

import "errors"

// 领域错误定义
var (
	// ErrNotFound 资源未找到
	ErrNotFound = errors.New("resource not found")

	// ErrInvalidInput 无效输入
	ErrInvalidInput = errors.New("invalid input")

	// ErrUnauthorized 未授权
	ErrUnauthorized = errors.New("unauthorized")

	// ErrForbidden 禁止访问
	ErrForbidden = errors.New("forbidden")

	// ErrConflict 冲突
	ErrConflict = errors.New("conflict")

	// ErrInternal 内部错误
	ErrInternal = errors.New("internal error")

	// ErrInsufficientInventory 库存不足
	ErrInsufficientInventory = errors.New("insufficient inventory")

	// ErrInvalidCredentials 无效凭证
	ErrInvalidCredentials = errors.New("invalid credentials")

	// ErrDuplicateResource 资源重复
	ErrDuplicateResource = errors.New("duplicate resource")

	// ErrInvalidStatus 无效状态
	ErrInvalidStatus = errors.New("invalid status")
)

// DomainError 领域错误结构
type DomainError struct {
	Err     error
	Message string
	Code    string
}

// Error 实现error接口
func (e *DomainError) Error() string {
	return e.Message
}

// Unwrap 实现errors.Unwrap
func (e *DomainError) Unwrap() error {
	return e.Err
}

// NewDomainError 创建领域错误
func NewDomainError(err error, message, code string) *DomainError {
	return &DomainError{
		Err:     err,
		Message: message,
		Code:    code,
	}
}