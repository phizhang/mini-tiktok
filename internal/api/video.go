package api

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/phizhang/mini-tiktok/internal/configs"
	"github.com/phizhang/mini-tiktok/internal/db"
	"github.com/phizhang/mini-tiktok/internal/models"
	"github.com/phizhang/mini-tiktok/internal/storage"
)

// UploadVideoHandler handles video file uploads and stores metadata in Cassandra
func UploadVideoHandler(c *gin.Context) {
	log.Printf("UploadVideoHandler called by user: %s", c.GetHeader("X-User-ID"))
	userID := c.GetHeader("X-User-ID") // Example: parse user ID from header (replace with real auth)
	if userID == "" {
		log.Println("Missing user ID in request header")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing user ID"})
		return
	}

	title := c.PostForm("title")
	tag := c.PostForm("tag")
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		log.Println("Missing video file in form data")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing video file"})
		return
	}
	defer file.Close()

	videoUUID := uuid.New()
	videoID := videoUUID.String()
	log.Println("videoID: ", videoID)

	log.Printf("Storing video for user %s with videoID %s using storage type %s", userID, videoID, configs.GlobalEnv.StorageType)
	// Use global env config to decide storage backend
	var store storage.Storage
	if configs.GlobalEnv.StorageType == "s3" {
		// TODO: Upload file to S3 and get videoID (for now, use timestamp)
		store = &storage.S3Storage{
			Bucket: configs.GlobalEnv.S3Bucket}
	} else {
		path := fmt.Sprintf("%s/%s", configs.GlobalEnv.LocalPath, userID)
		store = &storage.LocalStorage{BasePath: path}
	}

	location, err := store.Upload(file, videoID)
	if err != nil {
		log.Printf("Failed to upload video: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to upload video: %v", err)})
		return
	}

	meta := models.VideoMeta{
		UserID:      userID,
		VideoID:     videoID,
		CreatedTime: time.Now(),
		Title:       title,
		Tag:         tag,
		Location:    location,
	}

	if db.Session == nil {
		log.Println("Cassandra not connected")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cassandra not connected"})
		return
	}
	query := `INSERT INTO video_meta (user_id, video_id, created_time, title, tag, location) VALUES (?, ?, ?, ?, ?, ?)`
	err = db.Session.Query(query, meta.UserID, meta.VideoID, meta.CreatedTime, meta.Title, meta.Tag, meta.Location).Exec()
	if err != nil {
		log.Printf("Failed to insert video meta: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Video uploaded successfully: %s", videoID)
	c.JSON(http.StatusOK, gin.H{"message": "Video uploaded", "video_id": videoID})
}

// GetUserVideosHandler returns the list of uploaded videos for a user
func GetUserVideosHandler(c *gin.Context) {
	userID := c.GetHeader("X-User-ID") // Example: parse user ID from header (replace with real auth)
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing user ID"})
		return
	}
	if db.Session == nil {
		log.Println("Cassandra not connected")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cassandra not connected"})
		return
	}
	query := `SELECT video_id, created_time, title, tag, location FROM video_meta WHERE user_id = ?`
	iter := db.Session.Query(query, userID).Iter()

	var videos []models.VideoMeta
	var videoID, title, tag, location string
	var createdTime time.Time
	for iter.Scan(&videoID, &createdTime, &title, &tag, &location) {
		videos = append(videos, models.VideoMeta{
			UserID:      userID,
			VideoID:     videoID,
			CreatedTime: createdTime,
			Title:       title,
			Tag:         tag,
			Location:    location,
		})
	}
	if err := iter.Close(); err != nil {
		log.Printf("Failed to fetch videos: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"videos": videos})
}
