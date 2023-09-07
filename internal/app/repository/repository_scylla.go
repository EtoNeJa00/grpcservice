package repository

import (
	"context"
	"fmt"

	"GRPCService/internal/models"
	"GRPCService/internal/models/generated/scyllat"

	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"github.com/scylladb/gocqlx/v2"
	"github.com/scylladb/gocqlx/v2/table"
)

//go:generate go run github.com/scylladb/gocqlx/v2/cmd/schemagen --cluster=localhost --keyspace=grpcservice --output=../../models/generated/scyllat --pkgname=scyllat

const keySpace = "grpcservice"

type repositoryScylla struct {
	ses gocqlx.Session
	st  *table.Table
}

func (r repositoryScylla) GetRecord(ctx context.Context, id uuid.UUID) (models.Record, error) {
	recordS := scyllat.RecordsStruct{
		Id: id,
	}

	err := r.ses.Query(r.st.Get()).BindStruct(recordS).WithContext(ctx).GetRelease(&recordS)
	if err != nil {
		return models.Record{}, fmt.Errorf("scylla select: %w", err)
	}

	return models.Record{
		ID:   recordS.Id,
		Data: recordS.Record,
	}, nil
}

func (r repositoryScylla) SetRecord(ctx context.Context, record models.Record) (models.Record, error) {
	recordS := scyllat.RecordsStruct{
		Id:     record.ID,
		Record: record.Data,
	}

	if record.ID != uuid.Nil {
		return r.update(ctx, recordS)
	}

	return r.insert(ctx, recordS)
}

func (r repositoryScylla) insert(ctx context.Context, recordS scyllat.RecordsStruct) (models.Record, error) {
	recordS.Id = uuid.New()

	err := r.ses.Query(r.st.Insert()).BindStruct(&recordS).WithContext(ctx).ExecRelease()
	if err != nil {
		return models.Record{}, err
	}

	return models.Record{
		ID:   recordS.Id,
		Data: recordS.Record,
	}, nil
}

func (r repositoryScylla) update(ctx context.Context, recordS scyllat.RecordsStruct) (models.Record, error) {
	err := r.st.UpdateQuery(r.ses, r.st.Metadata().Columns[1:]...).BindStruct(&recordS).WithContext(ctx).ExecRelease()
	if err != nil {
		return models.Record{}, err
	}

	return models.Record{
		ID:   recordS.Id,
		Data: recordS.Record,
	}, err
}

func (r repositoryScylla) DeleteRecord(ctx context.Context, id uuid.UUID) (models.Record, error) {
	rec, err := r.GetRecord(ctx, id)
	if err != nil {
		return models.Record{}, err
	}

	recordS := scyllat.RecordsStruct{
		Id: rec.ID,
	}

	if err = r.ses.Query(r.st.Delete()).BindStruct(recordS).WithContext(ctx).ExecRelease(); err != nil {
		return models.Record{}, fmt.Errorf("scylla delete: %w", err)
	}

	return models.Record{
		ID:   rec.ID,
		Data: rec.Data,
	}, nil
}

func NewScyllaRepository(ctx context.Context, addr string) (Repository, error) {
	cluster := gocql.NewCluster(addr)
	cluster.Consistency = gocql.Quorum
	cluster.Keyspace = keySpace

	session, err := gocqlx.WrapSession(cluster.CreateSession())
	if err != nil {
		return nil, err
	}

	go func() {
		<-ctx.Done()

		session.Close()
	}()

	return repositoryScylla{ses: session, st: scyllat.Records}, nil
}
