package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"

	logger "github.com/tal-tech/loggerX"
	"github.com/tal-tech/xtools/confutil"
)

var Conf map[string]PluginConfig

var updateChan chan string
var metadataChan chan string

var filePath string

func InitFromFile() {
	filePath = confutil.GetConf("ServiceConfig", "path")
	if len(filePath) == 0 {
		logger.W("Rpcx-plugin Conf", "get filepath fail")
		return
	}
	paths := strings.Split(filePath, "/")
	if len(paths) > 1 {
		dir := strings.Join(paths[0:len(paths)-1], "/")
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			logger.E("Rpcx-plugin Conf", "Could not create directory %s\n", dir, err)
		}
	}
	f, err := os.OpenFile(filePath, os.O_CREATE, 0664)
	if err != nil {
		logger.W("Rpcx-plugin Conf", "openfile err:%v", err)
		return
	}
	defer f.Close()
	conf, err := ioutil.ReadAll(f)
	if err != nil {
		logger.W("Rpcx-plugin Conf", "readfile err:%v", err)
		return
	}
	err = refreshConf([]byte(conf))
	if err != nil {
		logger.W("Rpcx-plugin Conf", "refreshConf err:%v", err)
	}
	return
}

func refreshConf(data []byte) error {
	err := json.Unmarshal(data, &Conf)
	if err != nil {
		logger.W("Rpcx-plugin Conf", "json err:%v", err)
		return err
	}
	logger.I("Rpcx-plugin Conf", "RefreshConf:%v", Conf)
	return nil
}

func SaveConfFile() error {
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0664)
	if err != nil {
		logger.W("Rpcx-plugin Conf", "openfile err:%v", err)
		return err
	}
	defer f.Close()
	data, err := json.MarshalIndent(Conf, "", "\t")
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	return err
}

func init() {
	updateChan = make(chan string, 0)
	metadataChan = make(chan string, 0)
	Conf = make(map[string]PluginConfig, 0)
}

func GetUpdateChan() chan string {
	return updateChan
}

func GetMetadataChan() chan string {
	return metadataChan
}
