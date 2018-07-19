// Copyright (c) 2013, 2014 The btcsuite developers
// Copyright (c) 2015 The Decred developers
// Copyright (c) 2018-2020 The Hcd developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package hcutil_test

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"reflect"
	"testing"

	"github.com/HcashOrg/hcd/chaincfg"
	"github.com/HcashOrg/hcd/chaincfg/chainec"
	"github.com/HcashOrg/hcd/wire"
	"github.com/HcashOrg/hcutil"
	"golang.org/x/crypto/ripemd160"
)

// invalidNet is an invalid network.
const invalidNet = wire.CurrencyNet(0xffffffff)

func TestAddresses(t *testing.T) {
	tests := []struct {
		name    string
		addr    string
		saddr   string
		encoded string
		valid   bool
		result  hcutil.Address
		f       func() (hcutil.Address, error)
		net     *chaincfg.Params
	}{
		// Positive P2PKH tests.
		{
			name:    "mainnet p2pkh",
			addr:    "DsUZxxoHJSty8DCfwfartwTYbuhmVct7tJu",
			encoded: "DsUZxxoHJSty8DCfwfartwTYbuhmVct7tJu",
			valid:   true,
			result: hcutil.TstAddressPubKeyHash(
				[ripemd160.Size]byte{
					0x27, 0x89, 0xd5, 0x8c, 0xfa, 0x09, 0x57, 0xd2, 0x06, 0xf0,
					0x25, 0xc2, 0xaf, 0x05, 0x6f, 0xc8, 0xa7, 0x7c, 0xeb, 0xb0},

				chaincfg.MainNetParams.PubKeyHashAddrID),
			f: func() (hcutil.Address, error) {
				pkHash := []byte{
					0x27, 0x89, 0xd5, 0x8c, 0xfa, 0x09, 0x57, 0xd2, 0x06, 0xf0,
					0x25, 0xc2, 0xaf, 0x05, 0x6f, 0xc8, 0xa7, 0x7c, 0xeb, 0xb0}
				return hcutil.NewAddressPubKeyHash(pkHash,
					&chaincfg.MainNetParams, chainec.ECTypeSecp256k1)
			},
			net: &chaincfg.MainNetParams,
		},
		{
			name:    "mainnet p2pkh 2",
			addr:    "DsU7xcg53nxaKLLcAUSKyRndjG78Z2VZnX9",
			encoded: "DsU7xcg53nxaKLLcAUSKyRndjG78Z2VZnX9",
			valid:   true,
			result: hcutil.TstAddressPubKeyHash(
				[ripemd160.Size]byte{
					0x22, 0x9e, 0xba, 0xc3, 0x0e, 0xfd, 0x6a, 0x69, 0xee, 0xc9,
					0xc1, 0xa4, 0x8e, 0x04, 0x8b, 0x7c, 0x97, 0x5c, 0x25, 0xf2},
				chaincfg.MainNetParams.PubKeyHashAddrID),
			f: func() (hcutil.Address, error) {
				pkHash := []byte{
					0x22, 0x9e, 0xba, 0xc3, 0x0e, 0xfd, 0x6a, 0x69, 0xee, 0xc9,
					0xc1, 0xa4, 0x8e, 0x04, 0x8b, 0x7c, 0x97, 0x5c, 0x25, 0xf2}
				return hcutil.NewAddressPubKeyHash(pkHash,
					&chaincfg.MainNetParams, chainec.ECTypeSecp256k1)
			},
			net: &chaincfg.MainNetParams,
		},
		{
			name:    "testnet p2pkh",
			addr:    "Tso2MVTUeVrjHTBFedFhiyM7yVTbieqp91h",
			encoded: "Tso2MVTUeVrjHTBFedFhiyM7yVTbieqp91h",
			valid:   true,
			result: hcutil.TstAddressPubKeyHash(
				[ripemd160.Size]byte{
					0xf1, 0x5d, 0xa1, 0xcb, 0x8d, 0x1b, 0xcb, 0x16, 0x2c, 0x6a,
					0xb4, 0x46, 0xc9, 0x57, 0x57, 0xa6, 0xe7, 0x91, 0xc9, 0x16},
				chaincfg.TestNet2Params.PubKeyHashAddrID),
			f: func() (hcutil.Address, error) {
				pkHash := []byte{
					0xf1, 0x5d, 0xa1, 0xcb, 0x8d, 0x1b, 0xcb, 0x16, 0x2c, 0x6a,
					0xb4, 0x46, 0xc9, 0x57, 0x57, 0xa6, 0xe7, 0x91, 0xc9, 0x16}
				return hcutil.NewAddressPubKeyHash(pkHash,
					&chaincfg.TestNet2Params, chainec.ECTypeSecp256k1)
			},
			net: &chaincfg.TestNet2Params,
		},

		// Negative P2PKH tests.
		{
			name:  "p2pkh wrong hash length",
			addr:  "",
			valid: false,
			f: func() (hcutil.Address, error) {
				pkHash := []byte{
					0x00, 0x0e, 0xf0, 0x30, 0x10, 0x7f, 0xd2, 0x6e, 0x0b, 0x6b,
					0xf4, 0x05, 0x12, 0xbc, 0xa2, 0xce, 0xb1, 0xdd, 0x80, 0xad,
					0xaa}
				return hcutil.NewAddressPubKeyHash(pkHash,
					&chaincfg.MainNetParams,
					chainec.ECTypeSecp256k1)
			},
		},
		{
			name:  "p2pkh bad checksum",
			addr:  "TsmWaPM77WSyA3aiQ2Q1KnwGDVWvEkhip23",
			valid: false,
			net:   &chaincfg.TestNet2Params,
		},

		// Positive P2SH tests.
		{
			// Taken from transactions:
			// output: 3c9018e8d5615c306d72397f8f5eef44308c98fb576a88e030c25456b4f3a7ac
			// input:  837dea37ddc8b1e3ce646f1a656e79bbd8cc7f558ac56a169626d649ebe2a3ba.
			name:    "mainnet p2sh",
			addr:    "DcuQKx8BES9wU7C6Q5VmLBjw436r27hayjS",
			encoded: "DcuQKx8BES9wU7C6Q5VmLBjw436r27hayjS",
			valid:   true,
			result: hcutil.TstAddressScriptHash(
				[ripemd160.Size]byte{
					0xf0, 0xb4, 0xe8, 0x51, 0x00, 0xae, 0xe1, 0xa9, 0x96, 0xf2,
					0x29, 0x15, 0xeb, 0x3c, 0x3f, 0x76, 0x4d, 0x53, 0x77, 0x9a},
				chaincfg.MainNetParams.ScriptHashAddrID),
			f: func() (hcutil.Address, error) {
				txscript := []byte{
					0x51, 0x21, 0x03, 0xaa, 0x43, 0xf0, 0xa6, 0xc1, 0x57, 0x30,
					0xd8, 0x86, 0xcc, 0x1f, 0x03, 0x42, 0x04, 0x6d, 0x20, 0x17,
					0x54, 0x83, 0xd9, 0x0d, 0x7c, 0xcb, 0x65, 0x7f, 0x90, 0xc4,
					0x89, 0x11, 0x1d, 0x79, 0x4c, 0x51, 0xae}
				return hcutil.NewAddressScriptHash(txscript, &chaincfg.MainNetParams)
			},
			net: &chaincfg.MainNetParams,
		},
		{
			// Taken from transactions:
			// output: b0539a45de13b3e0403909b8bd1a555b8cbe45fd4e3f3fda76f3a5f52835c29d
			// input: (not yet redeemed at time test was written)
			name:    "mainnet p2sh 2",
			addr:    "DcqgK4N4Ccucu2Sq4VDAdu4wH4LASLhzLVp",
			encoded: "DcqgK4N4Ccucu2Sq4VDAdu4wH4LASLhzLVp",
			valid:   true,
			result: hcutil.TstAddressScriptHash(
				[ripemd160.Size]byte{
					0xc7, 0xda, 0x50, 0x95, 0x68, 0x34, 0x36, 0xf4, 0x43, 0x5f,
					0xc4, 0xe7, 0x16, 0x3d, 0xca, 0xfd, 0xa1, 0xa2, 0xd0, 0x07},
				chaincfg.MainNetParams.ScriptHashAddrID),
			f: func() (hcutil.Address, error) {
				hash := []byte{
					0xc7, 0xda, 0x50, 0x95, 0x68, 0x34, 0x36, 0xf4, 0x43, 0x5f,
					0xc4, 0xe7, 0x16, 0x3d, 0xca, 0xfd, 0xa1, 0xa2, 0xd0, 0x07}
				return hcutil.NewAddressScriptHashFromHash(hash, &chaincfg.MainNetParams)
			},
			net: &chaincfg.MainNetParams,
		},
		{
			// Taken from bitcoind base58_keys_valid.
			name:    "testnet p2sh",
			addr:    "TccWLgcquqvwrfBocq5mcK5kBiyw8MvyvCi",
			encoded: "TccWLgcquqvwrfBocq5mcK5kBiyw8MvyvCi",
			valid:   true,
			result: hcutil.TstAddressScriptHash(
				[ripemd160.Size]byte{
					0x36, 0xc1, 0xca, 0x10, 0xa8, 0xa6, 0xa4, 0xb5, 0xd4, 0x20,
					0x4a, 0xc9, 0x70, 0x85, 0x39, 0x79, 0x90, 0x3a, 0xa2, 0x84},
				chaincfg.TestNet2Params.ScriptHashAddrID),
			f: func() (hcutil.Address, error) {
				hash := []byte{
					0x36, 0xc1, 0xca, 0x10, 0xa8, 0xa6, 0xa4, 0xb5, 0xd4, 0x20,
					0x4a, 0xc9, 0x70, 0x85, 0x39, 0x79, 0x90, 0x3a, 0xa2, 0x84}
				return hcutil.NewAddressScriptHashFromHash(hash, &chaincfg.TestNet2Params)
			},
			net: &chaincfg.TestNet2Params,
		},

		// Negative P2SH tests.
		{
			name:  "p2sh wrong hash length",
			addr:  "",
			valid: false,
			f: func() (hcutil.Address, error) {
				hash := []byte{
					0x00, 0xf8, 0x15, 0xb0, 0x36, 0xd9, 0xbb, 0xbc, 0xe5, 0xe9,
					0xf2, 0xa0, 0x0a, 0xbd, 0x1b, 0xf3, 0xdc, 0x91, 0xe9, 0x55,
					0x10}
				return hcutil.NewAddressScriptHashFromHash(hash, &chaincfg.MainNetParams)
			},
			net: &chaincfg.MainNetParams,
		},

		// Positive P2PK tests.
		{
			name:    "mainnet p2pk compressed (0x02)",
			addr:    "DsT4FDqBKYG1Xr8aGrT1rKP3kiv6TZ5K5th",
			encoded: "DsT4FDqBKYG1Xr8aGrT1rKP3kiv6TZ5K5th",
			valid:   true,
			result: hcutil.TstAddressPubKey(
				[]byte{
					0x02, 0x8f, 0x53, 0x83, 0x8b, 0x76, 0x39, 0x56, 0x3f, 0x27,
					0xc9, 0x48, 0x45, 0x54, 0x9a, 0x41, 0xe5, 0x14, 0x6b, 0xcd,
					0x52, 0xe7, 0xfe, 0xf0, 0xea, 0x6d, 0xa1, 0x43, 0xa0, 0x2b,
					0x0f, 0xe2, 0xed},
				hcutil.PKFCompressed, chaincfg.MainNetParams.PubKeyHashAddrID),
			f: func() (hcutil.Address, error) {
				serializedPubKey := []byte{
					0x02, 0x8f, 0x53, 0x83, 0x8b, 0x76, 0x39, 0x56, 0x3f, 0x27,
					0xc9, 0x48, 0x45, 0x54, 0x9a, 0x41, 0xe5, 0x14, 0x6b, 0xcd,
					0x52, 0xe7, 0xfe, 0xf0, 0xea, 0x6d, 0xa1, 0x43, 0xa0, 0x2b,
					0x0f, 0xe2, 0xed}
				return hcutil.NewAddressSecpPubKey(serializedPubKey, &chaincfg.MainNetParams)
			},
			net: &chaincfg.MainNetParams,
		},
		{
			name:    "mainnet p2pk compressed (0x03)",
			addr:    "DsfiE2y23CGwKNxSGjbfPGeEW4xw1tamZdc",
			encoded: "DsfiE2y23CGwKNxSGjbfPGeEW4xw1tamZdc",
			valid:   true,
			result: hcutil.TstAddressPubKey(
				[]byte{
					0x03, 0xe9, 0x25, 0xaa, 0xfc, 0x1e, 0xdd, 0x44, 0xe7, 0xc7,
					0xf1, 0xea, 0x4f, 0xb7, 0xd2, 0x65, 0xdc, 0x67, 0x2f, 0x20,
					0x4c, 0x3d, 0x0c, 0x81, 0x93, 0x03, 0x89, 0xc1, 0x0b, 0x81,
					0xfb, 0x75, 0xde},
				hcutil.PKFCompressed, chaincfg.MainNetParams.PubKeyHashAddrID),
			f: func() (hcutil.Address, error) {
				serializedPubKey := []byte{
					0x03, 0xe9, 0x25, 0xaa, 0xfc, 0x1e, 0xdd, 0x44, 0xe7, 0xc7,
					0xf1, 0xea, 0x4f, 0xb7, 0xd2, 0x65, 0xdc, 0x67, 0x2f, 0x20,
					0x4c, 0x3d, 0x0c, 0x81, 0x93, 0x03, 0x89, 0xc1, 0x0b, 0x81,
					0xfb, 0x75, 0xde}
				return hcutil.NewAddressSecpPubKey(serializedPubKey, &chaincfg.MainNetParams)
			},
			net: &chaincfg.MainNetParams,
		},
		// Hybrid, uncompressed and compressed key types are supported, hcd consensus rules require a compressed key type however.
		{
			name:    "mainnet p2pk uncompressed (0x04)",
			addr:    "DkM3EyZ546GghVSkvzb6J47PvGDyntqiDtFgipQhNj78Xm2mUYRpf",
			encoded: "DsfFjaADsV8c5oHWx85ZqfxCZy74K8RFuhK",
			valid:   true,
			saddr:   "0264c44653d6567eff5753c5d24a682ddc2b2cadfe1b0c6433b16374dace6778f0",
			result: hcutil.TstAddressPubKey(
				[]byte{
					0x04, 0x64, 0xc4, 0x46, 0x53, 0xd6, 0x56, 0x7e, 0xff, 0x57,
					0x53, 0xc5, 0xd2, 0x4a, 0x68, 0x2d, 0xdc, 0x2b, 0x2c, 0xad,
					0xfe, 0x1b, 0x0c, 0x64, 0x33, 0xb1, 0x63, 0x74, 0xda, 0xce,
					0x67, 0x78, 0xf0, 0xb8, 0x7c, 0xa4, 0x27, 0x9b, 0x56, 0x5d,
					0x21, 0x30, 0xce, 0x59, 0xf7, 0x5b, 0xfb, 0xb2, 0xb8, 0x8d,
					0xa7, 0x94, 0x14, 0x3d, 0x7c, 0xfd, 0x3e, 0x80, 0x80, 0x8a,
					0x1f, 0xa3, 0x20, 0x39, 0x04},
				hcutil.PKFUncompressed, chaincfg.MainNetParams.PubKeyHashAddrID),
			f: func() (hcutil.Address, error) {
				serializedPubKey := []byte{
					0x04, 0x64, 0xc4, 0x46, 0x53, 0xd6, 0x56, 0x7e, 0xff, 0x57,
					0x53, 0xc5, 0xd2, 0x4a, 0x68, 0x2d, 0xdc, 0x2b, 0x2c, 0xad,
					0xfe, 0x1b, 0x0c, 0x64, 0x33, 0xb1, 0x63, 0x74, 0xda, 0xce,
					0x67, 0x78, 0xf0, 0xb8, 0x7c, 0xa4, 0x27, 0x9b, 0x56, 0x5d,
					0x21, 0x30, 0xce, 0x59, 0xf7, 0x5b, 0xfb, 0xb2, 0xb8, 0x8d,
					0xa7, 0x94, 0x14, 0x3d, 0x7c, 0xfd, 0x3e, 0x80, 0x80, 0x8a,
					0x1f, 0xa3, 0x20, 0x39, 0x04}
				return hcutil.NewAddressSecpPubKey(serializedPubKey, &chaincfg.MainNetParams)
			},
			net: &chaincfg.MainNetParams,
		},
		{
			name:    "mainnet p2pk hybrid (0x06)",
			addr:    "DkM3EyZ546GghVSkvzb6J47PvGDyntqiDtFgipQhNj78Xm2mUYRpf",
			encoded: "DsfFjaADsV8c5oHWx85ZqfxCZy74K8RFuhK",
			valid:   true,
			saddr:   "0264c44653d6567eff5753c5d24a682ddc2b2cadfe1b0c6433b16374dace6778f0",
			result: hcutil.TstAddressPubKey(
				[]byte{
					0x06, 0x64, 0xc4, 0x46, 0x53, 0xd6, 0x56, 0x7e, 0xff, 0x57,
					0x53, 0xc5, 0xd2, 0x4a, 0x68, 0x2d, 0xdc, 0x2b, 0x2c, 0xad,
					0xfe, 0x1b, 0x0c, 0x64, 0x33, 0xb1, 0x63, 0x74, 0xda, 0xce,
					0x67, 0x78, 0xf0, 0xb8, 0x7c, 0xa4, 0x27, 0x9b, 0x56, 0x5d,
					0x21, 0x30, 0xce, 0x59, 0xf7, 0x5b, 0xfb, 0xb2, 0xb8, 0x8d,
					0xa7, 0x94, 0x14, 0x3d, 0x7c, 0xfd, 0x3e, 0x80, 0x80, 0x8a,
					0x1f, 0xa3, 0x20, 0x39, 0x04},
				hcutil.PKFHybrid, chaincfg.MainNetParams.PubKeyHashAddrID),
			f: func() (hcutil.Address, error) {
				serializedPubKey := []byte{
					0x06, 0x64, 0xc4, 0x46, 0x53, 0xd6, 0x56, 0x7e, 0xff, 0x57,
					0x53, 0xc5, 0xd2, 0x4a, 0x68, 0x2d, 0xdc, 0x2b, 0x2c, 0xad,
					0xfe, 0x1b, 0x0c, 0x64, 0x33, 0xb1, 0x63, 0x74, 0xda, 0xce,
					0x67, 0x78, 0xf0, 0xb8, 0x7c, 0xa4, 0x27, 0x9b, 0x56, 0x5d,
					0x21, 0x30, 0xce, 0x59, 0xf7, 0x5b, 0xfb, 0xb2, 0xb8, 0x8d,
					0xa7, 0x94, 0x14, 0x3d, 0x7c, 0xfd, 0x3e, 0x80, 0x80, 0x8a,
					0x1f, 0xa3, 0x20, 0x39, 0x04}
				return hcutil.NewAddressSecpPubKey(serializedPubKey, &chaincfg.MainNetParams)
			},
			net: &chaincfg.MainNetParams,
		},
		{
			name:    "mainnet p2pk hybrid (0x07)",
			addr:    "DkRKh2aTdwjKKL1mkCb2DFp2Hr7SqMyx3zWqNwyc37PYiGpKmGRsi",
			encoded: "DskEQZMCs4nifL7wx7iHYGWxMQvR9ThCBKQ",
			valid:   true,
			saddr:   "03348d8aeb4253ca52456fe5da94ab1263bfee16bb8192497f666389ca964f8479",
			result: hcutil.TstAddressPubKey(
				[]byte{
					0x07, 0x34, 0x8d, 0x8a, 0xeb, 0x42, 0x53, 0xca, 0x52, 0x45,
					0x6f, 0xe5, 0xda, 0x94, 0xab, 0x12, 0x63, 0xbf, 0xee, 0x16,
					0xbb, 0x81, 0x92, 0x49, 0x7f, 0x66, 0x63, 0x89, 0xca, 0x96,
					0x4f, 0x84, 0x79, 0x83, 0x75, 0x12, 0x9d, 0x79, 0x58, 0x84,
					0x3b, 0x14, 0x25, 0x8b, 0x90, 0x5d, 0xc9, 0x4f, 0xae, 0xd3,
					0x24, 0xdd, 0x8a, 0x9d, 0x67, 0xff, 0xac, 0x8c, 0xc0, 0xa8,
					0x5b, 0xe8, 0x4b, 0xac, 0x5d},
				hcutil.PKFHybrid, chaincfg.MainNetParams.PubKeyHashAddrID),
			f: func() (hcutil.Address, error) {
				serializedPubKey := []byte{
					0x07, 0x34, 0x8d, 0x8a, 0xeb, 0x42, 0x53, 0xca, 0x52, 0x45,
					0x6f, 0xe5, 0xda, 0x94, 0xab, 0x12, 0x63, 0xbf, 0xee, 0x16,
					0xbb, 0x81, 0x92, 0x49, 0x7f, 0x66, 0x63, 0x89, 0xca, 0x96,
					0x4f, 0x84, 0x79, 0x83, 0x75, 0x12, 0x9d, 0x79, 0x58, 0x84,
					0x3b, 0x14, 0x25, 0x8b, 0x90, 0x5d, 0xc9, 0x4f, 0xae, 0xd3,
					0x24, 0xdd, 0x8a, 0x9d, 0x67, 0xff, 0xac, 0x8c, 0xc0, 0xa8,
					0x5b, 0xe8, 0x4b, 0xac, 0x5d}
				return hcutil.NewAddressSecpPubKey(serializedPubKey, &chaincfg.MainNetParams)
			},
			net: &chaincfg.MainNetParams,
		},
		{
			name:    "testnet p2pk compressed (0x02)",
			addr:    "Tso9sQD3ALqRsmEkAm7KvPrkGbeG2Vun7Kv",
			encoded: "Tso9sQD3ALqRsmEkAm7KvPrkGbeG2Vun7Kv",
			valid:   true,
			result: hcutil.TstAddressPubKey(
				[]byte{
					0x02, 0x6a, 0x40, 0xc4, 0x03, 0xe7, 0x46, 0x70, 0xc4, 0xde,
					0x76, 0x56, 0xa0, 0x9c, 0xaa, 0x23, 0x53, 0xd4, 0xb3, 0x83,
					0xa9, 0xce, 0x66, 0xee, 0xf5, 0x1e, 0x12, 0x20, 0xea, 0xcf,
					0x4b, 0xe0, 0x6e},
				hcutil.PKFCompressed, chaincfg.TestNet2Params.PubKeyHashAddrID),
			f: func() (hcutil.Address, error) {
				serializedPubKey := []byte{
					0x02, 0x6a, 0x40, 0xc4, 0x03, 0xe7, 0x46, 0x70, 0xc4, 0xde,
					0x76, 0x56, 0xa0, 0x9c, 0xaa, 0x23, 0x53, 0xd4, 0xb3, 0x83,
					0xa9, 0xce, 0x66, 0xee, 0xf5, 0x1e, 0x12, 0x20, 0xea, 0xcf,
					0x4b, 0xe0, 0x6e}
				return hcutil.NewAddressSecpPubKey(serializedPubKey, &chaincfg.TestNet2Params)
			},
			net: &chaincfg.TestNet2Params,
		},
		{
			name:    "testnet p2pk compressed (0x03)",
			addr:    "TsWZ1EzypJfMwBKAEDYKuyHRGctqGAxMje2",
			encoded: "TsWZ1EzypJfMwBKAEDYKuyHRGctqGAxMje2",
			valid:   true,
			result: hcutil.TstAddressPubKey(
				[]byte{
					0x03, 0x08, 0x44, 0xee, 0x70, 0xd8, 0x38, 0x4d, 0x52, 0x50,
					0xe9, 0xbb, 0x3a, 0x6a, 0x73, 0xd4, 0xb5, 0xbe, 0xc7, 0x70,
					0xe8, 0xb3, 0x1d, 0x6a, 0x0a, 0xe9, 0xfb, 0x73, 0x90, 0x09,
					0xd9, 0x1a, 0xf5},
				hcutil.PKFCompressed, chaincfg.TestNet2Params.PubKeyHashAddrID),
			f: func() (hcutil.Address, error) {
				serializedPubKey := []byte{
					0x03, 0x08, 0x44, 0xee, 0x70, 0xd8, 0x38, 0x4d, 0x52, 0x50,
					0xe9, 0xbb, 0x3a, 0x6a, 0x73, 0xd4, 0xb5, 0xbe, 0xc7, 0x70,
					0xe8, 0xb3, 0x1d, 0x6a, 0x0a, 0xe9, 0xfb, 0x73, 0x90, 0x09,
					0xd9, 0x1a, 0xf5}
				return hcutil.NewAddressSecpPubKey(serializedPubKey, &chaincfg.TestNet2Params)
			},
			net: &chaincfg.TestNet2Params,
		},
		{
			name:    "testnet p2pk uncompressed (0x04)",
			addr:    "TkKmMiY5iDh4U3KkSopYgkU1AzhAcQZiSoVhYhFymZHGMi9LM9Fdt",
			encoded: "Tso9sQD3ALqRsmEkAm7KvPrkGbeG2Vun7Kv",
			valid:   true,
			saddr:   "026a40c403e74670c4de7656a09caa2353d4b383a9ce66eef51e1220eacf4be06e",
			result: hcutil.TstAddressPubKey(
				[]byte{
					0x04, 0x6a, 0x40, 0xc4, 0x03, 0xe7, 0x46, 0x70, 0xc4, 0xde,
					0x76, 0x56, 0xa0, 0x9c, 0xaa, 0x23, 0x53, 0xd4, 0xb3, 0x83,
					0xa9, 0xce, 0x66, 0xee, 0xf5, 0x1e, 0x12, 0x20, 0xea, 0xcf,
					0x4b, 0xe0, 0x6e, 0xd5, 0x48, 0xc8, 0xc1, 0x6f, 0xb5, 0xeb,
					0x90, 0x07, 0xcb, 0x94, 0x22, 0x0b, 0x3b, 0xb8, 0x94, 0x91,
					0xd5, 0xa1, 0xfd, 0x2d, 0x77, 0x86, 0x7f, 0xca, 0x64, 0x21,
					0x7a, 0xce, 0xcf, 0x22, 0x44},
				hcutil.PKFUncompressed, chaincfg.TestNet2Params.PubKeyHashAddrID),
			f: func() (hcutil.Address, error) {
				serializedPubKey := []byte{
					0x04, 0x6a, 0x40, 0xc4, 0x03, 0xe7, 0x46, 0x70, 0xc4, 0xde,
					0x76, 0x56, 0xa0, 0x9c, 0xaa, 0x23, 0x53, 0xd4, 0xb3, 0x83,
					0xa9, 0xce, 0x66, 0xee, 0xf5, 0x1e, 0x12, 0x20, 0xea, 0xcf,
					0x4b, 0xe0, 0x6e, 0xd5, 0x48, 0xc8, 0xc1, 0x6f, 0xb5, 0xeb,
					0x90, 0x07, 0xcb, 0x94, 0x22, 0x0b, 0x3b, 0xb8, 0x94, 0x91,
					0xd5, 0xa1, 0xfd, 0x2d, 0x77, 0x86, 0x7f, 0xca, 0x64, 0x21,
					0x7a, 0xce, 0xcf, 0x22, 0x44}
				return hcutil.NewAddressSecpPubKey(serializedPubKey, &chaincfg.TestNet2Params)
			},
			net: &chaincfg.TestNet2Params,
		},
		{
			name:    "testnet p2pk hybrid (0x06)",
			addr:    "TkKmMiY5iDh4U3KkSopYgkU1AzhAcQZiSoVhYhFymZHGMi9LM9Fdt",
			encoded: "Tso9sQD3ALqRsmEkAm7KvPrkGbeG2Vun7Kv",
			valid:   true,
			saddr:   "026a40c403e74670c4de7656a09caa2353d4b383a9ce66eef51e1220eacf4be06e",
			result: hcutil.TstAddressPubKey(
				[]byte{
					0x06, 0x6a, 0x40, 0xc4, 0x03, 0xe7, 0x46, 0x70, 0xc4, 0xde,
					0x76, 0x56, 0xa0, 0x9c, 0xaa, 0x23, 0x53, 0xd4, 0xb3, 0x83,
					0xa9, 0xce, 0x66, 0xee, 0xf5, 0x1e, 0x12, 0x20, 0xea, 0xcf,
					0x4b, 0xe0, 0x6e, 0xd5, 0x48, 0xc8, 0xc1, 0x6f, 0xb5, 0xeb,
					0x90, 0x07, 0xcb, 0x94, 0x22, 0x0b, 0x3b, 0xb8, 0x94, 0x91,
					0xd5, 0xa1, 0xfd, 0x2d, 0x77, 0x86, 0x7f, 0xca, 0x64, 0x21,
					0x7a, 0xce, 0xcf, 0x22, 0x44},
				hcutil.PKFHybrid, chaincfg.TestNet2Params.PubKeyHashAddrID),
			f: func() (hcutil.Address, error) {
				serializedPubKey := []byte{
					0x06, 0x6a, 0x40, 0xc4, 0x03, 0xe7, 0x46, 0x70, 0xc4, 0xde,
					0x76, 0x56, 0xa0, 0x9c, 0xaa, 0x23, 0x53, 0xd4, 0xb3, 0x83,
					0xa9, 0xce, 0x66, 0xee, 0xf5, 0x1e, 0x12, 0x20, 0xea, 0xcf,
					0x4b, 0xe0, 0x6e, 0xd5, 0x48, 0xc8, 0xc1, 0x6f, 0xb5, 0xeb,
					0x90, 0x07, 0xcb, 0x94, 0x22, 0x0b, 0x3b, 0xb8, 0x94, 0x91,
					0xd5, 0xa1, 0xfd, 0x2d, 0x77, 0x86, 0x7f, 0xca, 0x64, 0x21,
					0x7a, 0xce, 0xcf, 0x22, 0x44}
				return hcutil.NewAddressSecpPubKey(serializedPubKey, &chaincfg.TestNet2Params)
			},
			net: &chaincfg.TestNet2Params,
		},
		{
			name:    "testnet p2pk hybrid (0x07)",
			addr:    "TkQ5Ax2ieEZpBDA963VDH4y27KMpXtP8qyeykzwBNFocDc8ZKqTGz",
			encoded: "TsTFLdM32YVrYsEQDFxo2zmPuKFcFhH5ZT3",
			valid:   true,
			saddr:   "03edd40747de905a9becb14987a1a26c1adbd617c45e1583c142a635bfda9493df",
			result: hcutil.TstAddressPubKey(
				[]byte{
					0x07, 0xed, 0xd4, 0x07, 0x47, 0xde, 0x90, 0x5a, 0x9b, 0xec,
					0xb1, 0x49, 0x87, 0xa1, 0xa2, 0x6c, 0x1a, 0xdb, 0xd6, 0x17,
					0xc4, 0x5e, 0x15, 0x83, 0xc1, 0x42, 0xa6, 0x35, 0xbf, 0xda,
					0x94, 0x93, 0xdf, 0xa1, 0xc6, 0xd3, 0x67, 0x35, 0x97, 0x49,
					0x65, 0xfe, 0x7b, 0x86, 0x1e, 0x7f, 0x6f, 0xcc, 0x08, 0x7d,
					0xc7, 0xfe, 0x47, 0x38, 0x0f, 0xa8, 0xbd, 0xe0, 0xd9, 0xc3,
					0x22, 0xd5, 0x3c, 0x0e, 0x89},
				hcutil.PKFHybrid, chaincfg.TestNet2Params.PubKeyHashAddrID),
			f: func() (hcutil.Address, error) {
				serializedPubKey := []byte{
					0x07, 0xed, 0xd4, 0x07, 0x47, 0xde, 0x90, 0x5a, 0x9b, 0xec,
					0xb1, 0x49, 0x87, 0xa1, 0xa2, 0x6c, 0x1a, 0xdb, 0xd6, 0x17,
					0xc4, 0x5e, 0x15, 0x83, 0xc1, 0x42, 0xa6, 0x35, 0xbf, 0xda,
					0x94, 0x93, 0xdf, 0xa1, 0xc6, 0xd3, 0x67, 0x35, 0x97, 0x49,
					0x65, 0xfe, 0x7b, 0x86, 0x1e, 0x7f, 0x6f, 0xcc, 0x08, 0x7d,
					0xc7, 0xfe, 0x47, 0x38, 0x0f, 0xa8, 0xbd, 0xe0, 0xd9, 0xc3,
					0x22, 0xd5, 0x3c, 0x0e, 0x89}
				return hcutil.NewAddressSecpPubKey(serializedPubKey, &chaincfg.TestNet2Params)
			},
			net: &chaincfg.TestNet2Params,
		},
	}

	for _, test := range tests {
		// Decode addr and compare error against valid.
		decoded, err := hcutil.DecodeAddress(test.addr)
		if (err == nil) != test.valid {
			t.Errorf("%v: decoding test failed: %v", test.name, err)
			return
		}

		if err == nil {
			// Ensure the stringer returns the same address as the
			// original.
			if decodedStringer, ok := decoded.(fmt.Stringer); ok {
				if test.addr != decodedStringer.String() {
					t.Errorf("%v: String on decoded value does not match expected value: %v != %v",
						test.name, test.addr, decodedStringer.String())
					return
				}
			}

			// Encode again and compare against the original.
			encoded := decoded.EncodeAddress()
			if test.encoded != encoded {
				t.Errorf("%v: decoding and encoding produced different addressess: %v != %v",
					test.name, test.encoded, encoded)
				return
			}

			// Perform type-specific calculations.
			var saddr []byte
			switch d := decoded.(type) {
			case *hcutil.AddressPubKeyHash:
				saddr = hcutil.TstAddressSAddr(encoded)

			case *hcutil.AddressScriptHash:
				saddr = hcutil.TstAddressSAddr(encoded)

			case *hcutil.AddressSecpPubKey:
				// Ignore the error here since the script
				// address is checked below.
				saddr, err = hex.DecodeString(d.String())
				if err != nil {
					saddr, _ = hex.DecodeString(test.saddr)
				}

			case *hcutil.AddressEdwardsPubKey:
				// Ignore the error here since the script
				// address is checked below.
				saddr, _ = hex.DecodeString(d.String())

			case *hcutil.AddressSecSchnorrPubKey:
				// Ignore the error here since the script
				// address is checked below.
				saddr, _ = hex.DecodeString(d.String())
			}

			// Check script address, as well as the Hash160 method for P2PKH and
			// P2SH addresses.
			if !bytes.Equal(saddr, decoded.ScriptAddress()) {
				t.Errorf("%v: script addresses do not match:\n%x != \n%x",
					test.name, saddr, decoded.ScriptAddress())
				return
			}
			switch a := decoded.(type) {
			case *hcutil.AddressPubKeyHash:
				if h := a.Hash160()[:]; !bytes.Equal(saddr, h) {
					t.Errorf("%v: hashes do not match:\n%x != \n%x",
						test.name, saddr, h)
					return
				}

			case *hcutil.AddressScriptHash:
				if h := a.Hash160()[:]; !bytes.Equal(saddr, h) {
					t.Errorf("%v: hashes do not match:\n%x != \n%x",
						test.name, saddr, h)
					return
				}
			}

			// Ensure the address is for the expected network.
			if !decoded.IsForNet(test.net) {
				t.Errorf("%v: calculated network does not match expected",
					test.name)
				return
			}
		}

		if !test.valid {
			// If address is invalid, but a creation function exists,
			// verify that it returns a nil addr and non-nil error.
			if test.f != nil {
				_, err := test.f()
				if err == nil {
					t.Errorf("%v: address is invalid but creating new address succeeded",
						test.name)
					return
				}
			}
			continue
		}

		// Valid test, compare address created with f against expected result.
		addr, err := test.f()
		if err != nil {
			t.Errorf("%v: address is valid but creating new address failed with error %v",
				test.name, err)
			return
		}
		if !reflect.DeepEqual(addr.ScriptAddress(), test.result.ScriptAddress()) {
			t.Errorf("%v: created address does not match expected result \n "+
				"	got %x, expected %x",
				test.name, addr.ScriptAddress(), test.result.ScriptAddress())
			return
		}
	}
}
