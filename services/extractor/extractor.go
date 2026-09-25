package extractor

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
)

//go:embed assets/bin/ffmpeg-bin.exe
var ffmpegExeFile []byte

type ffmpegExtractor struct {
	ffmpegBinary []byte
	outPath      string
}

func NewFfmpegExtractor() *ffmpegExtractor {
	return &ffmpegExtractor{
		ffmpegBinary: ffmpegExeFile,
	}
}

func (fe *ffmpegExtractor) Extract(currentDir bool) error {
	filename := "ffmpeg.exe"

	if currentDir {
		wd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("Falha ao obter informação do diretório atual: %w", err)
		}

		filename = filepath.Join(wd, filename)
	} else {
		tempdir := os.TempDir()
		filename = filepath.Join(tempdir, filename)
	}

	fe.outPath = filename

	color.Yellow("Extraindo binário do ffmpeg no diretório: %s\n", filepath.Dir(filename))

	err := os.WriteFile(filename, fe.ffmpegBinary, 0755)
	if err != nil {
		return fmt.Errorf("Falha ao extrair %s no diretório atual: %w", filename, err)
	}

	return nil
}

func (fe *ffmpegExtractor) Remove() error {
	if fe.outPath == "" {
		return nil
	}

	err := os.Remove(fe.outPath)
	if err != nil {
		return fmt.Errorf("Falha ao remover o binário do ffmpeg: %w", err)
	}

	return nil
}

func (fe *ffmpegExtractor) GetPath() string {
	return fe.outPath
}
