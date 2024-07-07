package worker

import (
	"errors"
	"fmt"
	"log"

	gearman "github.com/mikespook/gearman-go/worker"
	"github.com/sentiweb/gearworkers/pkg/config"
)

type Manager struct {
	config  *config.AppConfig
	workers []*Worker
}

type Worker struct {
	worker   *gearman.Worker
	config   config.JobConfig
	executor Executor
}

type Executor interface {
	Run(gearman.Job) ([]byte, error)
	Init() error
}

func NewManager(cfg *config.AppConfig) *Manager {
	return &Manager{
		config:  cfg,
		workers: make([]*Worker, 0, len(cfg.Jobs)),
	}
}

func (m *Manager) Start() error {
	for idx, job := range m.config.Jobs {
		var executor Executor = nil
		if job.Type == "shell" {
			shellConfig := job.ShellConfig
			if shellConfig == nil {
				return fmt.Errorf("job %d %s : shell config is not provided", idx, job.Name)
			}
			executor = NewShellExecutor(job.Name, *shellConfig)
		}
		if job.Type == "http" {
			httpConfig := job.HttpConfig
			if httpConfig == nil {
				return fmt.Errorf("job %d %s : http config is not provided", idx, job.Name)
			}
			executor = NewHttpExecutor(job.Name, *httpConfig)
		}
		if executor == nil {
			return fmt.Errorf("job %d %s : Unknown type '%s'", idx, job.Name, job.Type)
		}
		worker := NewWorker(job, executor)
		err := worker.Register(m.config.GearmanServer)
		if err != nil {
			return errors.Join(fmt.Errorf("error registering job %d '%s'", idx, job.Name), err)
		}
		err = worker.Start()
		if err != nil {
			return errors.Join(fmt.Errorf("error starting job %d '%s'", idx, job.Name), err)
		}
		log.Printf("Job %s started %s", job.Name, worker.WorkerId())
		m.workers = append(m.workers, worker)
	}

	return nil
}

func NewWorker(job config.JobConfig, executor Executor) *Worker {
	return &Worker{
		worker:   gearman.New(job.Concurrency),
		config:   job,
		executor: executor,
	}
}

func (w *Worker) Register(server string) error {
	err := w.worker.AddServer("tcp", server)
	if err != nil {
		return err
	}
	w.worker.ErrorHandler = func(e error) {
		log.Printf("Job %s Worker %s : %s", w.config.Name, w.worker.Id, e)
	}
	w.worker.JobHandler = func(job gearman.Job) error {
		log.Printf("Job received %s %s", w.config.Name, job.UniqueId())
		return nil
	}

	timeout, err := timeoutToSeconds(w.config.Timeout)
	if err != nil {
		return err
	}
	err = w.worker.AddFunc(w.config.Name, w.executor.Run, timeout)
	if err != nil {
		return err
	}
	return nil
}

func (w *Worker) WorkerId() string {
	return w.worker.Id
}

func (w *Worker) Start() error {
	err := w.worker.Ready()
	if err != nil {
		return err
	}
	go w.worker.Work()
	return nil
}
