package utils

import (
	"fmt"
	"github.com/unknwon/goconfig"
	"reflect"
	"strconv"
	"strings"
)

type OpenJudge struct {
	Host     string
	Port     string
	DataPath string
}

type MysqlConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	DBName   string
}

type Path struct {
	Data string
}

type Config struct {
	Mysql     MysqlConfig
	OpenJudge OpenJudge
	Path      Path
}

var GlobalConfig *Config

func initSectionConfig(sr *reflect.Value, dict map[string]string) {
	for i := 0; i < sr.NumField(); i++ {
		Type := sr.Type().Field(i)
		value := sr.Field(i)
		name := strings.ToLower(Type.Name)
		switch value.Type().Kind() {
		case reflect.String:
			{
				value.SetString(dict[name])
			}
		case reflect.Int:
			{
				number, err := strconv.Atoi(dict[name])
				if err != nil {
					panic(err)
				}
				value.SetInt(int64(number))
			}
		}
		fmt.Println(name, dict[name])
	}
}

func initGlobalConfig(sr *reflect.Value, file *goconfig.ConfigFile) {
	for i := 0; i < sr.NumField(); i++ {
		Type := sr.Type().Field(i)
		value := sr.Field(i)
		name := strings.ToLower(Type.Name)
		if value.Type().Kind() == reflect.Struct {
			dict, err := file.GetSection(strings.ToLower(name))
			if err != nil {
				panic(err)
			}
			fmt.Println("struct")
			initSectionConfig(&value, dict)
		}
	}
}

func init() {
	file, err := goconfig.LoadConfigFile("openjudge-server.conf")
	if err != nil {
		panic(err)
	}
	GlobalConfig = &Config{}
	sr := reflect.ValueOf(GlobalConfig).Elem()
	initGlobalConfig(&sr, file)
}
