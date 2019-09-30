/*
HTMLへの展開
*/

package expandHTML

import (
	"log"
	"net/http"

	librt "../../RelayTableLibrary"
	"github.com/gin-gonic/gin"
)

func ExpandRelaytableToHTML(tunnels chan map[librt.ID]librt.Information) (err error) {
	//gin.SetMode(gin.ReleaseMode)
	rsAddress := "192.168.100.5"
	log.Printf("info: expand HTML... \n")

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

	err = router.Run(":8989") //サーバ起動(エラーが発生しない限り実行されるメソッド)
	if err != nil {
		log.Printf("error: %v\n", err)
		return
	}
}
