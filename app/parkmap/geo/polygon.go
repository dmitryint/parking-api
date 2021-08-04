// Also added other functions and some tests related to geo based polygons.

package geo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
)

// A Polygon is carved out of a 2D plane by a set of (possibly disjoint) contours.
// It can thus contain holes, and can be self-intersecting.
type Polygon struct {
	points []*Point
}

// Creates and returns a new pointer to a Polygon
// composed of the passed in points.  Points are
// considered to be in order such that the last point
// forms an edge with the first point.
func NewPolygon(points []*Point) *Polygon {
	return &Polygon{points: points}
}

// Returns the points of the current Polygon.
func (p *Polygon) Points() []*Point {
	return p.points
}

// Appends the passed in contour to the current Polygon.
func (p *Polygon) Add(point *Point) {
	p.points = append(p.points, point)
}

// Returns whether or not the polygon is closed.
// TODO:  This can obviously be improved, but for now,
//        this should be sufficient for detecting if points
//        are contained using the raycast algorithm.
func (p *Polygon) IsClosed() bool {
	if len(p.points) < 3 {
		return false
	}

	return true
}

// Returns whether or not the current Polygon contains the passed in Point.
func (p *Polygon) Contains(point *Point) bool {
	if !p.IsClosed() {
		return false
	}

	start := len(p.points) - 1
	end := 0

	contains := p.intersectsWithRaycast(point, p.points[start], p.points[end])

	for i := 1; i < len(p.points); i++ {
		if p.intersectsWithRaycast(point, p.points[i-1], p.points[i]) {
			contains = !contains
		}
	}

	return contains
}

// Using the raycast algorithm, this returns whether or not the passed in point
// Intersects with the edge drawn by the passed in start and end points.
// Original implementation: http://rosettacode.org/wiki/Ray-casting_algorithm#Go
func (p *Polygon) intersectsWithRaycast(point *Point, start *Point, end *Point) bool {
	// Always ensure that the the first point
	// has a y coordinate that is less than the second point
	if start.lng > end.lng {

		// Switch the points if otherwise.
		start, end = end, start

	}

	// Move the point's y coordinate
	// outside of the bounds of the testing region
	// so we can start drawing a ray
	for point.lng == start.lng || point.lng == end.lng {
		newLng := math.Nextafter(point.lng, math.Inf(1))
		point = NewPoint(point.lat, newLng)
	}

	// If we are outside of the polygon, indicate so.
	if point.lng < start.lng || point.lng > end.lng {
		return false
	}

	if start.lat > end.lat {
		if point.lat > start.lat {
			return false
		}
		if point.lat < end.lat {
			return true
		}

	} else {
		if point.lat > end.lat {
			return false
		}
		if point.lat < start.lat {
			return true
		}
	}

	raySlope := (point.lng - start.lng) / (point.lat - start.lat)
	diagSlope := (end.lng - start.lng) / (end.lat - start.lat)

	return raySlope >= diagSlope
}

func (p *Polygon) GetOuterLines() (error, [][]*Point) {
	if len(p.points) < 2 {
		return fmt.Errorf("no lines"), [][]*Point{}
	}
	var res [][]*Point

	start := p.points[0]
	current := start
	i := 1
	for {
		res = append(res, []*Point{current, p.points[i]})
		if p.points[i] == start || len(p.points) == 2 {
			break
		}
		current = p.points[i]

		if i+1 < len(p.points) {
			i++
		} else {
			i = 0
		}
	}
	return nil, res
}

func (p *Polygon) distanceToLine(point, l1, l2 *Point) float64 {
	k := ((l2.Lat()-l1.Lat())*(point.Lng()-l1.Lng()) - (l2.Lng()-l1.Lng())*(point.Lat()-l1.Lat())) / (math.Pow(l2.Lat()-l1.Lat(), 2) + math.Pow(l2.Lng()-l1.Lng(), 2))
	perpendicular := NewPoint(
		point.Lat()+k*(l2.Lng()-l1.Lng()),
		point.Lng()-k*(l2.Lat()-l1.Lat()),
	)

	lineLength := l1.GreatCircleDistance(l2)

	if perpendicular.GreatCircleDistance(l1) > lineLength || perpendicular.GreatCircleDistance(l2) > lineLength {
		dist1 := point.GreatCircleDistance(l1)
		dist2 := point.GreatCircleDistance(l2)

		if dist1 > dist2 {
			return dist1
		}
		return dist2
	}
	return perpendicular.GreatCircleDistance(point)
	// return point.GreatCircleDistance(l1)
}

func (p *Polygon) MinimumDistance(point *Point) float64 {
	isPoint, lines := p.GetOuterLines()
	if isPoint != nil {
		return point.GreatCircleDistance(p.points[0])
	}

	var distances []float64
	for _, line := range lines {
		distances = append(distances, p.distanceToLine(point, line[0], line[1]))
	}

	min := distances[0]
	for _, v := range distances {
		if v < min {
			min = v
		}
	}
	return min
}

func (p *Polygon) UnmarshalJSON(data []byte) error {
	// TODO throw an error if there is an issue parsing the body.
	dec := json.NewDecoder(bytes.NewReader(data))
	var values []*Point
	err := dec.Decode(&values)

	if err != nil {
		log.Print(err)
		return err
	}

	*p = *NewPolygon(values)

	return nil
}
