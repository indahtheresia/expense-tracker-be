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
	ErrorGetCategories      = fmt.Errorf("failed to get categories. please try again")
	ErrorAddExpense         = fmt.Errorf("failed to add expense. please try again")
	ErrorParsingDate        = fmt.Errorf("invalid date format, use YYYY-MM-DD")
	ErrorUpdateExpense      = fmt.Errorf("failed to edit expense. please try again")
	ErrorDeleteExpense      = fmt.Errorf("failed to delete expense. please try again")
	ErrorExpenseNotFound    = fmt.Errorf("invalid expense id. please enter valid expense id")
	ErrorGetExpenses        = fmt.Errorf("failed to get all expenses. please try again")
)
