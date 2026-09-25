package discover

import (
	"fmt"
	"strings"
	"video-converter/types"

	"github.com/StackExchange/wmi"
)

var (
	ErrNoDedicatedGPUFound = fmt.Errorf("no dedicated GPU found")
)

type win32_VideoController struct {
	DeviceID   string
	Name       string
	AdapterRAM uint32
	Status     string
	SystemName string
}

type GpuInfo struct {
	win32_VideoController
	Type types.GpuType
}

type GpuDiscover struct{}

func (gd *GpuDiscover) Discover() (GpuInfo, error) {
	var infoList []win32_VideoController
	query := "SELECT DeviceID, Name, AdapterRAM, Status, SystemName FROM Win32_VideoController"

	err := wmi.Query(query, &infoList)
	if err != nil {
		return GpuInfo{}, fmt.Errorf("Falha ao rodar query WMI para identificação de GPU no ambiente: %v", err)
	}

	for _, info := range infoList {
		name := strings.ToLower(info.Name)

		if strings.Contains(name, "nvidia") {
			return GpuInfo{
				win32_VideoController: info,
				Type:                  types.GpuTypeNvidia,
			}, nil
		}

		if strings.Contains(name, "amd") {
			return GpuInfo{
				win32_VideoController: info,
				Type:                  types.GpuTypeAMD,
			}, nil
		}
	}

	return GpuInfo{}, ErrNoDedicatedGPUFound
}
