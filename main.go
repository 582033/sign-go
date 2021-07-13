package main

import (
	"sign/libs"
	"sync"
)

func main() {
	configList := libs.LoadConf("/Users/yjiang/go/src/sign/conf.yml")

	var wg sync.WaitGroup
	for _, config := range configList {
		wg.Add(1)
		go libs.Sign(config, func() {
			defer wg.Done()
		})
	}
	wg.Wait()
}
