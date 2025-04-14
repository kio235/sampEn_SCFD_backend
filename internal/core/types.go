package core

type batteryRecord struct {
	voltage     []float64 // 电压数据
	cell        []string  // 电池名称
	cellCount   int       // 电池数量
	recordCount int       // 数据数量
}

type FaultInfo struct {
	FaultTime int `json:"fault_time"`
	FaultCell int `json:"fault_cell"`
}
