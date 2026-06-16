package errcode

import (
	"testing"
)

func TestNoDuplicateCodes(t *testing.T) {
	seen := map[int]string{}
	all := []*BizError{
		ErrProductNotFound,
		ErrInvalidParams,
		ErrPaginationQuery,
		ErrUnauthorized,
		ErrUserNotFound,
		ErrOrderNotFound,
		ErrDuplicateOrder,
		ErrPaymentFailed,
		ErrInvalidCredentials,
		ErrNotFound,
		ErrAccountDisabled,
		ErrWechatClientNotConfigured,
		ErrUsernameAlreadyExists,
		ErrUnsupportedProvider,
		ErrIdentityAlreadyBound,
		ErrInvalidToken,
		ErrTokenRevoked,
		ErrGenerateAccessToken,
		ErrGenerateRefreshToken,
		ErrSaveRefreshToken,
		ErrUnexpectedSigningMethod,
		ErrParseToken,
		ErrDuplicateSKU,
		ErrPermissionNotFound,
		ErrInsufficientPermissions,
		ErrCannotModifySystemRole,
		ErrCannotDeleteSystemRole,
		ErrReviewNotFound,
		ErrReviewDuplicate,
		ErrReviewNotPurchased,
		ErrReviewInvalidRating,
		ErrReviewMediaLimitExceed,
		ErrReviewNotOwner,
		ErrReviewPendingModeration,
	}
	for _, e := range all {
		if name, ok := seen[e.Code]; ok {
			t.Errorf("duplicate code %d: %s and %s", e.Code, name, e.Message)
		}
		seen[e.Code] = e.Message
	}
}
