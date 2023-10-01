package domain

type ShortRepositoryInerface interface {
	Save(short Short)
	Get(slug string) (string, bool)
}
