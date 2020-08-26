package alicloud

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"

	"github.com/allanhung/ali-ecs-tag-update/pkg/log"
	"github.com/allanhung/ali-ecs-tag-update/pkg/types"
)

type QueryFlags struct {
	InstanceId   string
	InstanceName string
	PageSize     int
	Tag          types.ArgList
	RegName      types.ArgList
	NoTagKey     types.ArgList
	NoTagValue   types.ArgList
	Cron         string
}

func QueryECS(aliClient *AliClient, queryFlags QueryFlags) ([]ecs.Instance, error) {
	allInstances := make([]ecs.Instance, 0)
	remaining := 1
	pageNumber := 1
	pageSize := queryFlags.PageSize

	tags := []ecs.DescribeInstancesTag{}
	for _, tag := range queryFlags.Tag {
		tagList := strings.Split(tag, "=")
		instanceTag := ecs.DescribeInstancesTag{
			Key:   tagList[0],
			Value: tagList[1],
		}
		tags = append(tags, instanceTag)
	}
	for remaining > 0 {
		request := ecs.CreateDescribeInstancesRequest()
		request.Tag = &tags
		if queryFlags.InstanceName != "" {
			request.InstanceName = queryFlags.InstanceName
		}
		request.RegionId = aliClient.RegionID
		request.PageSize = requests.NewInteger(pageSize)
		request.PageNumber = requests.NewInteger(pageNumber)
		response, err := aliClient.EcsClient.DescribeInstances(request)
		if err != nil {
			return nil, fmt.Errorf("failed to get ECS information: %v", err)
		}
		for _, instance := range response.Instances.Instance {
			// name match
			nameMatch := (len(queryFlags.RegName) == 0)
			for _, regRule := range queryFlags.RegName {
				match, _ := regexp.MatchString(regRule, instance.InstanceName)
				if match {
					nameMatch = true
					log.Logger.Debugf("instance %s is include by name rule: %s", instance.InstanceName, regRule)
					break
				}
			}

			if nameMatch {
				// tag key match
				noTagKeyMatch := true
				for _, regRule := range queryFlags.NoTagKey {
					r, _ := regexp.Compile(regRule)
					for _, tag := range instance.Tags.Tag {
						if r.MatchString(tag.TagKey) {
							noTagKeyMatch = false
							log.Logger.Debugf("instance %s is exclude by no tag key rule: %s", instance.InstanceName, regRule)
							break
						}
					}
					if !noTagKeyMatch {
						break
					}
				}
				if noTagKeyMatch {
					// tag value match
					noTagValueMatch := true
					for _, regRule := range queryFlags.NoTagValue {
						r, _ := regexp.Compile(regRule)
						for _, tag := range instance.Tags.Tag {
							if r.MatchString(tag.TagValue) {
								noTagValueMatch = false
								log.Logger.Debugf("instance %s is exclude by no tag value rule: %s", instance.InstanceName, regRule)
								break
							}
						}
						if !noTagValueMatch {
							break
						}
					}
					if noTagValueMatch {
						allInstances = append(allInstances, instance)
					}
				}
			}
		}
		remaining = response.TotalCount - pageNumber*pageSize
		pageNumber++
	}
	return allInstances, nil
}

func AddInstanceTags(aliClient *AliClient, ecsInstance ecs.Instance, ecsTags []ecs.AddTagsTag) error {
	request := ecs.CreateAddTagsRequest()
	request.Scheme = "https"

	request.ResourceType = "instance"
	request.ResourceId = ecsInstance.InstanceId
	request.Tag = &ecsTags

	response, err := aliClient.EcsClient.AddTags(request)
	if err != nil {
		return err
	}
	log.Logger.Infof("instance: %s (%s) tags added", ecsInstance.InstanceId, ecsInstance.InstanceName)
	log.Logger.Debugf("response: %v", response)
	return nil
}

func QueryVpc(aliClient *AliClient, pageSize int) ([]ecs.Vpc, error) {
	remaining := 1
	pageNumber := 1

	allVpcs := make([]ecs.Vpc, 0)

	for remaining > 0 {
		request := ecs.CreateDescribeVpcsRequest()
		request.RegionId = aliClient.RegionID
		request.PageSize = requests.NewInteger(pageSize)
		request.PageNumber = requests.NewInteger(pageNumber)
		response, err := aliClient.EcsClient.DescribeVpcs(request)
		if err != nil {
			return allVpcs, err
		}
		for _, Vpc := range response.Vpcs.Vpc {
			allVpcs = append(allVpcs, Vpc)
		}
		remaining = response.TotalCount - pageNumber*pageSize
		pageNumber++
	}

	return allVpcs, nil
}
