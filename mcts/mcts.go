package mcts

import (
	"math"
	"math/rand"

	"github.com/Div9851/warehouse-sim/action"
	"github.com/Div9851/warehouse-sim/env"
	"github.com/Div9851/warehouse-sim/greedy"
	"github.com/Div9851/warehouse-sim/state"
)

//MCTS モンテカルロ木探索で行動を決定する
func MCTS(id int, startState *state.State, env *env.Env, rnd *rand.Rand) []int {
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

	var dfs func(int, int) float64
	dfs = func(stateID int, depth int) float64 {
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
				nxt, _, _, rewards := state.NextState(now, actionLists, env, rnd)
				now = nxt
				r += k * rewards[id]
				k *= env.DiscountFactor
				depth++
			}
			return r
		}
		var bestScore float64
		var bestActions []int
		for act := 0; act < action.NUM; act++ {
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
			actionLists := greedy.Greedy(states[stateID], env, rnd)
			actionLists[id] = []int{chosen}
			nxt, _, _, rewards := state.NextState(states[stateID], actionLists, env, rnd)
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
		r += env.DiscountFactor * dfs(to, depth+1)
		sumReward[stateID][chosen] += r
		totalCount[stateID]++
		counts[stateID][chosen]++
		return r
	}
	for i := 0; i < env.NumOfIter; i++ {
		dfs(0, 1)
	}
	bestScore := math.Inf(-1)
	var bestActions []int
	for act := 0; act < action.NUM; act++ {
		var score float64
		if counts[0][act] == 0 {
			score = math.Inf(-1)
		} else {
			score = sumReward[0][act] / float64(counts[0][act])
		}
		if bestScore < score {
			bestScore = score
			bestActions = []int{act}
		} else if bestScore == score {
			bestActions = append(bestActions, act)
		}
	}
	return []int{bestActions[rnd.Intn(len(bestActions))]}
}
