package repository

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/paul-ss/http-proxy/internal/database"
	"github.com/paul-ss/http-proxy/internal/domain"
	"log"
)

type Database struct {
	conn *pgx.Conn
}

func NewDatabase() *Database {
	conn, err := database.NewPostgresConn()
	if err != nil {
		log.Fatal("NewDatabase: " + err.Error())
	}

	return &Database{conn: conn}
}

func (d *Database) StoreRequest(req *domain.StoreRequest) (*domain.Request, error) {
	var id int32
	err := d.conn.QueryRow(context.Background(),
		"INSERT INTO requests (method, host, path, request) " +
		"VALUES ($1, $2, $3, $4) " +
		"RETURNING id ", req.Method, req.Host, req.Path, req.Req).Scan(&id)

	if err != nil {
		return nil, err
	}

	return &domain.Request{
		Id: id,
		Method: req.Method,
		Path: req.Path,
		Req: req.Req,
	}, nil
}

func (d *Database) GetShortRequests() ([]domain.RequestShort, error) {
	rows, err := d.conn.Query(context.Background(),
		"SELECT id, method, host, path "+
			"FROM requests "+
			"ORDER BY id ")

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reqs []domain.RequestShort
	for rows.Next() {
		req := domain.RequestShort{}
		if err := rows.Scan(&req.Id, &req.Method, &req.Host, &req.Path); err != nil {
			return nil, err
		}

		reqs = append(reqs, req)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return reqs, err
}

func (d *Database) GetRequestById(id int32) (*domain.Request, error) {
	req := domain.Request{}
	err := d.conn.QueryRow(context.Background(),
		"SELECT id, method, host, path, request "+
			"FROM requests " +
			"WHERE id = $1 ", id).
		Scan(&req.Id, &req.Method, &req.Host, &req.Path, &req.Req)

	if err != nil {
		return nil, err
	}

	return &req, nil
}

