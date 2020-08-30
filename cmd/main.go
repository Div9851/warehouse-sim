package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"sync"

	"github.com/Div9851/warehouse-sim/env"
	"github.com/Div9851/warehouse-sim/sim"
)

func main() {
	envPath := flag.String("env", "", "環境設定ファイルのパス")
	outputPath := flag.String("output", "", "結果を出力するディレクトリのパス")
	concurrent := flag.Int("concurrent", 1, "並行して実行するシミュレーションの数")
	total := flag.Int("total", 1, "実行するシミュレーションの数")
	verbose := flag.Bool("verbose", false, "シミュレーションの詳細を出力するかどうか")
	flag.Parse()
	env, err := env.Load(*envPath)
	if err != nil {
		panic(err)
	}
	for i := 0; i < *total; i += *concurrent {
		k := *concurrent
		if (*total - i) < k {
			k = *total - i
		}
		wg := &sync.WaitGroup{}
		for j := 0; j < k; j++ {
			wg.Add(1)
			go func() {
				seed := rand.Int63()
				sim := sim.New(env, seed)
				jsonData := sim.Do(*verbose)
				file, err := os.Create(filepath.Join(*outputPath, env.Name+"_"+fmt.Sprint(seed)+".json"))
				if err != nil {
					panic(err)
				}
				defer file.Close()
				file.WriteString(jsonData)
				wg.Done()
			}()
		}
		wg.Wait()
	}
}
