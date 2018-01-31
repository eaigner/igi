package node

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/eaigner/igi/storage"

	"github.com/eaigner/igi/hash"
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
	hashSizeTrits       = 243
	udpPacketBytes      = 1650
	hashTrailerBytes    = 46
	txnPacketBytes      = udpPacketBytes - hashTrailerBytes
	hashesInvalidBefore = 1508760000
)

var (
	errMessageTooShort    = errors.New("message too short")
	errTxAlreadyExists    = errors.New("transaction already exists")
	errInvalidTxTimestamp = errors.New("invalid transaction timestamp")
	errInvalidTxValue     = errors.New("invalid transaction value")
	errInvalidTxHash      = errors.New("invalid transaction hash")
	errInvalidTxAddress   = errors.New("invalid transaction address")
)

type Message struct {
	TxBytes           []byte // Raw transaction bytes
	TxTrits           []int8 // Raw transaction trits
	Address           []int8 // Address trits
	Trunk             []int8 // Trunk address trits
	Branch            []int8 // Branch address trits
	Bundle            []int8 // Bundle address trits
	Tag               []int8 // Tag
	ObsoleteTag       []int8
	Nonce             []int8 // Nonce
	ValueTrailer      []int8 // Trits after usable value
	AttachmentTs      int64  // Attachment timestamp
	AttachmentTsUpper int64  // Attachment timestamp upper bound
	AttachmentTsLower int64  // Attachment timestamp lower bound
	Value             int64  // Transaction value
	Ts                int64
	CurrentIndex      int64
	LastIndex         int64
	Trailer           []byte // UDP packet trailer. Only set if message was read with ParseUdpBytes.

	txHash      []int8
	trailerHash []int8
	digest      []byte // SHA-256 digest of the transaction packet, excluding trailer bytes
}

// ParseUdpBytes parses a transaction including UDP trailer bytes.
func ParseUdpBytes(b []byte) (*Message, error) {
	if len(b) != udpPacketBytes {
		return nil, errMessageTooShort
	}
	m, err := ParseTxBytes(b[:txnPacketBytes])
	if err != nil {
		return nil, err
	}
	m.Trailer = b[txnPacketBytes:]
	return m, nil
}

// ParseTxBytes parses transaction bytes without an UDP trailer.
func ParseTxBytes(b []byte) (*Message, error) {
	if len(b) != txnPacketBytes {
		return nil, errMessageTooShort
	}

	t := make([]int8, trinary.LenTrits(len(b)))
	_, err := trinary.Trits(t, b)
	if err != nil {
		return nil, err
	}

	m := new(Message)
	m.TxBytes = b
	m.TxTrits = t
	m.Address = chunk(t, addressTrinaryOffset, addressTrinarySize)
	m.Trunk = chunk(t, trunkTransactionTrinaryOffset, trunkTransactionTrinarySize)
	m.Branch = chunk(t, branchTransactionTrinaryOffset, branchTransactionTrinarySize)
	m.Bundle = chunk(t, bundleTrinaryOffset, bundleTrinarySize)
	m.Tag = chunk(t, tagTrinaryOffset, tagTrinarySize)
	m.ObsoleteTag = chunk(t, obsoleteTagTrinaryOffset, obsoleteTagTrinarySize)
	m.Nonce = chunk(t, nonceTrinaryOffset, nonceTrinarySize)
	m.ValueTrailer = chunk(t, valueTrinaryOffset+valueUsableTrinarySize, valueTrinarySize-valueUsableTrinarySize)
	m.AttachmentTs = chunkInt64(t, attachmentTimestampTrinaryOffset, attachmentTimestampTrinarySize)
	m.AttachmentTsUpper = chunkInt64(t, attachmentTimestampUpperBoundTrinaryOffset, attachmentTimestampUpperBoundTrinarySize)
	m.AttachmentTsLower = chunkInt64(t, attachmentTimestampLowerBoundTrinaryOffset, attachmentTimestampLowerBoundTrinarySize)
	m.Value = chunkInt64(t, valueTrinaryOffset, valueUsableTrinarySize)
	m.Ts = chunkInt64(t, timestampTrinaryOffset, timestampTrinarySize)
	m.CurrentIndex = chunkInt64(t, currentIndexTrinaryOffset, currentIndexTrinarySize)
	m.LastIndex = chunkInt64(t, lastIndexTrinaryOffset, lastIndexTrinarySize)

	return m, nil
}

