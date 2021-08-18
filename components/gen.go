package components

// AUTO GENERATED CODE DO NOT EDIT REFER TO gen/codegen

func (a *AnimeComponent) GetAnimeComponent() *AnimeComponent {
	return a
}

type AnimeFace interface {
	GetAnimeComponent() *AnimeComponent
}

func (c *CollisionComponent) GetCollisionComponent() *CollisionComponent {
	return c
}

type CollisionFace interface {
	GetCollisionComponent() *CollisionComponent
}

func (i *ImageComponent) GetImageComponent() *ImageComponent {
	return i
}

type ImageFace interface {
	GetImageComponent() *ImageComponent
}

func (m *MainGamePlayerComponent) GetMainGamePlayerComponent() *MainGamePlayerComponent {
	return m
}

type MainGamePlayerFace interface {
	GetMainGamePlayerComponent() *MainGamePlayerComponent
}

func (m *MovementComponent) GetMovementComponent() *MovementComponent {
	return m
}

type MovementFace interface {
	GetMovementComponent() *MovementComponent
}

func (s *ScrollableComponent) GetScrollableComponent() *ScrollableComponent {
	return s
}

type ScrollableFace interface {
	GetScrollableComponent() *ScrollableComponent
}

func (s *SoundComponent) GetSoundComponent() *SoundComponent {
	return s
}

type SoundFace interface {
	GetSoundComponent() *SoundComponent
}

func (t *TextComponent) GetTextComponent() *TextComponent {
	return t
}

type TextFace interface {
	GetTextComponent() *TextComponent
}

func (t *TransformComponent) GetTransformComponent() *TransformComponent {
	return t
}

type TransformFace interface {
	GetTransformComponent() *TransformComponent
}

func (v *VelocityComponent) GetVelocityComponent() *VelocityComponent {
	return v
}

type VelocityFace interface {
	GetVelocityComponent() *VelocityComponent
}
