package progress_v1

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// 全局变量1： 日志文件路径
var logPath string // 日志文件路径

// 启动函数
// 示例：
//  var c *string // 运行方式
//  c = flag.String("c", "run", "运行方式(status/run/start/stop/kill)")
//  flag.Parse()
//  progress_v1.Exec(fStart, *c, "logs/aaa.txt")
func Exec(fStart func(), cType string, LogPath string, fStop ...func()) {
	logPath = LogPath

	// 参数个数校验
	switch cType {
	case "status":
		runPids := GetMyProcessIdList()
		if len(runPids) > 0 {
			fmt.Println("运行中:", runPids)
		} else {
			fmt.Println("程序未启动")
		}

	case "run":
		fStart()

	case "start":
		startProcess(fStart)

	case "stop":
		if len(fStop) > 0 {
			fStop[0]()
		} else {
			killProcess()
		}

	case "kill":
		killProcess()

	default:
		fStart()
	}

}

/**-------------------------
// 名称：获取应用程序所在目录
***-----------------------*/
func GetProcessDir() string {

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	return dir

}

/**-------------------------
// 名称：查找自身程序运行的进程id
// 返回：相关进程id数组
***-----------------------*/
func GetMyProcessIdList() []string {

	var (
		fields  []string // 进程id列表
		absFile string   // 程序文件绝对路径
		err     error    // 错误
	)

	// 1、获取程序文件地址
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return []string{}
	}
	absFile, err = filepath.Abs(file)
	if err != nil {
		log.Error("The eacePath failed: %s\n", err.Error())
	}

	// 2、扫描相关进程
	output, _ := exec.Command("/bin/sh", "-c", `ps -ef|awk '{if($8=="`+absFile+`"){print $2}}'`).Output()
	data := string(output)
	fields = strings.Fields(data)

	// 3、过滤掉自身进程
	pid := fmt.Sprintf("%d", os.Getpid())
	for k, v := range fields {
		if v == pid {
			fields = append(fields[:k], fields[k+1:]...)
			break
		}
	}

	// 4、返回进程列表
	return fields

}

/**-------------------------
// 名称：kill掉自身进程
***-----------------------*/
func killProcess() int {

	// 获取进程id列表
	idList := GetMyProcessIdList()

	// 查找到了
	if len(idList) > 0 {
		fmt.Println("kill ", len(idList), "个进程")
		// 逐个kill掉
		for _, pid := range idList {
			_, _ = exec.Command("/bin/sh", "-c", `kill -9 `+pid).Output()
		}
	} else {
		fmt.Println("没有运行的进程")
	}

	return len(idList)

}

/**-------------------------
// 名称：启动自身进程
***-----------------------*/
func startProcess(run func()) bool {

	// 获取进程id列表
	idList := GetMyProcessIdList()

	// 有程序运行
	if len(idList) > 0 {
		fmt.Println("进程已经在运行了", idList)

		return false

	} else {
		// 判 断当其是否是子进程，当父进程return之后，子进程会被 系统1 号进程接管
		if os.Getppid() != 1 {
			// 日志文件所在目录检查
			logDir := filepath.Dir(logPath)

			if _, err := os.Stat(logDir); os.IsNotExist(err) {
				log.Panic("启动失败，日志文件目录不存在:", logDir)
			}
			f, err := os.OpenFile(logPath, os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_APPEND, 0644)
			if err != nil {
				log.Panic("启动失败，创建日志文件失败", logPath)
			}

			// 将命令行参数中执行文件路径转换成可用路径
			appPath, _ := filepath.Abs(os.Args[0])
			cmd := exec.Command(appPath, os.Args[1:]...)

			// 将其他命令传入生成出的进程
			cmd.Stdin = os.Stdin // 给新进程设置文件描述符，可以重定向到文件中
			cmd.Stdout = f
			cmd.Stderr = f
			_ = cmd.Start() // 开始执行新进程，不等待新进程退出

			fmt.Println("守护进程启动成功,进程id:", cmd.Process.Pid)
			return true

		}

		// /////////////////////////////////////////////正式运行
		// fmt.Println("开始运行，")
		fmt.Println("开始运行,进程id:", os.Getpid())
		run()

	}

	return false

}
