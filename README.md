# go-progress
go-progress用来控制程序的自启动start，及停止stop

## 安装
go get github.com/loudbund/go-progress

## 引入
```golang
import "github.com/loudbund/go-progress/progress_v1"
```

## 关键代码
```golang
// 参数1: 程序入口函数
// 参数2：status/run/start/stop/kill
// 参数3: 程序输出日志写入文件路径
// 参数4：可选重写stop函数，用于安全stop程序，不重写stop参数则直接kill进程
func Exec(fStart func(), cType string, LogPath string, fStop ...func()) 
```

## 示例代码
```golang
package main

import (
    "flag"
    "fmt"
    "github.com/loudbund/go-progress/progress_v1"
    "time"
)

func main() {
    // 使用flag获取启动参数
    var c *string // 运行方式
    c = flag.String("c", "run", "运行方式(status/run/start/stop/kill)")
    flag.Parse()

    // 
    progress_v1.Exec(fStart, *c, "log.txt")
}

// 主程序入口移到这里了
func fStart() {
    // 模拟应用程序运行
    T := time.NewTicker(time.Second)
    for {
        select {
        case <-T.C:
            fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
        }
    }
}

```

## 示例的构建和使用
```shell script
go build xxx.go -o xxx
./xxx -c start
```