package utils

import (
	"NUMParser/db/models"
	"github.com/agnivade/levenshtein"
	"runtime"
	"runtime/debug"
	"strings"
)

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func Filter[T any](arr []T, fn func(i int, e T) bool) []T {
	var list []T
	for i, t := range arr {
		if !fn(i, t) {
			list = append(list, t)
		}
	}
	return list
}

func IsEqTorrKP(t *models.TorrentDetails, kp *models.KPDetail) bool {
	if Abs(t.Year-kp.Year) > 1 {
		return false
	}
	ruHash := ClearStr(t.Name)
	ruLev := levenshtein.ComputeDistance(ruHash, ClearStr(strings.ToLower(kp.NameRu)))
	isRu := false
	isEn := false
	if len(ruHash) == 0 || kp.NameRu == "" {
		isRu = true
	} else if len([]rune(ruHash)) > 5 {
		isRu = ruLev < len([]rune(ruHash))/2
	} else {
		isRu = ruLev < len([]rune(ruHash))/3
	}

	if !isRu {
		return false
	}

	if t.GetNames() == "" || kp.NameEn+kp.NameOriginal == "" {
		isEn = true
	}

	if !isEn {
		for _, name := range t.Names {
			lev := levenshtein.ComputeDistance(ClearStr(name), ClearStr(kp.NameEn+kp.NameOriginal))
			if len([]rune(ClearStr(name))) > 5 {
				isEn = lev < len([]rune(ClearStr(name)))/2
			} else {
				isEn = lev < len([]rune(ClearStr(name)))/3
			}
			if isEn {
				break
			}
		}
	}

	return isEn && isRu
}

func FreeOSMemGC() {
	runtime.GC()
	debug.FreeOSMemory()
}
