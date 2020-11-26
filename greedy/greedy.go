package greedy

import (
	"math"
	"math/rand"
	"sort"

	"github.com/Div9851/warehouse-sim/action"
	"github.com/Div9851/warehouse-sim/env"
	"github.com/Div9851/warehouse-sim/pos"
	"github.com/Div9851/warehouse-sim/state"
)

//あるエージェントにとっての, ある点の価値を返す
func eval(id int, pos pos.Pos, state *state.State, env *env.Env) float64 {
	d := 1 + float64(env.MinDist[state.AgentPos[id]][pos])
	if pos == env.DepotPos {
		return float64(state.AgentItems[id]) * env.Reward / d
	}
	m := math.Min(float64(state.PosItems[pos]), float64(env.MaxItems-state.AgentItems[id]))
	return m * env.Reward / d
}

type tuple struct {
	ID        int
	Pos       pos.Pos
	Value     float64
	RandomVal float64
}

type tuples []tuple

func (t tuples) Len() int {
	return len(t)
}

func (t tuples) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t tuples) Less(i, j int) bool {
	if t[i].Value != t[j].Value {
		return t[i].Value < t[j].Value
	}
	if t[i].ID != t[j].ID {
		return t[i].ID < t[j].ID
	}
	return t[i].RandomVal < t[j].RandomVal
}

func makeTuple(id int, pos pos.Pos, value float64, randomVal float64) tuple {
	return tuple{ID: id, Pos: pos, Value: value, RandomVal: randomVal}
}

//Greedy 貪欲法で行動を決定する
func Greedy(state *state.State, env *env.Env, rnd *rand.Rand, plannerID int, plannerAction int) []int {
	reserved := make(map[pos.Pos]int)
	blocked := make(map[pos.Pos]bool)
	agentID := make(map[pos.Pos]int)
	decided := make([]bool, env.NumAgents)
	actions := make([]int, env.NumAgents)
	dest := make([]pos.Pos, env.NumAgents)
	ts := make(tuples, 0)
	if plannerID != -1 {
		decided[plannerID] = true
		actions[plannerID] = plannerAction
		nxt := pos.NextPos(state.AgentPos[plannerID], actions[plannerID], env.MapData)
		dest[plannerID] = nxt
		blocked[nxt] = true
	}
	for id := 0; id < env.NumAgents; id++ {
		if id == plannerID {
			continue
		}
		agentID[state.AgentPos[id]] = id
		for pos := range state.PosItems {
			ts = append(ts, makeTuple(id, pos, eval(id, pos, state, env), state.RandomValues[pos]))
		}
		ts = append(ts, makeTuple(id, env.DepotPos, eval(id, env.DepotPos, state, env), state.RandomValues[env.DepotPos]))
	}
	sort.Sort(sort.Reverse(ts))
	for _, t := range ts {
		if t.Value == 0 {
			break
		}
		//すでに行動が決まっているならスキップ
		if decided[t.ID] {
			continue
		}
		//すでにアイテム数と同じ数のエージェントが予約していたらダメ
		if t.Pos != env.DepotPos && reserved[t.Pos] == state.PosItems[t.Pos] {
			continue
		}
		//目的地にいるなら
		if state.AgentPos[t.ID] == t.Pos {
			decided[t.ID] = true
			if t.Pos == env.DepotPos {
				actions[t.ID] = action.CLEAR
			} else {
				actions[t.ID] = action.PICKUP
			}
			dest[t.ID] = t.Pos
			blocked[t.Pos] = true
			reserved[t.Pos]++
			continue
		}
		moves := []int{}
		validMoves := env.ValidMoves[state.AgentPos[t.ID]]
		for _, move := range validMoves {
			nxt := pos.NextPos(state.AgentPos[t.ID], move, env.MapData)
			//すでにブロックされているならダメ
			if blocked[nxt] {
				continue
			}
			//すれ違うような動き方はダメ
			otherID, exist := agentID[nxt]
			if exist && decided[otherID] && dest[otherID] == state.AgentPos[t.ID] {
				continue
			}
			//目的地に近づくなら
			if env.MinDist[state.AgentPos[t.ID]][t.Pos] > env.MinDist[nxt][t.Pos] {
				moves = append(moves, move)
			}
		}
		//目的地に近づく動き方がなければスキップ
		if len(moves) == 0 {
			continue
		}
		decided[t.ID] = true
		actions[t.ID] = moves[rnd.Intn(len(moves))]
		nxt := pos.NextPos(state.AgentPos[t.ID], actions[t.ID], env.MapData)
		dest[t.ID] = nxt
		blocked[nxt] = true
		reserved[t.Pos]++
	}
	for id := 0; id < env.NumAgents; id++ {
		//すでに行動が決まっているならスキップ
		if decided[id] {
			continue
		}
		moves := []int{}
		validMoves := env.ValidMoves[state.AgentPos[id]]
		for _, move := range validMoves {
			nxt := pos.NextPos(state.AgentPos[id], move, env.MapData)
			//すでにブロックされているならダメ
			if blocked[nxt] {
				continue
			}
			//すれ違うような動き方はダメ
			otherID, exist := agentID[nxt]
			if exist && decided[otherID] && dest[otherID] == state.AgentPos[id] {
				continue
			}
			moves = append(moves, move)
		}
		decided[id] = true
		if len(moves) == 0 {
			actions[id] = validMoves[rnd.Intn(len(validMoves))]
		} else {
			actions[id] = moves[rnd.Intn(len(moves))]
		}
		nxt := pos.NextPos(state.AgentPos[id], actions[id], env.MapData)
		dest[id] = nxt
		blocked[nxt] = true
	}
	return actions
}
