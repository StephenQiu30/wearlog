// Command server starts the Wearlog backend API and worker processes.
package main

import (
	"log/slog"
	"os"

	"github.com/StephenQiu30/wearlog-server/backend/internal/bootstrap"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	if err := bootstrap.Run(log); err != nil {
		log.Error("startup_or_runtime_failed", "reason", err.Error())
		os.Exit(1)
	}
}
