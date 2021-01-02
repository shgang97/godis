package config

import (
	"bufio"
	"godis/lib/logger"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
)

/*
@author: shg
@since: 2023/3/3 2:42 AM
@mail: shgang97@163.com
*/

type ServerProperties struct {
	Bind           string `cfg:"bind"`
	Port           int    `cfg:"port"`
	AppendOnly     bool   `cfg:"appendonly"`
	AppendFilename string `cfg:"appendfilename"`
	MaxClients     int    `cfg:"maxclients"`
}

// Properties 全剧配置
var Properties *ServerProperties

func init() {
	// 默认配置
	Properties = &ServerProperties{
		Bind:       "127.0.0.1",
		Port:       9379,
		AppendOnly: false,
	}
}

func SetupConfig(configFilename string) {
	file, err := os.Open(configFilename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	Properties = parse(file)
}

func parse(file io.Reader) *ServerProperties {
	serverProps := &ServerProperties{}

	configMap := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 && strings.TrimLeft(line, " ")[0] == '#' {
			continue
		}
		pivot := strings.IndexAny(line, " ")
		if pivot > 0 && pivot < len(line)-1 {
			key := line[0:pivot]
			value := strings.Trim(line[pivot+1:], " ")
			configMap[strings.ToLower(key)] = value
		}
	}
	if err := scanner.Err(); err != nil {
		logger.Fatal(err)
	}

	t := reflect.TypeOf(serverProps)
	v := reflect.ValueOf(serverProps)
	n := t.Elem().NumField()
	for i := 0; i < n; i++ {
		field := t.Elem().Field(i)
		fieldVal := v.Elem().Field(i)
		key, ok := field.Tag.Lookup("cfg")
		if !ok || strings.TrimLeft(key, " ") == "" {
			key = field.Name
		}
		value, ok := configMap[strings.ToLower(key)]
		if ok {
			switch field.Type.Kind() {
			case reflect.String:
				fieldVal.SetString(value)
			case reflect.Int:
				intValue, err := strconv.ParseInt(value, 10, 64)
				if err == nil {
					fieldVal.SetInt(intValue)
				}
			case reflect.Bool:
				boolValue := "yes" == value
				fieldVal.SetBool(boolValue)
			case reflect.Slice:
				if field.Type.Elem().Kind() == reflect.String {
					s := strings.Split(value, ",")
					fieldVal.Set(reflect.ValueOf(s))
				}

			}
		}
	}
	return serverProps
}
