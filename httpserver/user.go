package httpserver

type User struct {
	UserName 	string  `json:"user_name"`
	PasWord 	string  `json:"pas_word"`
	Status  	int		`json:"status"`
}

//用户列表
type UserList struct {
	Users  []User  `json:"users"`
}


