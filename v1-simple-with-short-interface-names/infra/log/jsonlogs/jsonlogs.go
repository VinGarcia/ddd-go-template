package jsonlogs

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/vingarcia/ddd-go-template/v1-simple-with-short-interface-names/infra/log"
	"github.com/vingarcia/ddd-go-template/v1-simple-with-short-interface-names/infra/maps"
)

// Client is the logger client, to instantiate it call `New()`
type Client struct {
	priorityLevel uint
	PrintlnFn     func(...interface{})

	ctxParsers []ContextParser
}

type ContextParser func(ctx context.Context) log.Body

// New builds a logger Client on the appropriate log level
func New(level string, parsers ...ContextParser) Client {
	var priority uint
	switch strings.ToUpper(level) {
	case "DEBUG":
		priority = 0
	case "INFO":
		priority = 1
	case "WARN":
		priority = 2
	case "ERROR":
		priority = 3
	default:
		priority = 1
	}

	return Client{
		priorityLevel: priority,
		PrintlnFn: func(args ...interface{}) {
			fmt.Println(args...)
		},
		ctxParsers: parsers,
	}
}

// Debug logs an entry on level "DEBUG" with the received title
// along with all the values collected from the input valueMaps and the context.
func (c Client) Debug(ctx context.Context, title string, valueMaps ...log.Body) {
	if c.priorityLevel > 0 {
		return
	}

	c.log(ctx, "DEBUG", title, valueMaps)
}

// Info logs an entry on level "INFO" with the received title
// along with all the values collected from the input valueMaps and the context.
func (c Client) Info(ctx context.Context, title string, valueMaps ...log.Body) {
	if c.priorityLevel > 1 {
		return
	}

	c.log(ctx, "INFO", title, valueMaps)
}

// Warn logs an entry on level "WARN" with the received title
// along with all the values collected from the input valueMaps and the context.
func (c Client) Warn(ctx context.Context, title string, valueMaps ...log.Body) {
	if c.priorityLevel > 2 {
		return
	}

	c.log(ctx, "WARN", title, valueMaps)
}

// Error logs an entry on level "ERROR" with the received title
// along with all the values collected from the input valueMaps and the context.
func (c Client) Error(ctx context.Context, title string, valueMaps ...log.Body) {
	if c.priorityLevel > 3 {
		return
	}

	c.log(ctx, "ERROR", title, valueMaps)
}

// Fatal logs an entry on level "ERROR" with the received title
// along with all the values collected from the input valueMaps and the context.
//
// After that it proceeds to exit the program with code 1.
func (c Client) Fatal(ctx context.Context, title string, valueMaps ...log.Body) {
	if c.priorityLevel > 3 {
		return
	}

	c.log(ctx, "ERROR", title, valueMaps)
	os.Exit(1)
}

func (c Client) log(ctx context.Context, level string, title string, valueMaps []log.Body) {
	body := log.Body{}
	for _, parser := range c.ctxParsers {
		maps.Merge(&body, parser(ctx))
	}
	maps.Merge(&body, valueMaps...)

	c.PrintlnFn(buildJSONString(level, title, body))
}

func buildJSONString(level string, title string, body log.Body) string {
	timestamp := time.Now().Format(time.RFC3339)

	// Remove reserved keys from the input map:
	delete(body, "level")
	delete(body, "title")
	delete(body, "timestamp")

	var separator = ""
	var bodyJSON = []byte("{}")
	var err error
	if len(body) > 0 {
		separator = ","

		bodyJSON, err = json.Marshal(body)
		if err != nil {
			// Marshalling this string is necessary for
			// escaping characters such as '"'
			bodyString, _ := json.Marshal(fmt.Sprintf("%#v", body))

			return fmt.Sprintf(
				`{"level":"ERROR","title":"could-not-marshal-log-body","timestamp":"%s","body":%s}%s`,
				time.Now().Format(time.RFC3339),
				bodyString,
				"\n",
			)
		}
	}

	titleJSON, _ := json.Marshal(title)
	return fmt.Sprintf(
		`{"level":"%s","title":%s,"timestamp":"%s"%s%s`,
		level,
		string(titleJSON),
		timestamp,
		separator,
		string(bodyJSON[1:]),
	)
}
