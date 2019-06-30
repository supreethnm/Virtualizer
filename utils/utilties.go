package utils

import (
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/clbanning/mxj"
)

func ToJsonBytes(data []byte) (jsonBytes []byte, err error) {
	// []byte to Map
	mapVal, err := mxj.NewMapXml(data)
	if err != nil {
		return nil, err
	}

	// Map to JSON
	jsonBytes, err = mapVal.Json()
	//jsonBytes, err = json.Marshal(tempXmlMap.)
	if err != nil {
		return nil, err
	}

	logrus.WithFields(logrus.Fields{"data": string(jsonBytes)}).Debug("data converted to JSON")
	return
}

func ToXmlBytes(data []byte) (xmlBytes []byte, err error) {
	// []byte to Map
	mapVal, err := mxj.NewMapJson(data)
	if err != nil {
		return nil, err
	}

	// Map to JSON
	xmlBytes, err = mapVal.Xml()
	if err != nil {
		return nil, err
	}

	logrus.WithFields(logrus.Fields{"data": string(xmlBytes)}).Debug("data converted to XML")
	return
}

func BytesToString(data []byte) string {
	return string(data[:])
}

func AddDelay(delay time.Duration, ch chan<- bool) {
	time.Sleep(delay * time.Second)
	ch <- true
}

func GetTempDir() string {
	return os.TempDir()
}

func GetConfigFilePath() string {
	configFilePath := GetTempDir() + "/virtualizer/config.toml"
	logrus.WithFields(logrus.Fields{"config file": configFilePath}).Debug("config.toml file location")
	return configFilePath
}
