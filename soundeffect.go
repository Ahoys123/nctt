package main

import (
	"log"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

var inited bool = false

type SoundEffect interface {
	Play()
}

type BeepSfx struct {
	buffer *beep.Buffer
}

func NewBeepSfx(filename string) (sfx *BeepSfx) {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	streamer, format, err := wav.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	if !inited {
		speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
		inited = true
	}

	sfx = &BeepSfx{beep.NewBuffer(format)}
	sfx.buffer.Append(streamer)
	streamer.Close()

	return sfx
}

func (sfx *BeepSfx) Play() {
	sound := sfx.buffer.Streamer(0, sfx.buffer.Len())
	speaker.Play(sound)
}
