package config

import (
	"log"
	"os"

	utils "security_chat_app/internal/utils/log"

	"gopkg.in/go-ini/ini.v1"
)

type ConfigList struct {
	Port           string
	LogFile        string
	Static         string
	ServiceKeyPath string
	ProjectId      string
	StorageBucket  string
}

var Config ConfigList

func init() {
	LoadConfig()
	utils.LoggingSettings(Config.LogFile)
}

func LoadConfig() {
	localConfig, err := ini.Load("config.local.ini")
	if err != nil {
		log.Println("config.local.iniが見つかりません。config.iniを使用します。")
		localConfig = nil
	}

	// config.iniを読み込む
	defaultConfig, err := ini.Load("config.ini")
	if err != nil {
		log.Fatalln("config.iniが見つかりません:", err)
	}

	Config = ConfigList{
		Port:           "8080",
		LogFile:        "",
		Static:         "",
		ServiceKeyPath: "",
		ProjectId:      "",
		StorageBucket:  "",
	}

	// config.local.iniから値を読み込む（存在する場合）
	if localConfig != nil {
		loadConfigValues(localConfig, &Config)
	}

	// 不足している値をconfig.iniから補完
	loadConfigValues(defaultConfig, &Config)

	// 必須項目の検証
	validateConfig(&Config)
}

// 設定ファイルから値を読み込む
func loadConfigValues(cfg *ini.File, config *ConfigList) {
	if port := cfg.Section("web").Key("port").String(); port != "" {
		config.Port = port
	}
	if logFile := cfg.Section("web").Key("logfile").String(); logFile != "" {
		config.LogFile = logFile
	}
	if static := cfg.Section("web").Key("static").String(); static != "" {
		config.Static = static
	}
	if serviceAccountKey := cfg.Section("firebase").Key("serviceKeyPath").String(); serviceAccountKey != "" {
		config.ServiceKeyPath = serviceAccountKey
	}
	if projectId := cfg.Section("firebase").Key("projectId").String(); projectId != "" {
		config.ProjectId = projectId
	}
	if storageBucket := cfg.Section("firebase").Key("storageBucket").String(); storageBucket != "" {
		config.StorageBucket = storageBucket
	}
}

// 設定値の検証
func validateConfig(config *ConfigList) {
	// 必須項目の検証
	if config.LogFile == "" {
		log.Println("警告: LogFileが設定されていません。デフォルト値を使用します。")
		config.LogFile = "debug.log"
	}
	if config.Static == "" {
		log.Println("警告: Staticが設定されていません。デフォルト値を使用します。")
		config.Static = "app/views"
	}
	if config.ServiceKeyPath == "" {
		log.Fatalln("エラー: serviceKeyPathが設定されていません。")
	}
	if config.ProjectId == "" {
		log.Fatalln("エラー: projectIdが設定されていません。")
	}
	if config.StorageBucket == "" {
		log.Fatalln("エラー: storageBucketが設定されていません。")
	}

	// ファイルの存在確認
	if _, err := os.Stat(config.ServiceKeyPath); os.IsNotExist(err) {
		log.Fatalf("エラー: serviceKeyPathファイルが見つかりません: %s", config.ServiceKeyPath)
	}
}
