package util

import (
	"context"
	"path/filepath"
	"time"

	checkpoint "github.com/solo-io/go-checkpoint"
)

func Telemetry(version string, t time.Time) {
	sigfile := filepath.Join(HomeDir(), ".glooctl.sig")
	configDir, err := ConfigDir()
	if err == nil {
		sigfile = filepath.Join(configDir, "telemetry.sig")
	}
	ctx := context.Background()
	report := &checkpoint.ReportParams{
		Product:       "glooctl",
		Version:       version,
		StartTime:     t,
		EndTime:       time.Now(),
		SignatureFile: sigfile,
	}
	checkpoint.Report(ctx, report)
}
