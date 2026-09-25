package flows

import (
	"bufio"
	"fmt"
	"slices"
	"strings"
	"video-converter/utils"
	"video-converter/utils/printer"
)

var ErrNoMatchinFlow error = fmt.Errorf("Nenhum fluxo encontrado para a opção selecionada")

type OperationTypeFlow struct {
	reader                                *bufio.Reader
	convertFlow, cutFlow, cutIntervalFlow Flow
	printer                               printer.Printer
}

func NewOperationTypeFlow(reader *bufio.Reader, printer printer.Printer, convert, cut, cutInterval Flow) Flow {
	return &OperationTypeFlow{
		reader:          reader,
		convertFlow:     convert,
		cutFlow:         cut,
		cutIntervalFlow: cutInterval,
		printer:         printer,
	}
}

func (op *OperationTypeFlow) Run() (string, error) {
	op.showBanner()
	choice := op.getOption()
	return op.runFlowByOption(choice)
}

func (op *OperationTypeFlow) showBanner() {
	op.printer.Info(true, "Escolha a operação que deseja efetuar dentre as listadas abaixo: ")
	fmt.Println()
	op.printer.Warn(true, "1. Converter vídeo completo")
	op.printer.Warn(true, "2. Cortar vídeo")
	op.printer.Warn(true, "3. Criar múltiplos cortes baseados em intervalo de tempo")
	fmt.Println()
}

func (op *OperationTypeFlow) getOption() string {
	validOptions := []string{"1", "2", "3"}

	for {
		op.printer.Info(false, "Sua opção: ")
		choice := utils.MustGetUserInput(op.reader)

		if slices.Index(validOptions, choice) == -1 {
			op.printer.Error(true, "Opção inválida! Opções válidas são: %s", strings.Join(validOptions, ", "))
			continue
		}

		return choice
	}
}

func (op *OperationTypeFlow) runFlowByOption(option string) (string, error) {
	switch option {
	case "1":
		return op.convertFlow.Run()
	case "2":
		return op.cutFlow.Run()
	case "3":
		return op.cutIntervalFlow.Run()
	}

	return "", ErrNoMatchinFlow
}
