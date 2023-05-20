package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()

	// Auth services
	r.GET("/req_pq", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})


	r.GET("/req_DH_params", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Biz services
	r.GET("/get_users", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/get_users_with_sql_inject", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.Run("0.0.0.0:6433") // listen and serve on 0.0.0.0:8080
}
