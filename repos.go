package main

import (
	"context"
	"database/sql"
	"echoserver/echolet"
	"time"
)

type Product struct {
	Id          int64  `json:"id,omitempty"`
	Title       string `json:"title,omitempty"`
	Price       int    `json:"price,omitempty"`
	DateCreated int64  `json:"date_created,omitempty"`
}

type Repos interface {
	list(ctx context.Context) ([]Product, error)
	create(ctx context.Context, p *Product) error
	updateById(ctx context.Context, id int, p *Product) (*Product, error)
	deleteById(ctx context.Context, id int) error
}

type reposInteranl struct {
	*sql.DB
}

func CreateRepos() Repos {
	return &reposInteranl{echolet.ConnectMySQL()}
}

func (repos *reposInteranl) list(ctx context.Context) ([]Product, error) {
	rows, err := repos.QueryContext(ctx, "select id,title,price,date_created from product order by id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := make([]Product, 0, 10)
	for rows.Next() {
		p := Product{}
		err := rows.Scan(&p.Id, &p.Title, &p.Price, &p.DateCreated)
		if err != nil {
			return list, nil
		}
		list = append(list, p)
	}
	return list, nil
}
func (repos *reposInteranl) create(ctx context.Context, p *Product) error {
	stmt, err := repos.PrepareContext(ctx, "insert into product(title,price,date_created) values(?,?,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	p.DateCreated = time.Now().UnixNano() / int64(time.Millisecond)
	rs, err := stmt.Exec(p.Title, p.Price, p.DateCreated)
	if err != nil {
		return err
	}
	lastId, err := rs.LastInsertId()
	if err != nil {
		return err
	}
	p.Id = lastId
	return nil
}

func (repos *reposInteranl) updateById(ctx context.Context, id int, p *Product) (*Product, error) {
	stmt, err := repos.PrepareContext(ctx, "update product set title=?, price=?, date_created=? where id=?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	p.DateCreated = time.Now().UnixNano() / int64(time.Millisecond)
	_, err = stmt.ExecContext(ctx, p.Title, p.Price, p.DateCreated, id)
	if err != nil {
		return nil, err
	}
	return repos.selectById(ctx, id)
}

func (repos *reposInteranl) deleteById(ctx context.Context, id int) error {
	stmt, err := repos.PrepareContext(ctx, "delete from product where id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, _ = stmt.ExecContext(ctx, id)
	return nil
}

func (repos *reposInteranl) selectById(ctx context.Context, id int) (*Product, error) {
	p := new(Product)
	err := repos.QueryRowContext(ctx, "select id,title,price,date_created from product where id=?", id).Scan(&p.Id, &p.Title, &p.Price, &p.DateCreated)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return p, nil
}
