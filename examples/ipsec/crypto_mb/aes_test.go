// Copyright 2019 Intel Corporation.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crypto_mb

import "testing"

//TODO Separate test package?

// Appendix B, C of FIPS 197: Cipher examples, Example vectors.
type CryptTest struct {
	key []byte
	in  []byte
	out []byte
}

var encryptTests = []CryptTest{
	{
		// Appendix B.
		[]byte{0x2b, 0x7e, 0x15, 0x16, 0x28, 0xae, 0xd2, 0xa6, 0xab, 0xf7, 0x15, 0x88, 0x09, 0xcf, 0x4f, 0x3c},
		[]byte{0x32, 0x43, 0xf6, 0xa8, 0x88, 0x5a, 0x30, 0x8d, 0x31, 0x31, 0x98, 0xa2, 0xe0, 0x37, 0x07, 0x34},
		[]byte{0x39, 0x25, 0x84, 0x1d, 0x02, 0xdc, 0x09, 0xfb, 0xdc, 0x11, 0x85, 0x97, 0x19, 0x6a, 0x0b, 0x32},
	},
	{
		// Appendix C.1.  AES-128
		[]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f},
		[]byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff},
		[]byte{0x69, 0xc4, 0xe0, 0xd8, 0x6a, 0x7b, 0x04, 0x30, 0xd8, 0xcd, 0xb7, 0x80, 0x70, 0xb4, 0xc5, 0x5a},
	},
}

// Test Cipher Encrypt method against FIPS 197 examples.
func TestCipherEncrypt(t *testing.T) {
	for i, tt := range encryptTests {
		c := NewAESMultiBlock(tt.key)
		out := make([][]byte, 0)
		in := make([][]byte, 0)
		for i := 0; i < c.VecSize(); i++ {
			in = append(in, tt.in)
			out = append(out, make([]byte, len(tt.in)))
		}
		c.EncryptMany(out, in)
		for s := 0; s < c.VecSize(); s++ {
			for j, v := range out[s] {
				if v != tt.out[j] {
					t.Errorf("EncryptMany slice %d  %d: out[%d] = %#x, want %#x", s, i, j, v, tt.out[j])
					break
				}
			}
		}
	}
}

func BenchmarkEncrypt(b *testing.B) {
	tt := encryptTests[0]
	c := NewAESMultiBlock(tt.key)
	out := make([][]byte, 0)
	in := make([][]byte, 0)
	for i := 0; i < c.VecSize(); i++ {
		in = append(in, tt.in)
		out = append(out, make([]byte, len(tt.in)))
	}
	b.SetBytes(int64(len(out[0]) * c.VecSize()))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.EncryptMany(out, in)
	}
}
