package repository

import (
	"context"

	"github.com/EtoNeJa00/GRPCService/internal/models"
	"github.com/EtoNeJa00/GRPCService/internal/pkg/inner_storage"

	"github.com/google/uuid"
)

func NewInnerStorageRepository(ctx context.Context) Repository {
	return &ISRRepository{iStorage: inner_storage.NewInnerStorage(ctx)}
}

type ISRRepository struct {
	iStorage inner_storage.InnerStorage
}

func (i *ISRRepository) GetRecord(_ context.Context, id uuid.UUID) (models.Record, error) {
	v, err := i.iStorage.Get(id)
	if err != nil {
		return models.Record{}, err
	}

	return models.Record{Id: id, Data: v}, nil
}

func (i *ISRRepository) SetRecord(_ context.Context, record models.Record) (models.Record, error) {
	if record.Id == uuid.Nil {

		id, v := i.iStorage.Create(record.Data)

		return models.Record{Id: id, Data: v}, nil
	}

	v, err := i.iStorage.Update(record.Id, record.Data)
	if err != nil {
		return models.Record{}, err
	}

	return models.Record{Id: record.Id, Data: v}, nil
}

func (i *ISRRepository) DeleteRecord(_ context.Context, id uuid.UUID) (models.Record, error) {
	v, err := i.iStorage.Delete(id)
	if err != nil {
		return models.Record{}, err
	}

	return models.Record{Id: id, Data: v}, nil
}
