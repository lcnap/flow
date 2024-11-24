package config

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func Test_loadConfig(t *testing.T) {

	f, err := os.Open("../config.json")
	if err != nil {
		t.Fatal(err)
	}

	decoder := json.NewDecoder(f)

	var config Conf
	err = decoder.Decode(&config)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v", config)

}
