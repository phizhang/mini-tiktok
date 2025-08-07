package configs

type StorageType string

const (
	StorageTypeS3    StorageType = "s3"
	StorageTypeLocal StorageType = "local"
)

type StorageConfig struct {
	Type  StorageType `json:"type"`
	S3    S3Config    `json:"s3"`
	Local LocalConfig `json:"local"`
}

type S3Config struct {
	Bucket string `json:"bucket"`
	Region string `json:"region"`
}

type LocalConfig struct {
	BasePath string `json:"base_path"`
}
