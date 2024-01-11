package ffmpeg

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

func GenerateSnapshot(ctx context.Context, inputPath, outputPath string, position time.Duration) error {
	// ffmpeg -n -i file:input.dav -ss 00:00:06.000 -vframes 1 output.jpg
	output, err := exec.Command(
		"ffmpeg",
		"-n",
		"-i", inputPath,
		"-ss", fmt.Sprintf("%s.000", time.Time{}.Add(position).Format(time.TimeOnly)),
		"-vframes", "1",
		outputPath,
	).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: %b", err, output)
	}

	return nil
}
