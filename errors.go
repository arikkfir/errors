package errors

import (
	"errors"
	"fmt"
	"maps"
	"slices"
)

type Tag string
type T = Tag

type Label struct {
	Key   string
	Value any
}

func L(key string, value any) Label {
	return Label{Key: key, Value: value}
}

func LabelOf(key string, value any) Label {
	return Label{Key: key, Value: value}
}

type ErrorWithMeta interface {
	error
	WithMeta(meta ...any) ErrorWithMeta
}

func New(message string) ErrorWithMeta {
	return withMetadata(errors.New(message))
}

func Errorf(format string, args ...any) ErrorWithMeta {
	return withMetadata(fmt.Errorf(format, args...))
}

func withMetadata(e error, metadata ...any) ErrorWithMeta {
	var tags []Tag
	var labels map[string]any
	for _, v := range metadata {
		if v != nil {
			switch m := v.(type) {
			case Tag:
				tags = append(tags, m)
			case Label:
				if labels == nil {
					labels = make(map[string]any)
				}
				labels[m.Key] = m.Value
			default:
				panic(fmt.Errorf("unknown type for metadata: %T", m))
			}
		}
	}
	return &errorWrapperWithMetadata{e: e, tags: tags, labels: labels}
}

type errorWrapperWithMetadata struct {
	e      error
	tags   []Tag
	labels map[string]any
}

func (e *errorWrapperWithMetadata) WithMeta(meta ...any) ErrorWithMeta {
	for _, v := range meta {
		if v != nil {
			switch m := v.(type) {
			case Tag:
				e.tags = append(e.tags, m)
			case Label:
				if e.labels == nil {
					e.labels = make(map[string]any)
				}
				e.labels[m.Key] = m.Value
			default:
				panic(fmt.Errorf("unknown type for metadata: %T", m))
			}
		}
	}
	return e
}

func (e *errorWrapperWithMetadata) Error() string {
	return e.e.Error()
}

func (e *errorWrapperWithMetadata) Unwrap() error {
	return e.e
}

func (e *errorWrapperWithMetadata) Tags() []Tag {
	return e.tags
}

func (e *errorWrapperWithMetadata) Labels() map[string]any {
	return e.labels
}

func HasTag(e error, tag Tag) bool {
	for e != nil {
		var em *errorWrapperWithMetadata
		if errors.As(e, &em) {
			if slices.Contains(em.tags, tag) {
				return true
			}
		}
		e = errors.Unwrap(e)
	}
	return false
}

func HasLabel(e error, key string) bool {
	for e != nil {
		var em *errorWrapperWithMetadata
		if errors.As(e, &em) {
			if _, ok := em.labels[key]; ok {
				return true
			}
		}
		e = errors.Unwrap(e)
	}
	return false
}

func GetLabel(e error, key string) any {
	for e != nil {
		var em *errorWrapperWithMetadata
		if errors.As(e, &em) {
			if v, ok := em.labels[key]; ok {
				return v
			}
		}
		e = errors.Unwrap(e)
	}
	return nil
}

func GetLabels(e error) map[string]any {
	var labels map[string]any
	for e != nil {
		var em *errorWrapperWithMetadata
		if errors.As(e, &em) {
			if labels == nil {
				labels = make(map[string]any)
			}
			maps.Copy(labels, em.labels)
		}
		e = errors.Unwrap(e)
	}
	return labels
}
