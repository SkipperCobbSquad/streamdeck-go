package streamdeck

import (
	"streamdeck/devices"

	"github.com/sstallion/go-hid"
)

type StreamDeckProductId uint16

const ElgatoVendorId uint16 = 0x0fd9

const (
	streamDeckMK2Id StreamDeckProductId = 0x80
)

func Init() {
	hid.Init()
}

func GetStreamDecks() []devices.StreamDeck {
	elgatoDevicesInfo := make([]*hid.DeviceInfo, 0, 5)
	hid.Enumerate(ElgatoVendorId, hid.ProductIDAny, func(info *hid.DeviceInfo) error {
		elgatoDevicesInfo = append(elgatoDevicesInfo, info)
		return nil
	})
	streamDecks := make([]devices.StreamDeck, 0, len(elgatoDevicesInfo))
	for _, elgatoInfo := range elgatoDevicesInfo {
		if streamDeck, err := detectElgatoDevice(elgatoInfo); streamDeck != nil && err == nil {
			streamDecks = append(streamDecks, streamDeck)
		}
	}
	return streamDecks

}

func detectElgatoDevice(info *hid.DeviceInfo) (devices.StreamDeck, error) {
	switch StreamDeckProductId(info.ProductID) {
	case streamDeckMK2Id:
		dev, err := hid.OpenPath(info.Path)
		if err != nil {
			return nil, err
		}
		return devices.NewStreamDeckMK2(dev), nil
	default:
		return nil, nil
	}
}

func Close() {
	hid.Exit()
}
