package aws

import (
	"fmt"
	"sync"
	"testing"
)

func TestGetLabel(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		regionName := GetRegions() // 在 Goroutine 中执行
		for _, name := range regionName {
			fmt.Printf("区域: %s\n", name)
			InstanceInfo := GetInstances(name) // 在 Goroutine 中执行
			fmt.Printf("实例信息: %v\n", InstanceInfo)

		}

	}()

	wg.Wait() // 等待所有 Goroutine 执行完成
}
