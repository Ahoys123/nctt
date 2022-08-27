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

type SoundEffect struct {
	buffer *beep.Buffer
}

func NewSoundEffect(filename string) (sfx *SoundEffect) {
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

	sfx = &SoundEffect{beep.NewBuffer(format)}
	sfx.buffer.Append(streamer)
	streamer.Close()

	return sfx
}

func (sfx *SoundEffect) Play() {
	sound := sfx.buffer.Streamer(0, sfx.buffer.Len())
	speaker.Play(sound)
}
