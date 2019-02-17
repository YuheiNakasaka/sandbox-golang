package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/YuheiNakasaka/sandbox-golang/util"
)

func dumpChunk(chunk io.Reader) {
	var length int32
	binary.Read(chunk, binary.BigEndian, &length)
	buffer := make([]byte, 4)
	chunk.Read(buffer)
	fmt.Printf("chunk '%v'(%d bytes)\n", string(buffer), length)
}

func readChunks(file *os.File) []io.Reader {
	var chunks []io.Reader
	file.Seek(8, 0)
	var offset int64 = 8
	const LegthSize int64 = 4
	const KindSize int64 = 4
	const CRCSize int64 = 4

	for {
		// int32 = 32bit = 4byte なので
		// binary.Readはfileのoffsetから4byte分だけ読み込んでくれる
		// つまり今回の場合はLengthを読み込んでいることになる
		var dataSize int32
		err := binary.Read(file, binary.BigEndian, &dataSize)
		if err == io.EOF {
			break
		}
		chunks = append(chunks, io.NewSectionReader(file, offset, LegthSize+int64(dataSize)+KindSize+CRCSize))
		// 現在のfileのoffsetはLengthを読み込んだところまでなので次のチャンクの先頭までSeekする
		offset, _ = file.Seek(int64(dataSize)+KindSize+CRCSize, 1)
	}
	return chunks
}

func main() {
	// NewSectionReader Example
	reader := strings.NewReader("Example of io.SectionReader\n")
	sectionReader := io.NewSectionReader(reader, 14, 7)
	io.Copy(os.Stdout, sectionReader)
	fmt.Println("")

	// ビッグエンディアン
	data := []byte{0x0, 0x0, 0x27, 0x10}
	var i int32
	binary.Read(bytes.NewReader(data), binary.BigEndian, &i)
	fmt.Printf("data: %d\n", i)

	// PNGのバイナリをチャンクごとに読み取る
	pjtpath := util.GetProjectPath()
	file, err := os.Open(pjtpath + "/lenna.png")
	if err != nil {
		panic(err)
	}
	chunks := readChunks(file)
	for _, chunk := range chunks {
		dumpChunk(chunk)
	}
}
