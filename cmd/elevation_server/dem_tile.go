package main

import (
	"github.com/wladich/elevation_server/pkg/constants"
	"math"
	"reflect"
	"unsafe"
)

const NoValue = -32768

type TileIndex struct {
	x, y int
}

type Tile struct {
	data  []int16
	index TileIndex
}

func bytesToInt16(b []byte) []int16 {
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sliceHeader.Cap /= 2
	sliceHeader.Len /= 2
	return *(*[]int16)(unsafe.Pointer(sliceHeader))
}

func (tile *Tile) getInterpolated(latlon LatLon) float64 {
	x := (latlon.lon*constants.HgtSplitParts - float64(tile.index.x)) * (constants.TileSize - 1)
	y := (latlon.lat*constants.HgtSplitParts - float64(tile.index.y)) * (constants.TileSize - 1)
	indX1 := int(math.Floor(x))
	indY1 := int(math.Floor(y))
	indX2 := indX1 + 1
	indY2 := indY1 + 1
	dx := x - float64(indX1)
	dy := y - float64(indY1)
	v1 := tile.data[indX1+indY1*constants.TileSize]
	v2 := tile.data[indX2+indY1*constants.TileSize]
	v3 := tile.data[indX1+indY2*constants.TileSize]
	v4 := tile.data[indX2+indY2*constants.TileSize]
	if v1 == NoValue || v2 == NoValue || v3 == NoValue || v4 == NoValue {
		return NoValue
	}
	return float64(v1)*(1-dx)*(1-dy) +
		float64(v2)*dx*(1-dy) +
		float64(v3)*(1-dx)*dy +
		float64(v4)*dx*dy
}

func tileIndexFromLatLon(latlon LatLon) TileIndex {
	x := int(math.Floor(latlon.lon * constants.HgtSplitParts))
	y := int(math.Floor(latlon.lat * constants.HgtSplitParts))
	return TileIndex{x, y}
}