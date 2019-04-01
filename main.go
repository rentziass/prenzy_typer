package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"time"

	termbox "github.com/nsf/termbox-go"
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

type Quote struct {
	Text string `json:"text"`
}

func main() {
	data, err := ioutil.ReadFile("./quotes.json")
	if err != nil {
		panic(err)
	}
	quotes := []*Quote{}
	err = json.Unmarshal(data, &quotes)
	if err != nil {
		panic(err)
	}

	err = termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc)

	rand.Seed(time.Now().Unix())
	t := NewText(quotes[rand.Intn(len(quotes))].Text)
	t.Draw()

mainloop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				break mainloop
			case termbox.KeyBackspace:
				t.Delete()
			case termbox.KeyBackspace2:
				t.Delete()
			case termbox.KeySpace:
				t.InsertRune(' ')
			default:
				if ev.Ch != 0 {
					t.InsertRune(ev.Ch)
				}
			}
		case termbox.EventError:
			panic(ev.Err)
		}
		t.Draw()
	}
}

//var targetText = `asdf`
var targetText = `Talk to me softly There's something in your eyes Don't hang your head in sorrow`
