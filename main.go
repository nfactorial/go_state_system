package state_system

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	jsonFile, err := os.Open("test.json")
	if err != nil {
		fmt.Println(err)
	}

	rawData, _ := ioutil.ReadAll(jsonFile)

	var stateDesc StateTreeDesc

	json.Unmarshal(rawData, &stateDesc)

	fmt.Println("State Tree Name: " + stateDesc.Name)
	fmt.Println("Main: " + stateDesc.Main)

	for i := 0; i < len(stateDesc.States); i++ {
		fmt.Println("State: " + stateDesc.States[i].Name)
	}

	defer jsonFile.Close()
}
