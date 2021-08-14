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

func (t *TransformComponent) GetTransformComponent() *TransformComponent {
	return t
}

type TransformFace interface {
	GetTransformComponent() *TransformComponent
}
