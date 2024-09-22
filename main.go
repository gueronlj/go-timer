package main

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"time"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

func main() {
	mainImpl()
}

func mainImpl() {
	//run termTimer to check out Version 1
	//go termTimer()

	go func() {
		w := new(app.Window)
		w.Option(app.Title("Timer"))
		w.Option(app.Size(unit.Dp(280), unit.Dp(240))) // Increased height to accommodate larger button

		err := run(w)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func run(window *app.Window) error {
	theme := material.NewTheme()
	//defining Operations form user interfacfe
	var ops op.Ops
	//buttons are actually widgets that are clickable
	var startButton widget.Clickable
	var runClock bool
	buttonText := "Start"

	seconds, minutes := 0, 0
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	go func() {
		for minutes <= 25 {
			<-ticker.C
			if runClock {
				seconds++
				if seconds == 60 {
					seconds = 0
					minutes++
				}
			}
			window.Invalidate()
		}
	}()

	// listen for events in the window.
	for {
		//Determine type of event
		switch e := window.Event().(type) {
		//DestroyEvent is when user closes the window
		case app.DestroyEvent:
			os.Exit(0)
			return e.Err
		// this is sent when the application should re-render.
		case app.FrameEvent:
			// This graphics context is used for managing the rendering state.
			gtx := app.NewContext(&ops, e)
			// Update the clock text with the current time
			time := fmt.Sprintf("%02d:%02d", minutes, seconds)

			for startButton.Clicked(gtx) {
				runClock = !runClock
			}

			if runClock {
				buttonText = "Pause"
			} else {
				buttonText = "Start"
			}

			layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.UniformInset(unit.Dp(20)).Layout(gtx, material.H1(theme, time).Layout)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Min.X = gtx.Constraints.Max.X
						button := material.Button(theme, &startButton, buttonText)
						button.Background = color.NRGBA{R: 60, G: 179, B: 113, A: 255}
						button.CornerRadius = unit.Dp(8)
						// button.Background = theme.ContrastBg
						button.TextSize = unit.Sp(20)
						return layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							gtx.Constraints.Min.Y = gtx.Dp(64) // Double the default height (32dp)
							return button.Layout(gtx)
						})
					}),
				)
			})
			// Pass the drawing operations to the GPU.
			e.Frame(gtx.Ops)
		}
	}
}

// This function is not used. It is here for history. This is how it started.
func termTimer() {
	miliseconds := 0
	seconds := 0
	minutes := 0
	for minutes <= 25 {
		//this weird magic clears the terminal as to keep printing in place
		fmt.Print("\r\033[K")
		fmt.Printf("%02d:%02d:%02d", minutes, seconds, miliseconds)
		miliseconds++
		if miliseconds == 100 {
			miliseconds = 0
			seconds++
			if seconds == 60 {
				seconds = 0
				minutes++
			}
		}
		time.Sleep(1 * time.Second / 1000)
	}
}
