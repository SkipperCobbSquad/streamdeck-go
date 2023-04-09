package devices

import (
	"math"

	"github.com/sstallion/go-hid"
)

const imageReportLength = 1024
const imageReportHeaderLength = 8

var streamDeckMK2Config StreamDeckConfig = StreamDeckConfig{
	KeyCount:                 15,
	KeyCols:                  5,
	KeyRows:                  3,
	KeyPixelWidth:            72,
	KeyPixelHeight:           72,
	KeyImageFormat:           Jpeg,
	ImageReportLength:        imageReportLength,
	ImageReportHeaderLength:  imageReportHeaderLength,
	ImageReportPayloadLength: imageReportLength - imageReportHeaderLength,
	FeatureLength:            32,
}

type streamDeckMK2 struct {
	device       *hid.Device
	config       *StreamDeckConfig
	stateChannel <-chan []byte
}

func NewStreamDeckMK2(device *hid.Device) StreamDeck {
	return &streamDeckMK2{device, &streamDeckMK2Config, nil}
}

func (s *streamDeckMK2) Init() error {
	if err := s.resetKeyStream(); err != nil {
		return err
	}
	s.stateChannel = s.makeStateChannel()
	return nil
}

func (s *streamDeckMK2) resetKeyStream() error {
	payload := make([]byte, s.config.ImageReportLength)
	payload[0] = 0x02
	if _, err := s.device.Write(payload); err != nil {
		return err
	}
	return nil
}

func (s *streamDeckMK2) Reset() error {
	payload := make([]byte, s.config.FeatureLength)
	payload[0] = 0x03
	payload[1] = 0x02
	if _, err := s.device.SendFeatureReport(payload); err != nil {
		return err
	}
	return nil
}

// Getters ============================>
func (s streamDeckMK2) GetKeyCount() int {
	return s.config.KeyCount
}

func (s streamDeckMK2) GetKeyImageFormat() StreamDeckImageFormats {
	return s.config.KeyImageFormat
}

func (s *streamDeckMK2) GetStateChannel() <-chan []byte {
	return s.stateChannel
}

//=====================================>

func (s *streamDeckMK2) SetBrightness(percent int) {
	percent = int(math.Min(math.Max(float64(percent), float64(0)), float64(100)))
	payload := make([]byte, s.config.FeatureLength)
	payload[0] = 0x03
	payload[1] = 0x08
	payload[2] = byte(percent)
	s.device.SendFeatureReport(payload)
}

func (s *streamDeckMK2) SetKeyImage(key byte, image []byte) {
	pageNumber := 0
	bytesRemaining := len(image)
	for bytesRemaining > 0 {
		thisLength := int(math.Min(float64(bytesRemaining), float64(s.config.ImageReportPayloadLength)))
		bytesSent := pageNumber * s.config.ImageReportPayloadLength
		next := 0
		if thisLength == bytesRemaining {
			next = 1
		}
		header := []byte{
			0x02,
			0x07,
			key,
			byte(next),
			byte(thisLength) & 0xFF,
			byte(thisLength >> 8),
			byte(pageNumber) & 0xFF,
			byte(pageNumber >> 8),
		}
		payload := append(header, image[bytesSent:bytesSent+thisLength]...)
		padding := make([]byte, s.config.ImageReportLength-len(payload))
		s.device.Write(append(payload, padding...))
		bytesRemaining -= thisLength
		pageNumber += 1
	}
}

func (s *streamDeckMK2) makeStateChannel() <-chan []byte {
	stateChan := make(chan []byte)
	buff := make([]byte, 4+s.config.KeyCount)
	go func() {
		for {
			if _, err := s.device.Read(buff); err != nil {
				close(stateChan)
				break
			}
			stateChan <- buff[4:]
		}
	}()
	return stateChan
}

func (s *streamDeckMK2) Close() {
	s.Reset()
	s.device.Close()
}
