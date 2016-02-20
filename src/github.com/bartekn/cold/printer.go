package cold

import (
	"fmt"
	"os"

	"github.com/fumiyas/go-tty"
	"github.com/fumiyas/qrc/lib"
	"github.com/kyokomi/emoji"
	"github.com/mattn/go-colorable"
	"github.com/qpliu/qrencode-go/qrencode"
)

func print(value string) {
	fmt.Print(emoji.Sprintf(":snowflake: " + value))
}

func printQrCode(value string) error {
	grid, err := qrencode.Encode(value, qrencode.ECLevelM)
	if err != nil {
		return err
	}

	da1, err := tty.GetDeviceAttributes1(os.Stdout)
	if err == nil && da1[tty.DA1_SIXEL] {
		qrc.PrintSixel(os.Stdout, grid, false)
	} else {
		stdout := colorable.NewColorableStdout()
		qrc.PrintAA(stdout, grid, false)
	}
	return nil
}
