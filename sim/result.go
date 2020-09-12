package sim

//Result シミュレーションの結果を表す構造体
type Result struct {
	TotalItems   int   `json:"total_items"`
	PickupCounts []int `json:"pickup_counts"`
	ClearCounts  []int `json:"clear_counts"`
}
