package gin

import (
	"github.com/gin-gonic/gin"
)

// Handler ...
func Handler() *gin.Engine {
	router := gin.Default()

	router.GET("/", root)

	return router
}

func root(ctx *gin.Context) {
	type svcInfo struct {
		JSONAPI struct {
			Version string `json:"version,omitempty"`
			Name    string `json:"name,omitempty"`
		} `json:"jsonapi"`
	}

	info := svcInfo{}
	info.JSONAPI.Version = "v1"
	info.JSONAPI.Name = "HRM API"

	ctx.JSON(200, info)
	// w.Write(jsonutil.Marshal(info))
}
