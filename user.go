package avito

type User struct {
	Id       int    `json:"-" db:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Is_admin bool   `json:"is_admin"`
}
