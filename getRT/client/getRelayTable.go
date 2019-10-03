/*
管理プロセス（リレーテーブルの取得，HTMLへの展開）
*/

package main

import (
	"log"
	"net"
	"os" //スライス比較
	"reflect"
	"strings"
	"time"

	librt "../RelayTableLibrary"
	"./RenderingHTML"
)

func main() {
	/*ソケット作成*/
	defer os.Remove(librt.SocketFilepath)               //プログラム終了時にファイルを削除する
	conn, err := net.Dial("unix", librt.SocketFilepath) //PIDの代わりにファイルパスを指定
	if err != nil {                                     //エラーによりソケットが作成できなかった場合
		log.Printf("Socketerror: %v\n", err)
		return
	}
	defer conn.Close() //プログラム終了時にソケットを閉じる（実行予約）
	log.Printf("info: Socket connected\n")

	/*受信した情報を格納するマップ型変数の定義*/
	var tunnels = make(map[librt.ID]librt.Information) //マップ型変数

	/*リレーテーブル用チャネル*/
	tuunelChan := make(chan map[librt.ID]librt.Information)

	/*リレーテーブルをHTMLに出力するスレッド*/
	go RenderingHTML.ExpandRelaytableToHTML(tuunelChan)

	defer conn.Write([]byte{49}) //client.go終了時に終了要求を送信

	/*リレーテーブル受信処理*/
	for range time.Tick(5 * time.Second) { //リレーテーブル取得周期（秒）
		log.Printf("info: Get Relaytable...\n")

		/*要求送信*/
		librt.SendMessage(conn, []byte{48})
		for {
			index := librt.ResiveMessage(conn)
			if reflect.DeepEqual(index[0:4], librt.EXIT) == true {
				log.Printf("info: finish for loop\n")
				break //exitの文字列を受信した場合forループ終了（リレーテーブルを全て受信した）
			}

			tunnelsIndex := [16]byte{}                 //マップ型変数の要素にアクセスするためのキー
			copy(tunnelsIndex[0:15], index[:])         //スライス（可変長配列）から配列（固定長）に変換
			tunnels[tunnelsIndex] = librt.Information{ //受信した情報をマップ型変数に追加
				Index: strings.Trim(string(index), string([]byte{0})),
				En1:   strings.Trim(string(librt.ResiveMessage(conn)), string([]byte{0})),
				En2:   strings.Trim(string(librt.ResiveMessage(conn)), string([]byte{0})),
			}
		} //for
		tuunelChan <- tunnels
	} //for range time.Tick(sec * time.Second)
}
