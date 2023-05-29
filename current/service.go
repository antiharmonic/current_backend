package current

type Service interface {
	ListMedia(string, string, string) ([]Media, error)
	ListRecentMedia(string, string) ([]Media, error)
}

type Storage interface {
	MediaStorage
}

type serviceImpl struct {
	db Storage
}

func New(storage Storage) Service {
	return &serviceImpl{
		db: storage,
	}
}