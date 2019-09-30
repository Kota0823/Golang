/*
管理プロセス（リレーテーブルの取得，HTMLへの展開）
*/

package main

import (
	"log"
	"net"
	"os"
	"reflect" //スライス比較
	"time"

	"../RelayTableLibrary"
	"./expandHTML"
)

func main() {
	/*ソケット作成*/
	defer os.Remove(RelayTableLibrary.SocketFilepath)               //プログラム終了時にファイルを削除する
	conn, err := net.Dial("unix", RelayTableLibrary.SocketFilepath) //PIDの代わりにファイルパスを指定
	if err != nil {                                                 //エラーによりソケットが作成できなかった場合
		log.Printf("error: %v\n", err)
		return
	}
	defer conn.Close() //プログラム終了時にソケットを閉じる（実行予約）
	log.Printf("info: Socket connected\n")

	/*受信した情報を格納するマップ型変数の定義*/
	var tunnels = make(map[RelayTableLibrary.ID]RelayTableLibrary.Information) //マップ型変数

	/*リレーテーブルをHTMLに出力*/
	log.Printf("info: expand HTML... \n")
	go expandHTML.ExpandRelaytableToHTML(tunnels)
	if err != nil {
		log.Printf("error: %v\n", err)
	}
	/*受信処理*/
	for range time.Tick(5 * time.Second) { //リレーテーブル取得周期（秒）
		conn.Write([]byte{48}) //要求送信

		log.Printf("info: Get Relaytable...")

		for {
			index := RelayTableLibrary.ResiveMessage(conn)
			if reflect.DeepEqual([]byte(index), RelayTableLibrary.EXIT) == true {
				break //exitの文字列を受信した場合forループ終了
			}

			tunnelsIndex := [16]byte{}                             //マップ型変数の要素にアクセスするためのキー
			copy([]byte(index), tunnelsIndex[0:15])                //スライス（可変長配列）から配列（固定長）に変換
			tunnels[tunnelsIndex] = RelayTableLibrary.Information{ //受信した情報をマップ型変数に追加
				Index: index,
				En1:   RelayTableLibrary.ResiveMessage(conn),
				En2:   RelayTableLibrary.ResiveMessage(conn),
			}
		} //for

	} //for range time.Tick(sec * time.Second)

	conn.Write([]byte{49}) //終了要求

}
