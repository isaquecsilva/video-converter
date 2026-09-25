package vcut

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
	"video-converter/utils"
)

type videoCutter struct {
	ffmpegPath     string
	progressWriter io.Writer
}

func NewVideoCutter(ffmpegPath string, progressWriter io.Writer) VideoCutterInterface {
	return &videoCutter{
		ffmpegPath:     ffmpegPath,
		progressWriter: progressWriter,
	}
}

func (vc *videoCutter) Cut(ctx context.Context, params CutParams) (string, error) {
	if err := vc.validateTimestampsForCut(params.From, params.To, params.TotalSeconds); err != nil {
		return "", err
	}

	filename := vc.filename(params.InputFile, params.From, params.To, params.Format.Ext())

	cmd := exec.CommandContext(ctx, vc.ffmpegPath,
		"-loglevel", os.Getenv("FFMPEG_LOG_LEVEL"),
		"-hide_banner",
		"-ss", params.From,
		"-to", params.To,
		"-i", params.InputFile,
		"-f", params.Format.Name(),
		"-vcodec", params.Codec.String(),
		"-progress", "-",
		filename,
		"-y",
	)

	if vc.progressWriter == nil {
		vc.progressWriter = os.Stdout
	}

	cmd.Stdout = vc.progressWriter
	cmd.Stderr = os.Stderr
	fmt.Println("Comando ffmpeg usado: ", cmd.String())

	return filename, cmd.Run()
}

func (vc *videoCutter) BatchCut(ctx context.Context, params CutInIntervalParams) ([]string, error) {
	if params.Interval > params.TotalSeconds {
		interval := time.Duration(params.Interval) * time.Second
		duration := time.Duration(params.TotalSeconds) * time.Second
		return nil, fmt.Errorf("Intervalo %s deve ser menor que a duração do vídeo, que é: %s", interval.String(), duration.String())
	}

	filename := vc.batchFilename(params.InputFile, params.Interval, params.TotalSeconds, params.Format.Ext())

	cmd := exec.CommandContext(ctx, vc.ffmpegPath,
		"-loglevel", os.Getenv("FFMPEG_LOG_LEVEL"),
		"-hide_banner",
		"-i", params.InputFile,
		"-f", "segment",
		"-segment_time", fmt.Sprint(params.Interval),
		"-reset_timestamps", "1",
		"-c:v", params.Codec.String(),
		"-preset", "fast",
		"-progress", "-",
		filename,
		"-y",
	)

	if vc.progressWriter == nil {
		vc.progressWriter = os.Stdout
	}

	cmd.Stdout = vc.progressWriter
	cmd.Stderr = os.Stderr

	fmt.Println("Comando ffmpeg usado: ", cmd.String())

	err := cmd.Run()

	if err != nil {
		return nil, err
	}

	return vc.getListOfAllCuttedFiles(params.InputFile, params.Format.Ext())
}

func (vc *videoCutter) getListOfAllCuttedFiles(inputFilename string, ext string) ([]string, error) {
	files, err := os.ReadDir(filepath.Dir(inputFilename))

	if err != nil {
		return nil, fmt.Errorf("Falha ao ler diretório contendo os arquivos cortados: %v", err)
	}

	basename := strings.Trim(filepath.Base(inputFilename), filepath.Ext(inputFilename))

	cutted := slices.Collect(func(yield func(string) bool) {
		re := regexp.MustCompile(fmt.Sprintf("%s_[0-9]+%s$", basename, strings.ReplaceAll(ext, ".", "\\.")))

		for _, f := range files {
			if strings.Contains(f.Name(), basename) {
				filename := filepath.Join(filepath.Dir(inputFilename), filepath.Join(f.Name()))

				if re.MatchString(f.Name()) && !yield(filename) {
					return
				}
			}
		}
	})

	return cutted, nil
}

func (vc *videoCutter) filename(inputFilename, from, to string, _ext string) string {
	basedir := filepath.Dir(inputFilename)
	basename := filepath.Base(inputFilename)
	ext := filepath.Ext(inputFilename)

	if ext != "" {
		basename = basename[:strings.Index(basename, ext)]
	}

	textEscapedTimestamps := make([]string, 2)

	// /home/user123/videos/myvideo_004215_021237.mp4

	for i, s := range []string{from, to} {
		textEscapedTimestamps[i] = strings.ReplaceAll(s, ":", "")
	}

	filename := filepath.Join(basedir, basename+"_"+strings.Join(textEscapedTimestamps, "_")+_ext)
	return filename
}

func (vc *videoCutter) batchFilename(inputFilename string, intervalSec, totalSec int64, ext string) string {
	basedir := filepath.Dir(inputFilename)
	basename := strings.Trim(filepath.Base(inputFilename), filepath.Ext(inputFilename))

	totalClipsNeeded := int(totalSec / intervalSec)

	ffmpegNumerationFormat := fmt.Sprintf("%%0%dd", len(strconv.Itoa(totalClipsNeeded)))

	return filepath.Join(basedir, fmt.Sprintf("%s_%s%s", basename, ffmpegNumerationFormat, ext))
}

func (vc *videoCutter) validateTimestampsForCut(from, to string, totalSeconds int64) error {
	f, err := utils.TimestampParser(from)
	if err != nil {
		return fmt.Errorf("Timestamp inválido: %v", err)
	}

	if f.Seconds() < 0 {
		return fmt.Errorf("Timestamp inválido: %v", f.String())
	}

	t, err := utils.TimestampParser(to)
	if err != nil {
		return fmt.Errorf("Timestamp inválido: %v", err)
	}

	if f.Seconds() >= t.Seconds() {
		return fmt.Errorf("O timestamp de inicío deve ser menor do que o final. O que foi fornecido é: %s - %s",
			from, to)
	}

	if t.Seconds() > float64(totalSeconds) {
		return fmt.Errorf("O timestamp final deve ser menor ou igual ao timestamp em que o vídeo acaba")
	}

	return nil
}
