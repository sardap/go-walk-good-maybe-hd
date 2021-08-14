package components

// AUTO GENERATED CODE DO NOT EDIT REFER TO gen/codegen

func (a *AnimeComponent) GetAnimeComponent() *AnimeComponent {
	return a
}

type AnimeFace interface {
	GetAnimeComponent() *AnimeComponent
}

func (i *ImageComponent) GetImageComponent() *ImageComponent {
	return i
}

type ImageFace interface {
	GetImageComponent() *ImageComponent
}

func (m *MovementComponent) GetMovementComponent() *MovementComponent {
	return m
}

type MovementFace interface {
	GetMovementComponent() *MovementComponent
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
