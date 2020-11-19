package state

import (
	"math/rand"

	"github.com/Div9851/warehouse-sim/env"
	"github.com/Div9851/warehouse-sim/pos"
)

//NextState 現在の状態, 各エージェントの行動, 環境設定, 乱数生成器を受け取り
//次の状態, エージェントが取った行動, アイテムが出現した場所, 各エージェントが得た報酬を返す
func NextState(state *State, actions []int, env *env.Env, rnd *rand.Rand) (*State, []int, *pos.Pos, []float64) {
	return NextStateOpt(state, actions, env, rnd, -1, 0.0)
}

//NextStateOpt あるエージェントを優先するようなNextState
func NextStateOpt(state *State, actions []int, env *env.Env, rnd *rand.Rand, favoredID int, opt float64) (*State, []int, *pos.Pos, []float64) {
	agentItems, posItems, successItems, rewards := nextItems(state, actions, env)
	nxtPos, successPos := nextPosOpt(state, actions, env, rnd, favoredID, opt)
	success := make([]bool, env.NumAgents)
	for i := 0; i < env.NumAgents; i++ {
		success[i] = successItems[i] || successPos[i]
	}
	var lastAppear *pos.Pos
	//与えられた確率で新しいアイテムを出現させる
	if rnd.Float64() < env.AppearProb {
		lastAppear = &env.AllPos[rnd.Intn(len(env.AllPos))]
		posItems[*lastAppear]++
	}
	return &State{Turn: state.Turn + 1, AgentItems: agentItems, AgentPos: nxtPos, PosItems: posItems, RandomValues: state.RandomValues, Success: success}, actions, lastAppear, rewards
}
