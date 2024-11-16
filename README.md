# js2go

Encode Go structs to `js.Value`s and viceversa based on tags.

[![Go Report Card](https://goreportcard.com/badge/github.com/adrianosela/js2go)](https://goreportcard.com/report/github.com/adrianosela/js2go)
[![GitHub issues](https://img.shields.io/github/issues/adrianosela/js2go.svg)](https://github.com/adrianosela/js2go/issues)
[![Documentation](https://godoc.org/github.com//adrianosela/js2go?status.svg)](https://godoc.org/github.com/adrianosela/js2go)
[![license](https://img.shields.io/github/license/adrianosela/js2go.svg)](https://github.com/adrianosela/js2go/blob/master/LICENSE)

## Usage

Assume we want to encode/decode a struct defined as follows:

```
type config struct {
	StringField  string   `js:"stringField"`
	IntegerField int      `js:"integerField"`
	StringsArray []string `js:"stringsArray"`
	NestedObject struct {
		InnerStringField  string   `js:"innerStringField"`
		InnerIntegerField int      `js:"innerIntegerField"`
		InnerStringsArray []string `js:"innerStringsArray"`
	} `js:"nestedObject"`
}
```

### Decode a `js.Value` onto a Go struct

```
var c config
err := js2go.Decode(val, &c)
if err != nil {
    // handle error
}
```

### Encode a Go struct as a `js.Value`:

```
jsValue, err := js2go.Encode(&config{
    StringField:  "hello world",
    IntegerField: 1234,
    StringsArray: []string{"hello", "world"},
})
if err != nil {
    // handle error
}
```
