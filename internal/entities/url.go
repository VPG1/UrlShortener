package entities

type URL struct {
	Id    int
	Url   string
	Alias string
}

func NewUrl(id int, url string, alias string) *URL {
	return &URL{Id: id, Url: url, Alias: alias}
}
