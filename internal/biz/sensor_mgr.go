package biz

import (
	"context"
	"sensor/internal/biz/entity"
	"time"
)

type SensorDataProcessing = entity.Sensor

type SensorRepository interface {
	GetSensorByID(ctx context.Context, id string) (*entity.Sensor, error)
	UpdateSensorStatus(ctx context.Context, id string, status string) error
	//if already have the rule just delete the relative of SetSensorThreshold
	SetSensorThreshold(ctx context.Context, id string, threshold float64) error
	UpdateSensor(ctx context.Context, user *User) error
}

type SensorManager struct {
	repo SensorRepository
}

// init SensorProcessingManager
func NewSensorProcessingManager(repo SensorRepository) *SensorManager {
	return &SensorManager{repo: repo}
}

// check whether have alarm
func (m *SensorManager) CheckAlarm(ctx context.Context, sensorID string, newValue float64) (bool, string, error) {
	sensor, err := m.repo.GetSensorByID(ctx, sensorID)
	if err != nil {
		return false, "", err
	}

	//if you have rule replace it to the below
	//if sensor.Status == "alarm" {
	//	return true, "Sensor is in alarm state!", nil
	//}
	//return false, "Sensor is operating normally.", nil
	//check whether excedding the threshold
	if abs(newValue-sensor.PreviousValue) > sensor.Threshold {
		return true, "Sensor data abnormal!", nil
	}

	// update the value
	sensor.PreviousValue = newValue
	sensor.LastUpdate = time.Now()
	if err := m.repo.UpdateSensor(ctx, sensor); err != nil {
		return false, "", err
	}

	return false, "Sensor data normal.", nil
}

// get the absolute number to check whether reach the threshold,you can change it by golang's own abs method
func abs(value float64) float64 {
	if value < 0 {
		return -value
	}
	return value
}

// interpolate data when the data missing
func (m *SensorManager) InterpolateData(ctx context.Context, sensorID string) (bool, string, error) {
	sensor, err := m.repo.GetSensorByID(ctx, sensorID)
	if err != nil {
		return false, "", err
	}

	//judge whether need to interpolate data
	if time.Since(sensor.LastUpdate) > 10*time.Minute {
		return true, "Interpolation data applied for missing entries.", nil
	}
	return false, "No interpolation needed.", nil
}
func (m *SensorManager) SetThreshold(ctx context.Context, sensorID string, threshold float64) error {
	return m.repo.SetSensorThreshold(ctx, sensorID, threshold)
}
