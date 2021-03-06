// Copyright 2017 syzkaller project authors. All rights reserved.
// Use of this source code is governed by Apache 2 LICENSE that can be found in the LICENSE file.

package openbsd

import (
	"github.com/google/syzkaller/prog"
	"github.com/google/syzkaller/sys/targets"
)

func InitTarget(target *prog.Target) {
	arch := &arch{
		unix: targets.MakeUnixSanitizer(target),
	}

	target.MakeMmap = targets.MakePosixMmap(target)
	target.SanitizeCall = arch.SanitizeCall
}

type arch struct {
	unix *targets.UnixSanitizer
}

func (arch *arch) SanitizeCall(c *prog.Call) {
	// Prevent vnodes of type VBAD from being created. Such vnodes will
	// likely trigger assertion errors by the kernel.
	pos := 1
	switch c.Meta.CallName {
	case "mknodat":
		pos = 2
		fallthrough
	case "mknod":
		mode := c.Args[pos].(*prog.ConstArg)
		if mode.Val&arch.unix.S_IFMT == arch.unix.S_IFMT {
			mode.Val &^= arch.unix.S_IFMT
			mode.Val |= arch.unix.S_IFCHR
		}
	default:
		arch.unix.SanitizeCall(c)
	}
}
