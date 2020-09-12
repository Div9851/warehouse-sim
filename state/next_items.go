package state

import (
	"math/rand"

	"github.com/Div9851/warehouse-sim/action"
	"github.com/Div9851/warehouse-sim/env"
	"github.com/Div9851/warehouse-sim/pos"
)

//nextItems 現在の状態, 各エージェントの行動, 環境設定, 乱数生成器を受け取り
//次の状態のAgentItems, PosItems, Success, アイテムが出現した場所, 各エージェントが得た報酬を返す
func nextItems(state *State, actions []int, env *env.Env, rnd *rand.Rand) ([][]int, map[pos.Pos][]int, []bool, *pos.Pos, []float64) {
	agentItems := make([][]int, env.NumAgents)
	for i, items := range state.AgentItems {
		agentItems[i] = make([]int, len(items))
		copy(agentItems[i], items)
	}
	posItems := make(map[pos.Pos][]int)
	for k, v := range state.PosItems {
		posItems[k] = make([]int, len(v))
		copy(posItems[k], v)
	}
	success := make([]bool, env.NumAgents)
	agentPos := state.AgentPos //AgentPosは更新しないのでコピーしない
	rewards := make([]float64, env.NumAgents)
	for i, pos := range agentPos {
		switch actions[i] {
		case action.PICKUP:
			//まだアイテムを拾うことが出来, かつそこにアイテムがあるなら
			if len(agentItems[i]) < env.MaxItems && len(posItems[pos]) > 0 {
				for id := 0; id < env.NumAgents; id++ {
					rewards[id] += env.Reward
				}
				success[i] = true
				rewards[i] += env.DIYBonus
				agentItems[i] = append(agentItems[i], posItems[pos][0])
				posItems[pos] = posItems[pos][1:]
				if len(posItems[pos]) == 0 {
					delete(posItems, pos)
				}
			}
		case action.CLEAR:
			//デポにいて, かつアイテムをもっているなら
			if pos == env.DepotPos && len(agentItems[i]) > 0 {
				for id := 0; id < env.NumAgents; id++ {
					rewards[id] += env.Reward * float64(len(agentItems[i]))
				}
				success[i] = true
				rewards[i] += env.DIYBonus * float64(len(agentItems[i]))
				agentItems[i] = []int{}
			}
		}
	}
	//タイムリミットが来たアイテムを消す
	for i, items := range agentItems {
		for len(items) > 0 && items[0] == state.Turn {
			items = items[1:]
		}
		agentItems[i] = items
	}
	for pos, items := range posItems {
		for len(items) > 0 && items[0] == state.Turn {
			items = items[1:]
		}
		posItems[pos] = items
		if len(posItems[pos]) == 0 {
			delete(posItems, pos)
		}
	}
	//env.AppearProbで与えられる確率で, 新たなアイテムを出現させる
	var lastAppear *pos.Pos
	if rnd.Float64() < env.AppearProb {
		pos := env.AllPos[rnd.Intn(len(env.AllPos))]
		posItems[pos] = append(posItems[pos], state.Turn+env.TimeLimit)
		lastAppear = &pos
	}
	return agentItems, posItems, success, lastAppear, rewards
}
