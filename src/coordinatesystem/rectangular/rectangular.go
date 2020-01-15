package rectangular

import "sync"

type Rectangular struct {
	positionMap sync.Map // 存储该坐标上的数据
}

func New() *Rectangular {
	var rec Rectangular

	return &rec
}

func (r *Rectangular) AngleBetweenXY() uint64 {
	return 90
}
func (r *Rectangular) AngleBetweenXZ() uint64 {
	return 90
}
func (r *Rectangular) AngleBetweenZY() uint64 {
	return 90
}

func (r *Rectangular) ListObjectsOnOnePosition(x, y, z int) interface{} {
	var pos Position
	pos.x = x
	pos.y = y
	pos.z = z

	value, ok := r.positionMap.Load(pos)
	if ok {
		return value
	} else {
		return nil
	}
}

func (r *Rectangular) SetObject(x, y, z int, obj interface{}) {

}
