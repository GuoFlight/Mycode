package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// PDFCracker PDF密码破解器（使用qpdf命令行工具）
type PDFCracker struct {
	filePath   string
	maxWorkers int
	found      int32
	password   string
	foundMutex sync.Mutex
	wg         sync.WaitGroup
	stopCh     chan struct{}
	attempted  int64
}

// NewPDFCracker 创建新的破解器
func NewPDFCracker(filePath string) *PDFCracker {
	return &PDFCracker{
		filePath:   filePath,
		maxWorkers: runtime.NumCPU() * 2,
		stopCh:     make(chan struct{}),
	}
}

// tryPassword 尝试单个密码 - 使用qpdf命令
func (c *PDFCracker) tryPassword(password string) bool {
	if atomic.LoadInt32(&c.found) == 1 {
		return true
	}

	// 构造qpdf命令
	// qpdf --password=xxx --decrypt input.pdf output.pdf
	// 使用NUL (Windows) 或 /dev/null (Linux/Mac) 作为输出
	nullDevice := "/dev/null"
	if runtime.GOOS == "windows" {
		nullDevice = "NUL"
	}

	cmd := exec.Command("qpdf",
		"--password="+password,
		"--decrypt",
		c.filePath,
		nullDevice,
	)

	// 运行命令，不输出任何信息
	err := cmd.Run()
	if err != nil {
		return false
	}

	// 密码正确！
	if atomic.CompareAndSwapInt32(&c.found, 0, 1) {
		c.foundMutex.Lock()
		c.password = password
		c.foundMutex.Unlock()
		close(c.stopCh)
		return true
	}
	return false
}

// worker 工作协程
func (c *PDFCracker) worker(passwords <-chan string) {
	defer c.wg.Done()
	for password := range passwords {
		select {
		case <-c.stopCh:
			return
		default:
			if c.tryPassword(password) {
				return
			}
			atomic.AddInt64(&c.attempted, 1)
		}
	}
}

// progressMonitor 进度监控
func (c *PDFCracker) progressMonitor(total int, stop <-chan struct{}) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	var lastAttempted int64
	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
			current := atomic.LoadInt64(&c.attempted)
			speed := current - lastAttempted
			lastAttempted = current
			percentage := float64(current) / float64(total) * 100
			fmt.Printf("\r🔍 进度: %.2f%% (%d/%d) | 速度: %d 密码/秒",
				percentage, current, total, speed)
		}
	}
}

// BruteForce 执行暴力破解
func (c *PDFCracker) BruteForce(passwordList []string) (string, error) {
	totalPasswords := len(passwordList)
	if totalPasswords == 0 {
		return "", fmt.Errorf("密码列表为空")
	}

	// 检查qpdf是否安装
	if _, err := exec.LookPath("qpdf"); err != nil {
		return "", fmt.Errorf("qpdf未安装，请先安装: sudo apt-get install qpdf (Linux) 或 brew install qpdf (Mac)")
	}

	fmt.Printf("📁 PDF文件: %s\n", c.filePath)
	fmt.Printf("📊 密码总数: %d\n", totalPasswords)
	fmt.Printf("💻 工作线程: %d\n", c.maxWorkers)
	fmt.Printf("⏰ 开始时间: %s\n\n", time.Now().Format("15:04:05"))

	startTime := time.Now()
	passwords := make(chan string, 10000)
	stopMonitor := make(chan struct{})
	go c.progressMonitor(totalPasswords, stopMonitor)

	for i := 0; i < c.maxWorkers; i++ {
		c.wg.Add(1)
		go c.worker(passwords)
	}

	go func() {
		for _, pwd := range passwordList {
			select {
			case <-c.stopCh:
				break
			default:
				passwords <- pwd
			}
		}
		close(passwords)
	}()

	c.wg.Wait()
	close(stopMonitor)
	elapsed := time.Since(startTime)
	fmt.Println()

	if atomic.LoadInt32(&c.found) == 1 {
		c.foundMutex.Lock()
		pwd := c.password
		c.foundMutex.Unlock()
		fmt.Printf("\n✅ 密码破解成功！\n")
		fmt.Printf("🔑 密码: %s\n", pwd)
		fmt.Printf("⏱️  耗时: %v\n", elapsed)
		fmt.Printf("📊 尝试密码: %d 个\n", atomic.LoadInt64(&c.attempted))
		return pwd, nil
	}

	fmt.Printf("\n❌ 密码破解失败\n")
	fmt.Printf("⏱️  耗时: %v\n", elapsed)
	fmt.Printf("📊 已尝试: %d/%d 个密码\n", atomic.LoadInt64(&c.attempted), totalPasswords)
	return "", fmt.Errorf("password not found in dictionary")
}

// LoadDictionary 从文件加载密码字典
func LoadDictionary(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("无法打开字典文件: %v", err)
	}
	defer file.Close()

	var passwords []string
	scanner := bufio.NewScanner(file)
	buf := make([]byte, 0, 1024*1024)
	scanner.Buffer(buf, 1024*1024)

	for scanner.Scan() {
		password := strings.TrimSpace(scanner.Text())
		if password == "" || strings.HasPrefix(password, "#") {
			continue
		}
		passwords = append(passwords, password)
	}
	return passwords, scanner.Err()
}

func printUsage() {
	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║     PDF密码暴力破解工具 (Go + qpdf)                       ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("前提: 请先安装 qpdf")
	fmt.Println("  Linux: sudo apt-get install qpdf")
	fmt.Println("  Mac:   brew install qpdf")
	fmt.Println("  Windows: 下载 https://github.com/qpdf/qpdf/releases")
	fmt.Println()
	fmt.Println("用法:")
	fmt.Println("  go run main.go <PDF文件路径> <字典文件路径>")
	fmt.Println()
	fmt.Println("示例:")
	fmt.Println("  go run main.go encrypted.pdf passwords.txt")
	fmt.Println()
	fmt.Println("⚠️  法律声明:")
	fmt.Println("  请确保您拥有该PDF文件的合法所有权或已获得授权！")
	os.Exit(1)
}

func main() {
	if len(os.Args) < 3 {
		printUsage()
	}

	pdfPath := os.Args[1]
	dictPath := os.Args[2]

	if _, err := os.Stat(pdfPath); os.IsNotExist(err) {
		fmt.Printf("❌ 错误: PDF文件不存在: %s\n", pdfPath)
		os.Exit(1)
	}
	if _, err := os.Stat(dictPath); os.IsNotExist(err) {
		fmt.Printf("❌ 错误: 字典文件不存在: %s\n", dictPath)
		os.Exit(1)
	}

	fmt.Println("📖 正在加载字典文件...")
	passwords, err := LoadDictionary(dictPath)
	if err != nil {
		fmt.Printf("❌ 加载字典失败: %v\n", err)
		os.Exit(1)
	}
	if len(passwords) == 0 {
		fmt.Println("❌ 错误: 字典文件中没有有效的密码")
		os.Exit(1)
	}
	fmt.Printf("✅ 成功加载 %d 个密码\n\n", len(passwords))

	cracker := NewPDFCracker(pdfPath)
	password, err := cracker.BruteForce(passwords)
	if err != nil {
		fmt.Printf("\n❌ 破解失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("\n🎉 恭喜！密码破解成功！密码: %s\n", password)
}

