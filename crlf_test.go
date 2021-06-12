package stdio

import (
	"bytes"
	"testing"
)

func testCrlfWriter(t *testing.T, i, v []byte) {
	b := bytes.Buffer{}
	w := CrlfWriter{&b}
	_, err := w.Write(i)
	if err != nil {
		t.Error(err)
	}

	o := b.Bytes()
	ok := bytes.Equal(o, v)
	if !ok {
		t.Errorf("%s != %s", o, v)
	}
}

func TestBeginning(t *testing.T) {
	testCrlfWriter(t, []byte("\nHello World!"), []byte("\r\nHello World!"))
}

func TestMiddle(t *testing.T) {
	testCrlfWriter(t, []byte("Hello\nWorld!"), []byte("Hello\r\nWorld!"))
}

func TestEnding(t *testing.T) {
	testCrlfWriter(t, []byte("Hello World!\n"), []byte("Hello World!\r\n"))
}

func TestNone(t *testing.T) {
	testCrlfWriter(t, []byte("Hello World!"), []byte("Hello World!"))
}

func TestMultiple(t *testing.T) {
	testCrlfWriter(t, []byte("\nHello\nWorld!\n"), []byte("\r\nHello\r\nWorld!\r\n"))
}

func TestConcat(t *testing.T) {
	testCrlfWriter(t, []byte("\nHello\n\n\nWorld!\n"), []byte("\r\nHello\r\n\r\n\r\nWorld!\r\n"))
}
