package utils

import (
	"encoding/hex"
	"errors"
	"hash"
	"io"
)

type CheckSum struct {
	io.Reader
	Hash hash.Hash
	comparator string
}

var ErrChecksum = errors.New("Error verifying checksum")

func (cs *CheckSum) Read(buf []byte) (int, error) {
	n, err := cs.Reader.Read(buf)
	if err == nil {
		cs.Hash.Write(buf[:n])
	}

	// Fail the read if we hit EOF without matching checksum
	if err == io.EOF {
		if cs.Sum() != cs.comparator {
			err = utils.ErrChecksum
		}
	}

	return n, err
}

func (cs *CheckSum) Sum() string {
	return hex.EncodeToString(cs.Hash.Sum(nil))
}
