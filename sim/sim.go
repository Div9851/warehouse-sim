package sim

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/Div9851/warehouse-sim/action"
	"github.com/Div9851/warehouse-sim/env"
	"github.com/Div9851/warehouse-sim/greedy"
	"github.com/Div9851/warehouse-sim/mcts"
	"github.com/Div9851/warehouse-sim/pos"
	"github.com/Div9851/warehouse-sim/state"
)

//Simulator シミュレータを表す構造体
type Simulator struct {
	Env          *env.Env
	State        *state.State
	LastActions  []int
	TotalRewards []float64
	LastRewards  []float64
	LastAppear   *pos.Pos
	PrevOpt      []float64
	TotalItems   int
	PickupCounts []int
	ClearCounts  []int
	SimRand      *rand.Rand
	Rands        []*rand.Rand
	Seed         int64
}

//New 環境設定とシード値を受け取り, シミュレータを返す
func New(env *env.Env, seed int64) *Simulator {
	totalRewards := make([]float64, env.NumAgents)
	agentItems := make([]int, env.NumAgents)
	agentPos := make([]pos.Pos, env.NumAgents)
	posItems := make(map[pos.Pos]int)
	pickupCounts := make([]int, env.NumAgents)
	clearCounts := make([]int, env.NumAgents)
	simRand := rand.New(rand.NewSource(seed))
	rands := make([]*rand.Rand, env.NumAgents)
	prevOpt := make([]float64, env.NumAgents)
	for i := range rands {
		rands[i] = rand.New(rand.NewSource(simRand.Int63()))
		agentPos[i] = env.AllPos[simRand.Intn(len(env.AllPos))]
	}
	randomValues := make(map[pos.Pos]float64)
	randomValues[env.DepotPos] = simRand.Float64()
	for _, pos := range env.AllPos {
		randomValues[pos] = simRand.Float64()
	}
	success := make([]bool, env.NumAgents)
	state := state.New(1, agentItems, agentPos, posItems, randomValues, success)
	return &Simulator{Env: env, State: state, TotalRewards: totalRewards, PrevOpt: prevOpt, PickupCounts: pickupCounts, ClearCounts: clearCounts, SimRand: simRand, Rands: rands, Seed: seed}
}

//Do シミュレーションを実行し, 実行時間を返す
func (sim *Simulator) Do(verbose bool) float64 {
	startTime := time.Now()
	for {
		if verbose {
			fmt.Printf("%v\n\n", sim.DumpState())
		}
		if !sim.Next() {
			break
		}
	}
	endTime := time.Now()
	processTime := endTime.Sub(startTime).Seconds()
	return processTime
}

//Next シミュレーションを1ステップ進める（すでに終了していればfalseを返す）
func (sim *Simulator) Next() bool {
	if sim.State.Turn == sim.Env.LastTurn {
		return false
	}
	actions := make([]int, sim.Env.NumAgents)
	wg := &sync.WaitGroup{}
	for i := 0; i < sim.Env.NumAgents; i++ {
		wg.Add(1)
		go func(id int) {
			switch sim.Env.Algorithms[id] {
			case "MCTS":
				actions[id], _ = mcts.MCTS(id, sim.State, sim.Env, sim.Rands[id], false, 0)
			case "MCTS_OPT":
				var nxtOpt float64
				actions[id], nxtOpt = mcts.MCTS(id, sim.State, sim.Env, sim.Rands[id], true, sim.PrevOpt[id])
				sim.PrevOpt[id] = nxtOpt
			default:
				actions[id] = greedy.Greedy(sim.State, sim.Env, sim.Rands[id], 0)[id]
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	nxtState, lastActions, lastAppear, lastRewards := state.NextState(sim.State, actions, sim.Env, sim.SimRand)
	for i, act := range lastActions {
		if act == action.CLEAR && nxtState.Success[i] {
			sim.ClearCounts[i]++
		}
		if act == action.PICKUP && nxtState.Success[i] {
			sim.PickupCounts[i]++
		}
	}
	if lastAppear != nil {
		sim.TotalItems++
	}
	sim.State = nxtState
	sim.LastActions = lastActions
	sim.LastRewards = lastRewards
	sim.LastAppear = lastAppear
	for i, r := range lastRewards {
		sim.TotalRewards[i] += r
	}
	return true
}

//DumpState 現在の状態を表す文字列を返す
func (sim *Simulator) DumpState() string {
	mapData := make([]string, len(sim.Env.MapData))
	copy(mapData, sim.Env.MapData)
	pos := sim.Env.DepotPos
	mapData[pos.Y] = mapData[pos.Y][:pos.X] + "D" + mapData[pos.Y][pos.X+1:]
	for pos := range sim.State.PosItems {
		mapData[pos.Y] = mapData[pos.Y][:pos.X] + "*" + mapData[pos.Y][pos.X+1:]
	}
	for i, pos := range sim.State.AgentPos {
		mapData[pos.Y] = mapData[pos.Y][:pos.X] + fmt.Sprint(i) + mapData[pos.Y][pos.X+1:]
	}
	var b strings.Builder
	fmt.Fprintf(&b, "[TURN %v]\n", sim.State.Turn)
	for _, row := range mapData {
		fmt.Fprintln(&b, row)
	}
	if sim.LastAppear != nil {
		fmt.Fprintln(&b, "[NEW ITEM]")
		fmt.Fprintf(&b, "(%v, %v)\n", sim.LastAppear.X, sim.LastAppear.Y)
	}
	if sim.LastActions != nil {
		fmt.Fprintln(&b, "[ACTIONS]")
		for i, act := range sim.LastActions {
			var actStr string
			switch act {
			case action.UP:
				actStr = "UP"
			case action.DOWN:
				actStr = "DOWN"
			case action.LEFT:
				actStr = "LEFT"
			case action.RIGHT:
				actStr = "RIGHT"
			case action.PICKUP:
				actStr = "PICKUP"
			case action.CLEAR:
				actStr = "CLEAR"
			case action.STAY:
				actStr = "STAY"
			}
			fmt.Fprintf(&b, "agent %v: %v ", i, actStr)
		}
		fmt.Fprintln(&b)
	}
	fmt.Fprintln(&b, "[ITEMS]")
	for i, items := range sim.State.AgentItems {
		fmt.Fprintf(&b, "agent %v: %v ", i, items)
	}
	fmt.Fprintln(&b)
	fmt.Fprintln(&b, "[REWARDS]")
	for i, r := range sim.TotalRewards {
		var lastReward float64
		if sim.LastRewards != nil {
			lastReward = sim.LastRewards[i]
		}
		var diff string
		if lastReward == 0 {
			diff = "±0"
		} else if lastReward > 0 {
			diff = "+" + fmt.Sprint(lastReward)
		} else {
			diff = fmt.Sprint(lastReward)
		}
		fmt.Fprintf(&b, "agent %v: %v (%v) ", i, r, diff)
	}
	return b.String()
}

//GetResult シミュレーションの結果を返す
func (sim *Simulator) GetResult() *Result {
	return &Result{TotalItems: sim.TotalItems, PickupCounts: sim.PickupCounts, ClearCounts: sim.ClearCounts}
}
