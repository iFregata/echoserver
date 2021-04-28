package main

import (
	"context"
	"database/sql"
	"echoserver/dbc"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Product struct {
	Id          int64  `json:"id,omitempty"`
	Title       string `json:"title,omitempty"`
	Price       int    `json:"price,omitempty"`
	DateCreated int64  `json:"date_created,omitempty"`
}

type Repos interface {
	list(ctx context.Context) ([]Product, error)
	findById(ctx context.Context, id int) (*Product, error)
	findByIdFuture(ctx context.Context, id int) <-chan Result
	create(ctx context.Context, p *Product) error
	updateById(ctx context.Context, id int, p *Product) (*Product, error)
	deleteById(ctx context.Context, id int) error
}

type reposInteranl struct {
	*pgxpool.Pool
}

func NewRepos() Repos {
	return &reposInteranl{dbc.ConnectPostgres()}
}

type Result struct {
	Error   error
	Product *Product
}

func (repos *reposInteranl) findByIdFuture(ctx context.Context, id int) <-chan Result {
	future := make(chan Result)
	go func() {
		defer close(future)
		select {
		case <-ctx.Done():
			return
		default:
		}
		p, err := repos.findById(ctx, id)
		future <- Result{Product: p, Error: err}
	}()
	return future
}

func (repos *reposInteranl) findById(ctx context.Context, id int) (*Product, error) {
	return repos.selectById(ctx, id)
}

func (repos *reposInteranl) list(ctx context.Context) ([]Product, error) {
	rows, err := repos.Query(ctx,
		"select id,title,price,date_created from product order by id")
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
	p.DateCreated = time.Now().UnixNano() / int64(time.Millisecond)
	err := repos.QueryRow(ctx,
		"insert into product(title,price,date_created) values($1,$2,$3) returning id",
		p.Title, p.Price, p.DateCreated).Scan(&p.Id)
	if err != nil {
		return err
	}
	return nil
}

func (repos *reposInteranl) updateById(ctx context.Context, id int, p *Product) (*Product, error) {
	p.DateCreated = time.Now().UnixNano() / int64(time.Millisecond)
	_, err := repos.Exec(ctx,
		"update product set title=$1, price=$2, date_created=$3 where id=$4",
		p.Title, p.Price, p.DateCreated, id)
	if err != nil {
		return nil, err
	}
	return repos.selectById(ctx, id)
}

func (repos *reposInteranl) deleteById(ctx context.Context, id int) error {
	_, err := repos.Exec(ctx,
		"delete from product where id=$1", id)
	if err != nil {
		return err
	}
	return nil
}

func (repos *reposInteranl) selectById(ctx context.Context, id int) (*Product, error) {
	p := new(Product)
	err := repos.QueryRow(ctx,
		"select id,title,price,date_created from product where id=$1",
		id).Scan(&p.Id, &p.Title, &p.Price, &p.DateCreated)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return p, nil
}
