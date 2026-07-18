package pierum


type ILinks[In, Out any] struct {
	Links []ILink[In, Out]
	Clean bool
}

type ILink[In, Out any] struct {
	Seq ISequence[In]
}
