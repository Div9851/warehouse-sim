package state

import (
	"math/rand"

	"github.com/Div9851/warehouse-sim/action"
	"github.com/Div9851/warehouse-sim/env"
	"github.com/Div9851/warehouse-sim/pos"
)

//nextItems 現在の状態, 各エージェントの行動, 環境設定, 乱数生成器を受け取り
//次の状態のAgentItems, PosItems, Success, アイテムが出現した場所, 各エージェントが得た報酬を返す
func nextItems(state *State, actions []int, env *env.Env, rnd *rand.Rand) ([]int, map[pos.Pos]int, []bool, *pos.Pos, []float64) {
	agentItems := make([]int, env.NumAgents)
	copy(agentItems, state.AgentItems)
	posItems := make(map[pos.Pos]int)
	for k, v := range state.PosItems {
		posItems[k] = v
	}
	success := make([]bool, env.NumAgents)
	agentPos := state.AgentPos //AgentPosは更新しないのでコピーしない
	rewards := make([]float64, env.NumAgents)
	for i, pos := range agentPos {
		switch actions[i] {
		case action.PICKUP:
			//まだアイテムを拾うことが出来, かつそこにアイテムがあるなら
			if agentItems[i] < env.MaxItems && posItems[pos] > 0 {
				for id := 0; id < env.NumAgents; id++ {
					rewards[id] += env.Reward
				}
				success[i] = true
				rewards[i] += env.DIYBonus
				agentItems[i]++
				posItems[pos]--
				if posItems[pos] == 0 {
					delete(posItems, pos)
				}
			}
		case action.CLEAR:
			//デポにいて, かつアイテムをもっているなら
			if pos == env.DepotPos && agentItems[i] > 0 {
				for id := 0; id < env.NumAgents; id++ {
					rewards[id] += env.Reward * float64(agentItems[i])
				}
				success[i] = true
				rewards[i] += env.DIYBonus * float64(agentItems[i])
				agentItems[i] = 0
			}
		}
	}
	//env.AppearProbで与えられる確率で, 新たなアイテムを出現させる
	var lastAppear *pos.Pos
	if rnd.Float64() < env.AppearProb {
		pos := env.AllPos[rnd.Intn(len(env.AllPos))]
		posItems[pos]++
		lastAppear = &pos
	}
	return agentItems, posItems, success, lastAppear, rewards
}
