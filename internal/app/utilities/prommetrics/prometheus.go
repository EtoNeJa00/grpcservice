package prommetrics

import (
	"context"
	"time"

	"GRPCService/internal/app/usecase"
	"GRPCService/internal/models"

	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
)

type promUsecase struct {
	requestLatency *prometheus.HistogramVec
	errorCount     *prometheus.CounterVec

	uc          usecase.Usecase
	storageType string
}

func (p promUsecase) GetRecord(ctx context.Context, id uuid.UUID) (models.Record, error) {
	start := time.Now()

	r, err := p.uc.GetRecord(ctx, id)
	if err != nil {
		p.errorCount.With(prometheus.Labels{"method": "get", "storage_type": p.storageType}).Inc()
	}

	p.requestLatency.With(prometheus.Labels{"method": "get", "storage_type": p.storageType}).Observe(p.getLatency(start))

	return r, err
}

func (p promUsecase) SetRecord(ctx context.Context, record models.Record) (models.Record, error) {
	start := time.Now()

	r, err := p.uc.SetRecord(ctx, record)
	if err != nil {
		p.errorCount.With(prometheus.Labels{"method": "set", "storage_type": p.storageType}).Inc()
	}

	p.requestLatency.With(prometheus.Labels{"method": "set", "storage_type": p.storageType}).Observe(p.getLatency(start))

	return r, err
}

func (p promUsecase) DeleteRecord(ctx context.Context, id uuid.UUID) (models.Record, error) {
	start := time.Now()

	r, err := p.uc.DeleteRecord(ctx, id)
	if err != nil {
		p.errorCount.With(prometheus.Labels{"method": "delete", "storage_type": p.storageType}).Inc()
	}

	p.requestLatency.With(prometheus.Labels{"method": "delete", "storage_type": p.storageType}).Observe(p.getLatency(start))

	return r, err
}

func (p promUsecase) getLatency(start time.Time) float64 {
	return float64(time.Since(start).Nanoseconds())
}

type Metrics struct {
	RequestLatency *prometheus.HistogramVec
	ErrorCount     *prometheus.CounterVec
}

func NewPrometheusMiddleware(uc usecase.Usecase, metrics *Metrics, storageType string) (usecase.Usecase, error) {
	return promUsecase{
		requestLatency: metrics.RequestLatency,
		errorCount:     metrics.ErrorCount,
		uc:             uc,
		storageType:    storageType,
	}, nil
}

func CreateMetrics(reg *prometheus.Registry) (*Metrics, error) {
	m := &Metrics{}

	m.RequestLatency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "processing_time",
		Help: "time of processing request",
	}, []string{"method", "storage_type"})

	m.ErrorCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "error_count",
		Help: "count of errors",
	}, []string{"method", "storage_type"})

	err := reg.Register(m.RequestLatency)
	if err != nil {
		return nil, err
	}

	err = reg.Register(m.ErrorCount)
	if err != nil {
		return nil, err
	}

	return m, err
}
