package main

import (
	// "encoding/json"
	"flag"
	"fmt"

	// "io/ioutil"
	"os"

	ctrdlog "github.com/containerd/containerd/log"
	log "github.com/sirupsen/logrus"

	fc "github.com/eth-easl/easyloader/internal/function"
	tc "github.com/eth-easl/easyloader/internal/trace"
)

var (
	debug    = flag.Bool("dbg", false, "Enable debug logging")
	rps      = flag.Int("rps", 100, "Request per second (default: 100)")
	duration = flag.Int("duration", 1, "Duration of the experiment (default: 1 min)")
)

func init() {

	flag.Parse()

	log.SetFormatter(&log.TextFormatter{
		TimestampFormat: ctrdlog.RFC3339NanoFixed,
		FullTimestamp:   true,
	})
	log.SetOutput(os.Stdout)
	if *debug {
		log.SetLevel(log.DebugLevel)
		log.Debug("Debug logging is enabled")
	} else {
		log.SetLevel(log.InfoLevel)
	}
}

func main() {
	// deploymentConcurrency := flag.Int("conc", 1, "Number of functions to deploy concurrently (for serving)")
	serviceConfigPath := "workloads/timed.yaml"
	// write the whole body at once

	traces := tc.ParseInvocationTrace("data/invocations_10.csv", *duration)
	tc.ParseDurationTrace(&traces, "data/durations_10.csv")
	tc.ParseMemoryTrace(&traces, "data/memory_10.csv")

	log.Info("Traces contain the following: ", len(traces.Functions), " functions")
	for _, function := range traces.Functions {
		fmt.Println("\t" + function.GetName())
	}

	/* Deployment */
	log.Info("Using service config file: ", serviceConfigPath)
	functions := fc.Deploy(traces.Functions, serviceConfigPath, 1) // TODO: Fixed number of functions per pod.

	/* Invokation */
	defer fc.Invoke(*rps, functions, traces.InvocationsPerMin, traces.TotalInvocationsEachMin)
}
