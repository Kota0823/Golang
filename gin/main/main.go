package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	type Information struct {
		Index string
		En1   string
		En2   string
	}
	var tunnels = make(map[[16]byte]Information) //マップ型変数

	tunnels[[16]byte{0}] = Information{
		Index: "1",
		En1:   "192.168.100.1",
		En2:   "192.168.100.2",
	}
	tunnels[[16]byte{1}] = Information{
		Index: "2",
		En1:   "192.168.200.3",
		En2:   "192.168.200.4",
	}

	/*HTMLファイルへレンダリング*/
	router := gin.Default()
	router.LoadHTMLGlob("templates/*.tmpl") //テンプレートファイル読み込み
	router.GET("/table", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"name":    "RelayTable(Gin)",
			"rsaddr":  "192.168.100.5",
			"tunnels": tunnels,
		})
	})

	router.Run(":8080") //サーバ起動
}
