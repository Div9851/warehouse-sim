package state

import (
	"github.com/Div9851/warehouse-sim/pos"
)

//State 状態を表す構造体
type State struct {
	Turn         int
	AgentItems   []int
	AgentPos     []pos.Pos
	PosItems     map[pos.Pos]int
	RandomValues map[pos.Pos]float64 //PosItemsのキーの順序を固定する
	Success      []bool              //PICKUPまたはCLEARを実行できた場合のみ真
}

//New 新しいStateへのポインタを返す
func New(turn int, agentItems []int, agentPos []pos.Pos, posItems map[pos.Pos]int, randomValues map[pos.Pos]float64, success []bool) *State {
	return &State{Turn: turn, AgentItems: agentItems, AgentPos: agentPos, PosItems: posItems, RandomValues: randomValues, Success: success}
}
