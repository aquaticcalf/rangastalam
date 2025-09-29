package ffmpeg

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/aquaticcalf/rangastalam/ast"
)

// Translator converts AST structures to ffmpeg commands
type Translator struct {
	project *ast.Project
}

// NewTranslator creates a new translator for the given project
func NewTranslator(project *ast.Project) *Translator {
	return &Translator{project: project}
}

// Translate converts the project to an ffmpeg command
func (t *Translator) Translate() (*Command, error) {
	if t.project == nil {
		return nil, fmt.Errorf("project cannot be nil")
	}
	
	cmd := &Command{}
	
	// Collect all inputs and create a mapping
	inputMap := make(map[string]int)
	inputIndex := 0
	
	// Sort tracks by Z-index for proper layering
	sortedTracks := make([]ast.Track, len(t.project.Tracks))
	copy(sortedTracks, t.project.Tracks)
	sort.Slice(sortedTracks, func(i, j int) bool {
		return sortedTracks[i].Z() < sortedTracks[j].Z()
	})
	
	var filterOutputs []string
	
	// Process each track
	for trackIndex, track := range sortedTracks {
		switch trackType := track.(type) {
		case ast.Video:
			outputs, err := t.translateVideoTrack(&trackType, cmd, &inputMap, &inputIndex, trackIndex)
			if err != nil {
				return nil, fmt.Errorf("failed to translate video track %s: %w", trackType.Name, err)
			}
			filterOutputs = append(filterOutputs, outputs...)
			
		case ast.Audio:
			outputs, err := t.translateAudioTrack(&trackType, cmd, &inputMap, &inputIndex, trackIndex)
			if err != nil {
				return nil, fmt.Errorf("failed to translate audio track %s: %w", trackType.Name, err)
			}
			filterOutputs = append(filterOutputs, outputs...)
			
		case ast.Image:
			outputs, err := t.translateImageTrack(&trackType, cmd, &inputMap, &inputIndex, trackIndex)
			if err != nil {
				return nil, fmt.Errorf("failed to translate image track %s: %w", trackType.Name, err)
			}
			filterOutputs = append(filterOutputs, outputs...)
			
		case ast.Text:
			outputs, err := t.translateTextTrack(&trackType, cmd, trackIndex)
			if err != nil {
				return nil, fmt.Errorf("failed to translate text track %s: %w", trackType.Name, err)
			}
			filterOutputs = append(filterOutputs, outputs...)
		}
	}
	
	// If we have multiple outputs, we need to overlay them
	if len(filterOutputs) > 1 {
		t.addOverlayFilters(cmd, filterOutputs)
	}
	
	// Add final output
	cmd.AddOutput("output.mp4", "-c:v", "libx264", "-c:a", "aac")
	
	return cmd, nil
}

// translateVideoTrack processes a video track
func (t *Translator) translateVideoTrack(video *ast.Video, cmd *Command, inputMap *map[string]int, inputIndex *int, trackIndex int) ([]string, error) {
	var outputs []string
	
	for clipIndex, clip := range video.Clips {
		// Add input if not already added
		if _, exists := (*inputMap)[clip.Path]; !exists {
			cmd.AddInput(clip.Path)
			(*inputMap)[clip.Path] = *inputIndex
			*inputIndex++
		}
		
		inputLabel := fmt.Sprintf("[%d:v]", (*inputMap)[clip.Path])
		outputLabel := fmt.Sprintf("[v%d_%d]", trackIndex, clipIndex)
		
		// Create filters for this clip
		filters := make(map[string]string)
		
		// Time trimming
		if clip.SrcStart > 0 || clip.SrcEnd > 0 {
			if clip.SrcEnd > 0 {
				filters["ss"] = fmt.Sprintf("%.3f", clip.SrcStart)
				filters["t"] = fmt.Sprintf("%.3f", clip.SrcEnd-clip.SrcStart)
			} else {
				filters["ss"] = fmt.Sprintf("%.3f", clip.SrcStart)
			}
		}
		
		// Scaling
		if clip.Size.Width != "" && clip.Size.Height != "" {
			filters["scale"] = fmt.Sprintf("%s:%s", clip.Size.Width, clip.Size.Height)
		}
		
		// Position (will be handled in overlay)
		
		// Add the filter
		if len(filters) > 0 {
			// For now, just add a simple scale filter as an example
			if scaleValue, hasScale := filters["scale"]; hasScale {
				cmd.AddFilter("scale", []string{inputLabel}, outputLabel, map[string]string{
					"": scaleValue,
				})
			} else {
				// Just pass through
				cmd.AddFilter("copy", []string{inputLabel}, outputLabel, nil)
			}
		}
		
		outputs = append(outputs, outputLabel)
	}
	
	return outputs, nil
}

