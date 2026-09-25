package converter

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	"video-converter/types/codecs"
	"video-converter/types/videoformat"
)

type VideoConverter struct {
	ffmpegPath     string
	progressWriter io.Writer
}

func NewVideoConverter(ffmpegPath string, progressWriter io.Writer) *VideoConverter {
	return &VideoConverter{
		ffmpegPath:     ffmpegPath,
		progressWriter: progressWriter,
	}
}

type ConvertParams struct {
	InputFile    string                  `json:"input_file"`
	VideoBitrate int                     `json:"bitrate"`
	Codec        codecs.VideoCodec       `json:"codec"`
	Format       videoformat.VideoFormat `json:"format"`
}

func (vc *VideoConverter) Convert(ctx context.Context, params ConvertParams) (string, error) {
	dir := filepath.Dir(params.InputFile)
	ext := filepath.Ext(params.InputFile)
	basename := strings.ReplaceAll(filepath.Base(params.InputFile), ext, "") +
		fmt.Sprint("_converted_", time.Now().Unix(), params.Format.Ext())
	filename := filepath.Join(dir, basename)

	args := []string{
		"-loglevel", os.Getenv("FFMPEG_LOG_LEVEL"),
		"-hide_banner",
		"-i", params.InputFile,
		"-b:v", fmt.Sprintf("%dk", params.VideoBitrate),
		"-c:v", params.Codec.String(),
		"-f", params.Format.Name(),
		"-preset", "fast",
		"-progress", "-",
		filename,
	}

	if vc.progressWriter == nil {
		vc.progressWriter = os.Stdout
	}

	cmd := exec.CommandContext(ctx, vc.ffmpegPath, args...)
	cmd.Stdout = vc.progressWriter
	cmd.Stderr = os.Stderr

	fmt.Println("Comando ffmpeg usado: ", cmd.String())

	return filename, cmd.Run()
}
