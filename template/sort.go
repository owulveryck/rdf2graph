package template

type bySubject []Current

func (a bySubject) Len() int      { return len(a) }
func (a bySubject) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a bySubject) Less(i, j int) bool {
	return a[i].Node.Subject.RawValue() < a[j].Node.Subject.RawValue()
}
