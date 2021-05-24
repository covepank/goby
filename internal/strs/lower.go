package strs

import (
	"github.com/sanbsy/gopkg/bufferpool"
)

func LowerName(name string) string {
	buf := bufferpool.Get()
	defer buf.Free()

	for idx := range name {
		if IsUpperLetter(name[idx]) && buf.Len() > 0 {
			if idx > 0 && !IsUpperLetter(name[idx-1]) || idx+1 < len(name) && !IsUpperLetter(name[idx+1]) {
				buf.WriteByte('_')
			}
		}
		buf.WriteByte(ToLower(name[idx]))
	}
	return buf.String()
}

func IsUpperLetter(c byte) bool {
	return c >= 'A' && c <= 'Z'
}

func IsLowerLetter(c byte) bool {
	return c >= 'a' && c <= 'z'
}

func ToLower(c byte) byte {
	if IsLowerLetter(c) {
		return c
	}
	if IsUpperLetter(c) {
		return c + 32
	}

	return c
}
