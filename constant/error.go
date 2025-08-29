package constant

import "fmt"

var (
	ErrorUnauthorized       = fmt.Errorf("unauthorized")
	ErrorGenerateHash       = fmt.Errorf("failed to generate hashed password")
	ErrorComparePassword    = fmt.Errorf("failed to compare password")
	ErrorHandleTokenExpired = fmt.Errorf("token expired")
	ErrorUnknownClaimsType  = fmt.Errorf("unknown claims type")
	ErrorGetClaimSubject    = fmt.Errorf("failed to get claim subject")
	ErrorInternalServer     = fmt.Errorf("internal server error")
)
