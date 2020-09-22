package env

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Div9851/warehouse-sim/action"
	"github.com/Div9851/warehouse-sim/pos"
)

//Env 環境設定
type Env struct {
	Name        string  `json:"name"`
	NumAgents   int     `json:"num_agents"`
	MaxItems    int     `json:"max_items"`
	LastTurn    int     `json:"last_turn"`
	MaxLen      int     `json:"max_len"` //行動のリストの最大長（この値より長い場合は無視される）
	Reward      float64 `json:"reward"`
	DIYBonus    float64 `json:"DIY_bonus"` //自分でアイテムを運んだとき/回収したときに追加で得られるReward
	MapDataPath string  `json:"map_data_path"`
	AppearProb  float64 `json:"appear_prob"`
	DepotPos    pos.Pos `json:"depot_pos"`
	Resolve     bool    `json:"resolve"`
	Algorithm   string  `json:"algorithm"` //GREEDY, MCTS

	DiscountFactor float64 `json:"mcts_discount_factor"`
	ExpandTheresh  int     `json:"mcts_expand_thresh"` //ノードを展開する閾値
	MaxChilds      int     `json:"mcts_max_childs"`    //遷移先の数の上限
	MaxDepth       int     `json:"mcts_max_depth"`
	NumOfIter      int     `json:"mcts_num_of_iter"`
	UCTparam       float64 `json:"uct_param"`
	MapData        []string
	MapDataH       int
	MapDataW       int
	MinDist        map[pos.Pos]map[pos.Pos]int
	AllPos         []pos.Pos         //壁でない全ての座標のスライス（デポを含まない）
	ValidMoves     map[pos.Pos][]int //その場所で選択できる行動のリスト
}

//Load 環境設定をJSONファイルから読み込む
func Load(path string) (*Env, error) {
	env, err := loadFromJSON(path)
	if err != nil {
		return nil, err
	}
	dir := filepath.Dir(path)
	env.MapData, err = loadMapData(filepath.Join(dir, env.MapDataPath))
	if err != nil {
		return nil, err
	}
	env.MapDataH = len(env.MapData)
	env.MapDataW = len(env.MapData[0])
	env.AllPos = getAllPos(env.MapData, env.DepotPos)
	env.ValidMoves = make(map[pos.Pos][]int)
	for _, p := range env.AllPos {
		env.ValidMoves[p] = getValidMoves(env.MapData, p)
	}
	env.ValidMoves[env.DepotPos] = getValidMoves(env.MapData, env.DepotPos)
	env.MinDist = make(map[pos.Pos]map[pos.Pos]int)
	for _, p := range env.AllPos {
		env.MinDist[p] = doBFS(env.MapData, p)
	}
	env.MinDist[env.DepotPos] = doBFS(env.MapData, env.DepotPos)
	return env, nil
}

//loadFromJSON 環境設定をJSONファイルから読み込む
func loadFromJSON(path string) (*Env, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("can't open `%s` (%s)", path, err)
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("can't read `%s` (%s)", path, err)
	}

	var env Env
	err = json.Unmarshal(b, &env)
	if err != nil {
		return nil, fmt.Errorf("can't decode `%s` (%s)", path, err)
	}
	return &env, nil
}

//loadMapData マップデータをテキストファイルから読み込む
func loadMapData(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("can't open `%s` (%s)", path, err)
	}
	defer f.Close()

	mapData := []string{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		mapData = append(mapData, scanner.Text())
	}
	if scanner.Err() != nil {
		return nil, fmt.Errorf("can't read `%s` (%s)", path, scanner.Err())
	}
	return mapData, nil
}

//getAllPos マップデータとデポの座標を受け取り, 壁でない全ての座標のスライスを返す（デポを含まない）
func getAllPos(mapData []string, depotPos pos.Pos) []pos.Pos {
	allPos := []pos.Pos{}
	for y, row := range mapData {
		for x, col := range row {
			if col != '#' && (x != depotPos.X || y != depotPos.Y) {
				allPos = append(allPos, pos.New(x, y))
			}
		}
	}
	return allPos
}

//getValidMoves マップデータと座標を受け取り, 選択できる行動のリストを返す
func getValidMoves(mapData []string, targetPos pos.Pos) []int {
	moves := []int{action.UP, action.DOWN, action.LEFT, action.RIGHT}
	validMoves := []int{}
	for _, move := range moves {
		nxt := pos.NextPos(targetPos, move, mapData)
		//動けるなら
		if nxt != targetPos {
			validMoves = append(validMoves, move)
		}
	}
	return validMoves
}

//doBFS マップデータと始点を受け取り, 各点までの最短距離のマップを返す
func doBFS(mapData []string, startPos pos.Pos) map[pos.Pos]int {
	H := len(mapData)
	W := len(mapData[0])
	minDist := make(map[pos.Pos]int)
	dx := []int{1, 0, -1, 0}
	dy := []int{0, 1, 0, -1}
	que := []pos.Pos{startPos}
	minDist[startPos] = 0
	for len(que) > 0 {
		now := que[0]
		que = que[1:]
		for i := 0; i < 4; i++ {
			nx := now.X + dx[i]
			ny := now.Y + dy[i]
			nxt := pos.New(nx, ny)
			_, visited := minDist[nxt]
			if 0 > nx || nx >= W || 0 > ny || ny >= H || mapData[ny][nx] == '#' || visited {
				continue
			}
			que = append(que, nxt)
			minDist[nxt] = minDist[now] + 1
		}
	}
	return minDist
}
