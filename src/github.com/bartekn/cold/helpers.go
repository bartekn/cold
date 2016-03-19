package cold

import (
	"bytes"
	"fmt"

	"github.com/stellar/go-stellar-base/hash"
	"github.com/stellar/go-stellar-base/xdr"
)

func transactionHash(tx xdr.Transaction, networkPassphrase string) ([32]byte, error) {
	var txBytes bytes.Buffer

	_, err := fmt.Fprintf(&txBytes, "%s", hash.Hash([]byte(networkPassphrase)))
	if err != nil {
		return [32]byte{}, err
	}

	_, err = xdr.Marshal(&txBytes, xdr.EnvelopeTypeEnvelopeTypeTx)
	if err != nil {
		return [32]byte{}, err
	}

	_, err = xdr.Marshal(&txBytes, tx)
	if err != nil {
		return [32]byte{}, err
	}

	return hash.Hash(txBytes.Bytes()), nil
}
