package current

import (
	_ "fmt"
	_ "encoding/json"
	"database/sql"
	"github.com/jackc/pgtype"
)

// the exotic types are because golang refuses to assign nulls to things like
// strings and time.Time. So instead you use a type that can be empty but has a 
// "valid" property to let you know if it was null or empty.
type Media struct {
	ID			NullInt	`db:"id" json:"id"`
	Title		NullString	`json:"title"`
	MediaType	NullInt		`db:"type" json:"type"`
	Weight		sql.NullFloat64	`json:"weight"`
	DateAdded	pgtype.Date	`db:"date_added" json:"date_added"`
	Referrer	NullString	`json:"referrer"`
	Removed		pgtype.Date	`json:"removed"`
	Started		pgtype.Date	`json:"started"`
	Priority	sql.NullBool	`json:"priority"`
	Genre		NullString	`json:"genre"`
}

type MediaQuery struct {
	Media
	Limit 			uint64
	OrderBy			NullString
	IncludeRemoved	sql.NullBool
}

type MediaStorage interface {
	ListMedia(int, int, string) ([]Media, error)
	ListRecentMedia(int, int) ([]Media, error)
	StartMedia(int) (*Media, error)
	SearchMedia(string, int) ([]Media, error)
	TopMedia(int) ([]Media, error)
	GetMediaByID(int) (*Media, error)
}

func (s serviceImpl) ListMedia(media_type int, limit int, genre string) ([]Media, error) {
	return s.db.ListMedia(media_type, limit, genre)
}

func (s serviceImpl) ListRecentMedia(media_type int, limit int) ([]Media, error) {
	return s.db.ListRecentMedia(media_type, limit)
}

func (s serviceImpl) StartMedia(id int) (*Media, error) {
	return s.db.StartMedia(id)
}

func (s serviceImpl) SearchMedia(title string, media_type int) ([]Media, error) {
	return s.db.SearchMedia(title, media_type)
}

func (s serviceImpl) TopMedia(media_type int) ([]Media, error) {
	return s.db.TopMedia(media_type)
}


func (s serviceImpl) GetMediaByID(id int) (*Media, error) {
	return s.db.GetMediaByID(id)
}