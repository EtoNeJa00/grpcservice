package repository

import (
	"context"

	"GRPCService/internal/models"
	"GRPCService/internal/pkg/innerstorage"

	"github.com/google/uuid"
)

func NewInnerStorageRepository() Repository {
	return &ISRRepository{iStorage: innerstorage.NewInnerStorage()}
}

type ISRRepository struct {
	iStorage innerstorage.InnerStorage
}

func (i *ISRRepository) GetRecord(_ context.Context, id uuid.UUID) (models.Record, error) {
	v, err := i.iStorage.Get(id)
	if err != nil {
		return models.Record{}, err
	}

	return models.Record{ID: id, Data: v}, nil
}

func (i *ISRRepository) SetRecord(_ context.Context, record models.Record) (models.Record, error) {
	if record.ID == uuid.Nil {
		id, v := i.iStorage.Create(record.Data)

		return models.Record{ID: id, Data: v}, nil
	}

	v, err := i.iStorage.Update(record.ID, record.Data)
	if err != nil {
		return models.Record{}, err
	}

	return models.Record{ID: record.ID, Data: v}, nil
}

func (i *ISRRepository) DeleteRecord(_ context.Context, id uuid.UUID) (models.Record, error) {
	v, err := i.iStorage.Delete(id)
	if err != nil {
		return models.Record{}, err
	}

	return models.Record{ID: id, Data: v}, nil
}
