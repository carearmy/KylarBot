package main

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type KConfig struct {
	Roles   []KRole   `yaml:"roles"`
	Members []KMember `yaml:"members"`
}

func dummyConfig() {
	config := KConfig{
		Roles:   []KRole{{ID: "general", Permissions: []string{"chat.write"}}},
		Members: []KMember{{ID: "2t23t", Permissions: []string{"chat.write"}, Roles: []string{"general"}}},
	}

	var res []byte
	var err error

	if res, err = yaml.Marshal(config); err != nil {
		fmt.Println(err.Error())
	}

	ioutil.WriteFile("config.yml", res, 0644)
}

func loadConfig() {

}

func saveConfig() {

}
