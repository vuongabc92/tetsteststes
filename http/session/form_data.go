package session

// Form data carries old data of form after submit
// The data will be clear after read
// The purpose when create this is holding data of a form to show again
// when validate form return error.
type FormData struct {
	Data map[string][]string
}

func (f *FormData) AddData(key string, value []string) {
	f.Data[key] = value
}

func (f *FormData) Has(k string) bool {
	_, has := f.Data[k]
	return has
}

func (f FormData) Get(k string) (d []string) {
	if f.Has(k) {
		d = f.Data[k]
	}
	return
}

func NewFormData() *FormData {
	return &FormData{Data: make(map[string][]string)}
}
