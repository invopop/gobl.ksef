package ksef_api

import (
	"github.com/KangSpace/gqrcode"
	"github.com/KangSpace/gqrcode/core/output"
)

func generateQRCode(env KSeFEnv, hash string, referenceNumber string, outFilePath string) (string, error) {
	url := env.Url + "/web/verify/" + referenceNumber + "/" + hash
	qr, err := gqrcode.NewQRCode(url)
	println(url)
	if err == nil {
		out := output.NewPNGOutput0()
		err = qr.Encode(out, outFilePath)
	}
	return url, err
}
