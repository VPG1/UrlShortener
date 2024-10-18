package entities

type URL struct {
	id    int
	url   string
	alias string
}

func NewUrl(id int, url string, alias string) *URL {
	return &URL{id: id, url: url, alias: alias}
}
