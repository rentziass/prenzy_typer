package main

import (
	"fmt"
	"time"
)

type Text struct {
	Chars    []*Char
	Position int

	Duration     float64
	TimerStarted bool
	Timer        chan bool

	Completed bool
	Accuracy  int
	WPM       int
}

func NewText(s string) *Text {
	t := &Text{}
	t.Timer = make(chan bool)

	for _, r := range s {
		t.Chars = append(t.Chars, &Char{Rune: r})
	}

	return t
}

func (t *Text) Draw() {
	var str string
	for i, c := range t.Chars {
		if i < t.Position {
			str = str + c.Format()
			continue
		}

		if i == t.Position {
			str = str + "\033[7m" + string(c.Rune) + "\033[0m"
			continue
		}

		str = str + string(c.Rune)
	}

	//fmt.Printf("\033[0;0H")
	fmt.Printf("\r%s", str)
	if t.Completed {
		fmt.Println()
		fmt.Println()
		fmt.Println("Well done")
		fmt.Printf("Completed in %f seconds.\n", t.Duration)

		t.CalcAccuracy()
		fmt.Printf("Accuracy: %d%% \n", t.Accuracy)

		t.CalcWPM()
		fmt.Printf("WPM: %d \n", t.WPM)
	}
}

func (t *Text) CalcAccuracy() {
	var correct int

	for _, c := range t.Chars {
		if c.Correct {
			correct++
		}
	}

	acc := float64(correct) / float64(len(t.Chars))

	t.Accuracy = int(acc * 100)
}

func (t *Text) CalcWPM() {
	wordsCount := float64(len(t.Chars)) / float64(5)

	t.WPM = int(wordsCount * 60 / t.Duration)
}

func (t *Text) Delete() {
	if t.Position > 0 {
		t.Position--
	}
}

func (t *Text) InsertRune(r rune) {
	if !t.TimerStarted {
		go t.StartTimer(t.Timer)
		t.TimerStarted = true
	}

	if t.Position > len(t.Chars)-1 {
		return
	}

	t.Chars[t.Position].Correct = r == t.Chars[t.Position].Rune
	t.Position++

	if t.Position > len(t.Chars)-1 {
		close(t.Timer)
		t.Completed = true
	}
}

type Char struct {
	Rune    rune
	Correct bool
}

func (c *Char) Format() string {
	format := "\033[32m"
	if !c.Correct {
		format = "\033[31m"
	}

	return format + string(c.Rune) + "\033[0m"
}

func (t *Text) StartTimer(timer chan bool) {
	now := time.Now()

	<-timer

	t.Duration = time.Since(now).Seconds()
}
