package config

import "sync"

type SmartAssistantModeType int

const (
	SmartAssistantModeLocal SmartAssistantModeType = iota + 1
	SmartAssistantModeCloud
)

var (
	mode SmartAssistantModeType = SmartAssistantModeLocal
	once sync.Once
)

func SetSmartAssistantMode(m SmartAssistantModeType) {
	once.Do(func() {
		mode = m
	})
}

func GetSmartAssistantMode() SmartAssistantModeType {
	once.Do(func() {
		mode = SmartAssistantModeLocal
	})

	return mode
}

func IsCloudSA() bool {
	return GetSmartAssistantMode() == SmartAssistantModeCloud
}
