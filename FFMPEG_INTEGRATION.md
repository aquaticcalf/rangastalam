# FFmpeg Integration Guide

This document explains how rangastalam integrates with FFmpeg using only Go's standard library.

## Architecture Overview

The FFmpeg integration consists of four main components:

### 1. Command Builder (`ffmpeg/command.go`)
Constructs FFmpeg commands with proper syntax for inputs, filters, and outputs.

**Key Features:**
- Type-safe command construction
- Filter complex generation
- Input/output management

**Example:**
```go
cmd := &ffmpeg.Command{}
cmd.AddInput("video.mp4")
cmd.AddFilter("scale", []string{"[0:v]"}, "[scaled]", map[string]string{"": "1920:1080"})
cmd.AddOutput("output.mp4", "-c:v", "libx264")
```

### 2. Executor (`ffmpeg/executor.go`)
Executes FFmpeg commands using `os/exec` with proper error handling and output capture.

**Key Features:**
- Context-based cancellation
- Real-time output streaming
- Dry-run mode for testing
- Timeout support
- Verbose logging

**Example:**
```go
executor := ffmpeg.NewExecutor()
executor.SetVerbose(true)
executor.SetDryRun(false)
err := executor.Execute(cmd)
```

### 3. Translator (`ffmpeg/translator.go`)
Converts AST structures to FFmpeg commands, handling complex filter chains and track composition.

**Key Features:**
- AST to FFmpeg translation
- Track layering and Z-index support
- Filter complex generation
- Time-based operations

**Translation Process:**
1. Collect all input files
2. Process tracks by Z-index (bottom to top)
3. Generate appropriate filters for each track type
4. Create overlay chains for composition
5. Set output parameters

### 4. Renderer (`ffmpeg/render.go`)
High-level interface that combines all components for easy project rendering.

**Key Features:**
- Simple project rendering
- Quality presets (high, medium, low)
- Format support (mp4, avi, mov)
- Context cancellation support

## Track Type Support

### Video Tracks
- **Source trimming**: `SrcStart`, `SrcEnd` parameters
- **Timeline placement**: `Start`, `End` parameters  
- **Scaling**: `Size.Width`, `Size.Height`
- **Positioning**: `Pos.X`, `Pos.Y` (for overlay)

**Generated Filters:**
- `trim`: For time-based cutting
- `scale`: For resolution changes
- `overlay`: For positioning multiple videos

### Audio Tracks
- **Volume control**: `Volume` parameter
- **Time trimming**: `SrcStart`, `SrcEnd`
- **Looping**: `Loop` parameter

**Generated Filters:**
- `volume`: For audio level adjustment
- `atrim`: For audio trimming
- `amix`: For mixing multiple audio sources

### Image Tracks
- **Duration**: Automatically looped for specified duration
- **Scaling**: `Size.Width`, `Size.Height`
- **Positioning**: `Pos.X`, `Pos.Y`

**Generated Filters:**
- `loop`: To extend image duration
- `scale`: For resizing
- `overlay`: For positioning

### Text Tracks
- **Content**: `Content` string
- **Styling**: `Font`, `Size`, `Color`, `Align`
- **Positioning**: `Pos.X`, `Pos.Y`
- **Timing**: `Start`, `End` with enable conditions

**Generated Filters:**
- `drawtext`: For text rendering with all styling options

## Usage Examples

### Basic Project Setup
```go
project := &ast.Project{
    Size: common.Size{Width: "1920", Height: "1080"},
    FPS: 30.0,
    Tracks: []ast.Track{
        ast.Video{
            Name: "main_video",
            Zindex: 1,
            Clips: []ast.VideoNode{
                {
                    ID: "clip1",
                    Path: "input.mp4",
                    SrcStart: 0.0,
                    SrcEnd: 10.0,
                    Start: 0.0,
                    End: 10.0,
                    Size: common.Size{Width: "1920", Height: "1080"},
                },
            },
        },
    },
}
```

### Rendering with Custom Options
```go
renderer := ffmpeg.NewRenderer()
options := &ffmpeg.RenderOptions{
    OutputPath: "output.mp4",
    Quality: "high",      // high, medium, low
    Format: "mp4",        // mp4, avi, mov
    Verbose: true,
    DryRun: false,
}

err := renderer.Render(project, options)
```

### Getting Command String (for debugging)
```go
cmdString, err := renderer.GetCommandString(project, options)
fmt.Println(cmdString)
// Output: ffmpeg -i input.mp4 -filter_complex '[0:v]scale=1920:1080[v0]' -c:v libx264 -preset slow -crf 18 output.mp4
```

## Command Generation Details

### Filter Complex Chains
The translator automatically generates `filter_complex` strings for operations that require them:

1. **Multiple inputs**: When combining video, audio, images, or text
2. **Overlays**: When layering multiple visual elements
3. **Time-based operations**: When timing doesn't align with simple trimming

### Quality Presets
- **High**: `-preset slow -crf 18` (best quality, slower encoding)
- **Medium**: `-preset medium -crf 23` (balanced)
- **Low**: `-preset ultrafast -crf 28` (fastest encoding, lower quality)

### Output Formats
- **MP4**: H.264 video + AAC audio (default)
- **AVI**: Compatible with older systems
- **MOV**: QuickTime format

## Error Handling

The integration provides comprehensive error handling:

1. **FFmpeg availability check**: Verifies FFmpeg is installed and working
2. **Command validation**: Ensures proper command structure
3. **Execution errors**: Captures and reports FFmpeg errors
4. **Timeout handling**: Prevents hanging on long operations
5. **Context cancellation**: Supports graceful shutdown

## Dependencies

**Standard Library Only:**
- `os/exec`: For running FFmpeg commands
- `context`: For cancellation support  
- `bufio`: For output parsing
- `fmt`, `strings`: For string manipulation
- `time`: For timeout handling

**No External Go Libraries Required** - The integration uses only Go's standard library as requested.

## Limitations and Considerations

1. **FFmpeg Installation**: Requires FFmpeg to be installed and available in PATH
2. **File Paths**: All input files must exist and be accessible
3. **Memory Usage**: Large projects may require significant memory for complex filter chains
4. **Performance**: Complex projects with many tracks may take time to render

## Future Enhancements

Potential improvements that could be added:

1. **Progress Reporting**: Parse FFmpeg output for progress information
2. **Hardware acceleration**: Support for GPU-accelerated encoding
3. **Advanced Filters**: More sophisticated video effects and transitions
4. **Streaming**: Support for streaming inputs/outputs
5. **Validation**: Pre-render validation of inputs and settings