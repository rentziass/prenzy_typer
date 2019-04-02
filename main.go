package main

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"time"

	termbox "github.com/nsf/termbox-go"
)

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
