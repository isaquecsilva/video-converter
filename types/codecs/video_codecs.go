package codecs

const (
	// h.264
	VideoCodecH264      VideoCodec = "h264"
	VideoCodecH264NvEnc VideoCodec = "h264_nvenc"
	VideoCodecH264Amf   VideoCodec = "h264_amf"

	// webm
	VideoCodecVP8      VideoCodec = "libvpx"
	VideoCodecVP9      VideoCodec = "libvpx-vp9"
	VideoCodecAV1      VideoCodec = "libaom-av1"
	VideoCodecAV1Amf   VideoCodec = "av1_amf"
	VideoCodecAV1NvEnc VideoCodec = "av1_nvenc"
)
