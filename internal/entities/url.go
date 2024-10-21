package entities

type URL struct {
	Id     uint64 `db:"id"`
	Url    string `db:"url"`
	Alias  string `db:"alias"`
	UserId uint64 `db:"user_id"`
}

func NewUrl(id uint64, url string, alias string, userId uint64) *URL {
	return &URL{Id: id, Url: url, Alias: alias, UserId: userId}
}
