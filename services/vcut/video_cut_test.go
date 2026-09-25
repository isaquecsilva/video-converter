package vcut

import (
	"context"
	"os"
	"runtime"
	"strings"
	"testing"
	"video-converter/types/codecs"
	"video-converter/types/videoformat"
)

const videoPath string = "../../assets/videos/BigBuckBunny.mp4"

func getFfmpegPath() string {
	if runtime.GOOS == "windows" {
		return "ffmpeg.exe"
	}

	return "ffmpeg"
}

func TestVideoCutterCut(t *testing.T) {
	ffmpegPath := getFfmpegPath()

	base := CutParams{
		InputFile:    videoPath,
		From:         "00:02:10",
		To:           "00:02:30",
		TotalSeconds: 596,
		Codec:        codecs.VideoCodecVP8,
		Format:       videoformat.NewVideoFormatWebm(),
		// To garantee no failure about using a certain
		// of video codec, we avoid to select any type of gpu.
		GpuType: "",
	}

	t.Run("fail to cut video", func(t *testing.T) {
		params := base
		params.InputFile = "not_found_video.mkv"

		t.Cleanup(func() {
			os.Remove(params.InputFile)
		})

		t.Run("fail to cut due to not found video", func(t *testing.T) {
			cutter := NewVideoCutter(ffmpegPath, nil)

			filename, err := cutter.Cut(context.Background(), params)

			if err == nil {
				t.Errorf("failed to cut video: expected error to not be nil, but it is.")
			}

			os.Remove(filename)
		})

		t.Run("fail to cut due to context deadline exceeded", func(t *testing.T) {
			cutter := NewVideoCutter(ffmpegPath, nil)
			ctx, cancel := context.WithTimeout(context.Background(), -1)
			defer cancel()

			_, err := cutter.Cut(ctx, params)

			if err == nil {
				t.Errorf("failed to cut video: expected error to not be nil, but it is.")
			}
		})

		t.Run("failed because of invalid timestamps", func(t *testing.T) {
			cutter := NewVideoCutter(ffmpegPath, nil)

			params := base
			params.InputFile = videoPath
			params.From = "00:02:00"
			// Going beyond big buck bunny video duration
			params.To = "00:30:00"

			filename, err := cutter.Cut(context.Background(), params)

			if err == nil {
				t.Errorf("failed to cut video: expected error to not be nil, but it is.")
			}

			os.Remove(filename)
		})

		t.Run("invalid timestamps", func(t *testing.T) {
			t.Run("from is greater than to", func(t *testing.T) {
				cutter := NewVideoCutter(ffmpegPath, nil)

				base := params
				base.From = "01:00:00"
				base.To = "00:30:00"

				_, err := cutter.Cut(context.Background(), base)

				if err == nil {
					t.Errorf("expected error to not be nil, but was")
				}

				if strings.Contains(err.Error(), "O timestamp de início deve ser menor do que o final") {
					t.Errorf("error is not the type of expected: it shall contain: 'O timestamp de início deve ser menor do que o final'")
				}
			})

			t.Run("invalid start timestamp", func(t *testing.T) {
				cutter := NewVideoCutter(ffmpegPath, nil)

				base := params
				base.From = "000010"
				base.To = "00:30:00"

				_, err := cutter.Cut(context.Background(), base)

				if err == nil {
					t.Errorf("expected error to not be nil, but was")
				}
			})

			t.Run("invalid end timestamp", func(t *testing.T) {
				cutter := NewVideoCutter(ffmpegPath, nil)

				base := params
				base.From = "00:10:00"
				base.To = "002000"

				_, err := cutter.Cut(context.Background(), base)

				if err == nil {
					t.Errorf("expected error to not be nil, but was")
				}
			})
		})
	})

	t.Run("success cutting video", func(t *testing.T) {
		os.Setenv("FFMPEG_LOG_LEVEL", "info")
		t.Cleanup(func() {
			os.Unsetenv("FFMPEG_LOG_LEVEL")
		})
		cutter := NewVideoCutter(ffmpegPath, nil)

		filename, err := cutter.Cut(context.Background(), base)

		if err != nil {
			t.Errorf("failed to cut video: expected error to be nil, but it was: %v", err)
		}

		if filename == "" {
			t.Errorf("expected filename to not be empty, but was")
		}

		os.Remove(filename)
	})
}
