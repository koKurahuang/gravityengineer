package rectangular

// 基本方向， 这里是立体坐标系，所以只有xyz
type baseDirection int

const (
	x baseDirection = iota
	y
	z
)

type Point struct {
	position int //位置坐标
	value    int //值todo沒想好有什么用

	direction baseDirection //这里其实是指坐标轴的朝向，对于基本的点，这个方向就是坐标轴
}

func newPoint(pos, value int, dir baseDirection) *Point {
	return &Point{pos, value, dir}
}
