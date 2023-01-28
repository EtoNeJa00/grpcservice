package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/EtoNeJa00/GRPCService/internal/models"
)

type Repository interface {
	GetRecord(ctx context.Context, id uuid.UUID) (models.Record, error)
	SetRecord(ctx context.Context, record models.Record) (models.Record, error)
	DeleteRecord(ctx context.Context, id uuid.UUID) (models.Record, error)
}
