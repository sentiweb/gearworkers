# Gearworkers

Simple gearman worker to handle simple case tasks like:

- Run a command in a external program (using shellexecutor)
- Run an http request as a task (task payload )

# Configuration

The worker configuration is expected in a yaml file:

```yaml
# The gearman server
gearman: "127.0.0.1:4730"
# The job configuration
jobs:
 - 
  name: "my_job_http"
  # Type define the executor to use
  # if "http", an "http_config" entry is expected
  # if "shell", a "shell_config" entry is expected
  type: "http"
  # Optional timeout, expressed as golang duration string, eg. "10s", "24h"
  timeout: "10s"
  http_config:
    url: "http://127.0.0.1:3202/test?action=do"
    # Http method to use, 
    method: 'GET'
    # Headers to provide to all requests
    headers:
      "X-Toto": "MyAwesomeHeader"
 - 
  name: "my_job_shell"
  type: "shell"
  # Optional timeout, expressed as golang duration string, eg. "10s", "24h"
  timeout: "1h"
  shell_config:
    # Command to run (without args)
    command: "/usr/bin/mycommand"
    args:
        - "-c"
        - "myconfig.yaml"
    # Working directory when the command is run
    working_dir: "/path/to/work"
    env:
      MY_ENV_VARIABLE: "value"
    log_file: "/path/to/logfile"
```

# Http Executor

The executor accepts a payload as a json object

```json
{
    "body":"",
    "query":{
        "action":"open"
     },
    "headers": {}
}
```

`query` object can be used to add entry in the query parameters of the URL
`body` can be used to provide a body (only for POST/PUT method)
`headers` can be used to add header to the query