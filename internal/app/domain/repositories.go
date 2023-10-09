package domain

type ShortRepositoryInerface interface {
	Save(short Short) error
	Get(slug string) (string, bool)
	Ping() error
}
