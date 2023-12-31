package system

type SystemVersion struct {
	Version string
}

type SystemModel struct {
	Model string
}

type SystemMemory struct {
	Type  string
	Total float64
	Used  float64
	Free  float64
}

type SystemCPU struct {
	Type string
	Used float64
	Idle float64
}

type SystemValue struct {
	isSet bool
	Value float64
}
