package cold

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/stellar/go-stellar-base/build"
	"github.com/stellar/go-stellar-base/keypair"
)

func init() {
	RootCmd.AddCommand(paymentCmd)
}

var paymentCmd = &cobra.Command{
	Use:   "payment",
	Short: "Creates a transaction with Payment operation",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)

		networkPassphrase, err := getNetworkPassphrase()
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Println(
			"Using network:",
			networkPassphrase,
			"(use --network to change)",
		)

		print("Sequence Number (+1): ")
		sequenceNumber, err := readUint64(reader)
		if err != nil {
			fmt.Println("Invalid value")
			return
		}

		print("Destination: ")
		destination := readString(reader)
		_, err = keypair.Parse(destination)
		if err != nil {
			fmt.Println("Invalid value")
			return
		}

		print("Amount: ")
		amount := readString(reader)

		print("Secret Seed: ")
		secretSeed := readString(reader)
		_, err = keypair.Parse(secretSeed)
		if err != nil {
			fmt.Println("Invalid value")
			return
		}

		tx := build.Transaction(
			build.SourceAccount{secretSeed},
			build.Sequence{sequenceNumber},
			build.Network{networkPassphrase},
			build.Payment(
				build.Destination{destination},
				build.NativeAmount{amount},
			),
		)

		if tx.Err != nil {
			fmt.Println("Error building transaction ", tx.Err)
			return
		}

		txe := tx.Sign(secretSeed)
		txeB64, err := txe.Base64()
		if err != nil {
			fmt.Println("Cannot encode transaction envelope ", err)
			return
		}

		if printText {
			fmt.Println(txeB64)
		} else {
			printQrCode(txeB64)
		}
	},
}

func readString(reader *bufio.Reader) string {
	line, _ := reader.ReadString('\n')
	return strings.TrimRight(line, "\n")
}

func readUint64(reader *bufio.Reader) (uint64, error) {
	line := readString(reader)
	return strconv.ParseUint(line, 10, 64)
}
