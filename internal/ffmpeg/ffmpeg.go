package ffmpeg

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os/exec"
	"time"
)

type VideoSnapshotConfig struct {
	Width    int
	Height   int
	Position time.Duration
}

func VideoSnapshot(ctx context.Context, inputPath, outputFormat string, outputWriter io.Writer, cfg VideoSnapshotConfig) error {
	var stderr bytes.Buffer

	// ffmpeg -hide_banner -i file:input.dav -ss 00:00:06.000 -vframes 1 pipe:1.jpg
	var args []string = []string{
		"-hide_banner",
		"-n",
		"-i", inputPath,
		"-ss", fmt.Sprintf("%s.000", time.Time{}.Add(cfg.Position).Format(time.TimeOnly)),
		"-vframes", "1",
	}
	if cfg.Width != 0 || cfg.Height != 0 {
		args = append(args, "-vf", fmt.Sprintf("scale=%d:%d", cfg.Width, cfg.Height))
	}
	args = append(args, fmt.Sprintf("pipe:1.%s", outputFormat))
	cmd := exec.Command(
		"ffmpeg",
		args...,
	)
	cmd.Stdout = outputWriter
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%w: %s", err, stderr.String())
	}

	return nil
}

type ImageSnapshotConfig struct {
	Width  int
	Height int
}

func ImageSnapshot(ctx context.Context, inputPath, outputFormat string, outputWriter io.Writer, cfg ImageSnapshotConfig) error {
	var stderr bytes.Buffer

	// ffmpeg -hide_banner file:input.jpg -vf scale=320:-1 pipe:1.jpg
	args := []string{
		"-hide_banner",
		"-i", inputPath,
	}
	if true {
		args = append(args, "-vf", fmt.Sprintf("scale=%d:%d", cfg.Width, cfg.Height))
	}
	args = append(args, fmt.Sprintf("pipe:1.%s", outputFormat))
	cmd := exec.Command(
		"ffmpeg",
		args...,
	)
	cmd.Stdout = outputWriter
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%w: %s", err, stderr.String())
	}

	return nil
}
