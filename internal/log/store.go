package log

import (
	"bufio"
	"encoding/binary"
	"os"
	"sync"
)

var (
	// レコード長とインデックスエントリを永続化するためのエンコーディング
	enc = binary.BigEndian
)

const (
	// レコード長を格納するために使うバイト数
	lenWidth = 8
)

type store struct {
	*os.File
	mu   sync.Mutex
	buf  *bufio.Writer
	size uint64
}

// 与えられたファイルに基づきストアを作成する
func newStore(f *os.File) (*store, error) {
	fi, err := os.Stat(f.Name())
	if err != nil {
		return nil, err
	}
	// サービス再起動時など、すでにデータを含むファイルから再作成する場合のためにサイズを取得
	size := uint64(fi.Size())
	return &store{
		File: f,
		size: size,
		buf:  bufio.NewWriter(f),
	}, nil
}

func (s *store) Append(p []byte) (n uint64, pos uint64, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	pos = s.size
	// レコードを読み出す時に何バイト読めばいいかがわかるようにレコード長を書き出す.
	// システムコールを減らしてパフォーマンス改善するため、ファイルに直接ではなくバッファ付きWriterに書き出す.
	if err := binary.Write(s.buf, enc, uint64(len(p))); err != nil {
		return 0, 0, err
	}
	w, err := s.buf.Write(p)
	if err != nil {
		return 0, 0, err
	}
	w += lenWidth
	s.size += uint64(w)
	return uint64(w), pos, nil
}

func (s *store) Read(pos uint64) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// バッファがまだディスクにFlushしていないレコードを読み出そうとしている場合に備えてまずFlushする
	if err := s.buf.Flush(); err != nil {
		return nil, err
	}
	// 何バイト読み出す必要があるか調べるため、`Append`の`binary.Write`で書き込んだ部分を読み取る
	size := make([]byte, lenWidth)
	if _, err := s.File.ReadAt(size, int64(pos)); err != nil {
		return nil, err
	}
	b := make([]byte, enc.Uint64(size))
	if _, err := s.File.ReadAt(b, int64(pos+lenWidth)); err != nil {
		return nil, err
	}
	return b, nil
}

func (s *store) ReadAt(p []byte, off int64) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.buf.Flush(); err != nil {
		return 0, err
	}
	return s.File.ReadAt(p, off)
}

func (s *store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	err := s.buf.Flush()
	if err != nil {
		return err
	}
	return s.File.Close()
}
