package repository

import (
	"context"
	"errors"
	"fmt"

	"GRPCService/internal/models"

	"github.com/gocql/gocql"
	"github.com/google/uuid"
)

const keySpace = "grpcservice"

var ErrNotFound = errors.New("not found")

type repositoryScylla struct {
	ses *gocql.Session
}

func NewScyllaRepository(ctx context.Context, addr string) (Repository, error) {
	cluster := gocql.NewCluster(addr)
	cluster.Consistency = gocql.Quorum
	cluster.Keyspace = keySpace

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}

	go func() {
		<-ctx.Done()

		session.Close()
	}()

	return repositoryScylla{ses: session}, nil
}

func (r repositoryScylla) GetRecord(ctx context.Context, id uuid.UUID) (rec models.Record, err error) {
	it := r.ses.Query("SELECT record FROM records WHERE id=?;", id.String()).WithContext(ctx).Iter()

	defer func() {
		errI := it.Close()
		if errI != nil {
			err = errI
		}
	}()

	var record *string

	it.Scan(&record)

	if record == nil {
		return models.Record{}, ErrNotFound
	}

	return models.Record{
		ID:   id,
		Data: *record,
	}, nil
}

func (r repositoryScylla) SetRecord(ctx context.Context, record models.Record) (models.Record, error) {
	if record.ID != uuid.Nil {
		return record, r.update(ctx, record)
	}

	return r.insert(ctx, record)
}

func (r repositoryScylla) insert(ctx context.Context, record models.Record) (models.Record, error) {
	record.ID = uuid.New()

	if err := r.ses.Query(`INSERT INTO records (id, record) VALUES (?, ?);`, record.ID.String(), record.Data).WithContext(ctx).Exec(); err != nil {
		return models.Record{}, err
	}

	return record, nil
}

func (r repositoryScylla) update(ctx context.Context, record models.Record) error {
	if err := r.ses.Query(`UPDATE records SET record=? WHERE id=?`, record.Data, record.ID.String()).WithContext(ctx).Exec(); err != nil {
		return err
	}

	return nil
}

func (r repositoryScylla) DeleteRecord(ctx context.Context, id uuid.UUID) (models.Record, error) {
	rec, err := r.GetRecord(ctx, id)
	if err != nil {
		return models.Record{}, err
	}

	if err = r.ses.Query(`DELETE record FROM records WHERE id = ?`, id.String()).WithContext(ctx).Exec(); err != nil {
		return models.Record{}, fmt.Errorf("scylla delete: %w", err)
	}

	return models.Record{
		ID:   rec.ID,
		Data: rec.Data,
	}, nil
}
