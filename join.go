package stdio

import (
	"io"

	"golang.org/x/sync/errgroup"
)

// Join connects two ReaderWriter interfaces to form a direct connection.
func Join(a io.ReadWriter, b io.ReadWriter) error {
	g := errgroup.Group{}

	g.Go(func() error {
		_, err := io.Copy(a, b)
		return err
	})

	g.Go(func() error {
		_, err := io.Copy(b, a)
		return err
	})

	return g.Wait()
}
