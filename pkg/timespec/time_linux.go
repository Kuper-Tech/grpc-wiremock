package timespec

import "syscall"

func CreationTime(info *syscall.Stat_t) syscall.Timespec {
	return info.Ctim
}
