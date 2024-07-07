package worker

import (
	"fmt"

	"github.com/VictoriaMetrics/metrics"
)

type ExecutorMetrics struct {
	counterTotal *metrics.Counter
}

func NewExecutorMetrics(executorType string, name string) *ExecutorMetrics {

	totalName := fmt.Sprintf(`gearworkers_executor_total{type="%s", name="%s"}`, executorType, name)

	requestsTotal := metrics.NewCounter(totalName)

	return &ExecutorMetrics{
		counterTotal: requestsTotal,
	}
}

func (m *ExecutorMetrics) IncTotalCounter() {
	m.counterTotal.Inc()
}
