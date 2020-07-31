package main

import (
	"fmt"
	"log"
	"path/filepath"

	"gitlab.com/gomidi/midi/writer"
)

// Note value as an int
type note = uint8

// Duration, currently in MIDI 'ticks'
type duration = uint32

const defaultVelocity = 255
const semiquaversIn16BarsOfFourFour = 4 * 4 * 16

func writeNote(wr *writer.SMF, n note, d duration) {
	writer.NoteOn(wr, n, defaultVelocity)
	wr.SetDelta(d)
	writer.NoteOff(wr, n)

}

func incrNotes(notes []note) (err error) {
	// The slice of notes, `notes` is little endian. So we incr the first one
	// and see if it raised by one. If it did, great. If not, incr the next one
	// and so on
	if len(notes) < 1 {
		return fmt.Errorf("not enough notes in slice")
	}

	for i := range notes {
		notes[i]++
		// If the note has not wrapped around, this is all the incrementing we
		// need to do right now. If not, increment the next note in the slice as
		// well.
		if notes[i] != 0 {
			break
		}
	}

	return nil
}

func musicMaker(wr *writer.SMF) error {
	const d = 120
	wr.SetChannel(1)
	notes := make([]note, 2)

	for i := range notes {
		notes[i] = 60
	}

	for i := 0; i < 500; i++ {

		for _, n := range notes {
			writeNote(wr, n, d)
		}
		wr.SetDelta(d)
		incrNotes(notes)
	}

	// A bit of empty space at the end to make sure wildmidi plays all of the tune
	wr.SetDelta(1200)

	// Apparently nothing can go wrong?
	// See https://github.com/gomidi/midi/blob/master/examples/smf/smf_test.go
	return nil
}

func main() {
	f := filepath.Join(".", "smf-test.mid")

	err := writer.WriteSMF(f, 1, musicMaker)
	if err != nil {
		log.Fatalf("Could not write SMF file %v", f)
	}

}
