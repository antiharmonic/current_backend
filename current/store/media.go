package store

import (
	"context"
	"github.com/georgysavva/scany/pgxscan"
	sq "github.com/Masterminds/squirrel"
	"github.com/antiharmonic/current_backend/current"
	"log"
	"fmt"
	"strings"
	"database/sql"
)

var (
	MediaTable = "current_media"
)

func (p postgres) GetMediaByID (id int) (*current.Media, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, err := psql.Select("*").From(MediaTable).Where(sq.Eq{"id": id}).ToSql()
	var m current.Media
	rows, err := p.pool.Query(context.Background(), sql, args...)
	if err != nil {
		return nil, err
	}
	err = pgxscan.ScanOne(&m, rows)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (p postgres) SearchMediaWrapper(m *current.MediaQuery) ([]current.Media, error){
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	builder := psql.Select("*").From(MediaTable)
	if m.Title.Valid && m.Title.String != "" {
		builder = builder.Where("lower(title) like ?", fmt.Sprint("%", strings.ToLower(m.Title.String), "%"))
	}
	if m.MediaType.Valid {
		builder = builder.Where("type = ?", m.MediaType.Int64)
	}
	if m.Genre.Valid && m.Genre.String != "" {	
		builder = builder.Where("lower(genre) like ?", fmt.Sprint("%", strings.ToLower(m.Genre.String), "%"))
	}
	//nlimit, err := strconv.ParseUint(limit, 10, 64)
	if m.Limit != 0 {
		builder = builder.Limit(m.Limit)
	}
	if m.OrderBy.Valid && m.OrderBy.String != "" {
		builder = builder.OrderBy(m.OrderBy.String)
	}
	if m.IncludeRemoved.Valid && m.IncludeRemoved.Bool == false {
		builder = builder.Where(sq.Eq{"removed": nil})
	} else if m.IncludeRemoved.Valid && m.IncludeRemoved.Bool == true {
		builder = builder.Where(sq.NotEq{"removed": nil})
	}
	sql, args, err := builder.ToSql()
	log.Println(sql)
	if err != nil {
		return nil, err
	}
	var results []current.Media
	err = pgxscan.Select(context.Background(), p.pool, &results, sql, args...)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (p postgres) ListMedia(media_type int, limit int, genre string) ([]current.Media, error) {
	nGenre, err := current.ParamToNullString(genre)
	if err != nil {
		return nil, err
	}
	nType, err := current.IntParamToNullInt(media_type, true)
	if err != nil {
		return nil, err
	}
	return p.SearchMediaWrapper(&current.MediaQuery{Media: current.Media{Genre: nGenre, MediaType: nType}, Limit: uint64(limit)})
}

func (p postgres) ListRecentMedia(media_type int, limit int) ([]current.Media, error) {
	nType, err := current.IntParamToNullInt(media_type, true)
	if err != nil {
		return nil, err
	}
	return p.SearchMediaWrapper(&current.MediaQuery{
		Media: current.Media{MediaType: nType}, 
		Limit: uint64(limit), 
		OrderBy: current.NullString{NullString: sql.NullString{String: "id desc", Valid:true}},
	} )//"", media_type, limit, "", "id desc", false)
}

func (p postgres) StartMedia(id int) (*current.Media, error) {
	stmnt := sq.Update(MediaTable).PlaceholderFormat(sq.Dollar).Where("id = ?", id).Set("started", sq.Expr("now()")).Suffix("returning *")
	sql, args, err := stmnt.ToSql()
	if err != nil {
		return nil, err
	}
	log.Println(sql, args)
	var m current.Media
	rows, err := p.pool.Query(context.Background(), sql, args...)
	if err != nil {
		return nil, err
	}
	err = pgxscan.ScanOne(&m, rows)
	//err = pgxscan.Select(context.Background(), p.pool, &m, sql, args...)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (p postgres) prioritizeMedia(id int, priority bool) (*current.Media, error) {
	stmnt := sq.Update(MediaTable).PlaceholderFormat(sq.Dollar).Where("id = ?", id).Set("priority", priority).Suffix("returning *")
	sql, args, err := stmnt.ToSql()
	if err != nil {
		return nil, err
	}
	log.Println(sql, args)
	var m current.Media
	rows, err := p.pool.Query(context.Background(), sql, args...)
	if err != nil {
		return nil, err
	}
	err = pgxscan.ScanOne(&m, rows)
	//err = pgxscan.Select(context.Background(), p.pool, &m, sql, args...)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (p postgres) UpgradeMedia(id int) (*current.Media, error) {
	return p.prioritizeMedia(id, true)
}

func (p postgres) DowngradeMedia(id int) (*current.Media, error) {
	return p.prioritizeMedia(id, false)
}

func (p postgres) SearchMedia(title string, media_type int) ([]current.Media, error) {
	nTitle, err := current.ParamToNullString(title)
	if err != nil {
		return nil, err
	}
	nType, err := current.IntParamToNullInt(media_type, true)
	if err != nil {
		return nil, err
	}
	return p.SearchMediaWrapper(&current.MediaQuery{Media: current.Media{Title: nTitle, MediaType: nType}})
}

func (p postgres) TopMedia(media_type int) ([]current.Media, error) {
	nType, err := current.IntParamToNullInt(media_type, true)
	if err != nil {
		return nil, err
	}

	query := current.MediaQuery{
		Media: current.Media{
			MediaType: nType,
			Priority: sql.NullBool{Bool: true, Valid: true},
			},
		IncludeRemoved: sql.NullBool{Bool: false, Valid: true},
		OrderBy: current.NullString{NullString: sql.NullString{String: "date_added, title", Valid: true}},
	}

	return p.SearchMediaWrapper(&query)
}