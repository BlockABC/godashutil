// Copyright (c) 2013-2014 The btcsuite developers
// Copyright (c) 2016 The Dash developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package godashutil_test

import (
	"bytes"
	"io"
	"reflect"
	"testing"

	"github.com/BlockABC/godash/wire"
	"github.com/BlockABC/godashutil"
	"github.com/davecgh/go-spew/spew"
)

// TestTx tests the API for Tx.
func TestTx(t *testing.T) {
	testTx := Block100000.Transactions[0]
	tx := godashutil.NewTx(testTx)

	// Ensure we get the same data back out.
	if msgTx := tx.MsgTx(); !reflect.DeepEqual(msgTx, testTx) {
		t.Errorf("MsgTx: mismatched MsgTx - got %v, want %v",
			spew.Sdump(msgTx), spew.Sdump(testTx))
	}

	// Ensure transaction index set and get work properly.
	wantIndex := 0
	tx.SetIndex(0)
	if gotIndex := tx.Index(); gotIndex != wantIndex {
		t.Errorf("Index: mismatched index - got %v, want %v",
			gotIndex, wantIndex)
	}

	// Hash for block 100,000 transaction 0.
	wantShaStr := "8c14f0db3df150123e6f3dbbf30f8b955a8249b62ac1d1ff16284aefa3d06d87"
	wantSha, err := wire.NewShaHashFromStr(wantShaStr)
	if err != nil {
		t.Errorf("NewShaHashFromStr: %v", err)
	}

	// Request the sha multiple times to test generation and caching.
	for i := 0; i < 2; i++ {
		sha := tx.Sha()
		if !sha.IsEqual(wantSha) {
			t.Errorf("Sha #%d mismatched sha - got %v, want %v", i,
				sha, wantSha)
		}
	}
}

// TestNewTxFromBytes tests creation of a Tx from serialized bytes.
func TestNewTxFromBytes(t *testing.T) {
	// Serialize the test transaction.
	testTx := Block100000.Transactions[0]
	var testTxBuf bytes.Buffer
	err := testTx.Serialize(&testTxBuf)
	if err != nil {
		t.Errorf("Serialize: %v", err)
	}
	testTxBytes := testTxBuf.Bytes()

	// Create a new transaction from the serialized bytes.
	tx, err := godashutil.NewTxFromBytes(testTxBytes)
	if err != nil {
		t.Errorf("NewTxFromBytes: %v", err)
		return
	}

	// Ensure the generated MsgTx is correct.
	if msgTx := tx.MsgTx(); !reflect.DeepEqual(msgTx, testTx) {
		t.Errorf("MsgTx: mismatched MsgTx - got %v, want %v",
			spew.Sdump(msgTx), spew.Sdump(testTx))
	}
}

// TestTxErrors tests the error paths for the Tx API.
func TestTxErrors(t *testing.T) {
	// Serialize the test transaction.
	testTx := Block100000.Transactions[0]
	var testTxBuf bytes.Buffer
	err := testTx.Serialize(&testTxBuf)
	if err != nil {
		t.Errorf("Serialize: %v", err)
	}
	testTxBytes := testTxBuf.Bytes()

	// Truncate the transaction byte buffer to force errors.
	shortBytes := testTxBytes[:4]
	_, err = godashutil.NewTxFromBytes(shortBytes)
	if err != io.EOF {
		t.Errorf("NewTxFromBytes: did not get expected error - "+
			"got %v, want %v", err, io.EOF)
	}
}
