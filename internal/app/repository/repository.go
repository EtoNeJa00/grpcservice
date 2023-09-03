package repository

import (
	"context"

	"GRPCService/internal/models"

	"github.com/google/uuid"
)

type Repository interface {
	GetRecord(ctx context.Context, id uuid.UUID) (models.Record, error)
	SetRecord(ctx context.Context, record models.Record) (models.Record, error)
	DeleteRecord(ctx context.Context, id uuid.UUID) (models.Record, error)
}
