/**
* @Author: kiosk
* @Mail: weijiaxiang007@foxmail.com
* @Date: 2020/6/25
**/
package global

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"zginx/utils"
	"zginx/zlog"
)

var ZConfig *ZginxEnv

type ZginxEnv struct {
	TCP 		TCPInfo
	Server 		ServerInfo
	Log 		LogInfo
}

type TCPInfo struct {
	Host 			string
	Port 			int
	Version			string
	MaxPacketSize 	uint32
	MaxConn			int
}

type ServerInfo struct {
	Name 				string
	Version 			string
	WorkerPoolSize  	uint32
	MaxWorkerTaskLen 	uint32
}

type LogInfo struct {
	LogFilePath		string
	DebugMode		bool
}


// Config Struct
type ZginxConfig struct {
	ConfFilePath	string
}

func (z ZginxConfig) ReadYaml() *ZginxEnv {
	data, _ := ioutil.ReadFile(z.ConfFilePath)

	var env = &ZginxEnv{}
	if err := yaml.Unmarshal(data, env); err != nil {
		panic("Read MFE Conf Error !")
	}

	zlog.Logger.Infof("Read Zginx Config From \"%s\" ", z.ConfFilePath)
	return env
}

func (z ZginxConfig) SetConfig() {
	if confFileExists, _ := utils.PathExists(z.ConfFilePath); confFileExists != true {
		zlog.Logger.Fatal("Config File ", z.ConfFilePath , " is not exist!!")
	}
	ZConfig = z.ReadYaml()

	Load()
}

func Load() {
	if ZConfig.Log.LogFilePath != "" {
		logfile, err := os.OpenFile(ZConfig.Log.LogFilePath,os.O_APPEND|os.O_CREATE|os.O_RDONLY,600)
		if err != nil {
			zlog.Logger.Fatal("Log File ",ZConfig.Log.LogFilePath, " open failed !! Error : ", err)
		}
		zlog.Logger.SetOutput(logfile)
	}

	if ZConfig.Log.DebugMode {
		// todo
	}

}





