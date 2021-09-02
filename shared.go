package main

import (
	"fmt"
	"os"
)

func ExistsFile(fileName string) bool {
	pathUrl := fmt.Sprintf("configs/%s.yaml", fileName)
	_, err := os.Stat(pathUrl)

	if err != nil {
		return os.IsExist(err)
	}

	return true
}
