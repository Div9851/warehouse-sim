package mcts

import (
	"math"
	"math/rand"
	"sort"

	"github.com/Div9851/warehouse-sim/action"
	"github.com/Div9851/warehouse-sim/env"
	"github.com/Div9851/warehouse-sim/greedy"
	"github.com/Div9851/warehouse-sim/state"
)

type tuple struct {
	ID    int
	Score float64
}

type tuples []tuple

func (t tuples) Len() int {
	return len(t)
}

func (t tuples) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t tuples) Less(i, j int) bool {
	if t[i].Score != t[j].Score {
		return t[i].Score < t[j].Score
	}
	return t[i].ID < t[j].ID
}

func makeTuple(id int, score float64) tuple {
	return tuple{ID: id, Score: score}
}

//MCTS モンテカルロ木探索で行動を決定する（選択した行動と, 更新された楽観度を返す）
func MCTS(id int, startState *state.State, env *env.Env, rnd *rand.Rand, opt bool, prevOpt float64) (int, float64) {
	states := []*state.State{startState}
	//ある状態に遷移したときに得られる報酬
	stateRewards := make([]float64, 1)
	//ある状態である行動を選んだときの遷移先
	childs := make([][][]int, 1)
	childs[0] = make([][]int, action.NUM)
	//ある状態である行動を選んだ回数
	counts := make([][]int, 1)
	counts[0] = make([]int, action.NUM)
	//ある状態で行動を選んだ回数の合計
	totalCount := make([]int, 1)
	//ある状態である行動を選んだときの報酬の和
	sumReward := make([][]float64, 1)
	sumReward[0] = make([]float64, action.NUM)
	//ある状態でロールアウトを行った回数
	simCount := []int{env.ExpandTheresh} //始点はすぐ展開

	var dfs func(int, int, float64) float64
	dfs = func(stateID int, depth int, opt float64) float64 {
		if states[stateID].Turn >= env.LastTurn || depth >= env.MaxDepth {
			return 0
		}
		if simCount[stateID] < env.ExpandTheresh {
			simCount[stateID]++
			var r float64
			var k float64 = 1
			//roll out
			now := states[stateID]
			for now.Turn < env.LastTurn && depth < env.MaxDepth {
				actionLists := greedy.Greedy(now, env, rnd)
				nxt, _, _, rewards := state.NextStateOpt(now, actionLists, env, rnd, id, opt)
				now = nxt
				r += k * rewards[id]
				k *= env.DiscountFactor
				depth++
			}
			return r
		}
		validActions := make([]int, len(env.ValidMoves[states[stateID].AgentPos[id]]))
		copy(validActions, env.ValidMoves[states[stateID].AgentPos[id]])
		if states[stateID].PosItems[states[stateID].AgentPos[id]] > 0 && states[stateID].AgentItems[id] < env.MaxItems {
			validActions = append(validActions, action.PICKUP)
		}
		if states[stateID].AgentPos[id] == env.DepotPos && states[stateID].AgentItems[id] > 0 {
			validActions = append(validActions, action.CLEAR)
		}
		var bestScore float64
		var bestActions []int
		for _, act := range validActions {
			var score float64
			if counts[stateID][act] == 0 {
				score = math.Inf(0)
			} else {
				//UCT
				score = sumReward[stateID][act] / float64(counts[stateID][act])
				score += math.Sqrt(env.UCTparam * math.Log(float64(totalCount[stateID])) / float64(counts[stateID][act]))
			}
			if bestScore < score {
				bestScore = score
				bestActions = []int{act}
			} else if bestScore == score {
				bestActions = append(bestActions, act)
			}
		}
		chosen := bestActions[rnd.Intn(len(bestActions))]
		var to int
		var r float64
		//遷移先の数が上限に達していたら
		if len(childs[stateID][chosen]) == env.MaxChilds {
			to = childs[stateID][chosen][rnd.Intn(len(childs[stateID][chosen]))]
		} else {
			actions := greedy.Greedy(states[stateID], env, rnd)
			actions[id] = chosen
			nxt, _, _, rewards := state.NextStateOpt(states[stateID], actions, env, rnd, id, opt)
			to = len(states)
			childs[stateID][chosen] = append(childs[stateID][chosen], to)
			states = append(states, nxt)
			stateRewards = append(stateRewards, rewards[id])
			childs = append(childs, make([][]int, action.NUM))
			counts = append(counts, make([]int, action.NUM))
			totalCount = append(totalCount, 0)
			sumReward = append(sumReward, make([]float64, action.NUM))
			simCount = append(simCount, 0)
		}
		r += stateRewards[to]
		r += env.DiscountFactor * dfs(to, depth+1, opt)
		sumReward[stateID][chosen] += r
		totalCount[stateID]++
		counts[stateID][chosen]++
		return r
	}
	nxtOpt := prevOpt
	if opt {
		if startState.Success[id] {
			nxtOpt = math.Min(nxtOpt*1.5, 1.0)
		} else {
			nxtOpt /= 2
		}
		for i := 0; i < env.NumOfIter; i++ {
			dfs(0, 1, nxtOpt)
		}
	} else {
		for i := 0; i < env.NumOfIter; i++ {
			dfs(0, 1, 0)
		}
	}
	ts := make(tuples, 0)
	validActions := make([]int, len(env.ValidMoves[states[0].AgentPos[id]]))
	copy(validActions, env.ValidMoves[states[0].AgentPos[id]])
	if states[0].PosItems[states[0].AgentPos[id]] > 0 && states[0].AgentItems[id] < env.MaxItems {
		validActions = append(validActions, action.PICKUP)
	}
	if states[0].AgentPos[id] == env.DepotPos && states[0].AgentItems[id] > 0 {
		validActions = append(validActions, action.CLEAR)
	}
	for _, act := range validActions {
		var score float64
		if counts[0][act] == 0 {
			score = math.Inf(-1)
		} else {
			score = sumReward[0][act] / float64(counts[0][act])
		}
		ts = append(ts, makeTuple(act, score))
	}
	sort.Sort(sort.Reverse(ts))
	return ts[0].ID, nxtOpt
}
