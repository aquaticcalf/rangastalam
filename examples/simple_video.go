package main

import (
	"fmt"
	"log"

	"github.com/aquaticcalf/rangastalam/ast"
	"github.com/aquaticcalf/rangastalam/common"
	"github.com/aquaticcalf/rangastalam/ffmpeg"
)

func main() {
	// Create a simple video project
	project := &ast.Project{
		Size: common.Size{
			Width:  "1280",
			Height: "720",
		},
		FPS: 24.0,
		Tracks: []ast.Track{
			// Main video track
			ast.Video{
				Name:   "main",
				Zindex: 1,
				Clips: []ast.VideoNode{
					{
						ID:       "clip1",
						Path:     "input.mp4",
						SrcStart: 5.0,   // Start from 5 seconds in source
						SrcEnd:   15.0,  // End at 15 seconds in source
						Start:    0.0,   // Place at beginning of timeline
						End:      10.0,  // 10 second duration in output
						Pos:      common.Vec2{X: 0, Y: 0},
						Size:     common.Size{Width: "1280", Height: "720"},
					},
				},
			},
		},
	}

	// Create renderer
	renderer := ffmpeg.NewRenderer()

	// Set render options
	options := &ffmpeg.RenderOptions{
		OutputPath: "simple_output.mp4",
		Quality:    "high",
		Format:     "mp4",
		Verbose:    true,
		DryRun:     true, // Set to false to actually render
	}

	// Show the command that would be generated
	cmdString, err := renderer.GetCommandString(project, options)
	if err != nil {
		log.Fatalf("Failed to generate command: %v", err)
	}

	fmt.Println("Generated command:")
	fmt.Println(cmdString)

	// Render the video
	if err := renderer.Render(project, options); err != nil {
		log.Fatalf("Render failed: %v", err)
	}

	fmt.Println("Done!")
}