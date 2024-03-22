package log

import (
	"io"
	"os"

	"github.com/tysonmote/gommap"
)

const (
	// レコードのオフセットを保持するフィールドの長さ (Bytes)
	offWidth uint64 = 4
	// レコードのファイル内での位置を保持するフィールドの長さ (Bytes)
	posWidth uint64 = 8
	// 各エントリーの長さ (Bytes)
	// オフセットが与えられた時にエントリーの位置 (offset * entWidth) に直接ジャンプするために使う
	entWidth = offWidth + posWidth
)

type index struct {
	file *os.File
	mmap gommap.MMap
	size uint64
}

func newIndex(f *os.File, c Config) (*index, error) {
	idx := &index{
		file: f,
	}
	fi, err := os.Stat(f.Name())
	if err != nil {
		return nil, err
	}
	// インデックスエントリを追加する際にファイル内のデータ量を管理可能にするため、現在のファイルサイズを保存
	idx.size = uint64(fi.Size())
	// メモリへマップする前に、ファイルを最大のインデックスサイズまで拡張
	if err = os.Truncate(
		f.Name(), int64(c.Segment.MaxIndexBytes),
	); err != nil {
		return nil, err
	}
	if idx.mmap, err = gommap.Map(
		idx.file.Fd(),
		gommap.PROT_READ|gommap.PROT_WRITE,
		gommap.MAP_SHARED,
	); err != nil {
		return nil, err
	}
	return idx, nil
}

func (i *index) Close() error {
	if err := i.mmap.Sync(gommap.MS_SYNC); err != nil {
		return err
	}
	if err := i.mmap.UnsafeUnmap(); err != nil {
		return err
	}
	if err := i.file.Sync(); err != nil {
		return err
	}
	// 永続化されたファイルのサイズを実際のデータ量に切り詰める
	if err := i.file.Truncate(int64(i.size)); err != nil {
		return err
	}
	return i.file.Close()
}

func (i *index) Read(in int64) (out uint32, pos uint64, err error) {
	if i.size == 0 {
		return 0, 0, io.EOF
	}
	if in == -1 {
		out = uint32((i.size / entWidth) - 1)
	} else {
		out = uint32(in)
	}
	pos = uint64(out) * entWidth
	if i.size < pos+entWidth {
		return 0, 0, io.EOF
	}
	out = enc.Uint32(i.mmap[pos : pos+offWidth])
	pos = enc.Uint64(i.mmap[pos+offWidth : pos+entWidth])
	return out, pos, nil
}

func (i *index) Write(off uint32, pos uint64) error {
	if i.isMaxed() {
		return io.EOF
	}
	enc.PutUint32(i.mmap[i.size:i.size+offWidth], off)
	enc.PutUint64(i.mmap[i.size+offWidth:i.size+entWidth], pos)
	i.size += uint64(entWidth)
	return nil
}

func (i *index) isMaxed() bool {
	return uint64(len(i.mmap)) < i.size+entWidth
}

func (i *index) Name() string {
	return i.file.Name()
}
