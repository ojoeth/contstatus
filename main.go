package main

import (
	"context"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	contstatusAuth := os.Getenv("CONTSTATUS_AUTH")
	if contstatusAuth == "" {
		panic("Auth env is empty. Please set CONTSTATUS_AUTH env var.")
	}
	router.Use(func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader != contstatusAuth {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	})
	router.GET("/getContainerStatus", getStatus)

	router.Run(":8080")
}

func getStatus(c *gin.Context) {
	apiClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	containers, err := apiClient.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		panic(err)
	}

	for _, ctr := range containers {
		for _, n := range ctr.Names {
			n = strings.Replace(n, "/", "", -1)
			if n == c.Query("service") {
				re := regexp.MustCompile(`Up\s+`)
				match := re.FindStringSubmatch(ctr.Status)
				if match != nil {
					c.IndentedJSON(http.StatusOK, map[string]string{"status": "up"})
				} else {
					c.IndentedJSON(http.StatusOK, map[string]string{"status": ctr.Status})
				}
				return
			}
		}
	}
	c.String(http.StatusNotFound, "not found\n")
}