// translateAudioTrack processes an audio track
func (t *Translator) translateAudioTrack(audio *ast.Audio, cmd *Command, inputMap *map[string]int, inputIndex *int, trackIndex int) ([]string, error) {
	var outputs []string
	
	for clipIndex, clip := range audio.Audios {
		// Add input if not already added
		if _, exists := (*inputMap)[clip.Path]; !exists {
			cmd.AddInput(clip.Path)
			(*inputMap)[clip.Path] = *inputIndex
			*inputIndex++
		}
		
		inputLabel := fmt.Sprintf("[%d:a]", (*inputMap)[clip.Path])
		outputLabel := fmt.Sprintf("[a%d_%d]", trackIndex, clipIndex)
		
		// Create audio filters
		filters := make(map[string]string)
		
		// Volume adjustment
		if clip.Volume != 0 {
			filters["volume"] = fmt.Sprintf("%.2f", clip.Volume)
		}
		
		// Time trimming
		if clip.SrcStart > 0 || clip.SrcEnd > 0 {
			if clip.SrcEnd > 0 {
				filters["ss"] = fmt.Sprintf("%.3f", clip.SrcStart)
				filters["t"] = fmt.Sprintf("%.3f", clip.SrcEnd-clip.SrcStart)
			}
		}
		
		// Add audio filter if needed
		if len(filters) > 0 {
			if volumeValue, hasVolume := filters["volume"]; hasVolume {
				cmd.AddFilter("volume", []string{inputLabel}, outputLabel, map[string]string{
					"": volumeValue,
				})
			}
		}
		
		outputs = append(outputs, outputLabel)
	}
	
	return outputs, nil
}

// translateImageTrack processes an image track
func (t *Translator) translateImageTrack(images *ast.Image, cmd *Command, inputMap *map[string]int, inputIndex *int, trackIndex int) ([]string, error) {
	var outputs []string
	
	for clipIndex, image := range images.Images {
		// Add input if not already added
		if _, exists := (*inputMap)[image.Path]; !exists {
			// For images, we might want to loop them for the duration
			duration := image.End - image.Start
			cmd.AddInput(image.Path, "-loop", "1", "-t", fmt.Sprintf("%.3f", duration))
			(*inputMap)[image.Path] = *inputIndex
			*inputIndex++
		}
		
		inputLabel := fmt.Sprintf("[%d:v]", (*inputMap)[image.Path])
		outputLabel := fmt.Sprintf("[i%d_%d]", trackIndex, clipIndex)
		
		// Scale the image if needed
		if image.Size.Width != "" && image.Size.Height != "" {
			cmd.AddFilter("scale", []string{inputLabel}, outputLabel, map[string]string{
				"": fmt.Sprintf("%s:%s", image.Size.Width, image.Size.Height),
			})
		}
		
		outputs = append(outputs, outputLabel)
	}
	
	return outputs, nil
}

// translateTextTrack processes a text track
func (t *Translator) translateTextTrack(text *ast.Text, cmd *Command, trackIndex int) ([]string, error) {
	var outputs []string
	
	for clipIndex, textNode := range text.Contents {
		outputLabel := fmt.Sprintf("[t%d_%d]", trackIndex, clipIndex)
		
		// Create drawtext filter
		params := map[string]string{
			"text":     textNode.Content,
			"x":        strconv.Itoa(textNode.Pos.X),
			"y":        strconv.Itoa(textNode.Pos.Y),
			"fontsize": fmt.Sprintf("%.0f", textNode.Style.Size),
		}
		
		if textNode.Style.Font != "" {
			params["font"] = textNode.Style.Font
		}
		if textNode.Style.Color != "" {
			params["fontcolor"] = textNode.Style.Color
		}
		
		// Enable/disable based on time
		if textNode.Start > 0 || textNode.End > 0 {
			params["enable"] = fmt.Sprintf("between(t,%.3f,%.3f)", textNode.Start, textNode.End)
		}
		
		cmd.AddFilter("drawtext", []string{"[base]"}, outputLabel, params)
		outputs = append(outputs, outputLabel)
	}
	
	return outputs, nil
}

// addOverlayFilters creates overlay filters to combine multiple video streams
func (t *Translator) addOverlayFilters(cmd *Command, inputs []string) {
	if len(inputs) < 2 {
		return
	}
	
	// Create a base black video if we don't have a video input
	if t.needsBaseVideo() {
		cmd.AddFilter("color", nil, "[base]", map[string]string{
			"c":    "black",
			"size": fmt.Sprintf("%sx%s", t.project.Size.Width, t.project.Size.Height),
			"r":    fmt.Sprintf("%.2f", t.project.FPS),
		})
		
		// Overlay the first input onto the base
		cmd.AddFilter("overlay", []string{"[base]", inputs[0]}, "[tmp0]", nil)
		
		// Overlay subsequent inputs
		for i := 1; i < len(inputs); i++ {
			prevLabel := fmt.Sprintf("[tmp%d]", i-1)
			nextLabel := fmt.Sprintf("[tmp%d]", i)
			if i == len(inputs)-1 {
				nextLabel = "[final]"
			}
			cmd.AddFilter("overlay", []string{prevLabel, inputs[i]}, nextLabel, nil)
		}
	} else {
		// Use the first input as base and overlay others
		for i := 1; i < len(inputs); i++ {
			var baseLabel string
			if i == 1 {
				baseLabel = inputs[0]
			} else {
				baseLabel = fmt.Sprintf("[tmp%d]", i-2)
			}
			
			nextLabel := fmt.Sprintf("[tmp%d]", i-1)
			if i == len(inputs)-1 {
				nextLabel = "[final]"
			}
			
			cmd.AddFilter("overlay", []string{baseLabel, inputs[i]}, nextLabel, nil)
		}
	}
}

// needsBaseVideo checks if we need to create a base video layer
func (t *Translator) needsBaseVideo() bool {
	for _, track := range t.project.Tracks {
		if _, isVideo := track.(ast.Video); isVideo {
			return false
		}
	}
	return true
}