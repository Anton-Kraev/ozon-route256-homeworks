package config

import (
	"log"

	"github.com/eschao/config"
)

type WorkersConfig struct {
	WorkersN int `json:"num_workers"`
	TasksN   int `json:"num_tasks"`
	HashesN  int `json:"num_hashes"`
}

func ParseWorkersConfig(configPath string) *WorkersConfig {
	workersConfig := &WorkersConfig{}
	if err := config.ParseConfigFile(workersConfig, configPath+"workers.json"); err != nil {
		log.Fatalln(err)
	}

	if workersConfig.WorkersN <= 0 {
		workersConfig.WorkersN = 1
	}

	if workersConfig.TasksN <= 0 {
		workersConfig.TasksN = 1
	}

	if workersConfig.HashesN <= 0 {
		workersConfig.HashesN = 1
	}

	return workersConfig
}
