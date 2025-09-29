package ast

import (
	"github.com/aquaticcalf/rangastalam/common"
)

type Project struct {
	Size   common.Size
	FPS    float64
	Tracks []Track
}

type Track interface {
	TrackName() string
	Z() int
}

type Video struct {
	Name  string
	Clips []VideoNode
}

func (v Video) TrackName() string { return v.Name }

type VideoNode struct {
	ID       string
	Path     string
	SrcStart float64
	SrcEnd   float64
	Start    float64
	End      float64
	Pos      common.Vec2
	Size     common.Size
}

type Text struct {
	Name     string
	Contents []TextNode
}

func (t Text) TrackName() string { return t.Name }

type TextNode struct {
	ID      string
	Content string
	Start   float64
	End     float64
	Pos     common.Vec2
	Style   TextStyle
}

type TextStyle struct {
	Font  string
	Size  float64
	Color string
	Align string
}
