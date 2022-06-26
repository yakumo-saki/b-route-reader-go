package config

import "os"

var B_ROUTE_ID = ""
var B_ROUTE_PASSWORD = ""
var SERIAL = ""

var EXEC_CMD = ""

var ACTIVE_SCAN_COUNT = 10

// EchonetLiteのGET時にリトライする回数
var MAX_ECHONET_GET_RETRY = 3

// 環境変数からconfigをセット
func Initialize() {
	B_ROUTE_ID = os.Getenv("B_ROUTE_ID")
	B_ROUTE_PASSWORD = os.Getenv("B_ROUTE_PASSWORD")
	SERIAL = os.Getenv("SERIAL")
	EXEC_CMD = os.Getenv("EXEC_CMD")

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

}
