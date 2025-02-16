package game

import (
	"image/color"

	"github.com/EngoEngine/ecs"
	"github.com/SolarLune/resolv"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/sardap/walk-good-maybe-hd/components"
	"github.com/sardap/walk-good-maybe-hd/entity"
)

type Resolvable interface {
	ecs.BasicFace
	components.TransformFace
	components.IdentityFace
	components.CollisionFace
}

type ResolvSystem struct {
	ents           map[uint64]Resolvable
	space          *resolv.Space
	overlay        *ebiten.Image
	OverlayEnabled bool
	debugInput     *entity.InputEnt
}

func CreateResolvSystem(space *resolv.Space, debugInput *entity.InputEnt) *ResolvSystem {
	return &ResolvSystem{
		space:      space,
		debugInput: debugInput,
	}
}

func (s *ResolvSystem) Priority() int {
	return int(systemPriorityResolvSystem)
}

func (s *ResolvSystem) New(world *ecs.World) {
	s.ents = make(map[uint64]Resolvable)
	s.overlay = ebiten.NewImage(windowWidth, windowHeight)
}

func (s *ResolvSystem) Update(dt float32) {
	const buffer = 1

	if s.debugInput != nil && s.debugInput.MovementComponent.InputJustReleased(components.InputKindToggleCollsionOverlay) {
		s.OverlayEnabled = !s.OverlayEnabled
	}

	for _, ent := range s.ents {
		colCom := ent.GetCollisionComponent()
		trans := ent.GetTransformComponent()
		colShape := colCom.CollisionShape

		colShape.X = trans.Postion.X
		colShape.Y = trans.Postion.Y
		colShape.W = trans.Size.X
		colShape.H = trans.Size.Y
	}

	for _, ent := range s.ents {
		colCom := ent.GetCollisionComponent()

		colShape := colCom.CollisionShape

		colCom.Collisions = nil

		colShape.X -= buffer
		colShape.Y -= buffer
		colShape.W += buffer * 2
		colShape.H += buffer * 2

		for _, collidingShape := range s.space.Collisions(colShape) {
			colCom.Collisions = append(colCom.Collisions, &components.CollisionEvent{
				Tags: collidingShape.GetTags(),
			})

			// apply damage
			if damage, ok := ent.(components.DamageFace); ok {
				damageCom := damage.GetDamageComponent()
				// That's right I broke the rules
				if other, ok := collidingShape.GetData().(components.LifeFace); ok {
					otherLifeCom := other.GetLifeComponent()
					otherLifeCom.DamageEvents = append(otherLifeCom.DamageEvents, &components.DamageEvent{
						Damage: damageCom.BaseDamage,
					})
				}
			}
		}

		colShape.X += buffer
		colShape.Y += buffer
		colShape.W -= buffer * 2
		colShape.H -= buffer * 2
	}
}

func (s *ResolvSystem) Render(cmds *RenderCmds) {
	s.overlay.Fill(color.RGBA{0, 0, 0, 0})

	if !s.OverlayEnabled {
		return
	}

	for _, ent := range s.ents {
		if !ent.GetCollisionComponent().Active {
			continue
		}

		x1 := ent.GetTransformComponent().Postion.X
		y1 := ent.GetTransformComponent().Postion.Y
		w := ent.GetTransformComponent().Size.X
		h := ent.GetTransformComponent().Size.Y

		clr := color.RGBA{255, 0, 0, 255}
		// Left Top to Right Top
		ebitenutil.DrawLine(s.overlay, x1, y1, x1+w, y1, clr)
		// Right Top to Right Bottom
		ebitenutil.DrawLine(s.overlay, x1+w, y1, x1+w, y1+h, clr)
		// Right Bottom to Left Bottom
		ebitenutil.DrawLine(s.overlay, x1+w, y1+h, x1, y1+h, clr)
		// Left Bottom to Left top
		ebitenutil.DrawLine(s.overlay, x1, y1+h, x1, y1, clr)
	}

	op := &ebiten.DrawImageOptions{}
	*cmds = append(*cmds, &RenderImageCmd{
		Image:   s.overlay,
		Options: op,
		Layer:   ImageLayerDebug,
	})
}

func (s *ResolvSystem) Add(r Resolvable) {
	s.ents[r.GetBasicEntity().ID()] = r

	if r.GetCollisionComponent().CollisionShape != nil {
		if !s.space.Contains(r.GetCollisionComponent().CollisionShape) {
			s.space.Add(r.GetCollisionComponent().CollisionShape)
		}
		return
	}

	trans := r.GetTransformComponent()
	ident := r.GetIdentityComponent()

	rectangle := resolv.NewRectangle(
		trans.Postion.X, trans.Postion.Y,
		trans.Size.X, trans.Size.Y,
	)

	rectangle.AddTags(ident.Tags...)

	// THIS DATA SHOULD ONLY BE USED FOR TESTING
	rectangle.Data = r

	s.space.Add(rectangle)

	r.GetCollisionComponent().CollisionShape = rectangle
}

func (s *ResolvSystem) Remove(e ecs.BasicEntity) {
	if ent, ok := s.ents[e.ID()]; ok {
		s.space.Remove(ent.GetCollisionComponent().CollisionShape)
	}

	delete(s.ents, e.ID())
}

func (s *ResolvSystem) AddByInterface(o ecs.Identifier) {
	s.Add(o.(Resolvable))
}
