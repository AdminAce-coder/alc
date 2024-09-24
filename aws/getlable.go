package aws

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lightsail"
)

// GetClientInit初始化并返回指定区域的kubectl客户端
func GetClientInit(region string) *lightsail.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		log.Printf("unable to load SDK config: %v", err)
		return nil
	}

	return lightsail.NewFromConfig(cfg)
}

// GetRegions检索所有可用区域
func GetRegions() []string {
	client := GetClientInit("us-east-1") // Use default region for fetching regions
	if client == nil {
		log.Printf("get client failed")
		return nil
	}

	resp, err := client.GetRegions(context.TODO(), &lightsail.GetRegionsInput{})
	if err != nil {
		log.Printf("get regions failed: %v", err)
		return nil
	}

	var regionNames []string
	for _, region := range resp.Regions {
		regionNames = append(regionNames, string(region.Name))
	}

	return regionNames
}

// InstanceInfo表示关于kubectl实例的信息
type InstanceInfo struct {
	Name       string
	Tags       map[string]string
	RegionName string
}

// GetInstances检索指定区域中的实例
func GetInstances(regionName string) []InstanceInfo {
	client := GetClientInit(regionName)
	if client == nil {
		log.Printf("get client failed")
		return nil
	}

	resp, err := client.GetInstances(context.TODO(), &lightsail.GetInstancesInput{})
	if err != nil {
		log.Printf("get instances failed: %v", err)
		return nil
	}

	var instanceList []InstanceInfo
	for _, instance := range resp.Instances {
		if instance.Name == nil {
			continue // Skip instances without a name
		}

		tags := make(map[string]string)
		for _, tag := range instance.Tags {
			if tag.Key != nil && tag.Value != nil {
				tags[*tag.Key] = *tag.Value
			}
		}

		instanceList = append(instanceList, InstanceInfo{
			Name:       *instance.Name,
			Tags:       tags,
			RegionName: regionName,
		})
	}

	return instanceList
}

// 停止实例
func StopInstance(regionName string, instanceName string) error {
	client := GetClientInit(regionName)
	if client == nil {
		log.Printf("get client failed")
		return nil
	}
	_, err := client.StopInstance(context.TODO(), &lightsail.StopInstanceInput{
		InstanceName: &instanceName,
	})
	if err != nil {
		log.Printf("stop instance failed: %v", err)
		return err
	}
	return nil
}

// 启动实例
func StartInstance(regionName string, instanceName string) error {
	client := GetClientInit(regionName)
	if client == nil {
		log.Printf("get client failed")
		return nil
	}
	_, err := client.StartInstance(context.TODO(), &lightsail.StartInstanceInput{
		InstanceName: &instanceName,
	})
	if err != nil {
		log.Printf("start instance failed: %v", err)
		return err
	}
	return nil
}

// 终止实例

func TerminateInstance(regionName string, instanceName string) error {
	client := GetClientInit(regionName)
	if client == nil {
		log.Printf("get client failed")
		return nil
	}
	_, err := client.DeleteInstance(context.TODO(), &lightsail.DeleteInstanceInput{
		InstanceName: &instanceName,
	})
	if err != nil {
		log.Printf("terminate instance failed: %v", err)
		return err
	}
	return nil
}
