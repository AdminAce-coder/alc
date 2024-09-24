package cmd

import (
	"bufio"
	"lightsailv2/aws"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var logFile *os.File

// termCmd represents the term command
var termCmd = &cobra.Command{
	Use:   "term",
	Short: "A brief description of your command",
	Long:  `终止实例.`,
	Run: func(cmd *cobra.Command, args []string) {
		// 获取区域
		region, _ := cmd.Flags().GetString("region")

		// 获取文件
		file, _ := cmd.Flags().GetString("file")
		// 读取文件内容
		f, err := os.Open(file)
		if err != nil {
			log.Fatalf("无法打开文件: %v", err)
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		total := 0 // 用于记录成功终止的实例数量
		for scanner.Scan() {
			// 获取实例ID
			instanceId := scanner.Text()
			// 终止实例
			err := aws.TerminateInstance(region, instanceId)
			if err != nil {
				log.Printf("终止实例 %s 失败: %v", instanceId, err)
			} else {
				log.Printf("终止实例 %s 成功", instanceId)
				total++ // 成功终止实例计数器增加
				log.Printf("成功终止的实例总数: %d", total)
			}
		}

		// 检查文件读取是否有错误
		if err := scanner.Err(); err != nil {
			log.Fatalf("读取文件时出错: %v", err)
		}

		// 打印成功终止的实例总数
		log.Printf("成功终止的实例总数: %d", total)
	},
}

func init() {
	rootCmd.AddCommand(termCmd)
	termCmd.Flags().StringP("region", "r", "", "指定区域")
	termCmd.Flags().StringP("file", "f", "", "指定停止实例的列表文件")

	// 设置日志输出到文件
	logFile = setupLogFile("terminate.log")
	log.SetOutput(logFile)
}

// setupLogFile创建并返回一个用于日志记录的文件
func setupLogFile(filePath string) *os.File {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("无法打开日志文件: %v", err)
	}
	return file
}
