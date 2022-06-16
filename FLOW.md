# 処理の流れ

コマンドリファレンスに書いてあるのをそのままなぞれば良い。

## 接続

SKVERコマンドを投げてOKが返ってくるか確かめる。

* -> SKVER
* <- EVER x.y.zz
* <- OK

## 受信パケット表示モードの確認

後ほどEchonetの通信をする際に、バイナリをそのままバイナリとして表示されると面倒なので
ASCIIモードに切り替える。この設定は、BP35A1のフラッシュに書き込まれるので毎度設定するとフラッシュの寿命が縮む。
このコマンドの影響は ERXUDP と ERXTCP のイベントの出力にのみ影響する

* -> ROPT
* <- OK 00 or <- OK 01

OK 01が返ってきた場合は次に。（WOPTコマンドは不要なので実行しない）
OK 00が返ってきた場合はバイナリモード（初期値）なのでASCIIモードに切り替える

* -> WOPT 01
* <- OK

## Bルートキー設定

これらのID、パスワードは電力会社から送付されてくるもの。
（IDはお手紙、パスワードはメール）

* -> SKSETPWD C <password>
* <- OK

* -> SKSETRBID 00112233445566778899AABBCCDDEEFF ※ ID
* <- OK

## アクティブスキャン

* -> SKSCAN 2 FFFFFFFF 6

## PAN ID からIPv6アドレスに変換

* -> SKLL64 <>
* <- FE80:~~~

## PAN通信設定のセット

S2は通信チャンネル（周波数） S3はPAN ID

* -> SKSREG S2 nn
* <- OK

* SKSREG S3 nnnn
* <- OK

## PAN認証

* -> SKJOIN FE80:~~~~
* <- wait EVENT 25

## Echonet通信

* -> SKSENDTO 1 FE80:xxxx 0E1A 0 0005 01234

SKSENDTO <HANDLE> <IPADDR> <PORT> <SECURE> <DATALEN> <DATA>

| -------- | ------------------- |
| HANDLE | とりあえず1でよさそう |
| IPADDR | SKLL64で求めたスマートメーターのIPアドレス |
| PORT   |: 16進数で指定。Echonet liteは 3610 = 0E1A |
| SECURE |  0（常に平文）| 
| DATALEN |  送信するデータ長。この数だけDATAを必ず指定するようにと注意書きがある。 |
| DATA   | 16進数のASCII表記 | 

## 切断

* -> SKTERM
* <- OK