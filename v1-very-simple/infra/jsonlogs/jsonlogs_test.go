package jsonlogs

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/vingarcia/ddd-go-template/v1-very-simple/domain"
	tt "github.com/vingarcia/ddd-go-template/v1-very-simple/infra/testtools"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		desc             string
		level            string
		expectedPriority uint
	}{
		{
			desc:             "should work for debug level",
			level:            "DEBUG",
			expectedPriority: 0,
		},
		{
			desc:             "should work for info level",
			level:            "INFO",
			expectedPriority: 1,
		},
		{
			desc:             "should work for warn level",
			level:            "WARN",
			expectedPriority: 2,
		},
		{
			desc:             "should work for error level",
			level:            "ERROR",
			expectedPriority: 3,
		},
		{
			desc:             "should default to info when input is unexpected",
			level:            "unexpected input",
			expectedPriority: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			instance := NewClient(test.level)
			assert.Equal(t, test.expectedPriority, instance.priorityLevel)
		})
	}
}

func TestBuildJSONString(t *testing.T) {
	tests := []struct {
		desc           string
		level          string
		title          string
		body           domain.LogBody
		expectedOutput map[string]interface{}
	}{
		{
			desc:  "should work with empty bodies",
			level: "fake-level",
			title: "fake-title",
			body:  domain.LogBody{},
			expectedOutput: map[string]interface{}{
				"level": "fake-level",
				"title": "fake-title",
				// "timestamp": "can't compare timestamps",
			},
		},
		{
			desc:  "should work with non empty bodies",
			level: "fake-level",
			title: "fake-title",
			body: domain.LogBody{
				"fake-key": "fake-timestamp",
			},
			expectedOutput: map[string]interface{}{
				"level": "fake-level",
				"title": "fake-title",
				// "timestamp": "can't compare timestamps",
				"fake-key": "fake-timestamp",
			},
		},
		{
			desc:  "should ignore reserved fields on body",
			level: "fake-level",
			title: "fake-title",
			body: domain.LogBody{
				"level":     "fake-level2",
				"title":     "fake-title2",
				"timestamp": "fake-timestamp2",
			},
			expectedOutput: map[string]interface{}{
				"level": "fake-level",
				"title": "fake-title",
				// "timestamp": "can't compare timestamps",
			},
		},
		{
			desc:  "should output an error log when unable to marshal the body",
			level: "fake-level",
			title: "fake-title",
			body: domain.LogBody{
				"value": CannotBeMarshaled{},
			},
			expectedOutput: map[string]interface{}{
				"level": "ERROR",
				"title": "could-not-marshal-log-body",
				"body":  fmt.Sprintf("%#v", domain.LogBody{"value": CannotBeMarshaled{}}),
				// "timestamp": "can't compare timestamps",
			},
		},
	}

	for _, test := range tests {
		jsonString := buildJSONString(test.level, test.title, test.body)
		fmt.Println("String:", jsonString)

		var output map[string]interface{}
		err := json.Unmarshal([]byte(jsonString), &output)
		assert.Nil(t, err)

		timestring, ok := output["timestamp"].(string)
		assert.True(t, ok)

		_, err = time.Parse(time.RFC3339, timestring)
		assert.Nil(t, err)

		delete(output, "timestamp")
		assert.Equal(t, test.expectedOutput, output)
	}
}

