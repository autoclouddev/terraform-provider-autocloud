package blueprint_config_references

import (
	"encoding/json"
	"fmt"
	"sync"
)

type Data struct {
	values map[string]string
}

var (
	instance *Data
	once     sync.Once
)

func GetInstance() *Data {
	once.Do(func() {
		instance = &Data{
			values: make(map[string]string),
		}
	})
	return instance
}

func (d *Data) SetValue(key string, value string) {
	d.values[key] = value
}

func (d *Data) GetValue(key string) string {
	return d.values[key]
}

func (d *Data) ToString() string {
	jsonData, err := json.Marshal(d.values)
	if err != nil {
		fmt.Println("error reading references: ", err)
	}

	return string(jsonData)
}
