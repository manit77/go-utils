package goutils

import (
	_ "fmt"
)

type Configs struct {
	ConfigItems map[string]interface{}
}

type ConfigItem struct {
	Name  string
	Value string
}

func (c *Configs) LoadConfig(filepath string) error {
	var err error
	c.ConfigItems, err = ParseJSON(filepath)
	if err != nil {
		return err
	}
	return nil
}

func (c *Configs) GetConfig(key string) string {
	if v, ok := c.ConfigItems[key]; ok {
		return v.(string)
	}
	return ""
}

func (c *Configs) GetConfigItem(key string) interface{} {
	if v, ok := c.ConfigItems[key]; ok {
		return v
	}
	return ""
}
