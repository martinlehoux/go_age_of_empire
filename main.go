package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"os"
	"runtime/pprof"

	"age_of_empires/engine"
	"age_of_empires/examples"
	"age_of_empires/selection"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/martinlehoux/kagamigo/kcore"
	"golang.org/x/exp/slog"
)

func listExamples() {
	fmt.Println("Available examples:")
	for name, entry := range examples.Registry {
		fmt.Printf("  %-12s %s\n", name, entry.Description)
	}
}

func main() {
	for _, arg := range os.Args {
		if arg == "--profile" {
			f, err := os.Create("cpu.prof")
			kcore.Expect(err, "could not create CPU profile")
			defer f.Close()
			kcore.Expect(pprof.StartCPUProfile(f), "could not start CPU profile")
			defer pprof.StopCPUProfile()
		}
	}
	logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	slog.SetDefault(slog.New(logHandler))
	ebiten.SetWindowSize(engine.ScreenWidth, engine.ScreenHeight)
	ebiten.SetWindowTitle("Age of Empire")
	icon, err := os.Open("icon-crop.jpg")
	kcore.Expect(err, "failed to open icon")
	defer icon.Close()
	iconImg, _, err := image.Decode(icon)
	kcore.Expect(err, "failed to decode icon")
	ebiten.SetWindowIcon([]image.Image{iconImg})
	fontData, err := os.ReadFile("fonts/Cinzel-VariableFont_wght.ttf")
	kcore.Expect(err, "failed to read font")
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fontData))
	kcore.Expect(err, "failed to create font source")

	baseUnitBuilder := engine.NewEntityBuilder().
		WithSolid(engine.NewFilledCircleImage(100, color.White)).
		WithSelectable("round", selection.Unit).
		WithMove().WithOrder().WithResourceGatherer(15)

	game := &engine.Game{
		FaceSource:  s,
		UnitBuilder: baseUnitBuilder,
	}

	deps := examples.BuilderDeps{
		FaceSource:  s,
		UnitBuilder: baseUnitBuilder,
	}

	if len(os.Args) >= 2 && os.Args[1] == "example" {
		if len(os.Args) >= 3 {
			name := os.Args[2]
			entry, ok := examples.Registry[name]
			if !ok {
				fmt.Fprintf(os.Stderr, "Unknown example %q\n\n", name)
				listExamples()
				os.Exit(1)
			}
			slog.Info("running example", slog.String("name", name))
			entry.Build(game, deps)
		} else {
			listExamples()
			os.Exit(0)
		}
	} else {
		examples.Registry["default"].Build(game, deps)
	}

	ebiten.SetRunnableOnUnfocused(true)
	kcore.Expect(ebiten.RunGame(game), "failed to run game")
}
