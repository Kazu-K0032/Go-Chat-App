package config

import (
	"log"
	"os"

	utils "security_chat_app/internal/utils/log"

	"gopkg.in/go-ini/ini.v1"
)

type ConfigList struct {
	Port            string
	LogFile         string
	Static          string
	DefaultIconDir  string
	ServiceKeyPath  string
	ProjectId       string
	StorageBucket   string
}

var Config ConfigList

func init() {
	LoadConfig()
	utils.LoggingSettings(Config.LogFile)
}

func LoadConfig() {
	Config = ConfigList{
		Port:            "8080",
		LogFile:         "",
		Static:          "",
		DefaultIconDir:  "",
		ServiceKeyPath:  "",
		ProjectId:       "",
		StorageBucket:   "",
	}

	// 環境変数から読み込む（優先度が最も高い）
	loadConfigFromEnv(&Config)

	localConfig, err := ini.Load("config.local.ini")
	if err != nil {
		log.Println("config.local.iniが見つかりません。config.iniを使用します。")
		localConfig = nil
	}

	// config.iniを読み込む
	defaultConfig, err := ini.Load("config.ini")
	if err != nil {
		log.Println("config.iniが見つかりません。環境変数のみを使用します。")
		defaultConfig = nil
	}

	// config.local.iniから値を読み込む（存在する場合、環境変数を上書きしない）
	if localConfig != nil {
		loadConfigValues(localConfig, &Config)
	}

	// 不足している値をconfig.iniから補完（環境変数を上書きしない）
	if defaultConfig != nil {
		loadConfigValues(defaultConfig, &Config)
	}

	// 必須項目の検証
	validateConfig(&Config)
}

// 環境変数から設定を読み込む（優先度が最も高い）
func loadConfigFromEnv(config *ConfigList) {
	if port := os.Getenv("PORT"); port != "" {
		config.Port = port
	}
	if logFile := os.Getenv("LOG_FILE"); logFile != "" {
		config.LogFile = logFile
	}
	if static := os.Getenv("STATIC_DIR"); static != "" {
		config.Static = static
	}
	if defaultIconDir := os.Getenv("DEFAULT_ICON_DIR"); defaultIconDir != "" {
		config.DefaultIconDir = defaultIconDir
	}
	if serviceKeyPath := os.Getenv("SERVICE_KEY_PATH"); serviceKeyPath != "" {
		config.ServiceKeyPath = serviceKeyPath
	}
	// Firebase認証情報を環境変数から直接読み込む場合（JSON文字列）
	if serviceKeyJSON := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS_JSON"); serviceKeyJSON != "" {
		// 一時ファイルに書き込む
		tmpFile := "/tmp/serviceAccountKey.json"
		if err := os.WriteFile(tmpFile, []byte(serviceKeyJSON), 0600); err != nil {
			log.Printf("警告: 一時ファイルの作成に失敗: %v", err)
		} else {
			config.ServiceKeyPath = tmpFile
		}
	}
	if projectId := os.Getenv("PROJECT_ID"); projectId != "" {
		config.ProjectId = projectId
	}
	if storageBucket := os.Getenv("STORAGE_BUCKET"); storageBucket != "" {
		config.StorageBucket = storageBucket
	}
}

// 設定ファイルから値を読み込む（環境変数で設定されていない場合のみ）
func loadConfigValues(cfg *ini.File, config *ConfigList) {
	if config.Port == "8080" || config.Port == "" {
		if port := cfg.Section("web").Key("port").String(); port != "" {
			config.Port = port
		}
	}
	if config.LogFile == "" {
		if logFile := cfg.Section("web").Key("logfile").String(); logFile != "" {
			config.LogFile = logFile
		}
	}
	if config.Static == "" {
		if static := cfg.Section("web").Key("static").String(); static != "" {
			config.Static = static
		}
	}
	if config.DefaultIconDir == "" {
		if defaultIconDir := cfg.Section("firebase").Key("defaultIconDir").String(); defaultIconDir != "" {
			config.DefaultIconDir = defaultIconDir
		}
	}
	if config.ServiceKeyPath == "" {
		if serviceAccountKey := cfg.Section("firebase").Key("serviceKeyPath").String(); serviceAccountKey != "" {
			config.ServiceKeyPath = serviceAccountKey
		}
	}
	if config.ProjectId == "" {
		if projectId := cfg.Section("firebase").Key("projectId").String(); projectId != "" {
			config.ProjectId = projectId
		}
	}
	if config.StorageBucket == "" {
		if storageBucket := cfg.Section("firebase").Key("storageBucket").String(); storageBucket != "" {
			config.StorageBucket = storageBucket
		}
	}
}

// 設定値の検証
func validateConfig(config *ConfigList) {
	if config.LogFile == "" {
		config.LogFile = "debug.log"
	}
	if config.Static == "" {
		config.Static = "app/views"
	}
	if config.DefaultIconDir == "" {
		config.DefaultIconDir = "internal/web/images/defaultIcon"
	}

	// ファイルの存在確認（空でない場合のみ）
	if config.ServiceKeyPath != "" {
		if _, err := os.Stat(config.ServiceKeyPath); os.IsNotExist(err) {
			log.Fatalf("エラー: serviceKeyPathファイルが見つかりません: %s", config.ServiceKeyPath)
		}
	}

	// Firebase設定の必須項目チェック
	if config.ProjectId == "" {
		log.Fatal("エラー: PROJECT_ID または projectId が設定されていません")
	}
	if config.StorageBucket == "" {
		log.Fatal("エラー: STORAGE_BUCKET または storageBucket が設定されていません")
	}
}
