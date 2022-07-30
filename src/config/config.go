package config

import (
	"os"
	"strings"
	"time"
)

var LOG_LEVEL = ""

var B_ROUTE_ID = ""
var B_ROUTE_PASSWORD = ""
var SERIAL = ""

var EXEC_CMD = ""

// SKSCAN を繰り返す回数
var ACTIVE_SCAN_COUNT = 10

// 瞬間消費電力値を取得する間隔（秒）
var NOW_CONSUMPTION_WAIT = 20 * time.Second

// 積算消費電力値を取得する間隔（秒）
var TOTAL_CONSUMPTION_WAIT = 180 * time.Second

// EchonetLiteのGET時にリトライする回数
var MAX_ECHONET_GET_RETRY = 3

// ログに日時を出力するか（systemdで動かす場合はfalse）
var LOG_NO_DATETIME = true

// 環境変数からconfigをセット
func Initialize() {
	LOG_LEVEL = os.Getenv("LOG_LEVEL")
	B_ROUTE_ID = os.Getenv("B_ROUTE_ID")
	B_ROUTE_PASSWORD = os.Getenv("B_ROUTE_PASSWORD")
	SERIAL = os.Getenv("SERIAL")
	EXEC_CMD = os.Getenv("EXEC_CMD")

	if LOG_LEVEL == "" {
		LOG_LEVEL = "INFO"
	} else {
		LOG_LEVEL = strings.ToUpper(LOG_LEVEL)
		switch LOG_LEVEL {
		case "DEBUG":
		case "INFO":
		case "WARN":
		case "ERROR":
		default:
			panic("LOG_LEVEL is not in [DEBUG|INFO|WARN|ERROR]")
		}
	}

	if B_ROUTE_ID == "" {
		panic("B_ROUTE_ID env value is not set")
	}
	if B_ROUTE_PASSWORD == "" {
		panic("B_ROUTE_PASSWORD env value is not set")
	}
	if SERIAL == "" {
		panic("SERIAL env value is not set")
	}
	if EXEC_CMD == "" {
		panic("EXEC_CMD env value is not set")
	}

	noDatetime := os.Getenv("LOG_NO_DATETIME")
	LOG_NO_DATETIME = (strings.ToLower(noDatetime) == "true")

}
