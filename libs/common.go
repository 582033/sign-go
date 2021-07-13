package libs

import (
	"os"
)

func CheckFile(file string) bool {
	isExist := true
	_, err := os.Stat(file)
	if err != nil {
		isExist = false
	}
	return isExist
}
