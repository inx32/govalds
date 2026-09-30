//go:build unix

package pidfile

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"

	"golang.org/x/sys/unix"
)

type Pid struct {
	filename string
	file     *os.File
}

func (p *Pid) Lock() error {
	if p.file != nil {
		return errors.New("already locked by this instance")
	}

	file, err := os.OpenFile(p.filename, os.O_RDWR|os.O_CREATE, 0o660)
	if err != nil {
		return fmt.Errorf("open %s: %w", p.filename, err)
	}
	p.file = file

	err = unix.Flock(int(file.Fd()), unix.LOCK_EX|unix.LOCK_NB)
	if err != nil {
		file.Close()
		p.file = nil
		return fmt.Errorf("another instance is already running (pidfile locked): %w", err)
	}

	if err := file.Truncate(0); err != nil {
		p.Release()
		return err
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		p.Release()
		return err
	}

	pid := strconv.Itoa(os.Getpid()) + "\n"
	if _, err := file.WriteString(pid); err != nil {
		p.Release()
		return fmt.Errorf("write pid: %w", err)
	}

	if err := file.Sync(); err != nil {
		p.Release()
		return fmt.Errorf("sync pid file: %w", err)
	}

	return nil
}

func (p *Pid) Release() error {
	if p.file == nil {
		return errors.New("file is nil")
	}

	if err := unix.Flock(int(p.file.Fd()), unix.LOCK_UN); err != nil {
		return fmt.Errorf("release lock: %w", err)
	}

	file := p.file
	p.file = nil

	if err := file.Close(); err != nil {
		return fmt.Errorf("close file: %w", err)
	}

	if err := os.Remove(p.filename); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("delete file: %w", err)
	}

	return nil
}

func New(filename string) *Pid {
	return &Pid{filename: filename}
}
