package worker

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	gearman "github.com/mikespook/gearman-go/worker"
	"github.com/sentiweb/gearworkers/pkg/config"
)

type ShellExecutor struct {
	config  config.ShellJobConfig
	timeout time.Duration
}

func NewShellExecutor(cfg config.ShellJobConfig) *ShellExecutor {
	return &ShellExecutor{
		config: cfg,
	}
}

func (executor *ShellExecutor) Init() error {
	if executor.config.Timeout != "" {
		timeout, err := ParseDuration(executor.config.Timeout, "s")
		if err != nil {
			return err
		}
		executor.timeout = timeout
	}
	return nil
}

func (executor *ShellExecutor) Run(job gearman.Job) ([]byte, error) {
	command := executor.config.Command
	args := executor.config.Args
	var cmd *exec.Cmd

	jobId := job.UniqueId()

	if executor.timeout != 0 {
		ctx, cancel := context.WithTimeout(context.Background(), executor.timeout)
		defer cancel()
		cmd = exec.CommandContext(ctx, command, args...)
	} else {
		cmd = exec.Command(command, args...)
	}

	env := make([]string, 0)

	env = append(env, os.Environ()...)
	if len(executor.config.Env) > 0 {
		for n, v := range executor.config.Env {
			e := fmt.Sprintf("%s=%s", n, v)
			env = append(env, e)
		}
	}

	if executor.config.WorkingDir != "" {
		cmd.Dir = executor.config.WorkingDir
	}

	cmd.Env = env
	var output ShellOutput
	if executor.config.LogFile != "" {
		output = NewShellOuputLogger(executor.config.LogFile)
	} else {
		output = NewShellOuputBuffer()
	}

	err := output.Configure(cmd)
	if err != nil {
		log.Printf("Job %s error configuring output : %v", jobId, err)
	}

	err = cmd.Run()
	quit := cmd.ProcessState.ExitCode()

	log.Printf("Job %s Command ended with code %d", jobId, quit)

	return output.Out(), err
}
