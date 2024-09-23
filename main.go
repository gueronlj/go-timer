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
		w.Option(app.Size(unit.Dp(280), unit.Dp(240)))

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
	var restartButton widget.Clickable

	var runClock bool
	buttonText := "Start"

	seconds, minutes := 0, 0
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	go func() {
		for minutes < 25 {
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
				if runClock {
					buttonText = "Pause"
				} else {
					buttonText = "Start"
				}
			}

			for restartButton.Clicked(gtx) {
				seconds, minutes = 0, 0
				runClock = false
				buttonText = "Start"
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
						return layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
								layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
									if !runClock && (minutes > 0 || seconds > 0) {
										// Show both buttons
										return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
											layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
												return layoutButton(gtx, theme, &startButton, buttonText, color.NRGBA{R: 60, G: 179, B: 113, A: 255})
											}),
											layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
												return layoutButton(gtx, theme, &restartButton, "Restart", color.NRGBA{R: 220, G: 20, B: 60, A: 255})
											}),
										)
									} else {
										// Show only start/pause button
										return layoutButton(gtx, theme, &startButton, buttonText, color.NRGBA{R: 60, G: 179, B: 113, A: 255})
									}
								}),
							)
						})
					}),
				)
			})

			e.Frame(gtx.Ops)
		}
	}
}

func layoutButton(gtx layout.Context, theme *material.Theme, button *widget.Clickable, text string, bg color.NRGBA) layout.Dimensions {
	btn := material.Button(theme, button, text)
	btn.Background = bg
	btn.CornerRadius = unit.Dp(8)
	btn.TextSize = unit.Sp(20)
	return layout.UniformInset(unit.Dp(4)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Min.Y = gtx.Dp(48)
		return btn.Layout(gtx)
	})
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
