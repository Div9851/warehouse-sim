package state

import (
	"math/rand"

	"github.com/Div9851/warehouse-sim/env"
	"github.com/Div9851/warehouse-sim/pos"
)

//NextState 現在の状態, 各エージェントの行動, 環境設定, 乱数生成器を受け取り
//次の状態, アイテムが出現した場所, 各エージェントが得た報酬を返す
func NextState(state *State, actions []int, env *env.Env, rnd *rand.Rand) (*State, *pos.Pos, []float64) {
	agentItems, posItems, success, lastAppear, rewards := nextItems(state, actions, env, rnd)
	agentPos := nextPos(state, actions, env, rnd)
	return &State{Turn: state.Turn + 1, AgentItems: agentItems, AgentPos: agentPos, PosItems: posItems, RandomValues: state.RandomValues, Success: success}, lastAppear, rewards
}
