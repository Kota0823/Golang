/*
リレーテーブルの情報を取得するスレッド
*/

package getrelaytable

import (
	"log"
	"net"
	"os"
	"strconv"

	"../../RelayTableLibrary"
)

func GetTable(channelTunnels chan map[RelayTableLibrary.ID]RelayTableLibrary.Information, control chan<- bool) {
	/*ソケット作成*/
	defer os.Remove(RelayTableLibrary.SocketFilepath)                     //プログラム終了時にファイルを削除する
	listener, err := net.Listen("unix", RelayTableLibrary.SocketFilepath) //PIDの代わりにファイルパスを指定
	log.Printf("info: Socket OK\n")

	conn, _ := listener.Accept() //クライアントと接続
	defer conn.Close()           //プログラム終了時にソケットを閉じる（実行予約）
	if err != nil {
		log.Printf("error: %v\n", err) //ソケット作成時エラーが発生した場合
		return
	}
	log.Printf("info: Socket connected\n")

	buf := make([]byte, 1)

	defer conn.Write([]byte("exit"))

	/*クライアントから命令を受信*/
	for {
		conn.Read(buf) //ソケットから受信した情報読み込み

		typeNumber, _ := strconv.Atoi(string(buf))
		log.Printf("recive:requestType = %d", typeNumber)

		switch typeNumber {
		case RelayTableLibrary.RequestRelayTable: //リレーテーブル送信処理
			for index, entity := range <-channelTunnels { //チャネルからリレーテーブルを取得
				conn.Write(index[:]) //配列([16]byte)をスライスに変換し，メッセージを送信
				log.Printf("send:Index %s", string(index[:]))
				conn.Read(make([]byte, 3)) //応答(ACK)を受信
				log.Printf("receive:ACK\n")

				conn.Write([]byte(entity.En1)) //メッセージを送信
				log.Printf("send:En1 %s", entity.En1)
				conn.Read(make([]byte, 3)) //応答(ACK)を受信
				log.Printf("receive:ACK\n")

				conn.Write([]byte(entity.En2)) //メッセージを送信
				log.Printf("send:En2 %s", entity.En2)
				conn.Read(make([]byte, 3)) //応答(ACK)を受信
				log.Printf("receive:ACK\n")
			}
			conn.Write([]byte("exit")) //すべて送信し終えた場合"exit"を送信
			conn.Read(make([]byte, 3)) //応答(ACK)を受信

		case RelayTableLibrary.Exit: //処理終了，main関数を終了させる
			control <- true
			break

		case RelayTableLibrary.Pass:
		}
	}
	//var tunnels = make(map[ID]Information)
	//tunnels <- channelTunnel

}