func TestLogFuncs(t *testing.T) {
	ctx := context.Background()

	t.Run("debug logs should produce logs correctly", func(t *testing.T) {
		var output string
		client := Client{
			priorityLevel: 0,
			PrintlnFn: func(args ...interface{}) {
				output = fmt.Sprintln(args...)
			},
			ctxParsers: []ContextParser{domain.GetCtxValues},
		}

		ctx = domain.CtxWithValues(ctx, domain.LogBody{
			"ctx_value1": "overwritten",
			"ctx_value2": "overwritten",
			"ctx_value3": "not-overwritten",
		})

		client.Debug(
			ctx,
			"fake-log-title",
			domain.LogBody{
				"ctx_value1":   "overwrites",
				"body1_value1": "overwritten",
				"body1_value2": "not-overwritten",
			},
			domain.LogBody{
				"ctx_value2":   "overwrites",
				"body1_value1": "overwrites",
				"body2_value1": "not-overwritten",
			},
		)

		var outputMap map[string]interface{}
		err := json.Unmarshal([]byte(output), &outputMap)
		tt.AssertNoErr(t, err)

		// Default log values:
		tt.AssertEqual(t, outputMap["level"], "DEBUG")
		tt.AssertEqual(t, outputMap["title"], "fake-log-title")
		tt.AssertNotEqual(t, outputMap["timestamp"], nil)

		// Final contextual values:
		tt.AssertEqual(t, outputMap["ctx_value1"], "overwrites")
		tt.AssertEqual(t, outputMap["ctx_value2"], "overwrites")
		tt.AssertEqual(t, outputMap["body1_value1"], "overwrites")
		tt.AssertEqual(t, outputMap["ctx_value3"], "not-overwritten")
		tt.AssertEqual(t, outputMap["body1_value2"], "not-overwritten")
		tt.AssertEqual(t, outputMap["body2_value1"], "not-overwritten")

		// No overwritten field should be present:
		tt.AssertFalse(t, strings.Contains(output, `"overwritten"`))
	})

	t.Run("debug logs should be ignored if priorityLevel > 0", func(t *testing.T) {
		var output string
		client := Client{
			priorityLevel: 1,
			PrintlnFn: func(args ...interface{}) {
				output = fmt.Sprintln(args...)
			},
		}

		client.Debug(ctx, "fake-log-title")

		tt.AssertEqual(t, "", output)
	})

	t.Run("info logs should produce logs correctly", func(t *testing.T) {
		var output string
		client := Client{
			priorityLevel: 0,
			PrintlnFn: func(args ...interface{}) {
				output = fmt.Sprintln(args...)
			},
			ctxParsers: []ContextParser{domain.GetCtxValues},
		}

		ctx = domain.CtxWithValues(ctx, domain.LogBody{
			"ctx_value1": "overwritten",
			"ctx_value2": "overwritten",
			"ctx_value3": "not-overwritten",
		})

		client.Info(
			ctx,
			"fake-log-title",
			domain.LogBody{
				"ctx_value1":   "overwrites",
				"body1_value1": "overwritten",
				"body1_value2": "not-overwritten",
			},
			domain.LogBody{
				"ctx_value2":   "overwrites",
				"body1_value1": "overwrites",
				"body2_value1": "not-overwritten",
			},
		)

		tt.AssertContains(t, output,
			// Default log values:
			`"level":"INFO"`,
			`"title":"fake-log-title"`,
			`"timestamp"`,

			// Final contextual values:
			`"ctx_value1":"overwrites"`,
			`"ctx_value2":"overwrites"`,
			`"body1_value1":"overwrites"`,
			`"ctx_value3":"not-overwritten"`,
			`"body1_value2":"not-overwritten"`,
			`"body2_value1":"not-overwritten"`,
		)

		// No overwritten field should be present:
		tt.AssertFalse(t, strings.Contains(output, `"overwritten"`))
	})

	t.Run("info logs should be ignored if priorityLevel > 1", func(t *testing.T) {
		var output string
		client := Client{
			priorityLevel: 2,
			PrintlnFn: func(args ...interface{}) {
				output = fmt.Sprintln(args...)
			},
		}

		client.Info(ctx, "fake-log-title")

		tt.AssertEqual(t, output, "")
	})

	t.Run("warn logs should produce logs correctly", func(t *testing.T) {
		var output string
		client := Client{
			priorityLevel: 0,
			PrintlnFn: func(args ...interface{}) {
				output = fmt.Sprintln(args...)
			},
			ctxParsers: []ContextParser{domain.GetCtxValues},
		}

		ctx = domain.CtxWithValues(ctx, domain.LogBody{
			"ctx_value1": "overwritten",
			"ctx_value2": "overwritten",
			"ctx_value3": "not-overwritten",
		})

		client.Warn(
			ctx,
			"fake-log-title",
			domain.LogBody{
				"ctx_value1":   "overwrites",
				"body1_value1": "overwritten",
				"body1_value2": "not-overwritten",
			},
			domain.LogBody{
				"ctx_value2":   "overwrites",
				"body1_value1": "overwrites",
				"body2_value1": "not-overwritten",
			},
		)

		tt.AssertContains(t, output,
			// Default log values:
			`"level":"WARN"`,
			`"title":"fake-log-title"`,
			`"timestamp"`,

			// Final contextual values:
			`"ctx_value1":"overwrites"`,
			`"ctx_value2":"overwrites"`,
			`"body1_value1":"overwrites"`,
			`"ctx_value3":"not-overwritten"`,
			`"body1_value2":"not-overwritten"`,
			`"body2_value1":"not-overwritten"`,
		)

		// No overwritten field should be present:
		tt.AssertFalse(t, strings.Contains(output, `"overwritten"`))
	})

	t.Run("warn logs should be ignored if priorityLevel > 2", func(t *testing.T) {
		var output string
		client := Client{
			priorityLevel: 3,
			PrintlnFn: func(args ...interface{}) {
				output = fmt.Sprintln(args...)
			},
		}

		client.Warn(ctx, "fake-log-title")

		tt.AssertEqual(t, output, "")
	})

	t.Run("error logs should produce logs correctly", func(t *testing.T) {
		var output string
		client := Client{
			priorityLevel: 0,
			PrintlnFn: func(args ...interface{}) {
				output = fmt.Sprintln(args...)
			},
			ctxParsers: []ContextParser{domain.GetCtxValues},
		}

		ctx = domain.CtxWithValues(ctx, domain.LogBody{
			"ctx_value1": "overwritten",
			"ctx_value2": "overwritten",
			"ctx_value3": "not-overwritten",
		})

		client.Error(
			ctx,
			"fake-log-title",
			domain.LogBody{
				"ctx_value1":   "overwrites",
				"body1_value1": "overwritten",
				"body1_value2": "not-overwritten",
			},
			domain.LogBody{
				"ctx_value2":   "overwrites",
				"body1_value1": "overwrites",
				"body2_value1": "not-overwritten",
			},
		)

		tt.AssertContains(t, output,
			// Default log values:
			`"level":"ERROR"`,
			`"title":"fake-log-title"`,
			`"timestamp"`,

			// Final contextual values:
			`"ctx_value1":"overwrites"`,
			`"ctx_value2":"overwrites"`,
			`"body1_value1":"overwrites"`,
			`"ctx_value3":"not-overwritten"`,
			`"body1_value2":"not-overwritten"`,
			`"body2_value1":"not-overwritten"`,
		)

		// No overwritten field should be present:
		tt.AssertFalse(t, strings.Contains(output, `"overwritten"`))
	})

	t.Run("error logs should be ignored if priorityLevel > 3", func(t *testing.T) {
		var output string
		client := Client{
			priorityLevel: 4,
			PrintlnFn: func(args ...interface{}) {
				output = fmt.Sprintln(args...)
			},
		}

		client.Error(ctx, "fake-log-title")

		tt.AssertEqual(t, output, "")
	})
}

type CannotBeMarshaled struct{}

func (c CannotBeMarshaled) MarshalJSON() ([]byte, error) {
	return nil, fmt.Errorf("fake-error-message")
}
