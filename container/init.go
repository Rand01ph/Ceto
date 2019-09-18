package container

import (
	"github.com/sirupsen/logrus"
	"os"
	"syscall"
)

func RunContainerInitProcess(command string, args []string) error {

	// systemd 加入linux之后, mount namespace 就变成 shared by default, 所以你必须显示
	//声明你要这个新的mount namespace独立。
	//在 syscall 当中就需要手动的以 private 的方式 mount 一遍根目录以达到效果（要在 chroot 之前）
	syscall.Mount("", "/", "", uintptr(syscall.MS_PRIVATE|syscall.MS_REC), "")

	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	argv := []string{command}
	if err := syscall.Exec(command, argv, os.Environ()); err != nil {
		logrus.Errorf(err.Error())
	}
	return nil
}
