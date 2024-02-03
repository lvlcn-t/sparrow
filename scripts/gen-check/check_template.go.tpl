package {{.PackageName}}

import (
    "context"
    "sync"
    "time"

    "github.com/caas-team/sparrow/internal/helper"
    "github.com/caas-team/sparrow/internal/logger"
    "github.com/caas-team/sparrow/pkg/checks"
    "github.com/getkin/kin-openapi/openapi3"
    "github.com/prometheus/client_golang/prometheus"
)

var (
	_ checks.Check   = (*{{.CheckStructName}})(nil)
	_ checks.Runtime = (*Config)(nil)
)

const CheckName = "{{.CheckName}}"

type {{.CheckStructName}} struct {
    checks.CheckBase
    config  Config
    metrics metrics
}

// NewCheck creates a new instance of the {{.CheckName}} check
func NewCheck() checks.Check {
	return &{{.CheckStructName}}{
		CheckBase: checks.CheckBase{
			Mu:      sync.Mutex{},
			CResult: nil,
			Done:    make(chan bool, 1),
		},
		config: Config{
			ConfigBase: checks.ConfigBase{
				Retry: checks.DefaultRetry,
			},
		},
		metrics: newMetrics(),
	}
}

// Config defines the configuration parameters for a {{.CheckName}} check
type Config struct {
	checks.ConfigBase
	Targets []string
    // Add configuration specific to this check
}

// result represents the result of a single {{.CheckName}} check for a specific target
type result struct {
	// Add results specific to this check
}

// metrics defines the metric collectors of the {{.CheckName}} check
type metrics struct {
    // Define metrics specific to this check
}

func (c *Config) For() string {
    return CheckName
}

// Run starts the {{.CheckName}} check
func (c *{{.CheckStructName}}) Run(ctx context.Context) error {
	ctx, cancel := logger.NewContextWithLogger(ctx)
	defer cancel()
	log := logger.FromContext(ctx)

	log.Info("Starting {{.CheckName}} check", "interval", c.config.Interval.String())
	for {
		select {
		case <-ctx.Done():
			log.Error("Context canceled", "err", ctx.Err())
			return ctx.Err()
		case <-c.Done:
			log.Debug("Soft shut down")
			return nil
		case <-time.After(c.config.Interval):
			res := c.check(ctx)
			errval := ""
			r := checks.Result{
				Data:      res,
				Err:       errval,
				Timestamp: time.Now(),
			}

			c.CResult <- r
			log.Debug("Successfully finished {{.CheckName}} check run")
		}
	}
}

// Startup is called once when the {{.CheckName}} check is registered
func (c *{{.CheckStructName}}) Startup(ctx context.Context, cResult chan<- checks.Result) error {
	log := logger.FromContext(ctx)
	log.Debug("Initializing {{.CheckName}} check")

	c.CResult = cResult
	return nil
}

// Shutdown is called once when the check is unregistered or sparrow shuts down
func (c *{{.CheckStructName}}) Shutdown(_ context.Context) error {
	c.Done <- true
	close(c.Done)

	return nil
}

// SetConfig sets the configuration for the {{.CheckName}} check
func (c *{{.CheckStructName}}) SetConfig(cfg checks.Runtime) error {
	if conf, ok := cfg.(*Config); ok {
		if len(conf.GlobalTargets) > 0 {
			conf.Targets = append(conf.Targets, conf.GlobalTargets...)
		}
		c.Mu.Lock()
		defer c.Mu.Unlock()
		c.config = *conf
		return nil
	}

	return checks.ErrConfigMismatch{
		Expected: CheckName,
		Current:  cfg.For(),
	}
}

// GetConfig returns the current configuration of the {{.CheckName}} check
func (c *{{.CheckStructName}}) GetConfig() checks.Runtime {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	return &c.config
}

// Name returns the name of the check
func (c *{{.CheckStructName}}) Name() string {
	return CheckName
}

// Schema provides the schema of the data that will be provided by the {{.CheckName}} check
func (c *{{.CheckStructName}}) Schema() (*openapi3.SchemaRef, error) {
	return checks.OpenapiFromPerfData[map[string]result](make(map[string]result))
}

// newMetrics initializes metric collectors of the {{.CheckName}} check
func newMetrics() metrics {
	return metrics{}
}

// GetMetricCollectors returns all metric collectors of check
func (c *{{.CheckStructName}}) GetMetricCollectors() []prometheus.Collector {
	return []prometheus.Collector{}
}

// check performs a {{.CheckName}} check using a retry function
// to get the {{.CheckName}} to all targets
func (c *{{.CheckStructName}}) check(ctx context.Context) map[string]result {
	log := logger.FromContext(ctx)
	log.Debug("Checking {{.CheckName}}")
	if len(c.config.Targets) == 0 {
		log.Debug("No targets defined")
		return map[string]result{}
	}
	log.Debug("Getting {{.CheckName}} status for each target in separate routine", "amount", len(c.config.Targets))

	var mu sync.Mutex
	var wg sync.WaitGroup
	results := map[string]result{}

	// Setup {{.CheckStructName}} check client

	for _, t := range c.config.Targets {
		target := t
		wg.Add(1)
		lo := log.With("target", target)

		get{{.CheckStructName}}Retry := helper.Retry(func(ctx context.Context) error {
			res, err := get{{.CheckStructName}}(ctx, client, target)
			mu.Lock()
			defer mu.Unlock()
			results[target] = res
			if err != nil {
				return err
			}
			return nil
		}, c.config.Retry)

		go func() {
			defer wg.Done()

			lo.Debug("Starting retry routine to get {{.CheckName}} status")
			if err := get{{.CheckStructName}}Retry(ctx); err != nil {
				lo.Error("Error while checking {{.CheckName}}", "error", err)
			}

			lo.Debug("Successfully got {{.CheckName}} status of target")
		}()
	}

	log.Debug("Waiting for all routines to finish")
	wg.Wait()

	log.Debug("Successfully got {{.CheckName}} status from all targets")
	return results
}

// get{{.CheckStructName}}
func get{{.CheckStructName}}(ctx context.Context, c *client, target string) (res result, err error) {
}