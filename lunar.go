package main

import (
	"fmt"
	"math"
	"net/http"
	"time"
)

var lander *Lunar

var flame *Flame

func main() {
	http.HandleFunc("/event", event)
	go http.ListenAndServe(":8090", nil)

	tmpFlame := NewFlame()
	flame = &tmpFlame
	flame.Active = true

	ground := NewGround()
	ground.Active = true

	tmpLander := NewLunar()
	lander = &tmpLander
	lander.SetPosition(600, 600)
	lander.Active = true

	lander.AccelerationX = 0
	lander.AccelerationY = 0

	var gameObjects []GameObject
	gameObjects = append(gameObjects, lander)
	gameObjects = append(gameObjects, flame)
	gameObjects = append(gameObjects, &ground)

	fmt.Printf("r=2000\n")
	fmt.Printf("e=1\n")

	go func() {
		for true {
			updateState()
			time.Sleep(time.Millisecond * 100)
		}
	}()

	for true {
		draw(gameObjects)
	}
}

type Position struct {
	X float64
	Y float64
}

type Line struct {
	X1 float64
	Y1 float64
	X2 float64
	Y2 float64
}

type GameObject interface {
	GetLines() []Line
	GetPosition() Position
	GetActive() bool
	//GetAngle() float64
	//SetAngle(float64)
}

func draw(objects []GameObject) {
	for _, object := range objects {
		if object.GetActive() {
			position := object.GetPosition()
			command := fmt.Sprintf("s=%v,%v,%v,%v,%v,%v\n", int(object.GetLines()[0].X1+position.X), int(object.GetLines()[0].Y1+position.Y), "0", "0", "0", "0")
			for _, line := range object.GetLines() {
				command = command + fmt.Sprintf("s=%v,%v,%v,%v,%v,%v\n", int(line.X1+position.X), int(line.Y1+position.Y), "4095", "4095", "1", "1")
				command = command + fmt.Sprintf("s=%v,%v,%v,%v,%v,%v\n", int(line.X2+position.X), int(line.Y2+position.Y), "4095", "4095", "1", "1")
			}
			fmt.Printf(command)
		}
	}
}

func event(w http.ResponseWriter, req *http.Request) {
	event, _ := req.URL.Query()["event"]

	key, _ := req.URL.Query()["key"]

	if event[0] == "down" && key[0] == "68" {
		turningLeft = true
	}
	if event[0] == "down" && key[0] == "65" {
		turningRight = true
	}
	if event[0] == "down" && key[0] == "87" {
		thrust = true
	}
	if event[0] == "up" && key[0] == "68" {
		turningLeft = false
	}
	if event[0] == "up" && key[0] == "65" {
		turningRight = false
	}
	if event[0] == "up" && key[0] == "87" {
		thrust = false
	}
}

var turningRight = false

var turningLeft = false

var thrust = false

func updateState() {
	if turningLeft {
		lander.Angle = lander.Angle - 10
	}
	if turningRight {
		lander.Angle = lander.Angle + 10
	}
	if lander.Angle < 0 {
		lander.Angle = lander.Angle + 360
	}
	if lander.Angle > 360 {
		lander.Angle = lander.Angle - 360
	}
	flame.Angle = lander.Angle

	if thrust {
		flame.Active = true

		val := math.Cos(lander.Angle * math.Pi / 180)
		lander.AccelerationY = lander.AccelerationY + (val * THRUST)

		val = math.Sin(lander.Angle * math.Pi / 180)
		lander.AccelerationX = lander.AccelerationX + (val * THRUST * -1)
	} else {
		flame.Active = false
		lander.AccelerationY = lander.AccelerationY + GRAVITY
	}
	landerY := lander.Position.Y + lander.AccelerationY
	landerX := lander.Position.X + lander.AccelerationX

	if landerY < 100 {
		landerY = 100
		if lander.AccelerationY < 0 {
			lander.AccelerationY = 0
		}
		lander.AccelerationX = 0
	}
	if landerY > 2100 {
		landerY = 2100
	}

	lander.SetPosition(landerX, landerY)
}

//var MAX_SPEED float64  = 250

//var DRAG float64  = 0

var THRUST float64 = .75

var GRAVITY float64 = -.5

type Ground struct {
	Position Position
	Lines    []Line
	Center   Position
	Angle    float64
	Active   bool
}

