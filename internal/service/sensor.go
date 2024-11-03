package service

import (
	"context"
	v1 "sensor/api/sensor/v1"
	"sensor/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type SensorService struct {
	v1.UnimplementedManagementServer
	mgr *biz.SensorManager
	log *log.Helper
}

func NewSensorService(mgr *biz.SensorManager, logger log.Logger) *SensorService {
	return &SensorService{
		mgr: mgr,
	}
}

// SetThreshold
func (s *SensorService) SetThreshold(ctx context.Context, req *v1.SetThresholdRequest) (*v1.SetThresholdResponse, error) {

	err := s.mgr.SetThreshold(ctx, req.SensorId, req.Threshold)
	if err != nil {
		s.log.Errorf("Error setting threshold: %v", err)
		return nil, err
	}
	return &v1.SetThresholdResponse{Success: true}, nil
}

// Check alarm status
func (s *SensorService) CheckAlarm(ctx context.Context, req *v1.AlarmRequest) (*v1.AlarmResponse, error) {

	sensorID := req.SensorId
	newValue := req.NewValue

	triggered, message, err := s.mgr.CheckAlarm(ctx, sensorID, newValue)
	if err != nil {
		s.log.Errorf("Error triggering alarm: %v", err)
		return nil, err
	}

	return &v1.AlarmResponse{AlarmTriggered: triggered, Message: message}, nil
}
