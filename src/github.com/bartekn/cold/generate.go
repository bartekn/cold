package cold

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stellar/go/keypair"
)

func init() {
	RootCmd.AddCommand(generateCmd)
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a new random key pair",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		random, err := keypair.Random()

		if err != nil {
			fmt.Println("Error generating key pair")
			return
		}

		print("Public Key:  " + random.Address() + "\n")
		print("Private Key: " + random.Seed() + "\n")
	},
}
