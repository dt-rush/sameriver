package main

type Lake struct {
	Vertices []*MapVertex
}

func (l *Lake) containsPoint2D(p Point2D) bool {
	polyVerts := make(Polygon, len(l.Vertices))
	for i, vertex := range l.Vertices {
		polyVerts[i] = vertex.pos
	}
	return Point2DInPolygon(p, polyVerts)
}
