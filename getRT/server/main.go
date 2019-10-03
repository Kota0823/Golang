package main

import (
	"time"

	"../RelayTableLibrary"
	"./DumpRelayTable"
)

func main() {
	/*ダミーデータ*/
	var tunnels = make(map[RelayTableLibrary.ID]RelayTableLibrary.Information) //マップ型変数

	/*リレーテーブルをスレッドに送るためのチャネル（スレッド間通信）*/
	RTChan := make(chan map[RelayTableLibrary.ID]RelayTableLibrary.Information, 100)

	/*ダミーデータ1*/
	tunnels[[16]byte{49}] = RelayTableLibrary.Information{
		En1: "192.168.100.1",
		En2: "192.168.100.2",
	}
	tunnels[[16]byte{50}] = RelayTableLibrary.Information{
		En1: "192.168.100.30",
		En2: "192.168.100.40",
	}
	RTChan <- tunnels //スレッド間通信（スレッドに送るトンネル情報をプッシュ）

	/*並行処理によるリレーテーブルの取得，送信*/
	go DumpRelayTable.GetTable(RTChan)

	time.Sleep(30 * time.Second)

	/*ダミーデータ2*/
	tunnels[[16]byte{51}] = RelayTableLibrary.Information{
		En1: "192.168.200.120",
		En2: "192.168.200.111",
	}

	tunnels[[16]byte{52}] = RelayTableLibrary.Information{
		En1: "192.168.200.6",
		En2: "192.168.200.87",
	}
	RTChan <- tunnels //スレッド間通信

	for {

	}

}
