package cmd

import (
	"strconv"
	"strings"
)

type CustomStringSlice struct {
	lines        []string
	column       int
	isNumeric    bool
	ignoreBlanks bool
}

func (cs CustomStringSlice) Len() int {
	return len(cs.lines)
}

func (cs CustomStringSlice) Swap(i, j int) {
	cs.lines[i], cs.lines[j] = cs.lines[j], cs.lines[i]
}

func (cs CustomStringSlice) Less(i, j int) bool {

	var valI, valJ string
	if cs.column == -1 {
		valI = cs.lines[i]
		valJ = cs.lines[j]
	} else {
		strI := strings.Split(cs.lines[i], "\t")
		strJ := strings.Split(cs.lines[j], "\t")

		if len(strI) > cs.column {
			valI = strI[cs.column]
		}
		if len(strJ) > cs.column {
			valJ = strJ[cs.column]
		}
	}

	if cs.ignoreBlanks {
		valI = strings.TrimRight(valI, " \t")
		valJ = strings.TrimRight(valJ, " \t")
	}

	if cs.isNumeric {
		nI, errI := strconv.Atoi(valI)
		nJ, errJ := strconv.Atoi(valJ)

		if errI == nil && errJ == nil {
			return nI < nJ
		} else if errI == nil {
			return true
		} else if errJ == nil {
			return false
		} else {
			return strings.Compare(valI, valJ) < 0
		}
	}

	return strings.Compare(valI, valJ) < 0
}
