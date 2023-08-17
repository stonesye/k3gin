package ws_api

import "io"

// WSApi 所有的API都需要实现它
type WSApi interface {
	Read(p []byte) (n int, err error)
	Write(p []byte) (n int, err error)
	Copy(io.Writer, io.Reader) (int64, error)
}
