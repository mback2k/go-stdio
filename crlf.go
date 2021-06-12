package stdio

import (
	"bytes"
	"io"
)

type CrlfWriter struct {
	W io.Writer
}

func (w *CrlfWriter) Write(p []byte) (int, error) {
	total := 0
	for _, split := range bytes.SplitAfter(p, []byte("\n")) {
		l := len(split)
		if l > 0 {
			if split[l-1] == '\n' {
				if l < 2 || split[l-2] != '\r' {
					wrt, err := w.W.Write(split[:l-1])
					total += wrt
					if err != nil {
						return total, err
					}
				}
				wrt, err := w.W.Write([]byte("\r\n"))
				total += wrt
				if err != nil {
					return total, err
				}
			} else {
				wrt, err := w.W.Write(split)
				total += wrt
				if err != nil {
					return total, err
				}
			}
		}
	}
	return total, nil
}
