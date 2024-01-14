package Backend

import (
	"github.com/APCS20-Thesis/Backend/api"
	"github.com/gin-gonic/gin"
	"net/http"
)

// getAlbums responds with the list of all albums as JSON.
func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, api.CommonResponse{
		Code:    0,
		Message: "Success",
	})
}

func main() {
	router := gin.Default()
	router.GET("/health", getAlbums)

	err := router.Run("localhost:8080")
	if err != nil {
		return
	}
}
