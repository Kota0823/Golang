/*
リレーテーブルの情報を取得するスレッド
*/

package getrelaytable

import (
	"log"
	"net"
	"os"
	"strconv"

	librt "../../RelayTableLibrary"
)

func GetTable(channelTunnels chan map[librt.ID]librt.Information) {
	/*ソケット作成*/
	defer os.Remove(librt.SocketFilepath)                     //プログラム終了時にファイルを削除する
	listener, err := net.Listen("unix", librt.SocketFilepath) //PIDの代わりにファイルパスを指定
	log.Printf("info: Socket OK\n")

	conn, _ := listener.Accept() //クライアントと接続
	defer conn.Close()           //プログラム終了時にソケットを閉じる（実行予約）
	if err != nil {
		log.Printf("error: %v\n", err) //ソケット作成時エラーが発生した場合
		return
	}
	log.Printf("info: Socket connected\n")

	buf := make([]byte, 1)

	/*クライアントから命令を受信*/
	for {
		conn.Read(buf) //ソケットから受信した情報読み込み

		typeNumber, _ := strconv.Atoi(string(buf))
		log.Printf("recive:requestType = %d", typeNumber)

		switch typeNumber {
		case librt.RequestRelayTable: //リレーテーブル送信処理
			for index, entity := range <-channelTunnels { //チャネルからリレーテーブルを取得
				librt.SendMessage(conn, index[:])
				librt.SendMessage(conn, []byte(entity.En1))
				librt.SendMessage(conn, []byte(entity.En2))
			}
			librt.SendMessage(conn, []byte("exit")) //すべて送信し終えた場合"exit"を送信

		case librt.Exit: //処理終了，main関数を終了させる
			break

		case librt.Pass:
		} //switch
	} //for
}
