package xerr

var (
	ErrInternalErr     = NewInternalErr(1001, "internal error")
	ErrSafeCall        = NewInternalErr(1002, "safe call error")
	ErrNoValidQuery    = NewBadRequestErrWithCode(1003, "no query param")
	ErrNoRows          = NewInternalErr(1004, "No records found")
	ErrNoRowsAffected  = NewInternalErr(1005, "No rows affected")
	ErrInvalidServerId = NewInternalErr(1006, "Invalid server id")

	/// auth
	ErrOauthInfo = NewInternalErr(1012, "invalid oauth info")
)
