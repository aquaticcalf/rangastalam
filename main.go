package main

import (
	"fmt"
	"log"

	"github.com/aquaticcalf/rangastalam/ast"
	"github.com/aquaticcalf/rangastalam/common"
	"github.com/aquaticcalf/rangastalam/ffmpeg"
)

func main() {
	fmt.Println("rangastalam - Command Line Video Editor")
	
	// Example usage: Create a simple project
	project := createExampleProject()
	
	// Create a renderer
	renderer := ffmpeg.NewRenderer()
	
	// Get the command that would be executed (for debugging)
	options := ffmpeg.DefaultRenderOptions()
	options.Verbose = true
	options.DryRun = true // Set to false to actually run ffmpeg
	
	cmdString, err := renderer.GetCommandString(project, options)
	if err != nil {
		log.Fatalf("Failed to generate command: %v", err)
	}
	
	fmt.Println("\nGenerated FFmpeg command:")
	fmt.Println(cmdString)
	
	// Render the project (currently in dry run mode)
	if err := renderer.Render(project, options); err != nil {
		log.Fatalf("Render failed: %v", err)
	}
	
	fmt.Println("\nRender completed successfully!")
}

// createExampleProject demonstrates how to build a project using the AST
func createExampleProject() *ast.Project {
	project := &ast.Project{
		Size: common.Size{
			Width:  "1920",
			Height: "1080",
		},
		FPS: 30.0,
		Tracks: []ast.Track{
			// Video track with a sample clip
			ast.Video{
				Name:   "main_video",
				Zindex: 1,
				Clips: []ast.VideoNode{
					{
						ID:       "clip1",
						Path:     "input_video.mp4",
						SrcStart: 0.0,
						SrcEnd:   10.0,
						Start:    0.0,
						End:      10.0,
						Pos:      common.Vec2{X: 0, Y: 0},
						Size: common.Size{
							Width:  "1920",
							Height: "1080",
						},
					},
				},
			},
			
			// Audio track
			ast.Audio{
				Name:   "background_music",
				Zindex: 0,
				Audios: []ast.AudioNode{
					{
						ID:       "audio1",
						Path:     "background.mp3",
						SrcStart: 0.0,
						SrcEnd:   10.0,
						Start:    0.0,
						End:      10.0,
						Volume:   0.5,
						Loop:     false,
					},
				},
			},
			
			// Text overlay
			ast.Text{
				Name:   "title_text",
				Zindex: 2,
				Contents: []ast.TextNode{
					{
						ID:      "title1",
						Content: "Welcome to Rangastalam",
						Start:   1.0,
						End:     4.0,
						Pos:     common.Vec2{X: 100, Y: 100},
						Style: ast.TextStyle{
							Font:  "Arial",
							Size:  48.0,
							Color: "white",
							Align: "left",
						},
					},
				},
			},
			
			// Image overlay
			ast.Image{
				Name:   "logo",
				Zindex: 3,
				Images: []ast.ImageNode{
					{
						ID:    "logo1",
						Path:  "logo.png",
						Start: 0.0,
						End:   2.0,
						Pos:   common.Vec2{X: 50, Y: 50},
						Size: common.Size{
							Width:  "200",
							Height: "100",
						},
					},
				},
			},
		},
	}
	
	return project
}
