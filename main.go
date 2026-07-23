package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Post struct {
	ID          int       `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Title       string    `gorm:"size:200;not null;column:title" json:"title"`
	Content     string    `gorm:"type:text;not null;column:content" json:"content"`
	Category    string    `gorm:"size:100;not null;column:category" json:"category"`
	CreatedDate time.Time `gorm:"column:created_date;autoCreateTime" json:"created_date"`
	UpdatedDate time.Time `gorm:"column:updated_date;autoUpdateTime" json:"updated_date"`
	Status      string    `gorm:"size:100;not null;column:status" json:"status"`
}

func (Post) TableName() string {
	return "posts"
}

type ArticleRequest struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	Category string `json:"category"`
	Status   string `json:"status"`
}

type ArticleResponse struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Category string `json:"category"`
	Status   string `json:"status"`
}

func bindJSON(c *gin.Context, req interface{}) error {
	if err := c.ShouldBindJSON(req); err != nil {
		if errors.Is(err, io.EOF) {
			return errors.New("request body tidak boleh kosong, sertakan payload JSON")
		}
		return fmt.Errorf("invalid JSON payload: %w", err)
	}
	return nil
}

func validateRequest(req *ArticleRequest) error {
	if strings.TrimSpace(req.Title) == "" {
		return errors.New("title is required")
	}
	if strings.TrimSpace(req.Content) == "" {
		return errors.New("content is required")
	}
	if strings.TrimSpace(req.Category) == "" {
		return errors.New("category is required")
	}
	st := strings.ToLower(strings.TrimSpace(req.Status))
	if st != "publish" && st != "draft" && st != "thrash" {
		return errors.New("status must be one of: 'publish', 'draft', or 'thrash'")
	}
	return nil
}

func main() {
	_ = godotenv.Load()

	dbHost := getEnv("DB_HOST", getEnv("MYSQLHOST", "127.0.0.1"))
	dbPort := getEnv("DB_PORT", getEnv("MYSQLPORT", "3306"))
	dbUser := getEnv("DB_USER", getEnv("MYSQLUSER", "root"))
	dbPass := getEnv("DB_PASSWORD", getEnv("MYSQLPASSWORD", ""))
	dbName := getEnv("DB_NAME", getEnv("MYSQLDATABASE", "article"))
	serverPort := getEnv("PORT", getEnv("SERVER_PORT", "8080"))

	dsnWithoutDB := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort)

	tempDB, err := gorm.Open(mysql.Open(dsnWithoutDB), &gorm.Config{})
	if err == nil {
		_ = tempDB.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`;", dbName)).Error
		sqlDB, _ := tempDB.DB()
		if sqlDB != nil {
			_ = sqlDB.Close()
		}
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Peringatan: Gagal terhubung ke MySQL database '%s': %v\n", dbName, err)
	} else {
		log.Println("Berhasil terhubung ke database MySQL!")
		_ = db.AutoMigrate(&Post{})
	}

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "Backend Microservice Sharing Vision 2023 is running"})
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.POST("/article/", func(c *gin.Context) {
		if db == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "database connection error"})
			return
		}

		var req ArticleRequest
		if err := bindJSON(c, &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := validateRequest(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		post := Post{
			Title:    req.Title,
			Content:  req.Content,
			Category: req.Category,
			Status:   strings.ToLower(strings.TrimSpace(req.Status)),
		}

		if err := db.Create(&post).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{})
	})

	r.GET("/article/:id/:offset", func(c *gin.Context) {
		if db == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "database connection error"})
			return
		}

		limit, err1 := strconv.Atoi(c.Param("id"))
		offset, err2 := strconv.Atoi(c.Param("offset"))
		if err1 != nil || err2 != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "parameter limit & offset harus berupa angka"})
			return
		}

		var posts []Post
		if err := db.Limit(limit).Offset(offset).Order("id asc").Find(&posts).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var responses []ArticleResponse
		for _, p := range posts {
			responses = append(responses, ArticleResponse{
				ID:       p.ID,
				Title:    p.Title,
				Content:  p.Content,
				Category: p.Category,
				Status:   p.Status,
			})
		}

		if responses == nil {
			responses = []ArticleResponse{}
		}

		c.JSON(http.StatusOK, responses)
	})

	r.GET("/article/:id", func(c *gin.Context) {
		if db == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "database connection error"})
			return
		}

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "parameter id harus berupa angka"})
			return
		}

		var post Post
		if err := db.First(&post, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "article tidak ditemukan"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, ArticleResponse{
			ID:       post.ID,
			Title:    post.Title,
			Content:  post.Content,
			Category: post.Category,
			Status:   post.Status,
		})
	})

	updateHandler := func(c *gin.Context) {
		if db == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "database connection error"})
			return
		}

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "parameter id harus berupa angka"})
			return
		}

		var req ArticleRequest
		if err := bindJSON(c, &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := validateRequest(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var post Post
		if err := db.First(&post, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "article tidak ditemukan"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		post.Title = req.Title
		post.Content = req.Content
		post.Category = req.Category
		post.Status = strings.ToLower(strings.TrimSpace(req.Status))

		if err := db.Save(&post).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{})
	}

	r.PUT("/article/:id", updateHandler)
	r.PATCH("/article/:id", updateHandler)
	r.POST("/article/:id", updateHandler)

	r.DELETE("/article/:id", func(c *gin.Context) {
		if db == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "database connection error"})
			return
		}

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "parameter id harus berupa angka"})
			return
		}

		result := db.Delete(&Post{}, id)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}

		if result.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "article tidak ditemukan"})
			return
		}

		c.JSON(http.StatusOK, gin.H{})
	})

	log.Printf("Server berjalan di port :%s\n", serverPort)
	if err := r.Run(":" + serverPort); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}
