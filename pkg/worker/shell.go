package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	gearman "github.com/mikespook/gearman-go/worker"
	"github.com/sentiweb/gearworkers/pkg/config"
	"github.com/sentiweb/gearworkers/pkg/types"
)

type ShellExecutor struct {
	config  config.ShellJobConfig
	metrics *ExecutorMetrics
	timeout time.Duration
}

func NewShellExecutor(name string, cfg config.ShellJobConfig) *ShellExecutor {
	return &ShellExecutor{
		config:  cfg,
		metrics: NewExecutorMetrics("shell", name),
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
	var err error
	jobId := job.UniqueId()

	payload := types.ShellJobPayload{}

	data := job.Data()
	if len(data) > 0 {
		err = json.Unmarshal(data, &payload)
		if err != nil {
			log.Printf("[Job %s] Unable to parse job data : %s", jobId, err)
			return nil, err
		}
	}

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

	if len(payload.EnvParams) > 0 {
		for n, v := range payload.EnvParams {
			e := fmt.Sprintf("%s=%s", n, v)
			env = append(env, e)
		}
	}

	if executor.config.WorkingDir != "" {
		cmd.Dir = executor.config.WorkingDir
	}

	cmd.Env = env
	var output ShellOutput

	logfile := executor.config.LogFile
	if payload.LogFile != "" {
		// Logfile can be overriden
		logfile = payload.LogFile
	}

	if logfile != "" {
		output = NewShellOuputLogger(logfile)
	} else {
		output = NewShellOuputBuffer()
	}

	err = output.Configure(cmd)
	if err != nil {
		log.Printf("Job %s error configuring output : %v", jobId, err)
	}

	executor.metrics.IncTotalCounter()

	err = cmd.Run()
	quit := cmd.ProcessState.ExitCode()

	log.Printf("Job %s Command ended with code %d", jobId, quit)

	return output.Out(), err
}
