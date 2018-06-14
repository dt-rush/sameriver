package main

func (w *World) UpdateEntityVel() {

	var sigma = func(x float64) float64 {
		return 1 / ((x*x)/(16*16) + 1)
	}

	if w.e != nil && w.e.moveTarget != nil {

		toward := VecFromPoints(w.e.pos, *w.e.moveTarget)
		force := toward.Unit().Scale(MOVESPEED)

		for _, o := range w.obstacles {
			ov := VecFromPoints(w.e.pos, o)
			pl := ov.PerpendicularUnit().Scale(-10)
			pr := ov.PerpendicularUnit().Scale(10)
			_, _, d := Distance(w.e.pos, o)
			if force.Project(pr) > 0 {
				force = force.Add(pr.Scale(sigma(d)))
			} else {
				force = force.Add(pl.Scale(sigma(d)))
			}
		}
		w.e.vel = force
	}
}

func (w *World) MoveEntity() {
	if w.e != nil {
		w.e.pos.X += w.e.vel.X
		w.e.pos.Y += w.e.vel.Y
		if w.e.moveTarget != nil {
			_, _, d := Distance(w.e.pos, *w.e.moveTarget)
			if d < MOVESPEED*2 {
				w.e.pos = *w.e.moveTarget
				w.e.vel = Vec2D{0, 0}
				w.e.moveTarget = nil
				return
			}
		}
	}
}
