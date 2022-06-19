package bp35a1

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
)

func sendReset() error {
	err := sendCommand("SKRESET")
	if err != nil {
		log.Error().Msg("Send SKRESET command fail.")
		return err
	}

	ret, err := waitForResult()
	if err != nil {
		return err
	}

	if !containsInResult(ret, RET_OK) {
		return fmt.Errorf("response is not OK")
	}

	return nil
}

// 接続が正しく確立して、コマンドが通るかテストする
// return がnilならOK
func connectionTest() error {
	err := sendCommand("SKVER")
	if err != nil {
		return err
	}

	_, err = waitForResult()
	if err != nil {
		return err
	}

	return nil
}

func isAsciiMode() (bool, error) {
	err := sendCommand("ROPT")
	if err != nil {
		return false, err
	}

	responses, err := waitForResult()
	if err != nil {
		return false, err
	}

	for _, v := range responses {
		if strings.HasPrefix(v, "OK 01") {
			log.Debug().Msg("ROPT OK. Already ASCII mode.")
			return true, nil
		}
	}

	log.Debug().Msg("ROPT OK. BINARY mode.")
	return false, nil
}

func setBrouteId(id string) error {
	err := sendCommand(fmt.Sprintf("SKSETRBID %s", id))
	if err != nil {
		return err
	}

	ret, err := waitForResult()
	if err != nil {
		return err
	}
	if endWithResult(ret, RET_OK) {
		return nil
	}

	return fmt.Errorf("response is not %s", RET_OK)
}

func setBroutePassword(password string) error {
	err := sendCommand(fmt.Sprintf("SKSETPWD C %s", password))
	if err != nil {
		return err
	}

	ret, err := waitForResult()
	if err != nil {
		return err
	}
	if endWithResult(ret, RET_OK) {
		return nil
	}

	return fmt.Errorf("response is not %s", RET_OK)
}

// PAN ADDR -> IPv6アドレス変換
func convertPanIdToIpv6(panAddr string) (string, error) {
	err := sendCommand(fmt.Sprintf("SKLL64 %s", panAddr))
	if err != nil {
		return "", err
	}

	ret, err := waitForResultSKLL64()
	if err != nil {
		return "", err
	}

	for _, v := range ret {
		if strings.HasPrefix(v, "FE80") {
			return v, nil
		}
	}

	return "", fmt.Errorf("command response %s is not expected", strings.Join(ret, ":"))
}
