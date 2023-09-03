package infrastructure

type Repository struct {
	urlShortenerMap map[string]string
}

func (r Repository) Save(url string, slug string) {
	r.urlShortenerMap[slug] = url
}

func (r Repository) Get(slug string) (string, bool) {
	value, ok := r.urlShortenerMap[slug]
	return value, ok
}

var RAMRepository = Repository{
	urlShortenerMap: map[string]string{},
}
