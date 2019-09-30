/*
外部パッケージで変数を宣言する場合は
変数名を大文字スタートで！

変数の宣言
*/

package RelayTableLibrary

import (
	"log"
	"net"
	"os"
	"path/filepath"
)

/*リレーテーブル*/
type ID [16]byte //マップ型変数で使用するキー
type Information struct {
	Index string
	En1   string
	En2   string
}

/*ソケット作成に使うファイルパス*/
var SocketFilepath = filepath.Join(os.TempDir(), "unixdomaisocketsample")

/*制御コード*/
const (
	RequestRelayTable int = iota //0
	Exit                         //1
	Pass                         //2
)

var EXIT = []byte{101, 120, 105, 116} //"exit"の文字コード

func SendMessage(conn net.Conn, message []byte) (err error) {
	_, err = conn.Write(message) //配列([16]byte)をスライスに変換し，メッセージを送信
	log.Printf("send: %s", string(message))
	if err != nil {
		return
	}

	_, err = conn.Read(make([]byte, 3)) //応答(ACK)を受信
	log.Printf("receive:ACK\n")
	return
}

func ResiveMessage(conn net.Conn) (message string) {
	buf := make([]byte, 20)
	_, err := conn.Read(buf) //配列([16]byte)をスライスに変換し，メッセージを送信
	if err != nil {
		log.Printf("error: %v\n", err) //ソケット作成時エラーが発生した場合
		return ""
	}
	log.Printf("receive: %s", string(buf))

	_, err = conn.Write([]byte("ACK")) //応答
	if err != nil {
		log.Printf("error: %v\n", err) //ソケット作成時エラーが発生した場合
		return ""
	}
	log.Printf("send: ACK\n")

	return string(buf)
}
