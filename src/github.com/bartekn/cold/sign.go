package cold

import (
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/stellar/go-stellar-base/keypair"
	"github.com/stellar/go-stellar-base/xdr"
)

func init() {
	RootCmd.AddCommand(signCmd)
}

var signCmd = &cobra.Command{
	Use:   "sign",
	Short: "Add a signature to transaction envelope",
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

		if !printText {
			fmt.Println("Transaction Envelope will be printed as QR code (use -s to print as text)")
		}

		print("Transaction Envelope (base64): ")
		envelope := readString(reader)
		var envelopeXdr xdr.TransactionEnvelope
		err = xdr.SafeUnmarshalBase64(envelope, &envelopeXdr)
		if err != nil {
			fmt.Println("Invalid value", err)
			return
		}

		print("Secret Seed: ")
		secretSeed := readString(reader)
		kp, err := keypair.Parse(secretSeed)
		if err != nil {
			fmt.Println("Invalid value")
			return
		}

		hash, err := transactionHash(envelopeXdr.Tx, networkPassphrase)
		if err != nil {
			fmt.Println("Error calculating transaction hash")
			return
		}

		sig, err := kp.SignDecorated(hash[:])
		if err != nil {
			fmt.Println("Error signing a transaction")
			return
		}

		envelopeXdr.Signatures = append(envelopeXdr.Signatures, sig)

		envelopeBase64, err := xdr.MarshalBase64(envelopeXdr)
		if err != nil {
			fmt.Println("Cannot encode transaction envelope")
			return
		}

		if printText {
			fmt.Println(envelopeBase64)
		} else {
			printQrCode(envelopeBase64)
		}
	},
}
