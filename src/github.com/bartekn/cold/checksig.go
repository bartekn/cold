package cold

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/fatih/color"
	"github.com/kyokomi/emoji"
	"github.com/spf13/cobra"
	"github.com/stellar/go-stellar-base/keypair"
	"github.com/stellar/go-stellar-base/xdr"
)

func init() {
	RootCmd.AddCommand(checkSigCmd)
}

var checkSigCmd = &cobra.Command{
	Use:   "checksig",
	Short: "Check if signatures added to transaction are valid",
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

		print("Transaction Envelope (base64): ")
		envelope := readString(reader)
		var envelopeXdr xdr.TransactionEnvelope
		err = xdr.SafeUnmarshalBase64(envelope, &envelopeXdr)
		if err != nil {
			fmt.Println("Invalid value", err)
			return
		}

		hash, err := transactionHash(envelopeXdr.Tx, networkPassphrase)
		if err != nil {
			fmt.Println("Error calculating transaction hash")
			return
		}

		var publicKeys []string

		print("Enter expected public keys. One per line, leave empty to continue.\n")

		for {
			publicKey := readString(reader)
			if publicKey == "" {
				break
			}
			publicKeys = append(publicKeys, publicKey)
		}

		header := color.New(color.FgWhite, color.Bold).SprintfFunc()
		valid := color.New(color.FgGreen).SprintFunc()(emoji.Sprintf(":white_check_mark: valid"))
		notfound := color.New(color.FgRed, color.Bold).SprintFunc()(emoji.Sprintf(":x: not found"))

		w := tabwriter.NewWriter(os.Stdout, 20, 56, 2, '\t', 0)
		headerLine := header("Index") + "\t" +
			header("Signature Hint") + "\t" +
			header("Public Key") + "\t" +
			header("Status")
		fmt.Fprintln(w, headerLine)

		for i, signature := range envelopeXdr.Signatures {
			hint := base64.StdEncoding.EncodeToString(signature.Hint[:])
			status := notfound
			foundPublicKey := strings.Repeat(" ", 8) // Table formatting

			for _, publicKey := range publicKeys {
				kp, err := keypair.Parse(publicKey)
				if err != nil {
					fmt.Println("Invalid value", err)
					return
				}

				err = kp.Verify(hash[:], signature.Signature)
				if err != nil {
					continue
				}

				status = valid
				foundPublicKey = publicKey[0:8]
				break
			}

			line := strconv.Itoa(i) + "\t" +
				hint + "\t" +
				foundPublicKey + "\t" +
				status
			fmt.Fprintln(w, line)
		}

		fmt.Fprintln(w)
		w.Flush()
	},
}
