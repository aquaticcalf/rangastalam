package ast

type Project struct {
	Width  int
	Height int
	FPS    float64
	Tracks []Track
}

type Track interface {
	Name() string
	Z() int
}
