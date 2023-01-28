package usecase

import (
	"context"

	"github.com/EtoNeJa00/GRPCService/internal/app/repository"

	"github.com/EtoNeJa00/GRPCService/internal/models"

	"github.com/google/uuid"
)

type Usecase interface {
	GetRecord(ctx context.Context, id uuid.UUID) (models.Record, error)
	SetRecord(ctx context.Context, record models.Record) (models.Record, error)
	DeleteRecord(ctx context.Context, id uuid.UUID) (models.Record, error)
}

type usecase struct {
	repo repository.Repository
}

func NewUsecase(r repository.Repository) Usecase {
	return usecase{
		repo: r,
	}
}

func (u usecase) GetRecord(ctx context.Context, id uuid.UUID) (models.Record, error) {
	return u.repo.GetRecord(ctx, id)
}

func (u usecase) SetRecord(ctx context.Context, record models.Record) (models.Record, error) {
	return u.repo.SetRecord(ctx, record)
}

func (u usecase) DeleteRecord(ctx context.Context, id uuid.UUID) (models.Record, error) {
	return u.repo.DeleteRecord(ctx, id)
}
