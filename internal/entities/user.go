package entities

type User struct {
	Id           uint64 `db:"id"`
	Name         string `db:"name"`
	UserName     string `db:"username"`
	PasswordHash string `db:"password_hash"`
}

func NewUser(id uint64, name string, UserName string, PasswordHash string) *User {
	return &User{Id: id, Name: name, UserName: UserName, PasswordHash: PasswordHash}
}
