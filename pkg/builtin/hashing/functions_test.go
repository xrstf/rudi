// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package hashing

import (
	"testing"

	"go.xrstf.de/rudi/pkg/testutil"
)

func TestSha1Function(t *testing.T) {
	testcases := []testutil.Testcase{
		{
			Expression: `(sha1)`,
			Invalid:    true,
		},
		{
			Expression: `(sha1 "too" "many")`,
			Invalid:    true,
		},
		{
			Expression: `(sha1 true)`,
			Invalid:    true,
		},
		{
			Expression: `(sha1 1)`,
			Invalid:    true,
		},
		{
			// strict coalescing allows null to turn into ""
			Expression: `(sha1 null)`,
			Expected:   "da39a3ee5e6b4b0d3255bfef95601890afd80709",
		},
		{
			Expression: `(sha1 "")`,
			Expected:   "da39a3ee5e6b4b0d3255bfef95601890afd80709",
		},
		{
			Expression: `(sha1 " ")`,
			Expected:   "b858cb282617fb0956d960215c8e84d1ccf909c6",
		},
	}

	for _, testcase := range testcases {
		testcase.Functions = Functions
		t.Run(testcase.String(), testcase.Run)
	}
}

func TestSha256Function(t *testing.T) {
	testcases := []testutil.Testcase{
		{
			Expression: `(sha256)`,
			Invalid:    true,
		},
		{
			Expression: `(sha256 "too" "many")`,
			Invalid:    true,
		},
		{
			Expression: `(sha256 true)`,
			Invalid:    true,
		},
		{
			Expression: `(sha256 1)`,
			Invalid:    true,
		},
		{
			// strict coalescing allows null to turn into ""
			Expression: `(sha256 null)`,
			Expected:   "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
		{
			Expression: `(sha256 "")`,
			Expected:   "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
		{
			Expression: `(sha256 " ")`,
			Expected:   "36a9e7f1c95b82ffb99743e0c5c4ce95d83c9a430aac59f84ef3cbfab6145068",
		},
	}

	for _, testcase := range testcases {
		testcase.Functions = Functions
		t.Run(testcase.String(), testcase.Run)
	}
}

func TestSha512Function(t *testing.T) {
	testcases := []testutil.Testcase{
		{
			Expression: `(sha512)`,
			Invalid:    true,
		},
		{
			Expression: `(sha512 "too" "many")`,
			Invalid:    true,
		},
		{
			Expression: `(sha512 true)`,
			Invalid:    true,
		},
		{
			Expression: `(sha512 1)`,
			Invalid:    true,
		},
		{
			// strict coalescing allows null to turn into ""
			Expression: `(sha512 null)`,
			Expected:   "cf83e1357eefb8bdf1542850d66d8007d620e4050b5715dc83f4a921d36ce9ce47d0d13c5d85f2b0ff8318d2877eec2f63b931bd47417a81a538327af927da3e",
		},
		{
			Expression: `(sha512 "")`,
			Expected:   "cf83e1357eefb8bdf1542850d66d8007d620e4050b5715dc83f4a921d36ce9ce47d0d13c5d85f2b0ff8318d2877eec2f63b931bd47417a81a538327af927da3e",
		},
		{
			Expression: `(sha512 " ")`,
			Expected:   "f90ddd77e400dfe6a3fcf479b00b1ee29e7015c5bb8cd70f5f15b4886cc339275ff553fc8a053f8ddc7324f45168cffaf81f8c3ac93996f6536eef38e5e40768",
		},
	}

	for _, testcase := range testcases {
		testcase.Functions = Functions
		t.Run(testcase.String(), testcase.Run)
	}
}
