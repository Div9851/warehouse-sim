package main

import (
	"flag"
	"fmt"

	"github.com/Div9851/warehouse-sim/env"
	"github.com/Div9851/warehouse-sim/sim"
)

func main() {
	envPath := flag.String("env", "", "環境設定ファイルのパス")
	seed := flag.Int64("seed", 123, "乱数のシード値")
	verbose := flag.Bool("verbose", false, "シミュレーションの詳細を出力するかどうか")
	flag.Parse()
	env, err := env.Load(*envPath)
	if err != nil {
		panic(err)
	}
	sim := sim.New(env, *seed)
	jsonData := sim.Do(*verbose)
	fmt.Println(jsonData)
}
