package errors

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"slices"
	"testing"
)

func TestWithMetadata(t *testing.T) {
	t.Run("Errorf", func(t *testing.T) {
		t1 := Tag("t1")
		t2 := T("t2")
		l1 := L("k1", "v1")
		l2 := LabelOf("k2", "v2")
		e := Errorf("error with args: %d, %5s", 1, "s").WithMeta(t1, l1, t2, nil, l2)
		//goland:noinspection GoTypeAssertionOnErrors
		err := e.(*errorWrapperWithMetadata)

		if err.Error() != "error with args: 1,     s" {
			t.Fatalf("incorrect message, expected 'error with args: 1, s', got: %s", e.Error())
		} else if len(err.tags) != 2 {
			t.Fatalf("incorrect length of tags, expect %d, got %d", 2, len(err.tags))
		} else if !slices.Contains(err.tags, t1) {
			t.Fatalf("expected error tags to contain '%s', but it did not: %+v", t1, err.tags)
		} else if !slices.Contains(err.tags, t2) {
			t.Fatalf("expected error tags to contain '%s', but it did not: %+v", t2, err.tags)
		}

		if len(err.labels) != 2 {
			t.Fatalf("incorrect length of labels, expect %d, got %d", 2, len(err.labels))
		} else if v, ok := err.labels[l1.Key]; !ok || v != l1.Value {
			t.Fatalf("expected labels to contain '%v', but it did not: %+v", l1, err.labels)
		} else if v, ok := err.labels["k2"]; !ok || v != "v2" {
			t.Fatalf("expected labels to contain '%v', but it did not: %+v", l2, err.labels)
		}
	})
	t.Run("New", func(t *testing.T) {
		t1 := Tag("t1")
		t2 := T("t2")
		l1 := L("k1", "v1")
		l2 := LabelOf("k2", "v2")
		e := New("error").WithMeta(t1, l1, t2, nil, l2)
		//goland:noinspection GoTypeAssertionOnErrors
		err := e.(*errorWrapperWithMetadata)

		if err.Error() != "error" {
			t.Fatalf("incorrect message, expected 'error', got: %s", e.Error())
		} else if len(err.tags) != 2 {
			t.Fatalf("incorrect length of tags, expect %d, got %d", 2, len(err.tags))
		} else if !slices.Contains(err.tags, t1) {
			t.Fatalf("expected error tags to contain '%s', but it did not: %+v", t1, err.tags)
		} else if !slices.Contains(err.tags, t2) {
			t.Fatalf("expected error tags to contain '%s', but it did not: %+v", t2, err.tags)
		}

		if len(err.labels) != 2 {
			t.Fatalf("incorrect length of labels, expect %d, got %d", 2, len(err.labels))
		} else if v, ok := err.labels[l1.Key]; !ok || v != l1.Value {
			t.Fatalf("expected labels to contain '%v', but it did not: %+v", l1, err.labels)
		} else if v, ok := err.labels["k2"]; !ok || v != "v2" {
			t.Fatalf("expected labels to contain '%v', but it did not: %+v", l2, err.labels)
		}
	})

	t.Run("Verify data structures", func(t *testing.T) {
		t1 := Tag("t1")
		t2 := T("t2")
		l1 := L("k1", "v1")
		l2 := LabelOf("k2", "v2")
		e := withMetadata(errors.New("src"), t1, nil, t2, l1, l2)
		//goland:noinspection GoTypeAssertionOnErrors
		err := e.(*errorWrapperWithMetadata)

		if len(err.tags) != 2 {
			t.Fatalf("incorrect length of tags, expect %d, got %d", 2, len(err.tags))
		} else if !slices.Contains(err.tags, t1) {
			t.Fatalf("expected error tags to contain '%s', but it did not: %+v", t1, err.tags)
		} else if !slices.Contains(err.tags, t2) {
			t.Fatalf("expected error tags to contain '%s', but it did not: %+v", t2, err.tags)
		}

		if len(err.labels) != 2 {
			t.Fatalf("incorrect length of labels, expect %d, got %d", 2, len(err.labels))
		} else if v, ok := err.labels[l1.Key]; !ok || v != l1.Value {
			t.Fatalf("expected labels to contain '%v', but it did not: %+v", l1, err.labels)
		} else if v, ok := err.labels["k2"]; !ok || v != "v2" {
			t.Fatalf("expected labels to contain '%v', but it did not: %+v", l2, err.labels)
		}
	})

	t.Run("Error()", func(t *testing.T) {
		src := errors.New("src")
		e := withMetadata(src)
		//goland:noinspection GoTypeAssertionOnErrors
		err := e.(*errorWrapperWithMetadata)
		expected := src.Error()
		actual := err.Error()
		if actual != expected {
			t.Fatalf("incorrect Error() result, expect %s, got %s", expected, actual)
		}
	})

	t.Run("Unwrap()", func(t *testing.T) {
		src := errors.New("src")
		e := withMetadata(src)
		//goland:noinspection GoTypeAssertionOnErrors
		err := e.(*errorWrapperWithMetadata)
		expected := src
		actual := err.Unwrap()

		//goland:noinspection GoDirectComparisonOfErrors
		if actual != expected {
			t.Fatalf("incorrect Unwrap() result, expect %v, got %v", expected, actual)
		}
	})

	t.Run("Tags()", func(t *testing.T) {
		t1 := Tag("t1")
		t2 := T("t2")
		e := withMetadata(errors.New("src"), t1, nil, t2)
		//goland:noinspection GoTypeAssertionOnErrors
		err := e.(*errorWrapperWithMetadata)
		expected := []Tag{t1, t2}
		actual := err.Tags()

		if !cmp.Equal(expected, actual) {
			t.Fatalf("incorrect Tags() result: %s", cmp.Diff(expected, actual))
		}
	})

	t.Run("Labels()", func(t *testing.T) {
		l1 := L("k1", "v1")
		l2 := LabelOf("k2", "v2")
		e := withMetadata(errors.New("src"), l1, nil, l2)
		//goland:noinspection GoTypeAssertionOnErrors
		err := e.(*errorWrapperWithMetadata)
		expected := map[string]any{l1.Key: l1.Value, l2.Key: l2.Value}
		actual := err.Labels()

		if !cmp.Equal(expected, actual) {
			t.Fatalf("incorrect Labels() result: %s", cmp.Diff(expected, actual))
		}
	})

	t.Run("HasTag()", func(t *testing.T) {
		t1 := Tag("t1")
		t2 := T("t2")
		l1 := L("k1", "v1")
		l2 := LabelOf("k2", "v2")
		e := withMetadata(errors.New("src"), t1, nil, t2, l1, l2)
		if HasTag(e, "t3") {
			t.Fatalf("incorrect HasTag(e, \"t3\") result, expected false, got true")
		}
		if HasTag(e, "") {
			t.Fatalf("incorrect HasTag(e, \"\") result, expected false, got true")
		}
		if !HasTag(e, "t1") {
			t.Fatalf("incorrect HasTag(e, \"t1\") result, expected true, got false")
		}
		if !HasTag(e, "t2") {
			t.Fatalf("incorrect HasTag(e, \"t2\") result, expected true, got false")
		}
	})

	t.Run("HasLabel()", func(t *testing.T) {
		t1 := Tag("t1")
		t2 := T("t2")
		l1 := L("k1", "v1")
		l2 := LabelOf("k2", "v2")
		e := withMetadata(errors.New("src"), t1, nil, t2, l1, l2)
		if HasLabel(e, "k3") {
			t.Fatalf("incorrect HasLabel(e, \"l3\") result, expected false, got true")
		}
		if HasLabel(e, "") {
			t.Fatalf("incorrect HasLabel(e, \"\") result, expected false, got true")
		}
		if !HasLabel(e, "k1") {
			t.Fatalf("incorrect HasLabel(e, \"k1\") result, expected true, got false")
		}
		if !HasLabel(e, "k2") {
			t.Fatalf("incorrect HasLabel(e, \"k2\") result, expected true, got false")
		}
	})

	t.Run("GetLabel()", func(t *testing.T) {
		t1 := Tag("t1")
		t2 := T("t2")
		l1 := L("k1", "v1")
		l2 := LabelOf("k2", "v2")
		e := withMetadata(errors.New("src"), t1, nil, t2, l1, l2)
		if r := GetLabel(e, "k3"); r != nil {
			t.Fatalf("incorrect GetLabel(e, \"k3\") result, expected nil, got %v", r)
		}
		if r := GetLabel(e, ""); r != nil {
			t.Fatalf("incorrect GetLabel(e, \"\") result, expected nil, got %v", r)
		}
		if r := GetLabel(e, l1.Key); r != l1.Value {
			t.Fatalf("incorrect GetLabel(e, \"%s\") result, expected %v, got %v", l1.Key, l1.Value, r)
		}
		if r := GetLabel(e, l2.Key); r != l2.Value {
			t.Fatalf("incorrect GetLabel(e, \"%s\") result, expected %v, got %v", l2.Key, l2.Value, r)
		}
	})

	t.Run("GetLabels()", func(t *testing.T) {
		t1 := Tag("t1")
		t2 := T("t2")
		l1 := L("k1", "v1")
		l2 := LabelOf("k2", "v2")
		e := withMetadata(errors.New("src"), t1, nil, t2, l1, l2)
		expected := map[string]any{l1.Key: l1.Value, l2.Key: l2.Value}
		if r := GetLabels(e); !cmp.Equal(expected, r) {
			t.Fatalf("incorrect GetLabels(e) result: %s", cmp.Diff(expected, r))
		}
	})
}
