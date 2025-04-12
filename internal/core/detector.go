package core

type Detector interface {
	Check(record []BatteryRecord) []FaultInfo
}

type ThresholdDetector struct {
	voltageThreshold float64
	tempThreshold    float64
}

func NewThresholdDetector(voltage, temp float64) *ThresholdDetector {
	return &ThresholdDetector{
		voltageThreshold: voltage,
		tempThreshold:    temp,
	}
}
