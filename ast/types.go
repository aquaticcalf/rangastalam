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
	Name   string
	Zindex int
	Clips  []VideoNode
}

func (v Video) TrackName() string { return v.Name }
func (v Video) Z() int            { return v.Zindex }

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
	Zindex   int
	Contents []TextNode
}

func (t Text) TrackName() string { return t.Name }
func (t Text) Z() int            { return t.Zindex }

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

type Image struct {
	Name   string
	Zindex int
	Images []ImageNode
}

func (i Image) TrackName() string { return i.Name }
func (i Image) Z() int            { return i.Zindex }

type ImageNode struct {
	ID    string
	Path  string
	Start float64
	End   float64
	Pos   common.Vec2
	Size  common.Size
}

type Audio struct {
	Name   string
	Zindex int
	Audios []AudioNode
}

func (a Audio) TrackName() string { return a.Name }
func (a Audio) Z() int            { return a.Zindex }

type AudioNode struct {
	ID       string
	Path     string
	SrcStart float64
	SrcEnd   float64
	Start    float64
	End      float64
	Volume   float64
	Loop     bool
}
