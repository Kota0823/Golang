/*
リレーテーブルの情報を取得するスレッド
*/

package main

import (
	"log"
	"net"
	"os"
	"strconv"

	librt "../RelayTableLibrary"
)

func GetTable() {
	/*リレーテーブル（ダンプ後）*/
	var dumpedTunnels = make(map[ID]Information) //マップ型変数

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

	/*クライアントから命令を受信*/
	for {
		typeNumber, _ := strconv.Atoi(string(librt.ResiveMessage(conn))) //ソケットから受信した情報読み込み
		log.Printf("recive:requestType = %d", typeNumber)

		switch typeNumber {
		case librt.RequestRelayTable: //リレーテーブル送信処理
			/*mainよりリレーテーブルの取得*/
			mutex.RLock() //読み込みロック（ロック取れるまでブロック）
			dumpedTunnels = tunnels
			mutex.RUnlock() //読み込みアンロック

			for index, entity := range dumpedTunnels { //リレーテーブルを取得
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
