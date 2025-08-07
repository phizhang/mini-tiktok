package configs

import (
	"os"
)

type EnvConfig struct {
	StorageType string
	S3Bucket    string
	S3Region    string
	LocalPath   string
	DBHost      string
}

var GlobalEnv EnvConfig

func LoadEnv() {
	GlobalEnv = EnvConfig{
		StorageType: os.Getenv("VIDEO_STORAGE_TYPE"),
		S3Bucket:    os.Getenv("S3_BUCKET"),
		S3Region:    os.Getenv("S3_REGION"),
		LocalPath:   os.Getenv("LOCAL_STORAGE_PATH"),
		DBHost:      os.Getenv("DB_HOSTS"),
	}
	if GlobalEnv.LocalPath == "" {
		GlobalEnv.LocalPath = "./videos"
	}
	if GlobalEnv.StorageType == "" {
		GlobalEnv.StorageType = "local" // Default to local storage if not set
	}
}
