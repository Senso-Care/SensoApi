package data

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/Senso-Care/SensoApi/internal/config"
	"github.com/Senso-Care/SensoApi/internal/models"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	log "github.com/sirupsen/logrus"
)

type InfluxService struct {
	InfluxClient influxdb2.Client
	WriteApi     api.WriteAPIBlocking
	QueryApi     api.QueryAPI
	bucket       string
}

type TimePointMap map[string]map[time.Time]*models.DataPoint

// NewDefaultApiService creates a default api service
func NewInfluxService(configuration *config.DatabaseConfiguration) InfluxServicer {
	client := influxdb2.NewClient(configuration.ConnectionUri, fmt.Sprintf("%s:%s", configuration.Username, configuration.Password))
	bucket := configuration.DbName + "/" + configuration.RetentionPolicy
	writeAPI := client.WriteAPIBlocking("", bucket)
	queryAPI := client.QueryAPI("")
	service := InfluxService{
		InfluxClient: client,
		WriteApi:     writeAPI,
		QueryApi:     queryAPI,
		bucket:       bucket,
	}
	return &service
}

func (service *InfluxService) GetSensors(ctx context.Context, timeRange string) ([]string, error) {
	result, err := service.QueryApi.Query(ctx, fmt.Sprintf(`from(bucket: "%s")|>range(start:-%s)|>group(columns:["sensor"])|>distinct(column:"sensor")`, service.bucket, timeRange))
	sensors := make([]string, 0)
	if err == nil {
		for result.Next() {
			sensor, ok := result.Record().ValueByKey("sensor").(string)
			if !ok {
				continue
			}
			sensors = append(sensors, sensor)
		}
		if result.Err() != nil {
			log.Printf("Query error: %s\n", result.Err().Error())
			return nil, err
		}
	} else {
		log.Error(err)
		return nil, err
	}
	return sensors, nil
}

func (service *InfluxService) GetMetrics(ctx context.Context, timeRange string) ([]string, error) {
	result, err := service.QueryApi.Query(ctx, fmt.Sprintf(`from(bucket: "%s")|>range(start:-%s)|>group(columns:["_measurement"])|>distinct(column:"_measurement")`, service.bucket, timeRange))
	measurements := make([]string, 0)
	if err == nil {
		for result.Next() {
			measurements = append(measurements, result.Record().Measurement())
		}
		if result.Err() != nil {
			log.Printf("Query error: %s\n", result.Err().Error())
			return nil, err
		}
	} else {
		log.Error(err)
		return nil, err
	}
	return measurements, nil
}

func (service *InfluxService) GetMetricsFromType(ctx context.Context, measurement, timeRange string) (*models.Metric, error) {
	result, err := service.QueryApi.Query(ctx, fmt.Sprintf(`from(bucket: "%s")|>range(start: -%s)|>group(columns:["sensor"])|>filter(fn: (r) => r._measurement == "%s")`, service.bucket, timeRange, measurement))
	data := make(TimePointMap)
	if err == nil {
		processQueryResultBySensor(result, data)
		if result.Err() != nil {
			log.Printf("Query error: %s\n", result.Err().Error())
			return nil, err
		}
	} else {
		log.Error(err)
		return nil, err
	}

	sensors := createSensorListFromMap(data)
	metric := &models.Metric{
		Type:    measurement,
		Sensors: sensors,
	}
	return metric, nil
}

func (service *InfluxService) GetLastMetrics(ctx context.Context, measurement, timeRange string) ([]models.SensorData, error) {
	result, err := service.QueryApi.Query(ctx, fmt.Sprintf(`from(bucket: "%s")|>range(start: -%s)|>group(columns:["sensor"])|>filter(fn: (r) => r._measurement == "%s")|>sort(columns:["time"], desc: true)|>limit(n:1)`, service.bucket, timeRange, measurement))
	data := make(TimePointMap)
	if err == nil {
		processQueryResultBySensor(result, data)
		if result.Err() != nil {
			log.Printf("Query error: %s\n", result.Err().Error())
			return nil, err
		}
	} else {
		log.Error(err)
		return nil, err
	}

	sensors := createSensorListFromMap(data)
	return sensors, nil
}

