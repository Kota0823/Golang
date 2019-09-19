package main

import (
	"log"
	"net"
	"os"
	"path/filepath"
)

func main() {
	/*ソケット作成*/
	file := filepath.Join(os.TempDir(), "unixdomaisocketsample") //ソケット作成に使うファイルパス
	defer os.Remove(file)                                        //プログラム終了時にファイルを削除する（実行予約）
	listener, err := net.Listen("unix", file)                    //PIDの代わりにファイルパスを指定
	conn, _ := listener.Accept()
	defer conn.Close() //プログラム終了時にソケットを閉じる（実行予約）
	if err != nil {
		log.Printf("error: %v\n", err) //ソケット作成時エラーが発生した場合
		return
	}

	/*クライアントからリクエスト受信*/
	buf := make([]byte, 20)
	_, err = conn.Read(buf) //ソケットから受信した情報(20字)読み込み
	log.Printf("recive:%s", buf)

	/*プロセス間通信によるリスト型変数送信*/
	type ID [16]byte
	type Information struct {
		En1 string
		En2 string
	}
	var tunnels = make(map[ID]Information) //マップ型変数

	tunnels[[16]byte{0}] = Information{
		En1: "192.168.100.1",
		En2: "192.168.100.2",
	}

	tunnels[[16]byte{32}] = Information{
		En1: "192.168.200.3",
		En2: "192.168.200.4",
	}

	/*送信処理*/
	for index := range tunnels {
		conn.Write(index[:]) //配列([16]byte)をスライスに変換し，メッセージを送信
		conn.Read(buf)

		conn.Write([]byte(tunnels[index].En1)) //メッセージを送信
		log.Printf("send:%s", tunnels[index].En1)
		conn.Read(buf)

		conn.Write([]byte(tunnels[index].En2)) //メッセージを送信
		log.Printf("send:%s", tunnels[index].En2)
		conn.Read(buf)
	}
	conn.Write([]byte("exit"))
}
