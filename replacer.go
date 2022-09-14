package tl

import "github.com/go-playground/validator/v10"

type replacerFunc func(raw string, f validator.FieldError) string

type replacer struct {
	tl      string
	replace replacerFunc
	next    *replacer
}

type replacerList struct {
	head *replacer
	tail *replacer
}

func newReplacer(rf replacerFunc) *replacer {
	return &replacer{
		replace: rf,
	}
}

func newReplacerList() *replacerList {
	return &replacerList{}
}

func registerReplacer() {
	rl = newReplacerList()
	for _, r := range replacers {
		rl.Add(r)
	}
}

func (r *replacer) reset() {
	r.tl = ""
}

func (rl *replacerList) Add(r *replacer) {
	if rl.head == nil {
		rl.head = r
		rl.tail = r
	} else {
		rl.tail.next = r
		rl.tail = r
	}
}

func (rl *replacerList) replace(raw string, f validator.FieldError) string {
	var translation string
	list := rl.head

	for list != nil {
		if list.tl == "" {
			list.tl = raw
		}

		if list.next != nil {
			list.next.tl = list.replace(list.tl, f)

			// reset current translation before moving to next replacer
			list.reset()

			list = list.next
		} else {
			list.tl = list.replace(list.tl, f)
			translation = list.tl

			// reset current translation before moving to next replacer
			list.reset()

			list = list.next
		}
	}

	return translation
}
