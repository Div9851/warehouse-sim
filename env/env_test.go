package env

import (
	"testing"

	"github.com/Div9851/warehouse-sim/pos"
)

func TestLoadFromJSON(t *testing.T) {
	env, err := loadFromJSON("testdata/example.json")
	if err != nil {
		t.Fatal(err)
	}
	if env.Name != "Example" {
		t.Fatalf("env.Name should be `Example`, but `%v`", env.Name)
	}
	if env.NumAgents != 3 {
		t.Fatalf("env.NumAgents should be `3`, but `%v`", env.NumAgents)
	}
	if env.MaxItems != 1 {
		t.Fatalf("env.MaxItems should be `1`, but `%v`", env.MaxItems)
	}
	if env.LastTurn != 100 {
		t.Fatalf("env.LastTurn should be `100`, but `%v`", env.LastTurn)
	}
	if env.Reward != 100 {
		t.Fatalf("env.Reward should be `100`, but `%v`", env.Reward)
	}
	if env.DIYBonus != 70 {
		t.Fatalf("env.DIYBonus should be `70`, but `%v`", env.DIYBonus)
	}
	if env.MapDataPath != "map_data.txt" {
		t.Fatalf("env.MapDataPath should be `map_data.txt`, but `%v`", env.MapDataPath)
	}
	if env.AppearProb != 0.8 {
		t.Fatalf("env.AppearProb should be `0.8`, but `%v`", env.AppearProb)
	}
	if env.DepotPos.X != 0 || env.DepotPos.Y != 3 {
		t.Fatalf("env.DepotPos should be `(0, 3)`, but `(%v, %v)`", env.DepotPos.X, env.DepotPos.Y)
	}
	if !env.Resolve {
		t.Fatalf("env.Resolve should be `true`, but `%v`", env.Resolve)
	}
	if env.Algorithm != "MCTS" {
		t.Fatalf("env.Algorithm should be `MCTS`, but `%v`", env.Algorithm)
	}
	if env.DiscountFactor != 0.9 {
		t.Fatalf("env.DiscountFactor should be `0.9`, but `%v`", env.DiscountFactor)
	}
	if env.ExpandTheresh != 1 {
		t.Fatalf("env.ExpandThresh should be `1`, but `%v`", env.ExpandTheresh)
	}
	if env.MaxChilds != 5 {
		t.Fatalf("env.MaxChilds should be `5`, but `%v`", env.MaxChilds)
	}
	if env.MaxDepth != 60 {
		t.Fatalf("env.MaxDepth should be `60`, but `%v`", env.MaxDepth)
	}
	if env.NumOfIter != 20000 {
		t.Fatalf("env.NumOfIter should be `20000`, but `%v`", env.NumOfIter)
	}
	if env.UCTparam != 2 {
		t.Fatalf("env.UCTparam should be `2`, but `%v`", env.UCTparam)
	}
}

func TestLoadMapData(t *testing.T) {
	mapData, err := loadMapData("testdata/map_data.txt")
	if err != nil {
		t.Fatal(err)
	}
	expected := []string{
		"...#...",
		".#.#.#.",
		".#.#.#.",
		".......",
		".##.##.",
		".##.##.",
		".##.##.",
	}
	if len(mapData) != len(expected) {
		t.Fatalf("len(mapData) should be `%v`, but `%v`", len(expected), len(mapData))
	}
	for i := range mapData {
		if mapData[i] != expected[i] {
			t.Fatalf("mapData[%v] should be `%v`, but `%v`", i, expected[i], mapData[i])
		}
	}
}

func TestGetAllPos(t *testing.T) {
	mapData, err := loadMapData("testdata/map_data.txt")
	if err != nil {
		t.Fatal(err)
	}
	allPos := getAllPos(mapData, pos.New(0, 3))
	if len(allPos) != 29 {
		t.Fatalf("len(allPos) should be `29`, but `%v`", len(allPos))
	}
}

func TestDoBFS(t *testing.T) {
	mapData, err := loadMapData("testdata/map_data.txt")
	if err != nil {
		t.Fatal(err)
	}
	minDist := doBFS(mapData, pos.New(0, 3))
	var to pos.Pos
	to = pos.New(2, 0)
	if minDist[to] != 5 {
		t.Fatalf("minDist[(2, 0)] should be `5`, but `%v`", minDist[to])
	}
	to = pos.New(6, 6)
	if minDist[to] != 9 {
		t.Fatalf("minDist[(6, 6)] should be `9`, but `%v`", minDist[to])
	}
}

func TestLoad(t *testing.T) {
	env, err := Load("testdata/example.json")
	if err != nil {
		t.Fatal(err)
	}
	if env.MapDataH != 7 {
		t.Fatalf("env.MapDataH should be `7`, but `%v`", env.MapDataH)
	}
	if env.MapDataW != 7 {
		t.Fatalf("env.MapDataW should be `7`, but `%v`", env.MapDataW)
	}
	if len(env.AllPos) != 29 {
		t.Fatalf("len(env.AllPos) should be `29`, but `%v`", len(env.AllPos))
	}
	var (
		from pos.Pos
		to   pos.Pos
	)
	from = pos.New(4, 0)
	to = pos.New(3, 5)
	if env.MinDist[from][to] != 6 {
		t.Fatalf("env.MinDist[(4, 0)][(3, 5)] should be `6`, but `%v`", env.MinDist[from][to])
	}
	from = pos.New(6, 6)
	to = pos.New(0, 3)
	if env.MinDist[from][to] != 9 {
		t.Fatalf("env.MinDist[(6, 6)][(0, 3)] should be `9`, but `%v`", env.MinDist[from][to])
	}
}
