package coordinatesystem

import (
	e "common/error"
	"common/log"
)

var logger = log.GetLogger("CoordinateSystem")

//
type CoordinateSystem interface {
	AngleBetweenXY() uint64

	AngleBetweenXZ() uint64

	AngleBetweenZY() uint64

	SetObject(x, y, z int, obj interface{})

	ListObjectsOnOnePosition(x, y, z int) interface{}
}

func NewCoordinateSystem(sysType string) CoordinateSystem {
	switch sysType {
	case "rectangular", "rec", "RECTANGULAR":

		return
	default:
		//todo
		err := e.New()
		logger.Error(err)
		return nil
	}
}
