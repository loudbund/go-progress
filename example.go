package main

import (
	"flag"
	"fmt"
	"github.com/loudbund/go-progress/progress_v1"
	"time"
)

func main() {
	fmt.Println(progress_v1.GetProcessDir())

	var c *string // 运行方式
	c = flag.String("c", "run", "运行方式(status/run/start/stop/kill)")
	flag.Parse()
	progress_v1.Exec(fStart, *c, "log.txt")
}

func fStart() {
	T := time.NewTicker(time.Second)
	for {
		select {
		case <-T.C:
			fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
		}
	}
}
