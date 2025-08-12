package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var supportedFormats = []string{"mp4", "avi", "mkv", "mov", "wmv", "rmvb"}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // 生产环境需限制
	},
}

// 全局变量存储当前转换的 WebSocket 连接
var clients = make(map[*websocket.Conn]bool)

// 解析 FFmpeg 进度文件
func parseProgress(progressFile string, totalDuration float64) (float64, error) {
	file, err := os.Open(progressFile)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var outTimeMs float64
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "out_time_ms=") {
			timeStr := strings.TrimPrefix(line, "out_time_ms=")
			if ms, err := strconv.ParseFloat(timeStr, 64); err == nil {
				outTimeMs = ms
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return 0, err
	}
	if totalDuration == 0 {
		return 0, nil
	}
	progress := (outTimeMs / 1000000) / totalDuration * 100
	if progress > 100 {
		progress = 100
	}
	return progress, nil
}

// 获取视频总时长（秒）
func getDuration(inputPath string) float64 {
	cmd := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", inputPath)
	output, err := cmd.Output()
	if err != nil {
		log.Printf("获取时长失败: %v", err)
		return 0
	}
	duration, err := strconv.ParseFloat(strings.TrimSpace(string(output)), 64)
	if err != nil {
		log.Printf("解析时长失败: %v", err)
		return 0
	}
	return duration
}

func main() {
	r := gin.Default()

	r.MaxMultipartMemory = 10 << 30 // 10GB

	r.Static("/static", "./static")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"Formats": supportedFormats,
		})
	})

	r.GET("/ws", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("WebSocket 升级失败: %v", err)
			return
		}
		clients[conn] = true
		defer func() {
			delete(clients, conn)
			conn.Close()
		}()

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				log.Printf("WebSocket 读取错误: %v", err)
				break
			}
		}
	})

	r.POST("/convert", func(c *gin.Context) {
		file, err := c.FormFile("video")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "上传文件失败"})
			return
		}

		outputFormat := c.PostForm("format")
		if !contains(supportedFormats, outputFormat) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "不支持的格式"})
			return
		}

		uploadDir := "uploads"
		if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
			os.Mkdir(uploadDir, 0755)
		}
		inputPath := filepath.Join(uploadDir, file.Filename)
		if err := c.SaveUploadedFile(file, inputPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败"})
			return
		}

		inputInfo, _ := os.Stat(inputPath)
		inputSizeMB := float64(inputInfo.Size()) / (1024 * 1024)
		log.Printf("输入文件: %s, 大小: %.2f MB", file.Filename, inputSizeMB)

		base := strings.TrimSuffix(file.Filename, filepath.Ext(file.Filename))
		outputPath := filepath.Join("outputs", fmt.Sprintf("%s_converted.%s", base, outputFormat))
		outputDir := "outputs"
		if _, err := os.Stat(outputDir); os.IsNotExist(err) {
			os.Mkdir(outputDir, 0755)
		}

		// 获取总时长
		totalDuration := getDuration(inputPath)
		log.Printf("视频总时长: %.2f 秒", totalDuration)

		startTime := time.Now()
		progressFile := filepath.Join(outputDir, "progress.txt")
		cmdArgs := []string{
			"-i", inputPath,
			"-analyzeduration", "100M",
			"-probesize", "100M",
			"-progress", progressFile,
		}
		if outputFormat == "rmvb" {
			cmdArgs = append(cmdArgs, "-c:v", "rv40", "-b:v", "2000k", "-c:a", "aac", "-b:a", "128k")
		} else if outputFormat == "avi" {
			cmdArgs = append(cmdArgs, "-c:v", "libx264", "-crf", "18", "-preset", "fast", "-c:a", "mp3", "-b:a", "192k")
		} else {
			cmdArgs = append(cmdArgs, "-c:v", "copy", "-c:a", "copy", "-map", "0", "-f", outputFormat)
		}
		cmdArgs = append(cmdArgs, outputPath)

		cmd := exec.Command("ffmpeg", cmdArgs...)
		if err := cmd.Start(); err != nil {
			log.Printf("FFmpeg 启动失败: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "转换启动失败"})
			os.Remove(inputPath)
			return
		}

		// 进度更新
		go func() {
			for {
				progress, err := parseProgress(progressFile, totalDuration)
				if err != nil {
					continue
				}
				log.Printf("转换进度: %.2f%%", progress)
				// 推送进度到所有客户端
				for client := range clients {
					err := client.WriteJSON(map[string]float64{"progress": progress})
					if err != nil {
						log.Printf("WebSocket 推送错误: %v", err)
						delete(clients, client)
					}
				}
				if progress >= 100 {
					break
				}
				time.Sleep(1 * time.Second)
			}
		}()

		if err := cmd.Wait(); err != nil {
			log.Printf("FFmpeg 错误: %v", err)
			cmdArgs = []string{
				"-i", inputPath,
				"-analyzeduration", "100M",
				"-probesize", "100M",
				"-c:v", "libx264", "-crf", "18", "-preset", "fast",
				"-c:a", "mp3", "-b:a", "192k",
				outputPath,
			}
			cmd = exec.Command("ffmpeg", cmdArgs...)
			if err := cmd.Start(); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "转换重试失败"})
				os.Remove(inputPath)
				return
			}
			if err := cmd.Wait(); err != nil {
				log.Printf("FFmpeg 重试错误: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "转换失败"})
				os.Remove(inputPath)
				return
			}
		}

		duration := time.Since(startTime).Seconds()
		log.Printf("转换耗时: %.2f 秒", duration)

		fileInfo, err := os.Stat(outputPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "无法获取输出文件信息"})
			os.Remove(inputPath)
			return
		}
		fileSizeMB := float64(fileInfo.Size()) / (1024 * 1024)
		log.Printf("输出文件: %s, 大小: %.2f MB", outputPath, fileSizeMB)

		c.JSON(http.StatusOK, gin.H{
			"download_url":  "/download/" + filepath.Base(outputPath),
			"file_size_mb":  fmt.Sprintf("%.2f MB", fileSizeMB),
			"convert_time":  fmt.Sprintf("%.2f 秒", duration),
		})

		os.Remove(inputPath)
		os.Remove(progressFile)
	})

	r.GET("/download/:filename", func(c *gin.Context) {
		filename := c.Param("filename")
		outputPath := filepath.Join("outputs", filename)
		if _, err := os.Stat(outputPath); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
			return
		}

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Transfer-Encoding", "binary")
		c.Header("Content-Disposition", "attachment; filename="+filename)
		c.Header("Content-Type", "application/octet-stream")
		c.File(outputPath)
	})

	r.LoadHTMLGlob("templates/*")
	r.Run(":8080")
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}