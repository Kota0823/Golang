/*
リレーテーブルの情報を取得するスレッド
*/

package DumpRelayTable

import (
	"log"
	"net"
	"os"
	"strconv"

	librt "../../RelayTableLibrary"
)

/*リレーテーブル(グローバル変数)*/
type ID [16]byte //マップ型変数で使用するキー

type Information struct {
	Index string
	En1   string
	En2   string
}

/*ダミーデータ*/
var tunnels = make(map[ID]Information) //マップ型変数

func GetTable(channelTunnels chan map[ID]Information) {
	/*ソケット作成*/
	os.Remove(librt.SocketFilepath)                           //プログラム終了時にファイルを削除する
	listener, err := net.Listen("unix", librt.SocketFilepath) //PIDの代わりにファイルパスを指定
	log.Printf("info: Socket OK\n")

	conn, _ := listener.Accept() //クライアントと接続
	defer conn.Close()           //プログラム終了時にソケットを閉じる（実行予約）
	if err != nil {
		log.Printf("error: %v\n", err) //ソケット作成時エラーが発生した場合
		return
	}
	log.Printf("info: Socket connected\n")

	tunnel := <-channelTunnels

	/*クライアントから命令を受信*/
	for {
		typeNumber, _ := strconv.Atoi(string(librt.ResiveMessage(conn))) //ソケットから受信した情報読み込み
		log.Printf("recive:requestType = %d", typeNumber)

		switch typeNumber {
		case librt.RequestRelayTable: //リレーテーブル送信処理
			select {
			case tunnel = <-channelTunnels: //チャネルに情報が入っている場合
			default: //チャネルに情報が入っていない場合
				log.Println("info: no value")
			}
			for index, entity := range tunnel { //リレーテーブルを取得
				librt.SendMessage(conn, index[:])
				librt.SendMessage(conn, []byte(entity.En1))
				librt.SendMessage(conn, []byte(entity.En2))
			}
			librt.SendMessage(conn, []byte("exit")) //すべて送信し終えた場合"exit"を送信

		case librt.Exit: //処理終了
			break

		case librt.Pass:
		} //switch
	} //for
}
