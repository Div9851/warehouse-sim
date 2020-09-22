package state

import (
	"math/rand"

	"github.com/Div9851/warehouse-sim/env"
	"github.com/Div9851/warehouse-sim/pos"
)

//NextState 現在の状態, 各エージェントの行動のリスト, 環境設定, 乱数生成器を受け取り
//次の状態, エージェントが取った行動, アイテムが出現した場所, 各エージェントが得た報酬を返す
func NextState(state *State, actionLists [][]int, env *env.Env, rnd *rand.Rand) (*State, []int, *pos.Pos, []float64) {
	actions := make([]int, env.NumAgents)
	rank := make([]int, env.NumAgents)
	for id := 0; id < env.NumAgents; id++ {
		actions[id] = actionLists[id][0]
		rank[id] = 1
	}
	count := 0
	for {
		count++
		agentItems, posItems, success, rewards := nextItems(state, actions, env)
		nxtPos := nextPos(state, actions, env)
		update := false
		for id := 0; id < env.NumAgents; id++ {
			if nxtPos[id] != pos.NextPos(state.AgentPos[id], actions[id], env.MapData) && rank[id] < len(actionLists[id]) {
				actions[id] = actionLists[id][rank[id]]
				rank[id]++
				update = true
			}
		}
		if !update || count >= env.MaxLen {
			var lastAppear *pos.Pos
			//与えられた確率で新しいアイテムを出現させる
			if rnd.Float64() < env.AppearProb {
				lastAppear = &env.AllPos[rnd.Intn(len(env.AllPos))]
				posItems[*lastAppear]++
			}
			return &State{Turn: state.Turn + 1, AgentItems: agentItems, AgentPos: nxtPos, PosItems: posItems, RandomValues: state.RandomValues, Success: success}, actions, lastAppear, rewards
		}
	}
}
