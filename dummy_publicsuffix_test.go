// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cookiejarx_test

import "github.com/eientei/cookiejarx"

type dummypsl struct {
	List cookiejarx.PublicSuffixList
}

func (dummypsl) PublicSuffix(domain string) string {
	return domain
}

func (dummypsl) String() string {
	return "dummy"
}

var publicsuffix = dummypsl{}
