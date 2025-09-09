package entity

type InsertUserReq struct {
	Name     string
	Email    string
	Password string
}

type LoginReq struct {
	Email    string
	Password string
}