func (service *InfluxService) GetMetricsFromSensor(ctx context.Context, sensor string, timeRange string) (*models.SensorData, error) {
	result, err := service.QueryApi.Query(ctx, fmt.Sprintf(`from(bucket: "%s")|>range(start: -%s)|>group(columns:["sensor"])|>filter(fn: (r) => r.sensor == "%s")`, service.bucket, timeRange, sensor))
	dataMap := make(map[time.Time]*models.DataPoint)
	if err == nil {
		for result.Next() {
			processQueryResult(result, dataMap)
		}
		if result.Err() != nil {
			log.Printf("Query error: %s\n", result.Err().Error())
			return nil, err
		}
	} else {
		log.Error(err)
		return nil, err
	}
	data := make([]models.DataPoint, 0)
	for _, value := range dataMap {
		data = append(data, *value)
	}
	sensorData := &models.SensorData{
		Name:   sensor,
		Series: data,
	}

	return sensorData, nil
}

func (service *InfluxService) PostMetricsFromType(ctx context.Context, type_ string, point models.DataPoint) error {
	date, err := time.Parse(time.RFC3339, point.Date)
	if err != nil {
		return err
	}
	value := influxdb2.NewPointWithMeasurement(type_).
		AddTag("sensor", type_+"-web").
		AddField("v", fmt.Sprintf("%f", point.Value)).
		SetTime(date)
	if len(point.Info) > 0 {
		value = value.AddField("info", point.Info)
	}
	if err := service.WriteApi.WritePoint(ctx, value); err != nil {
		log.Printf("Cannot insert point: %s\n", err)
		return err
	}

	return nil
}

func (service *InfluxService) Close() {
	service.InfluxClient.Close()
}

func dataPointExists(dataMap map[time.Time]*models.DataPoint, date time.Time) *models.DataPoint {
	ptr := dataMap[date]
	if ptr == nil {
		dataMap[date] = &models.DataPoint{
			Date: date.Format(time.RFC3339),
		}
		ptr = dataMap[date]
	}
	return ptr
}

func processQueryResultBySensor(result *api.QueryTableResult, dataMap TimePointMap) {
	for result.Next() {
		sensor, ok := result.Record().ValueByKey("sensor").(string)
		if !ok {
			log.Println("Bad type for sensor")
			continue
		}
		if dataMap[sensor] == nil {
			dataMap[sensor] = make(map[time.Time]*models.DataPoint)
		}
		processQueryResult(result, dataMap[sensor])
	}
}

func processQueryResult(result *api.QueryTableResult, dataMap map[time.Time]*models.DataPoint) {
	date := result.Record().Time()
	dataPoint := dataPointExists(dataMap, date)

	var err_ error
	switch result.Record().Field() {
	case "v":
		value, ok := result.Record().Value().(string)
		if !ok {
			err_ = errors.New("bad type for value")
			break
		}
		fValue, err := strconv.ParseFloat(value, 32)
		if err != nil {
			err_ = err
			break
		}
		dataPoint.Value = float32(fValue)
	case "info":
		value, ok := result.Record().Value().(string)
		if !ok {
			err_ = errors.New("bad type for info")
			break
		}
		dataPoint.Info = value
	}

	if err_ != nil {
		log.Println(err_)
	}
}

func createSensorListFromMap(dataMap TimePointMap) []models.SensorData {
	sensors := make([]models.SensorData, 0)
	for sensor, sensorMap := range dataMap {
		var dataPoints []models.DataPoint
		for _, dataPoint := range sensorMap {
			dataPoints = append(dataPoints, *dataPoint)
		}
		val := models.SensorData{
			Name:   sensor,
			Series: dataPoints,
		}
		sensors = append(sensors, val)
	}
	return sensors
}
