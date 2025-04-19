package configuration

type MiningCondition struct {
	Source      string
	TargetValue float64
}

type MinerConfig struct {
	DeviceId  string
	Condition MiningCondition
}
