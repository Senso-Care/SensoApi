package data

import (
	"context"
	"github.com/Senso-Care/SensoApi/internal/models"
)

type InfluxServicer interface {
	GetMetrics(context.Context, string) ([]string, error)
	GetMetricsFromSensor(context.Context, string, string) (*models.SensorData, error)
	GetMetricsFromType(context.Context, string, string) (*models.Metric, error)
	GetSensors(context.Context, string) ([]string, error)
	GetLastMetrics(context.Context, string, string) ([]models.SensorData, error)
	PostMetricsFromType(context.Context, string, models.DataPoint) error
	Close()
}