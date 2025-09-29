package ffmpeg

import (
	"strings"
	"testing"
)

func TestCommand_String(t *testing.T) {
	cmd := &Command{}
	cmd.AddInput("input.mp4")
	cmd.AddFilter("scale", []string{"[0:v]"}, "[scaled]", map[string]string{"": "1920:1080"})
	cmd.AddOutput("output.mp4", "-c:v", "libx264")
	
	result := cmd.String()
	
	// Check that the command contains expected parts
	if !strings.Contains(result, "ffmpeg") {
		t.Error("Command should start with ffmpeg")
	}
	
	if !strings.Contains(result, "-i input.mp4") {
		t.Error("Command should contain input")
	}
	
	if !strings.Contains(result, "output.mp4") {
		t.Error("Command should contain output")
	}
	
	if !strings.Contains(result, "scale=1920:1080") {
		t.Error("Command should contain scale filter")
	}
}

func TestCommand_AddInput(t *testing.T) {
	cmd := &Command{}
	cmd.AddInput("test.mp4", "-ss", "10")
	
	if len(cmd.Inputs) != 1 {
		t.Errorf("Expected 1 input, got %d", len(cmd.Inputs))
	}
	
	input := cmd.Inputs[0]
	if input.Path != "test.mp4" {
		t.Errorf("Expected path 'test.mp4', got '%s'", input.Path)
	}
	
	if len(input.Options) != 2 || input.Options[0] != "-ss" || input.Options[1] != "10" {
		t.Error("Input options not set correctly")
	}
}

func TestCommand_AddFilter(t *testing.T) {
	cmd := &Command{}
	cmd.AddFilter("scale", []string{"[0:v]"}, "[out]", map[string]string{"w": "1920", "h": "1080"})
	
	if len(cmd.Filters) != 1 {
		t.Errorf("Expected 1 filter, got %d", len(cmd.Filters))
	}
	
	filter := cmd.Filters[0]
	if filter.Name != "scale" {
		t.Errorf("Expected filter name 'scale', got '%s'", filter.Name)
	}
	
	if filter.Output != "[out]" {
		t.Errorf("Expected output '[out]', got '%s'", filter.Output)
	}
}