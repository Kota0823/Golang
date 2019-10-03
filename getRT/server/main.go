package main

import (
	"sync" //排他制御用
	"time"
)

type ID [16]byte //マップ型変数で使用するキー

type Information struct {
	Index string
	En1   string
	En2   string
}

/*リレーテーブル(グローバル変数)*/
var tunnels = make(map[ID]Information) //マップ型変数

func main() {
	/*相互排他制御用変数*/
	mutex := new(sync.RWMutex)

	/*ダミーデータ1*/
	mutex.Lock() //書き込みロック
	tunnels[[16]byte{49}] = Information{
		En1: "192.168.100.1",
		En2: "192.168.100.2",
	}
	tunnels[[16]byte{50}] = Information{
		En1: "192.168.100.30",
		En2: "192.168.100.40",
	}
	mutex.Unlock() //書き込みアンロック

	/*並行処理によるリレーテーブルの取得，送信*/
	go GetTable()

	time.Sleep(30 * time.Second)

	/*ダミーデータ2*/
	mutex.Lock() //書き込みロック
	tunnels[[16]byte{51}] = Information{
		En1: "192.168.200.120",
		En2: "192.168.200.111",
	}
	tunnels[[16]byte{52}] = Information{
		En1: "192.168.200.6",
		En2: "192.168.200.87",
	}
	mutex.Unlock() //書き込みアンロック

	for {
		//main関数が死なないように無限ループ
	}

}
