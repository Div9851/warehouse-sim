package sim

//Result シミュレーションの結果を表す構造体
type Result struct {
	EnvName      string  `json:"env_name"`
	TotalItems   int     `json:"total_items"`
	PickupCounts []int   `json:"pickup_counts"`
	ClearCounts  []int   `json:"clear_counts"`
	Seed         int64   `json:"seed"`
	ProcessTime  float64 `json:"process_time"`
}
