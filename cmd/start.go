/*
Copyright © 2024
*/
package cmd

import (
	"fmt"

	"lightsailv2/aws"

	"github.com/spf13/cobra"
)

// stopCmd represents the stop command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "启动指定区域和标签的实例",
	Long:  `输入start命令，停止指定区域和标签的实例。`,
	Run: func(cmd *cobra.Command, args []string) {
		region, _ := cmd.Flags().GetString("region")
		tag, _ := cmd.Flags().GetString("tag")

		// 检查是否提供了区域和标签
		if region == "" {
			fmt.Println("请使用 --region 或 -r 指定区域如：ca-central-1")
			return
		}
		if tag == "" {
			fmt.Println("请使用 --tag 或 -t 指定标签，格式为 key=value")
			return
		}

		// 解析标签为键和值
		k, v, err := parseTag(tag)
		if err != nil {
			fmt.Printf("标签格式错误: %v\n", err)
			return
		}

		// 获取实例信息
		instanceInfos := aws.GetInstances(region)
		if instanceInfos == nil {
			fmt.Println("未能获取实例信息")
			return
		} else { // 打印实例信息
			fmt.Printf("获取到 %d 个实例\n", len(instanceInfos))
			for _, instance := range instanceInfos {
				fmt.Printf("实例 %s, 标签 %v\n", instance.Name, instance.Tags)

			}
		}
		// 遍历实例，停止匹配的实例
		for _, instance := range instanceInfos {
			if instance.Tags[k] == v {
				err := aws.StartInstance(region, instance.Name)
				if err != nil {
					fmt.Printf("启动实例 %s 失败: %v\n", instance.Name, err)
				} else {
					fmt.Printf("已启动实例 %s\n", instance.Name)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	// 注册标志位
	startCmd.Flags().StringP("region", "r", "", "指定区域")
	startCmd.Flags().StringP("tag", "t", "", "指定标签，格式为 key=value")
}
