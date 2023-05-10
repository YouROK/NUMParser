package db

import (
	"NUMParser/db/models"
	"bytes"
	"compress/flate"
	"encoding/json"
	"log"
	"os"
	"time"
)

func GetKPDetails() []*models.KPDetail {
	return kpds
}

func SetKPDetails(list []*models.KPDetail) {
	kpds = list
	isKPDSChange = true
}

func AddKPD(k *models.KPDetail) {
	isKPDSChange = true
	k.UpdateDate = time.Now()
	for i, d := range kpds {
		if d.KinopoiskID == k.KinopoiskID {
			kpds[i] = k
			return
		}
	}

	muKpds.Lock()
	defer muKpds.Unlock()
	kpds = append(kpds, k)
}

func SaveKPDs() {
	muKpds.Lock()
	var list []*models.KPDetail
	for _, kpd := range kpds {
		if time.Now().Before(kpd.UpdateDate.Add(time.Hour * 24 * 7)) {
			list = append(list, kpd)
		}
	}
	if len(list) != len(kpds) {
		kpds = list
		isKPDSChange = true
	}
	muKpds.Unlock()

	if !isKPDSChange || len(kpds) == 0 {
		return
	}

	log.Println("Save KP")
	buf, err := json.Marshal(kpds)
	if err != nil {
		return
	}

	var b bytes.Buffer
	w, err := flate.NewWriter(&b, flate.BestCompression)
	if err != nil {
		return
	}
	w.Write(buf)
	w.Close()

	err = os.WriteFile("kpds.ls", b.Bytes(), 0666)
	if err != nil {
		return
	}
	isKPDSChange = false
}
