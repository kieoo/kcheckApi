package helm

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func HelmChange(c *gin.Context) {

	// save helm chart
	up := Upload{}
	if err := up.UploadInit(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("upload error : %s", err)})
		return
	}
	defer up.UploadClose()

	if err := up.UploadSave(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("upload error : %s", err)})
		return
	}

	// chart template to deploy

	template, err := Chart2Deploy(up.ChartDir)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("chage error : %s", string(template))})
		return
	}

	c.PureJSON(http.StatusOK, gin.H{"msg": string(template)})

}
