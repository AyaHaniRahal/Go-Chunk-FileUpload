package main

import (
	"io"
	"net/http"
	"os"
	"strconv"

	"upload-gin/helper"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	r.LoadHTMLGlob("FE/*.html")

	// Serve static files from the "static" directory.
	//r.Static("/upload.html", "./FE")

	// Define a route to serve the upload.html file.
	r.GET("/", func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				// Handle the panic here, log it, and respond with an error message if needed.
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			}
		}()
		c.HTML(http.StatusOK, "upload.html", gin.H{})
	})

	r.POST("/upload", func(c *gin.Context) {
		// Retrieve the uploaded file from the request.
		file, header, err := c.Request.FormFile("file")
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		defer file.Close()

		// Specify the directory where you want to save the uploaded file.
		uploadDir := "./uploads"

		// Create the directory if it doesn't exist.
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			c.JSON(500, gin.H{"error": "Unable to create directory"})
			return
		}

		// Create a new file on the server to store the uploaded content.
		destination, err := os.Create(uploadDir + "/" + header.Filename)
		if err != nil {
			c.JSON(500, gin.H{"error": "Unable to create file on server"})
			return
		}
		defer destination.Close()

		// Copy the uploaded file to the newly created file on the server in chunks.
		if _, err := io.Copy(destination, file); err != nil {
			c.JSON(500, gin.H{"error": "Unable to copy file to server"})
			return
		}

		c.JSON(200, gin.H{"message": "File uploaded successfully"})
	})

	r.POST("/chunkUpload", func(c *gin.Context) {

		// Parse form data to get the current chunk number and total chunks.
		currentChunkStr := c.PostForm("currentChunk")
		totalChunksStr := c.PostForm("totalChunks")

		currentChunk, err := strconv.Atoi(currentChunkStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid currentChunk value"})
			return
		}

		totalChunks, err := strconv.Atoi(totalChunksStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid totalChunks value"})
			return
		}

		// Determine the filename for the destination file.
		filename := c.PostForm("fileName")

		// Open the destination file in append mode.
		destination, err := helper.OpenFileForAppend(filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer destination.Close()

		// Open the uploaded chunk from the form.
		file, _, err := c.Request.FormFile("chunk")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer file.Close()

		// Copy the chunk's contents to the destination file.
		_, err = io.Copy(destination, file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Check if all chunks have been received.
		if currentChunk == totalChunks-1 {
			c.JSON(http.StatusOK, gin.H{"message": "File upload complete!"})
		} else {
			c.JSON(http.StatusOK, gin.H{"message": "Chunk received."})
		}

	})

	r.Run(":8080")

}
