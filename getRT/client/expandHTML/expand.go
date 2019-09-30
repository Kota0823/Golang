/*
HTMLへの展開
*/

package expandHTML

import (
	"net/http"

	"../../RelayTableLibrary"
	"github.com/gin-gonic/gin"
)

func ExpandRelaytableToHTML(tunnels map[RelayTableLibrary.ID]RelayTableLibrary.Information) (err error) {
	gin.SetMode(gin.ReleaseMode)

	rsAddress := "192.168.100.5"

	/*HTMLファイルへレンダリング*/
	router := gin.Default()
	router.LoadHTMLGlob("expandHTML/templates/*.tmpl") //テンプレートファイル読み込み
	router.GET("/relaytable", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"name":    "\"github.com/gin-gonic/gin\"",
			"rsaddr":  rsAddress,
			"tunnels": tunnels,
		})
	})

	err = router.Run(":8989") //サーバ起動
	return
}