func NewGround() Ground {
	ground := Ground{
		Center: Position{
			X: 0,
			Y: 0,
		},
		Angle:    0,
		Position: Position{},
		Lines: []Line{{
			X1: 0,
			Y1: 93,
			X2: 2000,
			Y2: 93,
		},
		},
	}
	return ground
}

func (ground Ground) GetLines() []Line {
	return ground.Lines
}

func (ground *Ground) GetPosition() Position {
	return ground.Position
}

func (ground *Ground) GetActive() bool {
	return ground.Active
}

func (ground *Ground) SetPosition(x float64, y float64) {
	ground.Position.X = x
	ground.Position.Y = y
}

type Lunar struct {
	Position      Position
	Lines         []Line
	Center        Position
	Angle         float64
	Active        bool
	AccelerationX float64
	AccelerationY float64
}

func NewLunar() Lunar {
	lunar := Lunar{
		Center: Position{
			X: 40,
			Y: 50,
		},
		Angle:    0,
		Position: Position{},
		Lines: []Line{{
			X1: 0,
			Y1: 0,
			X2: 40,
			Y2: 100,
		}, {
			X1: 40,
			Y1: 100,
			X2: 80,
			Y2: 0,
		}, {
			X1: 80,
			Y1: 0,
			X2: 0,
			Y2: 0,
		},
		},
	}
	return lunar
}

func (lunar Lunar) GetLines() []Line {
	var rotatedLines []Line

	for _, line := range lunar.Lines {
		s := math.Sin(lunar.Angle * (math.Pi / 180))
		c := math.Cos(lunar.Angle * (math.Pi / 180))

		line.X1 = line.X1 - lunar.Center.X
		line.Y1 = line.Y1 - lunar.Center.Y
		line.X2 = line.X2 - lunar.Center.X
		line.Y2 = line.Y2 - lunar.Center.Y

		x1new := line.X1*c - line.Y1*s
		y1new := line.X1*s + line.Y1*c
		x2new := line.X2*c - line.Y2*s
		y2new := line.X2*s + line.Y2*c

		line.X1 = x1new + lunar.Center.X
		line.Y1 = y1new + lunar.Center.Y
		line.X2 = x2new + lunar.Center.X
		line.Y2 = y2new + lunar.Center.Y

		rotatedLines = append(rotatedLines, line)
	}
	return rotatedLines
}

func (lunar *Lunar) GetPosition() Position {
	return lunar.Position
}

func (lunar *Lunar) GetActive() bool {
	return lunar.Active
}

func (lunar *Lunar) SetPosition(x float64, y float64) {
	lunar.Position.X = x
	lunar.Position.Y = y
	flame.SetPosition(x, y)
}

type Flame struct {
	Position Position
	Lines    []Line
	Center   Position
	Angle    float64
	Active   bool
}

func NewFlame() Flame {
	flame := Flame{
		Center: Position{
			X: 40,
			Y: 50,
		},
		Angle:    0,
		Position: Position{},
		Lines: []Line{{
			X1: 20,
			Y1: 0,
			X2: 40,
			Y2: -40,
		}, {
			X1: 40,
			Y1: -40,
			X2: 60,
			Y2: 0,
		},
		},
	}
	return flame
}

func (flame Flame) GetLines() []Line {
	var rotatedLines []Line

	for _, line := range flame.Lines {
		s := math.Sin(flame.Angle * (math.Pi / 180))
		c := math.Cos(flame.Angle * (math.Pi / 180))

		line.X1 = line.X1 - flame.Center.X
		line.Y1 = line.Y1 - flame.Center.Y
		line.X2 = line.X2 - flame.Center.X
		line.Y2 = line.Y2 - flame.Center.Y

		x1new := line.X1*c - line.Y1*s
		y1new := line.X1*s + line.Y1*c
		x2new := line.X2*c - line.Y2*s
		y2new := line.X2*s + line.Y2*c

		line.X1 = x1new + flame.Center.X
		line.Y1 = y1new + flame.Center.Y
		line.X2 = x2new + flame.Center.X
		line.Y2 = y2new + flame.Center.Y

		rotatedLines = append(rotatedLines, line)
	}
	return rotatedLines
}

func (flame *Flame) GetPosition() Position {
	return flame.Position
}

func (flame *Flame) GetActive() bool {
	return flame.Active
}

func (flame *Flame) SetPosition(x float64, y float64) {
	flame.Position.X = x
	flame.Position.Y = y
}
