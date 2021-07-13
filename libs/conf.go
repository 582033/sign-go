package libs

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

func configFieldGet(config *viper.Viper, taskName string, field string) string {
	return fmt.Sprintf("%s", config.Get(taskName+"."+field))
}

func LoadConf(confFile string) []RequestInfo {
	config := viper.New()
	config.SetConfigFile(confFile)
	config.SetConfigType("yaml") //设置文件的类型
	//尝试进行配置读取
	if err := config.ReadInConfig(); err != nil {
		log.Println(err)
	}

	var configList []RequestInfo
	for taskName, _ := range config.AllSettings() {
		CookieFile := configFieldGet(config, taskName, "CookieFile")
		Method := configFieldGet(config, taskName, "Method")
		Url := configFieldGet(config, taskName, "Url")

		configList = append(configList, RequestInfo{
			CookieFile: CookieFile,
			Headers:    config.Get(taskName + ".Headers"),
			Method:     Method,
			Url:        Url,
			Data:       config.Get(taskName + ".Data"),
		})
	}

	return configList
}
