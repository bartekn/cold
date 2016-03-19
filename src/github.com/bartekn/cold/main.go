package cold

import (
	"fmt"

	"github.com/kyokomi/emoji"
	"github.com/spf13/cobra"
	n "github.com/stellar/go-stellar-base/network"
)

var network string
var printText bool
var RootCmd = &cobra.Command{
	Use: "cold",
	Short: emoji.Sprint(
		` 

   :snowman: STELLAR COLD WALLET :snowman:
  
:snowflake::snowflake::snowflake::snowflake::snowflake::snowflake::snowflake::snowflake::snowflake::snowflake::snowflake::snowflake:`),
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&network, "network", "n", "test", "[public|test] network to use")
	RootCmd.PersistentFlags().BoolVarP(&printText, "print-string", "s", false, "when set, transaction envelope will be printed as string instead of QR code")
}

func getNetworkPassphrase() (passphrase string, err error) {
	switch network {
	case "public":
		passphrase = n.PublicNetworkPassphrase
	case "test":
		passphrase = n.TestNetworkPassphrase
	default:
		err = fmt.Errorf("Unknown network: %s", network)
	}
	return
}
