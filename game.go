package main

import "fmt"

type Game struct {
	Text *Text

	Duration     float64
	TimerStarted bool
	Timer        chan bool

	Completed bool
}

func (g *Game) Accuracy() int {
	var correct int

	for _, c := range g.Text.Chars {
		if c.Correct {
			correct++
		}
	}

	acc := float64(correct) / float64(len(g.Text.Chars))

	return int(acc * 100)
}

func (g *Game) WPM() int {
	wordsCount := float64(len(g.Text.Chars)) / float64(5)

	return int(wordsCount * 60 / g.Text.Duration)
}

func (g *Game) Draw() {
	g.Text.Draw()

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