func (m *Message) TxDigest() []byte {
	if len(m.digest) < sha256.Size {
		d := sha256.Sum256(m.TxBytes)
		m.digest = d[:]
	}
	return m.digest
}

func (m *Message) TxHash() []int8 {
	if len(m.txHash) > 0 {
		return m.txHash
	}

	m.txHash = make([]int8, hash.SizeTrits)

	var curl hash.Curl
	curl.Reset(hash.CurlP81)
	curl.Absorb(m.TxTrits[:trinarySize])
	curl.Squeeze(m.txHash)

	return m.txHash
}

func (m *Message) TrailerHash() []int8 {
	if len(m.trailerHash) > 0 {
		return m.trailerHash
	}

	m.trailerHash = make([]int8, hash.SizeTrits)

	_, err := trinary.Trits(m.trailerHash, m.Trailer)
	if err != nil {
		return nil
	}
	return m.trailerHash
}

func (m Message) Validate(minWeightMag int) error {
	txHash := m.TxHash()

	// Check weight magnitude
	if hash.WeightMagnitude(txHash[:]) < minWeightMag {
		return errInvalidTxHash
	}

	// Every non-zero hash before 'Mon 23rd Oct 2017 12:00:00 PM' is invalid.
	// Taken from IRI.
	if m.Ts < hashesInvalidBefore && hash.ZeroInt8(txHash[:]) {
		return errInvalidTxTimestamp
	}

	// Check if trits after value are zero.
	for _, v := range m.ValueTrailer {
		if v != 0 {
			return errInvalidTxValue
		}
	}

	// Check if last address trit is zero.
	if m.Value != 0 && len(m.Address) == hash.SizeTrits && m.Address[hash.SizeTrits-1] != 0 {
		return errInvalidTxAddress
	}

	return nil
}

// Store stores the message in the tangle.
// Returns an error if storage failed or the transaction already exists.
func (m Message) Store(tangle storage.Store) error {
	if !hash.ValidInt8(m.TxHash()) {
		return errInvalidTxHash
	}

	txHash := hash.ToBytes(m.TxHash())

	// TODO(era): make exists and write check atomic
	exists, err := storage.Exists(tangle, txHash, storage.TransactionBucket)

	if err != nil {
		return err
	}
	if exists {
		return errTxAlreadyExists
	}

	return storage.Write(tangle, txHash, m.TxBytes, storage.TransactionBucket)
}

func (m Message) AddressTrytes() string {
	return toTryte(m.Address)
}

func (m Message) BundleTrytes() string {
	return toTryte(m.Bundle)
}

func (m Message) TrunkTrytes() string {
	return toTryte(m.Trunk)
}

func (m Message) BranchTrytes() string {
	return toTryte(m.Branch)
}

func (m Message) ObsoleteTagTrytes() string {
	return toTryte(m.ObsoleteTag)
}

func (m Message) TagTrytes() string {
	return toTryte(m.Tag)
}

func (m Message) NonceTrytes() string {
	return toTryte(m.Nonce)
}

func toTryte(t []int8) string {
	s, _ := trinary.Trytes(t) // ignore err
	return s
}

func (m *Message) Debug() string {
	return fmt.Sprintf("<Message address=%s trunk=%s branch=%s bundle=%s tag=%s otag=%s nonce=%s ats=%v atsh=%v atsl=%v value=%v ts=%v index=%v lindex=%v trailer=%s digest=%s>",
		trytes(m.Address),
		trytes(m.Trunk),
		trytes(m.Branch),
		trytes(m.Bundle),
		trytes(m.Tag),
		trytes(m.ObsoleteTag),
		trytes(m.Nonce),
		m.AttachmentTs,
		m.AttachmentTsUpper,
		m.AttachmentTsLower,
		m.Value,
		m.Ts,
		m.CurrentIndex,
		m.LastIndex,
		hex.EncodeToString(m.Trailer),
		hex.EncodeToString(m.TxDigest()),
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
