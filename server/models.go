package server

type UpdateUsername struct {
	Username string `json:"username"`
	NewName  string `json:"new_name"`
}
