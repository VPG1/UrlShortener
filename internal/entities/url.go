package entities

type URL struct {
	Id    uint64 `db:"id"`
	Url   string `db:"url"`
	Alias string `db:"alias"`
}

func NewUrl(id uint64, url string, alias string) *URL {
	return &URL{Id: id, Url: url, Alias: alias}
}
