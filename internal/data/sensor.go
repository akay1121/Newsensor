package data

import (
	"context"
	"sensor/internal/biz"
	"time"

	"sensor/internal/biz/entity"
	"sensor/internal/ent"

	"github.com/go-kratos/kratos/v2/log"
)

// SensorDataProcessingRepo defines methods for sensor data processing.
type SensorDataProcessingRepo struct {
	db  *Data
	log *log.Helper
}

// NewSensorDataProcessingRepo creates a new instance of SensorDataProcessingRepo.
func NewSensorDataProcessingRepo(db *Data, logger log.Logger) *SensorDataProcessingRepo {
	return &SensorDataProcessingRepo{
		db:  db,
		log: log.NewHelper(logger),
	}
}

// SetSensorThreshold sets a threshold for triggering alerts for a sensor.
func (r *SensorDataProcessingRepo) SetSensorThreshold(ctx context.Context, sensorID string, threshold float64) error {
	_, err := r.db.Client.Sensor.Update().
		Where(sensor.IDEQ(sensorID)).
		SetThreshold(threshold).
		Save(ctx)
	if err != nil {
		r.log.Error("failed to set threshold for sensor:", err)
		return err
	}
	r.log.Info("set threshold for sensor:", sensorID)
	return nil
}

// GetSensorByID retrieves the sensor configuration by ID.
func (r *SensorDataProcessingRepo) GetSensorByID(ctx context.Context, sensorID string) (*entity.SensorDataProcessing, error) {
	s, err := r.data.Client.Sensor.Get(ctx, sensorID)
	if err != nil {
		r.log.Error("sensor not found:", err)
		return nil, err
	}
	return &entity.SensorDataProcessing{
		ID:            s.ID,
		TypeID:        s.TypeID,
		Description:   s.Desc,
		RuleID:        s.RuleID,
		Threshold:     s.Threshold,
		PreviousValue: s.PreviousValue,
		LastUpdate:    s.LastUpdate,
	}, nil
}

// InterpolateMissingData performs interpolation for missing sensor data based on historical records.
func (r *SensorDataProcessingRepo) InterpolateMissingData(ctx context.Context, sensorID string, start, end time.Time) ([]*entity.SensorValue, error) {
	// Example: Fetch recent values and interpolate missing data
	history, err := r.db.Client.SensorStatus.Query().
		Where(sensorstatus.IDEQ(sensorID), sensorstatus.TimestampGTE(start), sensorstatus.TimestampLTE(end)).
		Order(ent.Asc("timestamp")).
		All(ctx)
	if err != nil {
		r.log.Error("failed to retrieve sensor data for interpolation:", err)
		return nil, err
	}
	// Implement basic linear interpolation logic (for demonstration purposes)
	var interpolatedData []*entity.SensorValue
	for i := 1; i < len(history); i++ {
		if history[i].Timestamp.Sub(history[i-1].Timestamp) > time.Hour {
			// Simple linear interpolation between history[i-1] and history[i]
			interpolatedValue := (history[i].Value + history[i-1].Value) / 2
			missingData := &entity.SensorValue{
				Timestamp: history[i-1].Timestamp.Add(time.Hour),
				Value:     interpolatedValue,
				SensorID:  sensorID,
			}
			interpolatedData = append(interpolatedData, missingData)
		}
	}
	return interpolatedData, nil
}
func (r *SensorDataProcessingRepo) UpdateSensor(ctx context.Context, u *biz.User) error {
	return r.db.Client.User.Update().Where(user.IDEQ(u.Id)).Exec(ctx)
}
