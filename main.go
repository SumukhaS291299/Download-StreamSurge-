package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func portAssignment() (port string) {
	port, present := os.LookupEnv("PORT")
	if !present {
		port = "8080"
	}
	return port
}

func hostAssignment() (host string) {
	host, present := os.LookupEnv("HOST")
	if !present {
		host = "localhost"
	}
	return host
}

func listfiles() (filesListstr []string) {
	pwd, err := os.Getwd()
	if err != nil {
		log.Printf("Error occured %s", err)
	}
	dir := filepath.Join(pwd, "Files")
	fileList, err := os.ReadDir(dir)
	if err != nil {
		log.Printf("Error occured %s", err)
	}
	filesListstr = make([]string, 0)
	for _, files := range fileList {
		filesListstr = append(filesListstr, files.Name())
	}
	return filesListstr

}

func liveness(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"Liveness": true})
}

func list(c *gin.Context) {
	fileList := listfiles()
	c.JSON(http.StatusOK, gin.H{"Dir": fileList})
}

func downloadFile(c *gin.Context) {
	fileName := c.Param("name")
	fmt.Println("File name: ", fileName)
	pwd, err := os.Getwd()
	if err != nil {
		log.Printf("Error occured %s", err)
	}
	filePath := filepath.Join(pwd, filepath.Join("Files", fileName))
	fmt.Println("File path: ", filePath)
	c.FileAttachment(filePath, fileName)
}

func home(c *gin.Context) {
	fileList := listfiles()
	var sb strings.Builder
	sb.WriteString("<html><body><H1>Available Download Files</H1>")
	for _, file := range fileList {
		sb.WriteString("<a href=\"/download/" + file + "\">" + file + "</a>")
		sb.WriteString("<br>")
	}
	sb.WriteString("</body></html>")
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write([]byte(sb.String()))
}

func main() {
	r := gin.Default()
	r.GET("/", home)
	r.GET("/list", list)
	r.GET("/live", liveness)
	r.GET("/download/:name", downloadFile)
	r.Run(hostAssignment() + ":" + portAssignment())
}

// Need to check if docker compose and nginx if it works....
// Pushing the code anyways....
