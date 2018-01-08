package msg

import (
	"errors"
	"fmt"

	"github.com/eaigner/igi/trinary"
)

const (
	signatureMessageFragmentTrinaryOffset      = 0
	signatureMessageFragmentTrinarySize        = 6561
	addressTrinaryOffset                       = signatureMessageFragmentTrinaryOffset + signatureMessageFragmentTrinarySize
	addressTrinarySize                         = 243
	valueTrinaryOffset                         = addressTrinaryOffset + addressTrinarySize
	valueTrinarySize                           = 81
	valueUsableTrinarySize                     = 33
	obsoleteTagTrinaryOffset                   = valueTrinaryOffset + valueTrinarySize
	obsoleteTagTrinarySize                     = 81
	timestampTrinaryOffset                     = obsoleteTagTrinaryOffset + obsoleteTagTrinarySize
	timestampTrinarySize                       = 27
	currentIndexTrinaryOffset                  = timestampTrinaryOffset + timestampTrinarySize
	currentIndexTrinarySize                    = 27
	lastIndexTrinaryOffset                     = currentIndexTrinaryOffset + currentIndexTrinarySize
	lastIndexTrinarySize                       = 27
	bundleTrinaryOffset                        = lastIndexTrinaryOffset + lastIndexTrinarySize
	bundleTrinarySize                          = 243
	trunkTransactionTrinaryOffset              = bundleTrinaryOffset + bundleTrinarySize
	trunkTransactionTrinarySize                = 243
	branchTransactionTrinaryOffset             = trunkTransactionTrinaryOffset + trunkTransactionTrinarySize
	branchTransactionTrinarySize               = 243
	tagTrinaryOffset                           = branchTransactionTrinaryOffset + branchTransactionTrinarySize
	tagTrinarySize                             = 81
	attachmentTimestampTrinaryOffset           = tagTrinaryOffset + tagTrinarySize
	attachmentTimestampTrinarySize             = 27
	attachmentTimestampLowerBoundTrinaryOffset = attachmentTimestampTrinaryOffset + attachmentTimestampTrinarySize
	attachmentTimestampLowerBoundTrinarySize   = 27
	attachmentTimestampUpperBoundTrinaryOffset = attachmentTimestampLowerBoundTrinaryOffset + attachmentTimestampLowerBoundTrinarySize
	attachmentTimestampUpperBoundTrinarySize   = 27
	nonceTrinaryOffset                         = attachmentTimestampUpperBoundTrinaryOffset + attachmentTimestampUpperBoundTrinarySize
	nonceTrinarySize                           = 81
	trinarySize                                = nonceTrinaryOffset + nonceTrinarySize
	essenceTrinaryOffset                       = addressTrinaryOffset
	essenceTrinarySize                         = addressTrinarySize + valueTrinarySize + obsoleteTagTrinarySize + timestampTrinarySize + currentIndexTrinarySize + lastIndexTrinarySize
)

const (
	hashSizeTrits = 243
)

var (
	errMessageTooShort = errors.New("message too short")
)

type Message struct {
	Raw               []int8
	Address           []int8
	Trunk             []int8
	Branch            []int8
	Bundle            []int8
	Tag               []int8
	ObsoleteTag       []int8
	AttachmentTs      int64
	AttachmentTsUpper int64
	AttachmentTsLower int64
	Value             int64
	Ts                int64
	CurrentIndex      int64
	LastIndex         int64
}

func ParseBytes(b []byte) (*Message, error) {
	t := make([]int8, trinary.LenTrits(b))
	_, err := trinary.Trits(t, b)
	if err != nil {
		return nil, err
	}
	return ParseTrits(t)
}

func ParseTrits(t []int8) (*Message, error) {
	if len(t) < trinarySize {
		return nil, errMessageTooShort
	}
	m := new(Message)
	m.Raw = t
	m.Address = chunk(t, addressTrinaryOffset, addressTrinarySize)
	m.Trunk = chunk(t, trunkTransactionTrinaryOffset, trunkTransactionTrinarySize)
	m.Branch = chunk(t, branchTransactionTrinaryOffset, branchTransactionTrinarySize)
	m.Bundle = chunk(t, bundleTrinaryOffset, bundleTrinarySize)
	m.Tag = chunk(t, tagTrinaryOffset, tagTrinarySize)
	m.ObsoleteTag = chunk(t, obsoleteTagTrinaryOffset, obsoleteTagTrinarySize)
	m.AttachmentTs = chunkInt64(t, attachmentTimestampTrinaryOffset, attachmentTimestampTrinarySize)
	m.AttachmentTsUpper = chunkInt64(t, attachmentTimestampUpperBoundTrinaryOffset, attachmentTimestampUpperBoundTrinarySize)
	m.AttachmentTsLower = chunkInt64(t, attachmentTimestampLowerBoundTrinaryOffset, attachmentTimestampLowerBoundTrinarySize)
	m.Value = chunkInt64(t, valueTrinaryOffset, valueUsableTrinarySize)
	m.Ts = chunkInt64(t, timestampTrinaryOffset, timestampTrinarySize)
	m.CurrentIndex = chunkInt64(t, currentIndexTrinaryOffset, currentIndexTrinarySize)
	m.LastIndex = chunkInt64(t, lastIndexTrinaryOffset, lastIndexTrinarySize)

	return m, nil
}

func (m *Message) Debug() string {
	return fmt.Sprintf("<Message address=%s trunk=%s branch=%s bundle=%s tag=%s otag=%s ats=%v atsh=%v atsl=%v value=%v ts=%v index=%v lindex=%v>",
		trytes(m.Address),
		trytes(m.Trunk),
		trytes(m.Branch),
		trytes(m.Bundle),
		trytes(m.Tag),
		trytes(m.ObsoleteTag),
		m.AttachmentTs,
		m.AttachmentTsUpper,
		m.AttachmentTsLower,
		m.Value,
		m.Ts,
		m.CurrentIndex,
		m.LastIndex,
	)
}

func trytes(t []int8) string {
	s, _ := trinary.Trytes(t) // ignore err, only debug output
	return s
}

func chunk(t []int8, offset, size int) []int8 {
	return t[offset : offset+size]
}

func chunkInt64(t []int8, offset, size int) int64 {
	return trinary.Int64(t[offset : offset+size])
}
