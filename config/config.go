package config

import (
	_ "fmt"
	"goutils/data"
)

type Configs struct {
	ConfigItems map[string]interface{}
}

type ConfigItem struct {
	Name  string
	Value string
}

func (c *Configs) LoadConfig(filepath string) error {

	results, err := data.ParseJSONFromFile(filepath)
	if err != nil {
		return err
	}
	c.ConfigItems = results.(map[string]interface{})
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
