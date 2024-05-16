// Generated by `wit-bindgen-wrpc-go` 0.1.0. DO NOT EDIT!
package poll

import (
	bytes "bytes"
	context "context"
	binary "encoding/binary"
	errors "errors"
	fmt "fmt"
	wrpc "github.com/wrpc/wrpc/go"
	errgroup "golang.org/x/sync/errgroup"
	io "io"
	slog "log/slog"
	math "math"
)

// Poll for completion on a set of pollables.
//
// This function takes a list of pollables, which identify I/O sources of
// interest, and waits until one or more of the events is ready for I/O.
//
// The result `list<u32>` contains one or more indices of handles in the
// argument list that is ready for I/O.
//
// If the list contains more elements than can be indexed with a `u32`
// value, this function traps.
//
// A timeout can be implemented by adding a pollable from the
// wasi-clocks API to the list.
//
// This function does not return a `result`; polling in itself does not
// do any I/O so it doesn't fail. If any of the I/O sources identified by
// the pollables has an error, it is indicated by marking the source as
// being reaedy for I/O.
func Poll(ctx__ context.Context, wrpc__ wrpc.Client, in []Pollable) (r0__ []uint32, close__ func() error, err__ error) {
	if err__ = wrpc__.Invoke(ctx__, "wasi:io/poll@0.2.0", "poll", func(w__ wrpc.IndexWriter, r__ wrpc.IndexReadCloser) error {
		close__ = r__.Close
		var buf__ bytes.Buffer
		writes__ := make(map[uint32]func(wrpc.IndexWriter) error, 1)
		write0__, err__ := func(v []Pollable, w interface {
			io.ByteWriter
			io.Writer
		}) (write func(wrpc.IndexWriter) error, err error) {
			n := len(v)
			if n > math.MaxUint32 {
				return nil, fmt.Errorf("list length of %d overflows a 32-bit integer", n)
			}
			if err = func(v int, w io.Writer) error {
				b := make([]byte, binary.MaxVarintLen32)
				i := binary.PutUvarint(b, uint64(v))
				slog.Debug("writing list length", "len", n)
				_, err = w.Write(b[:i])
				return err
			}(n, w); err != nil {
				return nil, fmt.Errorf("failed to write list length of %d: %w", n, err)
			}
			slog.Debug("writing list elements")
			writes := make(map[uint32]func(wrpc.IndexWriter) error, n)
			for i, e := range v {
				write, err := (func(wrpc.IndexWriter) error)(nil), func(any) error { return errors.New("writing borrowed handles not supported yet") }(e)
				if err != nil {
					return nil, fmt.Errorf("failed to write list element %d: %w", i, err)
				}
				if write != nil {
					writes[uint32(i)] = write
				}
			}
			if len(writes) > 0 {
				return func(w wrpc.IndexWriter) error {
					var wg errgroup.Group
					for index, write := range writes {
						w, err := w.Index(index)
						if err != nil {
							return fmt.Errorf("failed to index writer: %w", err)
						}
						write := write
						wg.Go(func() error {
							return write(w)
						})
					}
					return wg.Wait()
				}, nil
			}
			return nil, nil
		}(in, &buf__)
		if err__ != nil {
			return fmt.Errorf("failed to write `in` parameter: %w", err__)
		}
		if write0__ != nil {
			writes__[0] = write0__
		}
		_, err__ = w__.Write(buf__.Bytes())
		if err__ != nil {
			return fmt.Errorf("failed to write parameters: %w", err__)
		}
		r0__, err__ = func(r wrpc.IndexReader, path ...uint32) ([]uint32, error) {
			var x uint32
			var s uint
			for i := 0; i < 5; i++ {
				slog.Debug("reading list length byte", "i", i)
				b, err := r.ReadByte()
				if err != nil {
					if i > 0 && err == io.EOF {
						err = io.ErrUnexpectedEOF
					}
					return nil, fmt.Errorf("failed to read list length byte: %w", err)
				}
				if b < 0x80 {
					if i == 4 && b > 1 {
						return nil, errors.New("list length overflows a 32-bit integer")
					}
					x = x | uint32(b)<<s
					vs := make([]uint32, x)
					for i := range vs {
						slog.Debug("reading list element", "i", i)
						vs[i], err = func(r io.ByteReader) (uint32, error) {
							var x uint32
							var s uint
							for i := 0; i < 5; i++ {
								slog.Debug("reading u32 byte", "i", i)
								b, err := r.ReadByte()
								if err != nil {
									if i > 0 && err == io.EOF {
										err = io.ErrUnexpectedEOF
									}
									return x, fmt.Errorf("failed to read u32 byte: %w", err)
								}
								if b < 0x80 {
									if i == 4 && b > 1 {
										return x, errors.New("varint overflows a 32-bit integer")
									}
									return x | uint32(b)<<s, nil
								}
								x |= uint32(b&0x7f) << s
								s += 7
							}
							return x, errors.New("varint overflows a 32-bit integer")
						}(r)
						if err != nil {
							return nil, fmt.Errorf("failed to read list element %d: %w", i, err)
						}
					}
					return vs, nil
				}
				x |= uint32(b&0x7f) << s
				s += 7
			}
			return nil, errors.New("list length overflows a 32-bit integer")
		}(r__, []uint32{0}...)
		if err__ != nil {
			return fmt.Errorf("failed to read result 0: %w", err__)
		}
		return nil
	}); err__ != nil {
		err__ = fmt.Errorf("failed to invoke `poll`: %w", err__)
		return
	}
	return
}

type Pollable interface {
	// Return the readiness of a pollable. This function never blocks.
	//
	// Returns `true` when the pollable is ready, and `false` otherwise.
	Ready(ctx__ context.Context, wrpc__ wrpc.Client) (bool, func() error, error)
	// `block` returns immediately if the pollable is ready, and otherwise
	// blocks until ready.
	//
	// This function is equivalent to calling `poll.poll` on a list
	// containing only this pollable.
	Block(ctx__ context.Context, wrpc__ wrpc.Client) (func() error, error)
	Drop(ctx__ context.Context, wrpc__ wrpc.Client) error
}
