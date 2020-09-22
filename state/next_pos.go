package state

import (
	"github.com/Div9851/warehouse-sim/env"
	"github.com/Div9851/warehouse-sim/pos"
)

//nextPos 現在の状態, 各エージェントの行動, 環境設定を受け取って
//次の状態のAgentPosを返す
func nextPos(state *State, actions []int, env *env.Env) []pos.Pos {
	nxtPos := make([]pos.Pos, env.NumAgents)
	copy(nxtPos, state.AgentPos)
	//現在ある座標にいるエージェントのID
	currentID := make(map[pos.Pos]int)
	//次にある座標に行きたいエージェントのID
	nextID := make(map[pos.Pos][]int)
	for id, now := range state.AgentPos {
		currentID[now] = id
		nxtPos[id] = pos.NextPos(now, actions[id], env.MapData)
		if nxtPos[id] != now {
			nextID[nxtPos[id]] = append(nextID[nxtPos[id]], id)
		}
	}
	doDFS(state, env, nxtPos, currentID, nextID)
	return nxtPos
}

func doDFS(state *State, env *env.Env, nxtPos []pos.Pos, currentID map[pos.Pos]int, nextID map[pos.Pos][]int) {
	const (
		INIT int = iota
		PENDING
		DECIDED
	)
	status := make([]int, env.NumAgents)

	var dfs func(id int)
	dfs = func(id int) {
		status[id] = PENDING
		//その場にとどまる場合
		if nxtPos[id] == state.AgentPos[id] {
			status[id] = DECIDED
			return
		}
		curID, exist := currentID[nxtPos[id]]
		//進む先に現在エージェントがいる場合
		if exist {
			canMove := true
			if status[curID] == INIT {
				dfs(curID)
			}
			if status[curID] == PENDING || nxtPos[curID] == state.AgentPos[curID] {
				canMove = false
			}
			if !canMove {
				for _, nxtID := range nextID[nxtPos[id]] {
					status[nxtID] = DECIDED
					nxtPos[nxtID] = state.AgentPos[nxtID]
				}
				return
			}
		}
		//誰とも競合していないなら
		if len(nextID[nxtPos[id]]) == 1 {
			status[id] = DECIDED
			return
		}
		if env.Resolve {
			chosenID := -1
			bestScore := -1
			for _, nxtID := range nextID[nxtPos[id]] {
				score := state.AgentItems[nxtID]
				//デポにいるエージェントが最も優先される
				if state.AgentPos[nxtID] == env.DepotPos {
					score = env.MaxItems + 1
				}
				if score > bestScore {
					chosenID = nxtID
					bestScore = score
				} else if score == bestScore && chosenID < nxtID {
					chosenID = nxtID
				}
			}
			for _, nxtID := range nextID[nxtPos[id]] {
				if nxtID != chosenID {
					status[nxtID] = DECIDED
					nxtPos[nxtID] = state.AgentPos[nxtID]
				}
			}
		} else {
			for _, nxtID := range nextID[nxtPos[id]] {
				status[nxtID] = DECIDED
				nxtPos[nxtID] = state.AgentPos[nxtID]
			}
		}
	}
	for i := 0; i < env.NumAgents; i++ {
		if status[i] == INIT {
			dfs(i)
		}
	}
}
