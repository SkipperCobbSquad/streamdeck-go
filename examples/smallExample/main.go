package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"streamdeck"
)

func main() {
	streamdeck.Init()
	defer streamdeck.Close()
	decks := streamdeck.GetStreamDecks()
	streamDeck := decks[0]
	if err := streamDeck.Init(); err != nil {
		log.Fatal(err)
	}
	if err := streamDeck.Reset(); err != nil {
		log.Fatal(err)
	}
	// streamDeck.SetBrightness(50)
	defer streamDeck.Close()
	prevState := make([]byte, 15)
	for state := range streamDeck.GetStateChannel() {
		fmt.Printf("%X\n", state)
		for key := range state {
			if state[key] == 1 {
				prevState[key] = 1
				streamDeck.SetKeyImage(byte(key), smallImage(color.White))

			} else if state[key] == 0 && prevState[key] == 1 {
				prevState[key] = 0
				fmt.Println(key)
				streamDeck.SetKeyImage(byte(key), smallImage(color.Black))
				if key == 14 {
					return
				}
			}
		}
	}
}

func smallImage(color color.Color) []byte {
	img := image.NewRGBA(image.Rect(0, 0, 72, 72))
	b := img.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			img.Set(x, y, color)
		}
	}
	buff := &bytes.Buffer{}
	jpeg.Encode(buff, img, nil)
	return buff.Bytes()
}
