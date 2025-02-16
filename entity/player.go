package entity

import (
	"time"

	"github.com/EngoEngine/ecs"
	"github.com/sardap/walk-good-maybe-hd/assets"
	"github.com/sardap/walk-good-maybe-hd/components"
	"github.com/sardap/walk-good-maybe-hd/math"
)

type Player struct {
	ecs.BasicEntity
	*components.TransformComponent
	*components.AnimeComponent
	*components.CollisionComponent
	*components.DamageComponent
	*components.GravityComponent
	*components.IdentityComponent
	*components.InputComponent
	*components.LifeComponent
	*components.MainGamePlayerComponent
	*components.MovementComponent
	*components.ScrollableComponent
	*components.SoundComponent
	*components.TileImageComponent
	*components.VelocityComponent
}

func CreatePlayer() *Player {
	img, _ := assets.LoadEbitenImage(assets.ImageWhaleAirTileSet)

	tileMap := components.CreateTileMap(1, 1, img, assets.ImageWhaleAirTileSet.FrameWidth)

	result := &Player{
		BasicEntity: ecs.NewBasic(),
		TransformComponent: &components.TransformComponent{
			Size: math.Vector2{
				X: float64(assets.ImageWhaleAirTileSet.FrameWidth),
				Y: float64(assets.ImageWhaleAirTileSet.FrameWidth),
			},
		},
		AnimeComponent: &components.AnimeComponent{
			FrameDuration: 50 * time.Millisecond,
		},
		CollisionComponent: &components.CollisionComponent{
			Active: true,
		},
		DamageComponent: &components.DamageComponent{
			BaseDamage: 1,
		},
		GravityComponent: &components.GravityComponent{},
		IdentityComponent: &components.IdentityComponent{
			Tags: []int{TagPlayer},
		},
		InputComponent: &components.InputComponent{
			InputMode: components.InputModeKeyboard,
			Gamepad:   components.DefaultGamepadInputType(),
			Keyboard:  components.DefaultKeyboardInputType(),
		},
		LifeComponent: &components.LifeComponent{
			HP:                100,
			InvincibilityTime: 1 * time.Second,
		},
		MainGamePlayerComponent: &components.MainGamePlayerComponent{
			Speed:                700,
			JumpPower:            1,
			State:                components.MainGamePlayerStateFlying,
			ShootCooldown:        250 * time.Millisecond,
			AirHorzSpeedModifier: 0.5,
		},
		MovementComponent: components.CreateMovementComponent(),
		ScrollableComponent: &components.ScrollableComponent{
			Modifier: 1,
		},
		SoundComponent: &components.SoundComponent{
			Active: false,
		},
		TileImageComponent: &components.TileImageComponent{
			Active:  true,
			TileMap: tileMap,
		},
		VelocityComponent: &components.VelocityComponent{
			Vel: math.Vector2{},
		},
	}

	return result
}
