package entity

import "time"

// Sensor reprsents the entity of it
type Sensor struct {
	ID          string    // snsor ID
	TypeID      int64     // sensor type Id
	Description string    // description
	RuleID      int64     // rule ID
	Status      string    // sensor status
	LastUpdate  time.Time // last update time
	// I don't know whether need i to make the rule ,if don't just delete it
	Threshold     float64 // 设定的阈值
	PreviousValue float64 // 之前的传感器数据值

}
