package entity

type Position struct {
	Position Vector3f
	Rotation Vector3f
}

type Vector3f struct {
	X float32
	Y float32
	Z float32
}

type Vector2f struct {
	X float32
	Y float32
}

type Vector2 struct {
	X int
	Y int
}
