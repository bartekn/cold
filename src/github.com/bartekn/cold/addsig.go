package cold

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/xdr"
)

func init() {
	RootCmd.AddCommand(addsigCmd)
}

var addsigCmd = &cobra.Command{
	Use:   "addsig",
	Short: "Add existing signature to transaction envelope",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)

		if !printText {
			fmt.Println("Transaction Envelope will be printed as QR code (use -s to print as text)")
		}

		print("Transaction Envelope (base64): ")
		envelope := readString(reader)
		var envelopeXdr xdr.TransactionEnvelope
		err := xdr.SafeUnmarshalBase64(envelope, &envelopeXdr)
		if err != nil {
			fmt.Println("Invalid value", err)
			return
		}

		print("Public key: ")
		publicKey := readString(reader)
		kp, err := keypair.Parse(publicKey)
		if err != nil {
			fmt.Println("Invalid value")
			return
		}

		print("Signature (base64): ")
		signature := readString(reader)

		signatureBytes, err := base64.StdEncoding.DecodeString(signature)
		if err != nil {
			fmt.Println("Signature is not invalid: ", err)
			return
		}

		if len(signatureBytes) != 64 {
			fmt.Println("Invalid signature length")
			return
		}

		decoratedSignature := xdr.DecoratedSignature{
			Hint:      xdr.SignatureHint(kp.Hint()),
			Signature: signatureBytes,
		}

		envelopeXdr.Signatures = append(envelopeXdr.Signatures, decoratedSignature)

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
