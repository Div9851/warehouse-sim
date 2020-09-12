package main

import (
	"flag"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/Div9851/warehouse-sim/env"
	"github.com/Div9851/warehouse-sim/sim"
)

func main() {
	envPath := flag.String("env", "", "環境設定ファイルのパス")
	concurrent := flag.Int("concurrent", 1, "並行して実行するシミュレーションの数")
	total := flag.Int("total", 1, "実行するシミュレーションの数")
	verbose := flag.Bool("verbose", false, "シミュレーションの詳細を出力するかどうか")
	flag.Parse()
	env, err := env.Load(*envPath)
	if err != nil {
		panic(err)
	}
	var totalProcessTime float64
	var totalItems int
	var totalPickupCounts int
	var totalClearCounts int

	for done := 0; done < *total; done += *concurrent {
		startTime := time.Now()
		now := *concurrent
		if (*total - done) < now {
			now = *total - done
		}
		results := make([]*sim.Result, now)
		wg := &sync.WaitGroup{}
		for i := 0; i < now; i++ {
			wg.Add(1)
			go func(idx int) {
				seed := rand.Int63()
				sim := sim.New(env, seed)
				sim.Do(*verbose)
				results[idx] = sim.GetResult()
				wg.Done()
			}(i)
		}
		wg.Wait()
		for _, result := range results {
			totalItems += result.TotalItems
			for _, pickup := range result.PickupCounts {
				totalPickupCounts += pickup
			}
			for _, clear := range result.ClearCounts {
				totalClearCounts += clear
			}
		}
		endTime := time.Now()
		processTime := endTime.Sub(startTime).Seconds()
		fmt.Printf("process time %v sec\n", processTime)
		totalProcessTime += processTime
	}
	fmt.Printf("total process time %v sec\n", totalProcessTime)
	fmt.Printf("avg. items: %v\n", float64(totalItems)/float64(*total))
	fmt.Printf("avg. pickup: %v\n", float64(totalPickupCounts)/float64(*total))
	fmt.Printf("avg. clear: %v\n", float64(totalClearCounts)/float64(*total))
}
