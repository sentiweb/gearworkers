package worker

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
)

type ShellOutput interface {
	Out() []byte
	Configure(cmd *exec.Cmd) error
}

type ShellOuputLogger struct {
	file string
}

func NewShellOuputLogger(file string) *ShellOuputLogger {
	return &ShellOuputLogger{file: file}
}

func (o *ShellOuputLogger) Configure(cmd *exec.Cmd) error {
	outFile, err := os.OpenFile(o.file, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return errors.Join(fmt.Errorf("unable to open '%s'", o.file), err)
	}
	cmd.Stdout = outFile
	cmd.Stderr = outFile
	return nil
}

func (o *ShellOuputLogger) Out() []byte {
	return nil
}

type ShellOuputBuffer struct {
	buffer *bytes.Buffer
}

func NewShellOuputBuffer() *ShellOuputBuffer {
	return &ShellOuputBuffer{}
}

func (o *ShellOuputBuffer) Configure(cmd *exec.Cmd) error {
	var b bytes.Buffer
	cmd.Stdout = &b
	cmd.Stderr = &b
	o.buffer = &b
	return nil
}

func (o *ShellOuputBuffer) Out() []byte {
	return o.buffer.Bytes()
}
