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
	ID			int64	`db:"id" json:"id"`
	Title		string	`json:"title"`
	MediaType	int		`db:"type" json:"type"`
	Weight		float32	`json:"weight"`
	DateAdded	pgtype.Date	`db:"date_added" json:"date_added"`
	Referrer	sql.NullString	`json:"referrer"`
	Removed		pgtype.Date	`json:"removed"`
	Started		pgtype.Date	`json:"started"`
	Priority	bool	`json:"priority"`
	Genre		sql.NullString	`json:"genre"`
}

type MediaStorage interface {
	ListMedia(string, string, string) ([]Media, error)
	ListRecentMedia(string, string) ([]Media, error)
	StartMedia(int) (*Media, error)
}

func (s serviceImpl) ListMedia(media_type string, limit string, genre string) ([]Media, error) {
	return s.db.ListMedia(media_type, limit, genre)
}

func (s serviceImpl) ListRecentMedia(media_type string, limit string) ([]Media, error) {
	return s.db.ListRecentMedia(media_type, limit)
}

func (s serviceImpl) StartMedia(id int) (*Media, error) {
	return s.db.StartMedia(id)
}