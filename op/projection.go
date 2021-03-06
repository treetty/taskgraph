package op

// This implementaton is useful for
type oneParameter struct {
	value float32
	size  int
}

func (s *oneParameter) Get(index int) float32 {
	return s.value
}

func (s *oneParameter) Set(index int, value float32) {
	panic("can not set value")
}

func (s *oneParameter) Add(index int, value float32) {
	panic("can not add value")
}

func (s *oneParameter) Data() []float32 {
	panic("can not get data pointer")
	return nil
}

// This allow us to generate parameter with same width.
func (s *oneParameter) CloneWithoutCopy() Parameter {
	return &oneParameter{value: s.value, size: s.size}
}

// This allow one to enumerate through all parameters
func (s *oneParameter) IndexIterator() IndexIterator {
	return MakeRangeIndexIterator(s.size)
}

// This creates a parameter that has the same value on all dimensions
func NewAllTheSameParameter(v float32, s int) Parameter {
	return &oneParameter{value: v, size: s}
}

// Project is used to clip the parameter and gradient.
type Projection struct {
	upper_bound, lower_bound Parameter
}

// this creates a Project with specified upper and lower bound.
// NOTE(baigang): to fix scope visibility of `upper_bound` and `lower_bound` inside `Projection`.
func NewProjection(ub, lb Parameter) *Projection {
	return &Projection{
		upper_bound: ub,
		lower_bound: lb,
	}
}

// We assume the base and gradient are in the same dimensions. In another words,
// the IndexIterator will return the same from base and gradient.
func (p *Projection) ClipGradient(base, gradient Parameter) {
	for iter := base.IndexIterator(); iter.Next(); {
		i := iter.Index()
		// We clip gradient to zero if it is out of bound
		if base.Get(i) <= p.lower_bound.Get(i) {
			gradient.Set(i, min(gradient.Get(i), 0))
		}
		if base.Get(i) >= p.upper_bound.Get(i) {
			gradient.Set(i, max(gradient.Get(i), 0))
		}
	}
}

func (p *Projection) ClipPoint(vec Parameter) {
	for iter := vec.IndexIterator(); iter.Next(); {
		i := iter.Index()
		value := max(vec.Get(i), p.lower_bound.Get(i))
		vec.Set(i, min(value, p.upper_bound.Get(i)))
	}
}

func min(x, y float32) float32 {
	if x < y {
		return x
	} else {
		return y
	}
}

func max(x, y float32) float32 {
	if x > y {
		return x
	} else {
		return y
	}
}
