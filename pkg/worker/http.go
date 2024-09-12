package worker

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	gearman "github.com/mikespook/gearman-go/worker"
	"github.com/sentiweb/gearworkers/pkg/config"
	"github.com/sentiweb/gearworkers/pkg/types"
)

type HttpExecutor struct {
	config         config.HttpJobConfig
	metrics        *ExecutorMetrics
	timeout        time.Duration
	allowedQuery   Acceptable
	allowedHeaders Acceptable
}

func NewHttpExecutor(name string, cfg config.HttpJobConfig) *HttpExecutor {
	return &HttpExecutor{
		config:  cfg,
		metrics: NewExecutorMetrics("http", name),
	}
}

var ErrUnableToParseDuration = errors.New("unable to parse duration")
var ErrUnableToParseUrl = errors.New("unable to parse url")

func (executor *HttpExecutor) Init() error {
	if executor.config.Timeout != "" {
		timeout, err := ParseDuration(executor.config.Timeout, "s")
		if err != nil {
			return errors.Join(ErrUnableToParseDuration, err)
		}
		executor.timeout = timeout
	} else {
		executor.timeout = 10 * time.Second
	}

	_, err := url.Parse(executor.config.Url)
	if err != nil {
		return errors.Join(ErrUnableToParseUrl, err)
	}

	var acceptQuery Acceptable
	if len(executor.config.AllowedQueryParams) > 0 {
		acceptQuery = NewAcceptableValues(executor.config.AllowedQueryParams)
	} else {
		acceptQuery = NewAcceptableAll()
	}
	executor.allowedQuery = acceptQuery

	var allowedHeaders Acceptable
	if len(executor.config.AllowedHeaders) > 0 {
		allowedHeaders = NewAcceptableValues(executor.config.AllowedHeaders)
	} else {
		allowedHeaders = NewAcceptableAll()
	}
	executor.allowedHeaders = allowedHeaders

	return nil
}

func extractBody(body interface{}) (io.Reader, error) {
	s, ok := body.(string)
	if ok {
		return strings.NewReader(s), nil
	}
	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}

var UrlVariableRegex = regexp.MustCompile("{(\\w+)}")

func BuildUrl(template string, vars map[string]string) (string, error) {
	matches := UrlVariableRegex.FindAllStringSubmatch(template, -1)

	if len(matches) == 0 {
		return template, nil
	}

	url := template
	for _, match := range matches {
		varName := match[1]
		varExpr := match[0]

		value, ok := vars[varName]
		if !ok {
			return "", fmt.Errorf("missing expected url param '%s' ", varName)
		}
		url = strings.ReplaceAll(url, varExpr, value)
	}
	return url, nil
}

func (executor *HttpExecutor) Run(job gearman.Job) ([]byte, error) {

	jobId := job.UniqueId()

	payload := types.HttpJobPayload{}
	err := json.Unmarshal(job.Data(), &payload)
	if err != nil {
		log.Printf("[Job %s] Unable to parse job data : %s", jobId, err)
		return nil, err
	}

	svcUrl, err := BuildUrl(executor.config.Url, payload.UrlParams)
	if err != nil {
		log.Printf("[Job %s] Unable to parse url params : %s", jobId, err)
		return nil, err
	}

	// Url parsing is checked during init phase
	u, _ := url.Parse(svcUrl)

	if len(payload.QueryParams) > 0 {
		queryValues := u.Query()
		for n, v := range payload.QueryParams {
			if executor.allowedQuery.Accept(n) {
				queryValues.Set(n, v)
			} else {
				log.Printf("[Job %s]  : query param not allowed %s", jobId, n)
			}
		}
		u.RawQuery = queryValues.Encode()
	}

	method := executor.config.Method

	var body io.Reader

	if method == "POST" || method == "PUT" {
		body = nil
		if payload.Body != nil {
			body, err = extractBody(payload.Body)
			if err != nil {
				log.Printf("[Job %s] Unable to parse job body : %s", jobId, err)
			}
		}
	} else {
		body = nil
	}

	/*
		rc, cancel := context.WithTimeout(context.Background(), executor.timeout)
		defer cancel()

		req, err := http.NewRequestWithContext(rc, method, u.String(), body)
	*/
	log.Printf("[Job %s] fetching %s %s", jobId, method, u.String())

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		log.Printf("[Job %s] Error preparing request %s", jobId, err)
		return nil, err
	}

	populateHeaders(req, executor.config.Headers, NewAcceptableAll())
	populateHeaders(req, payload.Headers, executor.allowedHeaders)

	client := &http.Client{Timeout: executor.timeout}

	executor.metrics.IncTotalCounter()

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error during request %s", err)
	} else {
		defer resp.Body.Close()
		log.Printf("[Job %s] Response %d %s", jobId, resp.StatusCode, resp.Status)
	}
	return nil, nil
}

func populateHeaders(req *http.Request, headers map[string]string, accepter Acceptable) {
	if len(headers) > 0 {
		for n, v := range headers {
			if accepter.Accept(n) {
				req.Header.Set(n, v)
			}
		}
	}
}
