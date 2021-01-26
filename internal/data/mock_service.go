package data

import (
	"context"
	"fmt"
	"github.com/Senso-Care/SensoApi/internal/config"
	"github.com/Senso-Care/SensoApi/internal/models"
	"math/rand"
	"strings"
	"time"
)

type MockInfluxService struct {
	bucket string
}

func NewMockService(configuration *config.DatabaseConfiguration) InfluxServicer {
	bucket := configuration.DbName + "/" + configuration.RetentionPolicy
	service := MockInfluxService{
		bucket: bucket,
	}
	return &service
}

func metrics() []string {
	return []string{"temperature", "pressure", "vox2", "humidity"}
}

func sensors() []string {
	return []string{"temperature-bathroom", "temperature-kitchen", "temperature-living-room", "pressure-bathroom", "pressure-kitchen", "pressure-living-room", "vox2-bathroom", "vox2-kitchen", "vox2-living-room", "humidity-bathroom", "humidity-kitchen", "humidity-living-room"}
}

func sensorsFromMetric(metric string) []string {
	var sensorL []string
	for _, sensor := range sensors() {
		if metric == sensor[0:len(metric)] {
			sensorL = append(sensorL, sensor)
		}
	}
	return sensorL
}

func getRandF(sensor string) func() float32 {
	switch sensor {
	case "temperature":
		return func() float32 {
			return float32(rand.Int31n(24 - 20)) + 20.0 + rand.Float32()
		}
	case "pressure":
		return func() float32 {
			return float32(rand.Int31n(100-90)) + 90.0 + rand.Float32()
		}
	case "vox2":
		return func() float32 {
			return float32(rand.Int31n(99-95)) + 95 + rand.Float32()
		}
	case "humidity":
		return func() float32 {
			return float32(rand.Int31n(40)) + rand.Float32()
		}
	default:
		return func() float32 {
			return float32(rand.Int31n(40)) + rand.Float32()
		}
	}
}


func genDataPoint(sensor string, range_ int) []models.DataPoint {
	// 172800
	var dataPoints []models.DataPoint
	sensor = strings.Split(sensor, "-")[0]
	randF := getRandF(sensor)
	fmt.Printf("%s sensor\n", sensor)
	var unit int64 = 24  // 15 minutes per day
	timestamp := time.Now().Unix() - (int64(range_) * 24 * 60 * 60)
	for i := 0; i < range_ * int(unit); i++ {
		timestamp = timestamp + (60 * 60)
		dataPoint := models.DataPoint{
			Date:  time.Unix(timestamp, 0).Format(time.RFC3339),
			Value: randF(),
		}
		dataPoints = append(dataPoints, dataPoint)
	}
	return dataPoints
}

func (service *MockInfluxService) GetSensors(ctx context.Context, timeRange string) ([]string, error) {
	return sensors(), nil
}

func (service *MockInfluxService) GetMetrics(ctx context.Context, timeRange string) ([]string, error) {
	return metrics(), nil
}

func (service *MockInfluxService) GetMetricsFromType(ctx context.Context, measurement, timeRange string) (*models.Metric, error) {
	sensors := sensorsFromMetric(measurement)
	var sensorData []models.SensorData
	metric := &models.Metric{
		Type: measurement,
		Sensors: sensorData,
	}
	for _, sensor := range sensors {
		metric.Sensors = append(metric.Sensors, models.SensorData{
			Name:   sensor,
			Series: genDataPoint(sensor, 30),
		})

	}
	return metric, nil
}

func (service *MockInfluxService) GetLastMetrics(ctx context.Context, measurement, timeRange string) ([]models.SensorData, error) {
	sensorsL := sensorsFromMetric(measurement)
	sensors := make([]models.SensorData, 0)
	for _, sensor := range sensorsL {
		val := models.SensorData{
			Name:   sensor,
			Series: genDataPoint(sensor, 1),
		}
		sensors = append(sensors, val)
	}
	return sensors, nil
}

func (service *MockInfluxService) GetMetricsFromSensor(ctx context.Context, sensor string, timeRange string) (*models.SensorData, error) {
	sensorData := &models.SensorData{
		Name:   sensor,
		Series: genDataPoint(sensor, 30),
	}

	return sensorData, nil
}

func (service *MockInfluxService) PostMetricsFromType(ctx context.Context, s string, point models.DataPoint) error {
	//nothing to do
	return nil
}

func (service *MockInfluxService) Close() {
	// do nothing
}
