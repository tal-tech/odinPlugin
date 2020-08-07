package config

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/smallnest/rpcx/server"
	logger "github.com/tal-tech/loggerX"
)

func InitConfManager(port string) {
	s := server.NewServer()
	s.RegisterName("ConfManager", new(ConfManager), "")
	portInt, err := strconv.Atoi(port)
	if err != nil {
		return
	}
	managePort := strconv.Itoa(portInt + 10000)
	go s.Serve("tcp", ":"+managePort)
}

type ConfManager struct{}

type NameArgs struct {
	ServerName string
}

type NilArgs struct{}

type CommonArgs struct {
	ServerName string
	Data       []byte `json:"data"`
}

func (c *ConfManager) GetConf(ctx context.Context, args *NameArgs, reply *CommonArgs) error {
	conf := Conf[args.ServerName]
	data, err := json.MarshalIndent(conf, "", "\t")
	if err != nil {
		logger.W("ConfManager", "json err:%v", err)
		return err
	}
	reply.Data = data
	return nil
}

func (c *ConfManager) UpdateConf(ctx context.Context, args *CommonArgs, reply *NilArgs) error {
	data, _ := json.Marshal(Conf[args.ServerName])
	logger.I("ConfManager", "updateconf,server:%s,backup:%s", args.ServerName, string(data))
	var conf PluginConfig
	err := json.Unmarshal(args.Data, &conf)
	if err != nil {
		logger.W("ConfManager", "updateconf unmarshal err:%v", err)
		return err
	}
	Conf[args.ServerName] = conf
	updateChan <- args.ServerName
	return nil
}

func (c *ConfManager) ResetConf(ctx context.Context, args *CommonArgs, reply *NilArgs) error {
	var conf PluginConfig
	err := json.Unmarshal(args.Data, &conf)
	if err != nil {
		logger.W("ConfManager", "saveconf unmarshal err:%v", err)
		return err
	}
	Conf[args.ServerName] = conf
	err = SaveConfFile()
	if err != nil {
		logger.W("ConfManager", "saveconf to file err:%v", err)
		return err
	}
	logger.I("ConfManager", "saveconf,server:%s,backup:%s", args.ServerName, string(args.Data))
	return nil
}

func (c *ConfManager) UpdateMetadata(ctx context.Context, args *CommonArgs, reply *NilArgs) error {
	data, _ := json.Marshal(Conf[args.ServerName])
	logger.I("ConfManager", "updatemetadata,server:%s,backup:%s", args.ServerName, string(data))
	md := make(map[string]string, 0)
	items := strings.Split(string(args.Data), "&")
	for _, item := range items {
		kvs := strings.Split(item, "=")
		if len(kvs) == 2 {
			md[kvs[0]] = kvs[1]
		}
	}
	conf := Conf[args.ServerName]
	conf.MDConfigs = md
	Conf[args.ServerName] = conf
	metadataChan <- args.ServerName
	return nil
}

func (c *ConfManager) ActivateService(ctx context.Context, args *CommonArgs, reply *NilArgs) error {
	data, _ := json.Marshal(Conf[args.ServerName])
	logger.I("ConfManager", "updatemetadata,server:%s,backup:%s", args.ServerName, string(data))
	conf := Conf[args.ServerName]
	/*
		if conf.MDConfigs != nil {
			conf.MDConfigs["state"] = "active"
		} else {
	*/
	md := make(map[string]string, 0)
	md["state"] = "active"
	conf.MDConfigs = md
	//}
	Conf[args.ServerName] = conf
	metadataChan <- args.ServerName
	return nil
}

func (c *ConfManager) DeactivateService(ctx context.Context, args *CommonArgs, reply *NilArgs) error {
	data, _ := json.Marshal(Conf[args.ServerName])
	logger.I("ConfManager", "updatemetadata,server:%s,backup:%s", args.ServerName, string(data))
	conf := Conf[args.ServerName]
	/*
		if conf.MDConfigs != nil {
			conf.MDConfigs["state"] = "inactive"
		} else {
	*/
	md := make(map[string]string, 0)
	md["state"] = "inactive"
	conf.MDConfigs = md
	//}
	Conf[args.ServerName] = conf
	metadataChan <- args.ServerName
	return nil
}
