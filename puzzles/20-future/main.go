package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	unidecode "golang.org/x/text/encoding/unicode"

	"github.com/lthummus/i18n-puzzles/input"
)

func chunk20(x []int32) []byte {
	if len(x)%2 != 0 {
		panic("needs to be even")
	}

	var ret []byte
	for i := 0; i < len(x); i += 2 {
		ret = append(ret, byte(x[i]>>12))
		ret = append(ret, byte((x[i]&0b00000000111111110000)>>4))

		n := (x[i]&0b1111)<<4 | (x[i+1]&0b11110000000000000000)>>16
		ret = append(ret, byte(n))

		ret = append(ret, byte((x[i+1]&0b00001111111100000000)>>8))
		ret = append(ret, byte(x[i+1]&0xFF))
	}

	return ret
}

func chunk28(x []int32) []byte {
	if len(x)%2 != 0 {
		x = append(x, 0)
	}
	var ret []byte
	for i := 0; i < len(x); i += 2 {
		ret = append(ret, byte(x[i]>>20))
		ret = append(ret, byte((x[i]>>12)&0xFF))
		ret = append(ret, byte((x[i]>>4)&0xFF))

		n := byte(x[i]&0xF)<<4 | byte(x[i+1]>>24)
		ret = append(ret, n)

		ret = append(ret, byte((x[i+1]>>16)&0xFF))
		ret = append(ret, byte((x[i+1]>>8)&0xFF))
		ret = append(ret, byte(x[i+1]&0xFF))
	}
	return ret
}

func futureDecode(b []byte) []int32 {
	var ret []int32
	for i := 0; i < len(b); {
		firstByte := b[i]
		var v int32
		var size int
		if firstByte&0x80 == 0 {
			v = int32(firstByte)
			size = 1
		} else if firstByte&0xE0 == 0xC0 {
			v = int32(firstByte & 0x1F)
			size = 2
		} else if firstByte&0xF0 == 0xE0 {
			v = int32(firstByte & 0x0F)
			size = 3
		} else if firstByte&0xF8 == 0xF0 {
			v = int32(firstByte & 0x07)
			size = 4
		} else if firstByte&0xFC == 0xF8 {
			v = int32(firstByte & 0x03)
			size = 5
		} else if firstByte&0xFe == 0xFC {
			v = int32(firstByte & 0x01)
			size = 6
		} else {
			panic("invalid byte")
		}

		for j := 1; j < size; j++ {
			if i+j >= len(b) {
				panic("not enough bytes")
			}
			c := b[i+j]
			if c&0xC0 != 0x80 {
				panic("expected continuation byte")
			}
			v = (v << 6) | int32(c&0x3F)
		}

		if v == 0 {
			// null byte, end it
			break
		}

		ret = append(ret, v)
		i += size
	}

	return ret
}

func decrypt(x string) string {
	x = strings.Replace(x, "\n", "", -1)
	decoded, err := base64.StdEncoding.DecodeString(x)
	if err != nil {
		panic(err)
	}

	// there's a UTF-16 BOM at the front, so decode that too
	decoder := unidecode.UTF16(unidecode.LittleEndian, unidecode.ExpectBOM).NewDecoder()
	t, err := decoder.String(string(decoded))
	if err != nil {
		panic(err)
	}

	var chunks []int32

	for _, curr := range t {
		chunks = append(chunks, curr)
	}

	c := chunk20(chunks)

	fd := futureDecode(c)
	fd2 := chunk28(fd)
	c2 := futureDecode(fd2)

	return string(c2)
}

// This has some work to do to clean up -- I was able to decrypt enough of the message in order to figure out
// what the problem wanted 😏
func main() {
	in, err := input.GetInputUTF8(context.Background(), 20, input.RealInput)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", decrypt(in))

}
