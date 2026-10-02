package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"video-converter/flows"
	"video-converter/services/converter"
	"video-converter/services/discover"
	"video-converter/services/extractor"
	"video-converter/services/metadata"
	"video-converter/services/vcut"
	"video-converter/types/videoformat"
	"video-converter/utils"
	"video-converter/utils/printer"
)

var (
	currentDir = flag.Bool("current-dir", false, "Extrair ffmpeg.exe no diretório atual")
	openDir    = flag.Bool("opendir", false, "Abrir pasta de saída ao terminar conversão.")
	version    = flag.Bool("version", false, "Exibe a versão da aplicação.")
	logLevel   = flag.String("ffmpegLogLevel", "error", "Qual nível de logs usar para o ffmpeg")
)

func main() {
	flag.Parse()

	if !flag.Parsed() {
		log.Fatal("Falha ao realizar o parsing dos command-line arguments.")
	}

	if _, set := os.LookupEnv("FFMPEG_LOG_LEVEL"); !set {
		os.Setenv("FFMPEG_LOG_LEVEL", *logLevel)
	}

	if *version == true {
		fmt.Println("video-converter-0.3.1")
		return
	}

	ext := extractor.NewFfmpegExtractor()
	if err := ext.Extract(*currentDir); err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := ext.Remove(); err != nil {
			log.Fatal(err)
		}
	}()

	// Dependencies and Services
	stdin := bufio.NewReader(os.Stdin)
	md := metadata.NewVideoMetadataExtractor(ext.GetPath())
	gpuDiscover := new(discover.GpuDiscover)
	formatKit := videoformat.NewVideoFormatKitImpl()
	printer := printer.NewDefaultPrinter(os.Stdout, nil)
	progressWriter := utils.NewFfmpegProgressWriter()
	vc := converter.NewVideoConverter(ext.GetPath(), progressWriter)
	videoCutter := vcut.NewVideoCutter(ext.GetPath(), progressWriter)

	// Flows
	baseFlow := flows.NewBaseFlow(formatKit, gpuDiscover, printer, stdin)
	convertionFlow := flows.NewConvertionFlow(stdin, md, baseFlow, vc)
	cutFlow := flows.NewCutFlow(stdin, baseFlow, md, videoCutter, printer)
	cutIntervalFlow := flows.NewCutIntervalFlow(stdin, md, baseFlow, printer, videoCutter)
	operationTypeFlow := flows.NewOperationTypeFlow(stdin, printer, convertionFlow, cutFlow, cutIntervalFlow)

	filename, err := operationTypeFlow.Run()
	fmt.Println()

	if err != nil {
		printer.Error(true, err.Error())
		printer.Warn(true, "\n\nPressione Enter para sair...")
		bufio.NewReader(os.Stdin).ReadLine()
		return
	}

	if filename != "" && *openDir {
		utils.OpenFileInExplorer(filepath.Dir(filename))
	}

	printer.Warn(true, "Operação concluída! Pressione Enter para sair...")
	bufio.NewReader(os.Stdin).ReadLine()
}
