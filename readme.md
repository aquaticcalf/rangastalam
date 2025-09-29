## rangastalam

rangastalam is an experimental command line video editor being written in go

it uses a simple [DSL](https://wikipedia.org/wiki/Domain-specific_language) to describe timelines of video, audio, text, and images,  
and then translates that into [ffmpeg](https://ffmpeg.org) commands  

the long term goal is to make it easy to create video essays in a scriptable way :
- cut and arrange multiple clips
- add audio tracks and voiceovers
- overlay text, images, or additional videos
- apply transforms ( crop, scale, rotate, move ) and animations
- render everything reproducibly from a single script

right now this project is in early development, expect things to change a lot

## FFmpeg Integration

rangastalam now includes a complete FFmpeg integration built using **only Go's standard library**. No external Go dependencies are required beyond the standard library.

### Features

- **Command Generation**: Converts AST structures to proper FFmpeg commands
- **Execution**: Runs FFmpeg with proper error handling and output capture  
- **Multiple Track Types**: Support for video, audio, image, and text tracks
- **Layer Composition**: Automatic overlay generation based on Z-index
- **Quality Presets**: High, medium, and low quality encoding options
- **Format Support**: MP4, AVI, MOV output formats
- **Dry Run Mode**: Test command generation without executing
- **Verbose Output**: Detailed logging for debugging

### Quick Start

```go
package main

import (
    "github.com/aquaticcalf/rangastalam/ast"
    "github.com/aquaticcalf/rangastalam/common"
    "github.com/aquaticcalf/rangastalam/ffmpeg"
)

func main() {
    // Create a project
    project := &ast.Project{
        Size: common.Size{Width: "1920", Height: "1080"},
        FPS: 30.0,
        Tracks: []ast.Track{
            ast.Video{
                Name: "main",
                Zindex: 1,
                Clips: []ast.VideoNode{{
                    Path: "input.mp4",
                    Start: 0.0,
                    End: 10.0,
                }},
            },
        },
    }
    
    // Render it
    renderer := ffmpeg.NewRenderer()
    options := ffmpeg.DefaultRenderOptions()
    renderer.Render(project, options)
}
```

### Requirements

- Go 1.25+ (uses only standard library)
- FFmpeg installed and available in PATH

### Documentation

See [FFMPEG_INTEGRATION.md](FFMPEG_INTEGRATION.md) for detailed documentation of the integration architecture and usage examples.

### Examples

Check the `examples/` directory for more complete usage examples.

