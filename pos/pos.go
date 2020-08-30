package pos

import (
	"github.com/Div9851/warehouse-sim/action"
)

//Pos 座標を表す構造体
type Pos struct {
	X int `json:"x"`
	Y int `json:"y"`
}

//New 座標を受け取りPosを返す
func New(x int, y int) Pos {
	return Pos{X: x, Y: y}
}

//NextPos 現在の座標, 行動, マップデータを受け取り, 次の座標を返す
func NextPos(pos Pos, act int, mapData []string) Pos {
	H, W := len(mapData), len(mapData[0])
	nx, ny := pos.X, pos.Y
	switch act {
	case action.UP:
		ny--
	case action.DOWN:
		ny++
	case action.LEFT:
		nx--
	case action.RIGHT:
		nx++
	}
	if 0 > nx || nx >= W || 0 > ny || ny >= H || mapData[ny][nx] == '#' {
		nx, ny = pos.X, pos.Y
	}
	return New(nx, ny)
}
