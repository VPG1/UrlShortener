package entities

type URL struct {
	Id     uint64 `db:"id" json:"id"`
	Url    string `db:"url" json:"url"`
	Alias  string `db:"alias" json:"alias"`
	UserId uint64 `db:"user_id" json"user_id"`
}

func NewUrl(id uint64, url string, alias string, userId uint64) *URL {
	return &URL{Id: id, Url: url, Alias: alias, UserId: userId}
}
