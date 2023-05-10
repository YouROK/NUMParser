package db

import (
	"NUMParser/db/models"
)

func indexTorrs() {
	//if IsTorrsChange {
	//torrsearch.NewIndex(GetTorrs())
	//}
}

func SearchTorr(query string) []*models.TorrentDetails {
	//matchedIDs := torrsearch.Search(query)
	//if len(matchedIDs) == 0 {
	//	return nil
	//}
	////torrs := GetTorrs()
	//var list []*models.TorrentDetails
	//for _, id := range matchedIDs {
	//	list = append(list, torrs[id])
	//}
	//
	//hash := utils.ClearStr(query)
	//
	//sort.Slice(list, func(i, j int) bool {
	//	lhash := utils.ClearStr(strings.ToLower(list[i].Name+list[i].GetNames())) + strconv.Itoa(list[i].Year)
	//	lev1 := levenshtein.ComputeDistance(hash, lhash)
	//	lhash = utils.ClearStr(strings.ToLower(list[j].Name+list[j].GetNames())) + strconv.Itoa(list[j].Year)
	//	lev2 := levenshtein.ComputeDistance(hash, lhash)
	//	if lev1 == lev2 {
	//		return list[j].CreateDate.Before(list[i].CreateDate)
	//	}
	//	return lev1 < lev2
	//})
	//return list
	return nil
}
