# b-route-reader-go

ROHM BP35A1を使用してBルート経由でスマートメーターから消費電力を取得する。

## 動作環境

* Raspberry Pi Zero WH

開発はRPi Zeroで行っていますが、Linuxが動いてBP35A1とシリアル通信できる環境であれば、
動作すると思います。
## デバッグ・開発

BP35A1と接続された機器上で開発するのも手ですが、色々と面倒です。  
以下のコマンドを使うと、シリアルポートを自分のPCにリダイレクトすることができます。  

* シリアルポートがあるPCのIPアドレスを 10.1.0.191 として説明しています。
* 33444はポート番号です。何番でも構わないので好きな番号にしてください。

Run at serial port machine (eg: raspberry pi)
`socat -d -d /dev/ttyAMA0,echo=0 tcp-listen:33444`

running at your coding pc
`socat -d -d pty,link=$HOME/ttyAMA0,waitslave tcp-connect:10.1.0.191:33444`

### 蛇足

* socatは一度実行されると終了してしまう（なにかおかしい気がする）ので while で無限ループさせる(bash)

`while true; do <コマンド>; done`

### 蛇足２

* socatで転送する前に minicom でハングアップを行っておく必要があります。
* minicom を起動して ctrl-a -> h で `hangup line?` -> `yes` した後、ctrl-a -> q で終了します。（ctrl-a -> x ではない）
* このプログラムを実行すると、ローカルエコーがオフになるのでminicomには何も表示されませんが気にせずコマンドを入れることができます。