package current

type Service interface {
	ListMedia(int, int, string) ([]Media, error)
	ListRecentMedia(int, int) ([]Media, error)
	StartMedia(int) (*Media, error)
	SearchMedia(string, int) ([]Media, error)
	GetMediaByID(int) (*Media, error)
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