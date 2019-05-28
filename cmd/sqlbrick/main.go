// Copyright (c) 2018-present Anbillon Team (anbillonteam@gmail.com).
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package main

import (
	"log"
)

func main() {
	cmd := newRootCmd()
	if err := cmd.Execute(); err != nil {
		log.Printf("%v", err)
	}
}
