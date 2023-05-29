package store

import (
	"context"
	"github.com/georgysavva/scany/pgxscan"
	sq "github.com/Masterminds/squirrel"
	"github.com/antiharmonic/current_backend/current"
	"strconv"
	"log"
	"fmt"
	"strings"
)

var (
	MediaTable = "current_media"
)

// func (p postgres) CreateMedia(m current.Media) error {
// 	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
// 	stmnt, args, err := psql.Insert(MediaTable).
// 		Columns("id", "")
// }

func (p postgres) ListMediaWrapper(media_type string, limit string, genre string, orderby string, include_removed bool) ([]current.Media, error){
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	builder := psql.Select("*").From(MediaTable)
	if media_type != "" {
		log.Println("Adding media type", media_type)
		builder = builder.Where("type = ?", media_type)
	}
	if genre != "" {	
		builder = builder.Where("lower(genre) like ?", fmt.Sprint("%", strings.ToLower(genre), "%"))
	}
	nlimit, err := strconv.ParseUint(limit, 10, 64)
	if err == nil {
		builder = builder.Limit(nlimit)
	}
	if orderby != "" {
		builder = builder.OrderBy(orderby)
	}
	if include_removed == false {
		builder = builder.Where(sq.Eq{"removed": nil})
	}
	sql, args, err := builder.ToSql()
	log.Println(sql)
	if err != nil {
		return nil, err
	}

	var m  []current.Media
	err = pgxscan.Select(context.Background(), p.pool, &m, sql, args...)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (p postgres) ListMedia(media_type string, limit string, genre string) ([]current.Media, error) {
	return p.ListMediaWrapper(media_type, limit, genre, "title", true)
}

func (p postgres) ListRecentMedia(media_type string, limit string) ([]current.Media, error) {
	return p.ListMediaWrapper(media_type, limit, "", "id desc", false)
}