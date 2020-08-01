package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"unicode/utf8"

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

	fullyWrappedAroundAllNotes := true

	for i := range notes {
		notes[i]++
		// If the note has not wrapped around, this is all the incrementing we
		// need to do right now. If not, increment the next note in the slice as
		// well.
		if notes[i] != 0 {
			fullyWrappedAroundAllNotes = false
			break
		}
	}

	if fullyWrappedAroundAllNotes {
		return fmt.Errorf("fully wrapped around all notes")
	}

	return nil
}

// Generates music maker functions
func musicMakerGenerator(notes []note) func(*writer.SMF) error {

	return func(wr *writer.SMF) error {
		const d = 120
		wr.SetChannel(1)

		for _, n := range notes {
			writeNote(wr, n, d)
		}

		// A bit of empty space at the end to make sure wildmidi plays all of the tune
		wr.SetDelta(1200)

		// Apparently nothing can go wrong?
		// See https://github.com/gomidi/midi/blob/master/examples/smf/smf_test.go
		return nil
	}
}

// [97, 122]
func incrString(s string) (retString string, err error) {
	// First, we define a couple of bookends for the characters allowed in
	// filenames. This is just a-z.
	// TODO support A-Z as well
	const minASCII = 97
	const maxASCII = 122
	// Little endian again
	willIncrNextRune := false
	sLength := utf8.RuneCountInString(s)
	for i, r := range s {
		if r < minASCII {
			return retString, fmt.Errorf("Out of range rune found: %q", r)
		}
		r++
		if r > maxASCII {
			r = minASCII
			willIncrNextRune = true
		} else {
			willIncrNextRune = false
		}

		// If this is the last character in the string
		if sLength == i+1 {
			s = s[:i] + string(r)
		} else {
			s = s[:i] + string(r) + s[i+1:]
		}

		if !willIncrNextRune {
			break
		}
	}

	// If, after all that processing, the string is entirely composed of
	// `minASCII`, add another `minASCII` on the end as the filename has rolled over
	allRunesEqualToMinASCII := true
	for _, v := range s {
		if v != minASCII {
			allRunesEqualToMinASCII = false
			break
		}
	}
	if allRunesEqualToMinASCII {
		s = s + string(minASCII)
	}

	return s, nil
}

func main() {
	var err error

	dir := "music/all_semiquavers/"
	err = os.MkdirAll(dir, os.ModeDir|os.FileMode(0777))
	if err != nil {
		log.Fatalf("Something went real wrong with creating a dir: %q", err)
	}

	notes := make([]note, semiquaversIn16BarsOfFourFour)
	fname := "a"

	log.Printf("Checking starting point from args")
	args := os.Args[1:]

	for i, s := range args {
		note, err := strconv.ParseUint(s, 10, 8)
		if err != nil {
			log.Fatalf("Problem interpreting argument: %q", err)
		}
		notes[i] = uint8(note)
	}

	for notes[len(notes)-1] != math.MaxUint8 {
		musicMaker := musicMakerGenerator(notes)
		f := filepath.Join(dir, fname+".mid")

		// If the file does not exist, create it
		if _, statErr := os.Stat(f); os.IsNotExist(statErr) {
			err = writer.WriteSMF(f, 1, musicMaker)
			if err != nil {
				log.Fatalf("Could not write SMF file %v, err: %q", f, err)
			}

			if notes[0] == math.MaxUint8 && notes[1] == math.MaxUint8 {
				log.Printf("Processed: %v", notes)
			}
		} else {
			if notes[0] == math.MaxUint8 && notes[1] == math.MaxUint8 {
				log.Printf("Skipped: %v", notes)
			}

		}

		// Increment notes and filename
		err = incrNotes(notes)
		if err != nil {
			log.Fatalf("Problem incrementing notes: %q", err)
		}
		fname, err = incrString(fname)
		if err != nil {
			log.Fatalf("Problem iterating filename: %q", err)
		}
	}

}
