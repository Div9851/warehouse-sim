package state

import (
	"math/rand"

	"github.com/Div9851/warehouse-sim/env"
	"github.com/Div9851/warehouse-sim/pos"
)

//nextPos 現在の状態, 各エージェントの行動, 環境設定, 乱数生成器を受け取って
//次の状態のAgentPosを返す
func nextPos(state *State, actions []int, env *env.Env, rnd *rand.Rand) []pos.Pos {
	agentPos := make([]pos.Pos, len(state.AgentPos))
	copy(agentPos, state.AgentPos)
	//現在ある座標にいるエージェントのID
	former := make(map[pos.Pos]int)
	//次にある座標に行きたいエージェントのID
	candidates := make(map[pos.Pos][]int)
	for i, now := range agentPos {
		former[now] = i
		agentPos[i] = pos.NextPos(now, actions[i], env.MapData)
		candidates[agentPos[i]] = append(candidates[agentPos[i]], i)
	}
	doDFS(state, env, agentPos, former, candidates)
	return agentPos
}

func doDFS(state *State, env *env.Env, agentPos []pos.Pos, former map[pos.Pos]int, candidates map[pos.Pos][]int) {
	const (
		INIT int = iota
		PENDING
		DECIDED
	)
	status := make([]int, env.NumAgents)

	var dfs func(id int)
	dfs = func(id int) {
		status[id] = PENDING
		canMove := true
		formerID, exist := former[agentPos[id]]
		//進む先に現在エージェントがいる場合
		if exist && formerID != id {
			switch {
			case status[formerID] == INIT:
				dfs(formerID)
				if agentPos[formerID] == agentPos[id] {
					canMove = false
				}
			case status[formerID] == PENDING:
				canMove = false
			case status[formerID] == DECIDED && agentPos[formerID] == agentPos[id]:
				canMove = false
			}
		}
		if !canMove {
			for _, candID := range candidates[agentPos[id]] {
				status[candID] = DECIDED
				agentPos[candID] = state.AgentPos[candID]
			}
			return
		}
		//誰とも競合していないなら
		if len(candidates[agentPos[id]]) == 1 {
			status[id] = DECIDED
			return
		}
		switch env.Resolve {
		case "ALL STAY":
			for _, candID := range candidates[agentPos[id]] {
				status[candID] = DECIDED
				agentPos[candID] = state.AgentPos[candID]
			}
		case "DEADLINE BASE":
			winner := -1
			winnerDeadline := 1 << 30
			for _, candID := range candidates[agentPos[id]] {
				candDeadline := 1 << 30
				if len(state.AgentItems[candID]) > 0 {
					candDeadline = state.AgentItems[candID][0]
				}
				//今いる場所にとどまるエージェントが最優先
				if candID == formerID {
					candDeadline = -100
				}
				//現在デポにいるエージェントを優先
				if state.AgentPos[candID] == env.DepotPos {
					candDeadline = -99
				}
				if candDeadline < winnerDeadline {
					winner = candID
					winnerDeadline = candDeadline
				}
			}
			for _, candID := range candidates[agentPos[id]] {
				status[candID] = DECIDED
				if candID != winner {
					agentPos[candID] = state.AgentPos[candID]
				}
			}
		}
	}
	for i := 0; i < env.NumAgents; i++ {
		if status[i] == INIT {
			dfs(i)
		}
	}
}
