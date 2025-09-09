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
	ErrorUserEmailExists    = fmt.Errorf("email already exists. please try another email")
	ErrorUserEmailNotExists = fmt.Errorf("email not found. please check again")
	ErrorInsertNewUser      = fmt.Errorf("error inserting new user, please try again later")
)
