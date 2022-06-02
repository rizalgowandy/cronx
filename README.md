[![Go Doc](https://pkg.go.dev/badge/github.com/rizalgowandy/cronx?status.svg)](https://pkg.go.dev/github.com/rizalgowandy/cronx?tab=doc)
[![Release](https://img.shields.io/github/release/rizalgowandy/cronx.svg?style=flat-square)](https://github.com/rizalgowandy/cronx/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/rizalgowandy/cronx)](https://goreportcard.com/report/github.com/rizalgowandy/cronx)
[![Build Status](https://github.com/rizalgowandy/cronx/workflows/Go/badge.svg?branch=main)](https://github.com/rizalgowandy/cronx/actions?query=branch%3Amain)
[![Sourcegraph](https://sourcegraph.com/github.com/rizalgowandy/cronx/-/badge.svg)](https://sourcegraph.com/github.com/rizalgowandy/cronx?badge)

![logo](https://socialify.git.ci/rizalgowandy/cronx/image?description=1&language=1&pattern=Floating%20Cogs&theme=Light)

Cronx is a library to manage cron jobs, a cron manager library. It includes a live monitoring of current schedule and state of active jobs that can be outputted as JSON or HTML template.

## Installation

In order to install cronx package, you need to install Go and set your Go workspace first.

You first need Go installed (version 1.14+ is required), then you can use the below Go command to install cronx.

```shell
go get -v github.com/rizalgowandy/cronx
```

Import it in your code:

```shell
package main

import "github.com/rizalgowandy/cronx"
```

## Quick Start

Check the example [here](example/2-storage/main.go).

Run docker:

```shell
docker-compose up
```

Run the binary:

```shell
make run | jq -R -r '. as $line | try fromjson catch $line'
```

Then, browse to:

- http://localhost:9001 => see server health status.
- http://localhost:9001/jobs => see the current job status as UI response.
- http://localhost:9001/api/jobs => see the current job status as JSON response.
- http://localhost:9001/api/histories => see previous job run histories as JSON response.

![cronx](docs/screenshot/6_jobs_page.png)

## Available Status

* **Down** => Job fails to be registered.
* **Up** => Job has just been created.
* **Running** => Job is currently running.
* **Idle** => Job is waiting for next execution time.
* **Error** => Job fails on the last run.

## Schedule Specification Format

### Schedule

| Field name   | Mandatory? | Allowed values  | Allowed special characters |
|--------------|------------|-----------------|----------------------------|
| Seconds      | Optional   | 0-59            | * / , -                    |
| Minutes      | Yes        | 0-59            | * / , -                    |
| Hours        | Yes        | 0-23            | * / , -                    |
| Day of month | Yes        | 1-31            | * / , - ?                  |
| Month        | Yes        | 1-12 or JAN-DEC | * / , -                    |
| Day of week  | Yes        | 0-6 or SUN-SAT  | * / , - ?                  |

### Predefined schedules

| Entry                  | Description                                | Equivalent  |
|------------------------|--------------------------------------------|-------------|
| @yearly (or @annually) | Run once a year, midnight, Jan. 1st        | 0 0 0 1 1 * |
| @monthly               | Run once a month, midnight, first of month | 0 0 0 1 * * |
| @weekly                | Run once a week, midnight between Sat/Sun  | 0 0 0 * * 0 |
| @daily (or @midnight)  | Run once a day, midnight                   | 0 0 0 * * * |
| @hourly                | Run once an hour, beginning of hour        | 0 0 * * * * |

### Intervals

```
@every <duration>
```

For example, "@every 1h30m10s" would indicate a schedule that activates after 1 hour, 30 minutes, 10 seconds, and then every interval after that.

Please refer to this [link](https://pkg.go.dev/github.com/robfig/cron?readme=expanded#section-readme/) for more detail.

## Interceptor / Middleware

Interceptor or commonly known as middleware is an operation that commonly executed before any of other operation. This library has the capability to add multiple middlewares that will be executed before or after the real job. It means you can log the running job, send telemetry, or protect the application from going
down because of panic by adding middlewares. The idea of a middleware is to be declared once, and be executed on all registered jobs. Hence, reduce the code duplication on each job implementation.

### Adding Interceptor / Middleware

```go
package main

import (
	"github.com/rizalgowandy/cronx"
	"github.com/rizalgowandy/cronx/interceptor"
)

func main() {
	// Create cron middleware.
	// The order is important.
	// The first one will be executed first.
	middleware := cronx.Chain(
		interceptor.RequestID,           // Inject request id to context.
		interceptor.Recover(),           // Auto recover from panic.
		interceptor.Logger(),            // Log start and finish process.
		interceptor.DefaultWorkerPool(), // Limit concurrent running job.
	)

	cronx.NewManager(cronx.WithInterceptor(middleware))
}
```

Check all available interceptors [here](interceptor).

### Custom Interceptor / Middleware

```go
package main

import (
	"context"
	"time"

	"github.com/rizalgowandy/cronx"
)

// Sleep is a middleware that sleep a few second after job has been executed.
func Sleep() cronx.Interceptor {
	return func(ctx context.Context, job *cronx.Job, handler cronx.Handler) error {
		err := handler(ctx, job)
		time.Sleep(10 * time.Second)
		return err
	}
}
```

For more example check [here](interceptor).

## FAQ

### What are the available commands?

Here the list of commonly used commands.

```go
package main

import (
	"context"

	"github.com/rizalgowandy/cronx"
)

// Schedule sets a job to run at specific time.
// Example:
//  @every 5m
//  0 */10 * * * * => every 10m
func Schedule(spec string, job cronx.JobItf) error

// ScheduleFunc adds a func to the Cron to be run on the given schedule.
func ScheduleFunc(spec, name string, cmd func(ctx context.Context) error) error

// Schedules sets a job to run multiple times at specific time.
// Symbol */,-? should never be used as separator character.
// These symbols are reserved for cron specification.
//
// Example:
//  Spec		: "0 0 1 * * *#0 0 2 * * *#0 0 3 * * *
//  Separator	: "#"
//  This input schedules the job to run 3 times.
func Schedules(spec, separator string, job cronx.JobItf) error

// SchedulesFunc adds a func to the Cron to be run on the given schedules.
func SchedulesFunc(spec, separator, name string, cmd func(ctx context.Context) error) error
```

Go [here](cronx.go) to see the list of available commands.

### What are the available interceptors?

Go [here](interceptor) to see the available interceptors.

### Can I use my own router without starting the built-in router?

Yes, you can. This library is very modular.

Here's an example of using [gin](https://github.com/gin-gonic/gin).

```go
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rizalgowandy/cronx"
)

func main() {
	// Since we want to create custom HTTP server.
	// Do not forget to shut down the cron gracefully manually here.
	manager := cronx.NewManager()
	defer manager.Stop()

	// An example using gin as the router.
	r := gin.Default()
	r.GET("/custom-path", func(c *gin.Context) {
		c.JSON(http.StatusOK, manager.GetInfo())
	})

	// Start your own server.
	r.Run()
}
```

### Can I still get the built-in template if I use my own router?

Yes, you can.

```go
package main

import (
	"github.com/labstack/echo/v4"
	"github.com/rizalgowandy/cronx"
	"github.com/rizalgowandy/cronx/page"
)

func main() {
	// Since we want to create custom HTTP server.
	// Do not forget to shut down the cron gracefully manually here.
	manager := cronx.NewManager()
	defer manager.Stop()

	// An example using echo as the router.
	e := echo.New()
	index, _ := page.GetStatusTemplate()
	e.GET("/jobs", func(context echo.Context) error {
		// Serve the template to the writer and pass the current status data.
		return index.Execute(context.Response().Writer, manager.GetStatusData(ctx.QueryParam(cronx.QueryParamSort)))
	})
}
```

### Server is located in the US, but my user is in Jakarta, can I change the cron timezone?

Yes, you can. By default, the cron timezone will follow the server location timezone using `time.Local`. If you placed the server in the US, it will use the US timezone. If you placed the server in the SG, it will use the SG timezone.

```go
package main

import (
	"time"

	"github.com/rizalgowandy/cronx"
)

func main() {
	loc := func() *time.Location { // Change timezone to Jakarta.
		jakarta, err := time.LoadLocation("Asia/Jakarta")
		if err != nil {
			secondsEastOfUTC := int((7 * time.Hour).Seconds())
			jakarta = time.FixedZone("WIB", secondsEastOfUTC)
		}
		return jakarta
	}()

	// Create a custom config.
	cronx.NewManager(cronx.WithLocation(loc))
}
```

### My job requires certain information like current wave number, how can I get this information?

This kind of information is stored inside metadata, which stored automatically inside `context`.

```go
package main

import (
	"context"
	"errors"

	"github.com/rizalgowandy/cronx"
	"github.com/rizalgowandy/gdk/pkg/logx"
)

type subscription struct{}

func (subscription) Run(ctx context.Context) error {
	md, ok := cronx.GetJobMetadata(ctx)
	if !ok {
		return errors.New("cannot job metadata")
	}
	logx.INF(ctx, md, "subscription is running")
	return nil
}
```
