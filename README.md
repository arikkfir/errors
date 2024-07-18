# errors

Drop-in replacement for the Go SDK "errors" package with some extra sugar.

## Drop-in replacement

This package exports the exact same symbols as exported by the Go SDK standard library `errors`
package, with few additions:

- The `New` function returns an instance of `ErrorWithMeta` interface, which extends `error` with methods for adding
  tags and labels (see below)
- The `Errorf` function is added, a proxy to `fmt.Errorf` but also returns an instance of `ErrorwithMeta` but with
  methods for adding tags and labels (see below)
- The `Tag` and `Label` structs, used for adding metadata to errors
- The `LabelOf` function to create a `Label` instance.
- The `HasTag`, `HasLabel`, `GetLabel` and `GetLabels` functions for getting an `error` metadata (if any)

## Adding metadata

This package allows adding tags metadata to errors, so that later decisions can be made based on whether a certain error
has a certain tag or not (e.g. `user-facing`, `abort`, etc).

To add tags metadata to an error, do this:

```go
package myPackage

import "github.com/arikkfir/errors"

func failsWithErrorForTheUser() error {
	return errors.New("an error").WithMeta(errors.Tag("user-facing")) 
}

func failsWithCustomerLabel() error {
	return errors.New("an error").WithMeta(errors.LabelOf("file", "/var/db/bad-file"), errors.LabelOf("customer", "3h18ksk2")) 
}
```
