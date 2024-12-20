package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	// "github.com/AnomalyFi/hypersdk/chain"
	"github.com/ava-labs/avalanchego/ids"
)

type SEQTransaction struct {
	Namespace   string `json:"namespace"`
	Tx_id       string `json:"tx_id"`
	Index       uint64 `json:"tx_index"`
	Transaction []byte `json:"transaction"`
}

type SequencerBlock struct {
	StateRoot ids.ID                       `json:"state_root"`
	Prnt      ids.ID                       `json:"parent"`
	Tmstmp    int64                        `json:"timestamp"`
	Hght      uint64                       `json:"height"`
	Txs       map[string][]*SEQTransaction `json:"transactions"`
}

// A BigInt type which serializes to JSON a a hex string.
type U256 struct {
	big.Int
}

func NewU256() *U256 {
	return new(U256)
}

func (i *U256) SetBigInt(n *big.Int) *U256 {
	i.Int.Set(n)
	return i
}

func (i *U256) SetUint64(n uint64) *U256 {
	i.Int.SetUint64(n)
	return i
}

func (i *U256) SetBytes(buf [32]byte) *U256 {
	i.Int.SetBytes(buf[:])
	return i
}

func (i U256) MarshalJSON() ([]byte, error) {
	return json.Marshal(fmt.Sprintf("0x%s", i.Text(16)))
}

func (i *U256) UnmarshalJSON(in []byte) error {
	var s string
	if err := json.Unmarshal(in, &s); err != nil {
		return err
	}
	if _, err := fmt.Sscanf(s, "0x%x", &i.Int); err != nil {
		return err
	}
	return nil
}

type Header struct {
	Height           uint64  `json:"height"`
	Timestamp        uint64  `json:"timestamp"`
	L1Head           uint64  `json:"l1_head"`
	TransactionsRoot NmtRoot `json:"transactions_root"`
}

func (h *Header) UnmarshalJSON(b []byte) error {
	type Dec struct {
		Height           *uint64  `json:"height"`
		Timestamp        *uint64  `json:"timestamp"`
		L1Head           *uint64  `json:"l1_head"`
		TransactionsRoot *NmtRoot `json:"transactions_root"`
	}

	var dec Dec
	if err := json.Unmarshal(b, &dec); err != nil {
		return err
	}

	if dec.Height == nil {
		return errors.New("Field height of type Header is required")
	}
	h.Height = *dec.Height

	if dec.Timestamp == nil {
		return errors.New("Field timestamp of type Header is required")
	}
	h.Timestamp = *dec.Timestamp

	if dec.L1Head == nil {
		return errors.New("Field l1_head of type Header is required")
	}
	h.L1Head = *dec.L1Head

	if dec.TransactionsRoot == nil {
		return errors.New("Field transactions_root of type Header is required")
	}
	h.TransactionsRoot = *dec.TransactionsRoot

	return nil
}

func (self *Header) Commit() Commitment {
	return NewRawCommitmentBuilder("BLOCK").
		Uint64Field("height", self.Height).
		Uint64Field("timestamp", self.Timestamp).
		Uint64Field("l1_head", self.L1Head).
		Field("transactions_root", self.TransactionsRoot.Commit()).
		Finalize()
}

type NmtRoot struct {
	Root Bytes `json:"root"`
}

func (r *NmtRoot) UnmarshalJSON(b []byte) error {
	// Parse using pointers so we can distinguish between missing and default fields.
	type Dec struct {
		Root *Bytes `json:"root"`
	}

	var dec Dec
	if err := json.Unmarshal(b, &dec); err != nil {
		return err
	}

	if dec.Root == nil {
		return errors.New("Field root of type NmtRoot is required")
	}
	r.Root = *dec.Root

	return nil
}

func (self *NmtRoot) Commit() Commitment {
	return NewRawCommitmentBuilder("NMTROOT").
		VarSizeField("root", self.Root).
		Finalize()
}

type Bytes []byte
