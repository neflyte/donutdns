//go:build !linux

package donutdns

import (
	"github.com/shoenig/go-landlock"
)

var sysPaths []*landlock.Path
