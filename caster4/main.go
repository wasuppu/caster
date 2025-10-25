package main

type vec2[T float64 | int] struct {
	x, y T
}

func (v vec2[T]) add(o vec2[T]) vec2[T] {
	return vec2[T]{v.x + o.x, v.y + o.y}
}

func (v vec2[T]) sub(o vec2[T]) vec2[T] {
	return vec2[T]{v.x - o.x, v.y - o.y}
}

func (v vec2[T]) muln(n float64) vec2[T] {
	return vec2[T]{T(float64(v.x) * n), T(float64(v.y) * n)}
}

func (v vec2[T]) divn(n float64) vec2[T] {
	return vec2[T]{T(float64(v.x) / n), T(float64(v.y) / n)}
}

type vec2f = vec2[float64]
type vec2i = vec2[int]

type Engine interface {
	Init() error
	Done() bool
	CleanUp()
	Calculate()
	Render()
	HandleKeys()
	SetFrameTime()
}

func chooseEngine() Engine {
	return &Graphic{}
}

func main() {
	engine := chooseEngine()
	engine.Init()
	defer engine.CleanUp()

	for !engine.Done() {
		engine.Calculate()
		engine.SetFrameTime()
		engine.Render()
		engine.HandleKeys()
	}
}
