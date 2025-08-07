package models

import "time"

// VideoMeta represents metadata for a video stored in Cassandra
// Partition key: UserID
// Clustering keys: VideoID, CreatedTime
type VideoMeta struct {
	UserID      string    // partition key
	VideoID     string    // clustering key
	CreatedTime time.Time // clustering key
	Title       string
	Tag         string
	Location    string
}
