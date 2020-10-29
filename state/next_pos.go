package state

import (
	"math/rand"

	"github.com/Div9851/warehouse-sim/env"
	"github.com/Div9851/warehouse-sim/pos"
)

//nextPos 現在の状態, 各エージェントの行動, 環境設定を受け取って
//次の状態のAgentPosを返す
func nextPos(state *State, actions []int, env *env.Env, rnd *rand.Rand) []pos.Pos {
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
	doDFS(state, env, rnd, nxtPos, currentID, nextID, -1, 0.0)
	return nxtPos
}

//nextPosOpt あるエージェントを優先するようなnextPos
func nextPosOpt(state *State, actions []int, env *env.Env, rnd *rand.Rand, favoredID int, opt float64) []pos.Pos {
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
	doDFS(state, env, rnd, nxtPos, currentID, nextID, favoredID, opt)
	return nxtPos
}

func doDFS(state *State, env *env.Env, rnd *rand.Rand, nxtPos []pos.Pos, currentID map[pos.Pos]int, nextID map[pos.Pos][]int, favoredID int, opt float64) {
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
		//競合しているなら全員STAY（ただし優先されるエージェントがいる時は例外）
		for _, nxtID := range nextID[nxtPos[id]] {
			if nxtID == favoredID && rnd.Float64() < opt {
				status[nxtID] = DECIDED
			} else {
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
