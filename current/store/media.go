package store

import (
	"context"
	"github.com/georgysavva/scany/pgxscan"
	sq "github.com/Masterminds/squirrel"
	"github.com/antiharmonic/current_backend/current"
)

var (
	MediaTable = "current_media"
)

// func (p postgres) CreateMedia(m current.Media) error {
// 	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
// 	stmnt, args, err := psql.Insert(MediaTable).
// 		Columns("id", "")
// }

func (p postgres) ListMedia() ([]current.Media, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, _, err := psql.Select("*").From(MediaTable).ToSql()
	if err != nil {
		return nil, err
	}

	var m  []current.Media
	err = pgxscan.Select(context.Background(), p.pool, &m, sql)
	if err != nil {
		return nil, err
	}

	return m, nil
}