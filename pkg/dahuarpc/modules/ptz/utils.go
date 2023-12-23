package ptz

import (
	"regexp"
	"strconv"
	"strings"
)

var numbersReg *regexp.Regexp = regexp.MustCompile("[0-9]+")

func getSuffix(s string, count int) string {
	end := len(s)
	start := end - count
	if start < 0 {
		start = 0
	}
	return s[start:end]
}

func nextSeq(session string, id int) int {
	sessionNumberRaw := strings.Join(numbersReg.FindAllString(session, -1), "")
	sessionNumber, _ := strconv.Atoi(sessionNumberRaw)
	sessionBinary := getSuffix(strconv.FormatInt(int64(sessionNumber), 2), 24)
	sessionBinary2 := getSuffix("00000000"+strconv.FormatInt(int64(id), 2), 8)
	seq, _ := strconv.ParseInt(sessionBinary+sessionBinary2, 2, 32)
	return int(seq)
}

func nextID(id int) int {
	return (id + 1) % 256
}
