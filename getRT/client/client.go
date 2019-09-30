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

	librt "../RelayTableLibrary"
	"./expandHTML"
)

func main() {
	/*ソケット作成*/
	defer os.Remove(librt.SocketFilepath)               //プログラム終了時にファイルを削除する
	conn, err := net.Dial("unix", librt.SocketFilepath) //PIDの代わりにファイルパスを指定
	if err != nil {                                     //エラーによりソケットが作成できなかった場合
		log.Printf("error: %v\n", err)
		return
	}
	defer conn.Close() //プログラム終了時にソケットを閉じる（実行予約）
	log.Printf("info: Socket connected\n")

	/*受信した情報を格納するマップ型変数の定義*/
	var tunnels = make(map[librt.ID]librt.Information) //マップ型変数

	/*リレーテーブル用チャネル*/
	tuunelChan := make(chan map[librt.ID]librt.Information)

	/*リレーテーブルをHTMLに出力するスレッド*/
	go expandHTML.ExpandRelaytableToHTML(tuunelChan)

	defer conn.Write([]byte{49}) //client.go終了時に終了要求を送信

	/*リレーテーブル受信処理*/
	for range time.Tick(5 * time.Second) { //リレーテーブル取得周期（秒）
		/*要求送信*/
		conn.Write([]byte{48}) //0(文字コード)
		log.Printf("info: Get Relaytable...")

		for {
			index := librt.ResiveMessage(conn)
			if reflect.DeepEqual([]byte(index), librt.EXIT) == true {
				break //exitの文字列を受信した場合forループ終了（リレーテーブルを全て受信した）
			}

			tunnelsIndex := [16]byte{}                 //マップ型変数の要素にアクセスするためのキー
			copy([]byte(index), tunnelsIndex[0:15])    //スライス（可変長配列）から配列（固定長）に変換
			tunnels[tunnelsIndex] = librt.Information{ //受信した情報をマップ型変数に追加
				Index: index,
				En1:   librt.ResiveMessage(conn),
				En2:   librt.ResiveMessage(conn),
			}
		} //for
		tuunelChan <- tunnels
	} //for range time.Tick(sec * time.Second)
}
