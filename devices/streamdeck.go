package devices

type StreamDeckImageFormats string

const (
	Jpeg StreamDeckImageFormats = "jpeg"
)

type StreamDeckConfig struct {
	KeyCount                 int
	KeyCols                  int
	KeyRows                  int
	KeyPixelWidth            int
	KeyPixelHeight           int
	KeyImageFormat           StreamDeckImageFormats
	ImageReportLength        int
	ImageReportHeaderLength  int
	ImageReportPayloadLength int
	FeatureLength            int
}

type StreamDeck interface {
	Init() error
	Reset() error
	GetKeyCount() int
	GetKeyImageFormat() StreamDeckImageFormats
	GetStateChannel() <-chan []byte
	SetBrightness(percent int)
	SetKeyImage(key byte, image []byte)
	Close()
}
