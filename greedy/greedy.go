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
		return float64(len(state.AgentItems[id])) * env.ClearReward / d
	}
	m := math.Min(float64(len(state.PosItems[pos])), float64(env.MaxItems-len(state.AgentItems[id])))
	return m * env.PickupReward / d
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
func Greedy(state *state.State, env *env.Env, rnd *rand.Rand) []int {
	dest := make(map[int]pos.Pos)
	value := make(map[int]float64)
	reserved := make(map[pos.Pos]bool)
	ts := make(tuples, 0)
	for id := 0; id < env.NumAgents; id++ {
		for pos := range state.PosItems {
			ts = append(ts, makeTuple(id, pos, eval(id, pos, state, env), state.RandomValues[pos]))
		}
		ts = append(ts, makeTuple(id, env.DepotPos, eval(id, env.DepotPos, state, env), state.RandomValues[env.DepotPos]))
	}
	sort.Sort(sort.Reverse(ts))
	for _, t := range ts {
		//すでに目的地が決まっているならスキップ
		_, exist := dest[t.ID]
		if exist {
			continue
		}
		//デポは例外的に複数のエージェントが同時に目的地に出来る
		if t.Pos == env.DepotPos {
			dest[t.ID] = t.Pos
			value[t.ID] = t.Value
			continue
		}
		//目的地に出来るのは一人だけ
		if reserved[t.Pos] {
			continue
		}
		dest[t.ID] = t.Pos
		value[t.ID] = t.Value
		reserved[t.Pos] = true
	}
	actions := make([]int, env.NumAgents)
	for id := 0; id < env.NumAgents; id++ {
		switch {
		case state.AgentPos[id] == dest[id] && value[id] > 0: //目的地にいて, 価値が非0なら
			if dest[id] == env.DepotPos {
				actions[id] = action.CLEAR
			} else {
				actions[id] = action.PICKUP
			}
		case state.AgentPos[id] != dest[id] && value[id] > 0: //目的地にいなくて, 価値が非0なら
			moves := []int{action.UP, action.DOWN, action.LEFT, action.RIGHT}
			goodMoves := []int{}
			for move := range moves {
				nxt := pos.NextPos(state.AgentPos[id], move, env.MapData)
				//目的地に近づくなら
				if env.MinDist[state.AgentPos[id]][dest[id]] > env.MinDist[nxt][dest[id]] {
					goodMoves = append(goodMoves, move)
				}
			}
			actions[id] = goodMoves[rnd.Intn(len(goodMoves))]
		default: //目的地の価値が0ならランダムに行動
			moves := []int{action.UP, action.DOWN, action.LEFT, action.RIGHT}
			goodMoves := []int{}
			for move := range moves {
				nxt := pos.NextPos(state.AgentPos[id], move, env.MapData)
				//動けるなら
				if nxt != state.AgentPos[id] {
					goodMoves = append(goodMoves, move)
				}
			}
			actions[id] = goodMoves[rnd.Intn(len(goodMoves))]
		}
	}
	return actions
}
