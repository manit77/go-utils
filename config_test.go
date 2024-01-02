package goutils

import (
	"fmt"
	"testing"
)

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestConfig(t *testing.T) {

	fmt.Printf("%v \n", "Begin TestConfig")

	var appConfig Configs
	err := appConfig.LoadConfig("config.json")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(appConfig.ConfigItems)

	connection_string := appConfig.GetConfigItem("connection_string").(string)
	if connection_string == "" {
		t.Fatal("connection_string is empty")
	}
	fmt.Printf("connection_string %v \n", connection_string)
	fmt.Printf("%v \n", "End TestConfig")
}
