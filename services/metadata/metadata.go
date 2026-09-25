package metadata

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
	"video-converter/utils"
)

type VideoMetadata struct {
	VideoBitrateInKbps int
	Duration           time.Duration
}

type VideoMetadataExtractor struct {
	ffmpegPath string
}

func NewVideoMetadataExtractor(ffmpegPath string) *VideoMetadataExtractor {
	return &VideoMetadataExtractor{ffmpegPath: ffmpegPath}
}

func (vme *VideoMetadataExtractor) ExtractMetadata(inputFile string) (*VideoMetadata, error) {
	inputFile = strings.Trim(inputFile, "\"'\r\n")

	buffer := new(bytes.Buffer)
	cmd := exec.Command(vme.ffmpegPath, "-i", inputFile, "-hide_banner", "-f", "ffmetadata", "-")
	cmd.Stderr = buffer
	err := cmd.Run()

	if err != nil {
		return nil, fmt.Errorf("Falha ao extrair metadados do vídeo: %v", err)
	}

	bitrateStr := strings.TrimSpace(
		strings.Replace(
			regexp.MustCompile(`\d+\skb/s`).FindString(buffer.String()),
			"kb/s", "", 1,
		),
	)

	durationStr := strings.TrimSpace(
		strings.TrimPrefix(
			regexp.MustCompile(`Duration:\s[0-9]{2}:[0-9]{2}:[0-9]{2}(\.[0-9]+)*`).FindString(buffer.String()),
			"Duration: ",
		),
	)

	bitrate, err := strconv.Atoi(bitrateStr)
	if err != nil {
		return nil, fmt.Errorf("Falha ao converter bitrate para inteiro: %v", err)
	}

	duration, err := utils.TimestampParser(durationStr)
	if err != nil {
		return nil, fmt.Errorf("Falha ao identificar a duração do vídeo: %v", err)
	}
	return &VideoMetadata{VideoBitrateInKbps: bitrate, Duration: duration}, nil
}
