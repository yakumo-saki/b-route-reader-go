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
* <- OK

その後、非同期で以下が飛んでくる。

```
EVENT 20 FE80:0000:0000:0000:~~~~
EPANDESC
  Channel:3A           ※ 後で使用する
  Channel Page:09         使わない
  Pan ID:A123          ※ 後で使用する
  Addr:00123456789ABCD ※ 後で使用する
  LQI:F5                  使わない
  PairID:01234567      ※ 後で使用する
EVENT 22 FE80:0000:0000:0000:~~~~
```

## PAN Addr からIPv6アドレスに変換

* -> SKLL64 00123456789ABCD
* <- FE80:~~~

## PAN通信設定のセット

S2は通信チャンネル（周波数） S3はPAN ID

* -> SKSREG S2 3A
* <- OK

* SKSREG S3 A123
* <- OK

## PAN認証

* -> SKJOIN FE80:~~~~
* <- EVENT 21 , EVENT 02 , ERXUDP が複数発生するが無視してよい
* <- wait EVENT 25

## Echonet通信

* -> SKSENDTO 1 FE80:xxxx 0E1A 0 0005 01234

SKSENDTO <HANDLE> <IPADDR> <PORT> <SECURE> <DATALEN> <DATA>

| -------- | ------------------- |
| HANDLE | とりあえず1でよさそう |
| IPADDR | SKLL64で求めたスマートメーターのIPアドレス |
| PORT   |: 16進数で指定。Echonet liteは 3610 = 0E1A |
| SECURE |  0（常に平文）| 
| DATALEN |  送信するデータ長。この数だけDATAを必ず指定するようにと注意書きがある。 16進数なので注意 |
| DATA   | 16進数のASCII表記 | 

## 切断

* -> SKTERM
* <- OK


## 参考

### BP35A1 コマンドリファレンス

* 公式では見れないようだ。
* "BP35A1 コマンドリファレンス" でWeb検索するとPDFが見つかる。

### ECHONET Lite 電文

第２部 ECHONET Lite 通信ミドルウェア仕様
https://echonet.jp/wp/wp-content/uploads/pdf/General/Standard/ECHONET_lite_V1_13_jp/ECHONET-Lite_Ver.1.13_02.pdf
第３章 電文構成（フレームフォーマット）

### 通信の流れ

低圧スマート電力量メータ・HEMS コントローラ間 アプリケーション通信インタフェース仕様書
https://echonet.jp/wp/wp-content/uploads/pdf/General/Standard/AIF/lvsm/lvsm_aif_ver1.01.pdf

### EPC(プロパティコード)の定義と値の意味

ECHONET SPECIFICATION APPENDIX ECHONET 機器オブジェクト詳細規定
３．３．２５ 低圧スマート電力量メータクラス規定
https://echonet.jp/wp/wp-content/uploads/pdf/General/Standard/Release/Release_Q/Appendix_Release_Q.pdf
