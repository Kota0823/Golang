/*
HTMLへの展開
*/

package RenderingHTML

import (
	"log"
	"net/http"
	"time"

	librt "../../RelayTableLibrary"
	"github.com/gin-gonic/gin"
)

func ExpandRelaytableToHTML(tunnelchan chan map[librt.ID]librt.Information) (err error) {
	gin.SetMode(gin.ReleaseMode)
	rsAddress := "192.168.100.5"

	/*リレーテーブル用チャネルから取得*/
	tunnel := <-tunnelchan

	updateTableTimeChan := time.Now()

	/*HTMLファイルへレンダリング*/
	router := gin.Default()
	router.LoadHTMLGlob("expandHTML/templates/*.tmpl") //テンプレートファイル読み込み
	router.GET("/relaytable", func(c *gin.Context) {
		/*リレーテーブル用キューに情報がある場合は取得*/
		select {
		case tunnel = <-tunnelchan: //キューに情報が入っている場合
			updateTableTimeChan = time.Now() //リレーテーブル取得時間を格納
		default: //キューに情報が入っていない場合
			log.Println("info: no value")
		}
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"time":    updateTableTimeChan,
			"rsaddr":  rsAddress,
			"tunnels": tunnel,
		})
	})

	log.Printf("info: expand HTML... \n")
	err = router.Run(":8989") //サーバ起動(エラーが発生しない限り実行されるメソッド)
	if err != nil {
		log.Printf("gin error: %v\n", err)
		return
	}
	return
}
