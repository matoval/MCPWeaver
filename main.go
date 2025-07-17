package main

import (
	"context"
	"embed"

	"MCPWeaver/internal/app"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	application := app.NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:            "MCPWeaver",
		Width:            1200,
		Height:           800,
		MinWidth:         800,
		MinHeight:        600,
		MaxWidth:         0,
		MaxHeight:        0,
		DisableResize:    false,
		Fullscreen:       false,
		Frameless:        false,
		StartHidden:      false,
		HideWindowOnClose: false,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        func(ctx context.Context) { application.OnStartup(ctx) },
		OnDomReady:       application.OnDomReady,
		OnBeforeClose:    application.OnBeforeClose,
		OnShutdown:       func(ctx context.Context) { application.OnShutdown(ctx) },
		WindowStartState: options.Normal,
		Bind: []interface{}{
			application,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
