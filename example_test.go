// SPDX-FileCopyrightText: © 2022 Jade Meskill <iamruinous@ruinous.social>
//
// SPDX-License-Identifier: MIT

package geminitest_test

import (
	geminitest "codeberg.org/iamruinous/geminitest"
	"context"
	"fmt"
	gemini "git.sr.ht/~adnano/go-gemini"
	"io"
)

func ExampleResponseRecorder() {
	handler := func(ctx context.Context, w gemini.ResponseWriter, r *gemini.Request) {
		io.WriteString(w, "# Hello World!")
	}

	req := geminitest.NewRequest("gemini://example.com/foo")
	w := geminitest.NewRecorder()
	ctx := context.Background()
	handler(ctx, w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	fmt.Println(int(resp.Status))
	fmt.Println(string(body))

	// Output:
	// 20
	// # Hello World!
}
