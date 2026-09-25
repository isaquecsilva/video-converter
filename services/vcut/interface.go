package vcut

import (
	"context"
	"video-converter/types"
	"video-converter/types/codecs"
	"video-converter/types/videoformat"
)

type CutParams struct {
	InputFile    string
	From         string
	To           string
	Bitrate      int
	TotalSeconds int64
	Format       videoformat.VideoFormat
	Codec        codecs.VideoCodec
	GpuType      types.GpuType
}

type CutInIntervalParams struct {
	InputFile    string
	Interval     int64
	TotalSeconds int64
	Bitrate      int
	Format       videoformat.VideoFormat
	Codec        codecs.VideoCodec
}

type VideoCutterInterface interface {
	Cut(ctx context.Context, params CutParams) (string, error)
	BatchCut(ctx context.Context, params CutInIntervalParams) ([]string, error)
}
