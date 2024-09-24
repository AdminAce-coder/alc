/*
Copyright © 2024
*/
package cmd

import (
	"fmt"
	"strings"

	"lightsailv2/aws"

	"github.com/spf13/cobra"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "停止指定区域和标签的实例",
	Long:  `输入stop命令，停止指定区域和标签的实例。`,
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
				err := aws.StopInstance(region, instance.Name)
				if err != nil {
					fmt.Printf("停止实例 %s 失败: %v\n", instance.Name, err)
				} else {
					fmt.Printf("已停止实例 %s\n", instance.Name)
				}
			}
		}
	},
}

func parseTag(tag string) (string, string, error) {
	parts := strings.SplitN(tag, "=", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("标签应为 key=value 格式")
	}
	return parts[0], parts[1], nil
}

func init() {
	rootCmd.AddCommand(stopCmd)
	// 注册标志位
	stopCmd.Flags().StringP("region", "r", "", "指定区域")
	stopCmd.Flags().StringP("tag", "t", "", "指定标签，格式为 key=value")

}
