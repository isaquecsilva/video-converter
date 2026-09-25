package utils

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/fatih/color"
)

type ffmpegProgressWriter struct {
	re *regexp.Regexp
}

func NewFfmpegProgressWriter() *ffmpegProgressWriter {
	return &ffmpegProgressWriter{
		re: regexp.MustCompile(`out_time=[0-9]*:[0-9]*:[0-9]*\.*[0-9]*`),
	}
}

func (pw *ffmpegProgressWriter) Write(p []byte) (int, error) {
	s := color.BlueString(pw.re.FindString(string(p)))
	fmt.Printf("\r%s", strings.Replace(s, "out_time=", "progresso=", 1))
	return len(p), nil
}

// frame=14989
// fps=116.05
// stream_0_0_q=-0.0
// bitrate=1625.6kbits/s
// total_size=127139888
// out_time_us=625706167
// out_time_ms=625706167
// out_time=00:10:25.706167
// dup_frames=0
// drop_frames=0
// speed=4.84x
// progress=continue
