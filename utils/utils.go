package utils

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"
)

func BytesToMebibytes(bytes int64) float64 {
	return float64(bytes) / (1024 * 1024)
}

func OpenFileInExplorer(path string) {
	exec.Command("explorer", path).Run()
}

func TimestampParser(timestamp string) (time.Duration, error) {
	parts := strings.Split(timestamp, ":")
	if len(parts) != 3 {
		return -1, errors.New("timestamp fornecido é inválido")
	}

	units := [3]string{"h", "m", "s"}
	pos := 0

	var sb strings.Builder

	for _, str := range parts {
		sb.WriteString(str + units[pos])
		pos++

		if pos == 3 {
			break
		}
	}

	return time.ParseDuration(sb.String())
}

func ValidateTimestampFormat(timestamp string) error {
	re := regexp.MustCompile(`[0-9]{2,}:[0-9]{2,}:[0-9]{2,}`)

	if re.MatchString(timestamp) != true {
		return fmt.Errorf("Timestamp %s é inválido", timestamp)
	}

	return nil
}

func GetFfmpegLogLevel() string {
	level, ok := os.LookupEnv("FFMPEG_LOG_LEVEL")
	if ok == false {
		level = "error"
	}

	return level
}

func CleanString(s string) string {
	return strings.Trim(s, "\r\n ")
}

func MustGetUserInput(reader *bufio.Reader) string {
	choice, _ := reader.ReadString('\n')
	choice = strings.Trim(choice, "\r\n ")
	return choice
}

func SignalContext() (context.Context, context.CancelFunc) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)
	return ctx, stop
}
