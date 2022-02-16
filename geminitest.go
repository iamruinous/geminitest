// SPDX-FileCopyrightText: © 2022 Jade Meskill <iamruinous@ruinous.social>
//
// SPDX-License-Identifier: MIT

package geminitest

import (
	"bufio"
	gemini "git.sr.ht/~adnano/go-gemini"
	"strings"
)

// NewRequest returns a new incoming server Request, suitable
// for passing to an gemini.Handler for testing.
//
// NewRequest panics on error for ease of use in testing, where a
// panic is acceptable.
func NewRequest(rawurl string) *gemini.Request {
	req, err := gemini.ReadRequest(bufio.NewReader(strings.NewReader(rawurl + "\r\n")))
	if err != nil {
		panic("invalid NewRequest arguments; " + err.Error())
	}

	return req
}
