package worker

func createSet(values []string) map[string]struct{} {
	set := make(map[string]struct{}, len(values))
	for _, v := range values {
		set[v] = struct{}{}
	}
	return set
}

type Acceptable interface {
	Accept(value string) bool
}

type AcceptableValues struct {
	values map[string]struct{}
}

func NewAcceptableValues(values []string) *AcceptableValues {
	set := createSet(values)
	return &AcceptableValues{values: set}
}

func (a *AcceptableValues) Accept(value string) bool {
	_, ok := a.values[value]
	return ok
}

type AcceptableAll struct {
}

func NewAcceptableAll() *AcceptableAll {
	return &AcceptableAll{}
}

func (a *AcceptableAll) Accept(value string) bool {
	return true
}
