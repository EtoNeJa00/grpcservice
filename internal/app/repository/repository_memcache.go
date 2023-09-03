package repository

import (
	"context"
	"errors"

	"GRPCService/internal/models"
	"GRPCService/internal/pkg/memcache"

	"github.com/google/uuid"
)

type repositoryMemcache struct {
	mc memcache.MemCache
}

func NewMemcacheRepository(addr string) (Repository, error) {
	mc, err := memcache.NewMemcache(addr)
	if err != nil {
		return nil, err
	}

	return &repositoryMemcache{
		mc: mc,
	}, nil
}

func (m *repositoryMemcache) GetRecord(_ context.Context, id uuid.UUID) (models.Record, error) {
	records, err := m.mc.Get(id)
	if err != nil {
		return models.Record{}, err
	}

	return models.Record{
		Id:   id,
		Data: string(records[0]),
	}, nil
}

func (m *repositoryMemcache) SetRecord(ctx context.Context, record models.Record) (models.Record, error) {
	if record.Id != uuid.Nil {
		res, err := m.GetRecord(ctx, record.Id)
		if errors.Is(err, memcache.ErrNotFound) {
			return m.createRecord(record)
		} else if err != nil {
			return models.Record{}, err
		}

		return m.updateRecord(res, record)
	}

	return m.createRecord(record)
}

func (m *repositoryMemcache) updateRecord(res models.Record, record models.Record) (models.Record, error) {
	err := m.mc.Delete(res.Id)
	if err != nil {
		return models.Record{}, err
	}

	return m.createRecord(record)
}

func (m *repositoryMemcache) createRecord(record models.Record) (models.Record, error) {
	id, err := m.mc.Set([]byte(record.Data))
	if err != nil {
		return models.Record{}, err
	}

	return models.Record{Id: id, Data: record.Data}, err
}

func (m *repositoryMemcache) DeleteRecord(ctx context.Context, id uuid.UUID) (models.Record, error) {
	res, err := m.GetRecord(ctx, id)
	if err != nil {
		return models.Record{}, err
	}

	err = m.mc.Delete(id)
	if err != nil {
		return models.Record{}, err
	}

	return res, nil
}
