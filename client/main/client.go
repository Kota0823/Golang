package main

import (
	"log"
	"net"
	"os"
	"path/filepath"
	"reflect" //スライス比較
)

func main() {
	/*ソケット作成*/
	file := filepath.Join(os.TempDir(), "unixdomaisocketsample") //ソケット作成に使うファイルパス
	defer os.Remove(file)
	conn, err := net.Dial("unix", file) //PIDの代わりにファイルパスを指定
	defer conn.Close()                  //プログラム終了時にファイルを削除する（実行予約）
	if err != nil {
		log.Printf("error: %v\n", err)
		return
	}
	/*リクエスト送信*/
	conn.Write([]byte("request"))
	log.Printf("send request")

	/*プロセス間通信によるリスト型変数取得*/
	type ID [16]byte
	type Information struct {
		En1 string
		En2 string
	}
	var tunnels = make(map[ID]Information) //マップ型変数

	/*受信処理*/
	indexbuf := make([]byte, 32)
	En1buf := make([]byte, 15)
	En2buf := make([]byte, 15)
	exit := []byte{101, 120, 105, 116} //101,120,105,116
	for {
		conn.Read(indexbuf) //メッセージを受信
		if reflect.DeepEqual(indexbuf[0:4], exit) == true {
			break
		}
		log.Printf("receive: %s\n", indexbuf)
		conn.Write([]byte("ACK"))

		conn.Read(En1buf) //メッセージを受信
		log.Printf("receive: %s\n", En1buf)
		conn.Write([]byte("ACK"))

		conn.Read(En2buf) //メッセージを受信
		log.Printf("receive: %s\n", En2buf)

		tunnelsIndex := [16]byte{}
		copy(indexbuf[:], tunnelsIndex[0:15]) //スライスから配列に変換
		tunnels[tunnelsIndex] = Information{
			En1: string(En1buf),
			En2: string(En2buf),
		}
		conn.Write([]byte("ACK"))
	}
}
