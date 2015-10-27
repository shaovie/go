package main

import (
	"errors"
	"fmt"
	"os"
	"path"
	"runtime"
	"strconv"
	"syscall"
)

var pidFile *os.File

func outputPid(pidPath string) error {
	if pidPath == "" {
		return errors.New("pid path is empty")
	}

	err := os.MkdirAll(path.Dir(pidPath), 0755)
	if err != nil {
		return err
	}
	pidFile, err = os.OpenFile(pidPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	// Luanch only one instance.
	if err := syscall.Flock(int(pidFile.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		pidFile.Close()
		return err
	}
	pidFile.Truncate(0)
	pidFile.Write([]byte(strconv.Itoa(os.Getpid())))
	return nil
}

func daemon() error {
	isDarwin := runtime.GOOS == "darwin"

	// already a daemon
	if syscall.Getppid() == 1 {
		return nil
	}

	// fork off the parent process
	ret, ret2, errno := syscall.RawSyscall(syscall.SYS_FORK, 0, 0, 0)
	if errno != 0 {
		return errors.New(fmt.Sprintf("fork error! [errno=%d]", errno))
	}

	// failure
	if ret2 < 0 {
		os.Exit(1)
	}

	// handle exception for darwin
	if isDarwin && ret2 == 1 {
		ret = 0
	}

	// if we got a good PID, then we call exit the parent process.
	if int(ret) > 0 {
		os.Exit(0)
	}

	syscall.Umask(0)

	// create a new SID for the child process
	_, err := syscall.Setsid()
	if err != nil {
		return err
	}

	//os.Chdir("/")

	f, err := os.OpenFile("/dev/null", os.O_RDWR, 0)
	if err == nil {
		fd := f.Fd()
		syscall.Dup2(int(fd), int(os.Stdin.Fd()))
		syscall.Dup2(int(fd), int(os.Stdout.Fd()))
		syscall.Dup2(int(fd), int(os.Stderr.Fd()))
	}
	return nil
}
