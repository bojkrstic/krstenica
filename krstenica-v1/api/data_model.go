package api

import "krstenica/krstenica-v1/dao"

type HramWo struct {
	NazivHrama string `json:"naziv_hrama"`
}
type HramUpdateData struct {
	NazivHrama *string `json:"naziv_hrama"`
	Status     *string `json:"status"`
}
type EparhijaWo struct {
	NazivEparhije string `json:"naziv_eparhije"`
}

type EparhijaUpdateWo struct {
	NazivEparhije *string `json:"naziv_eparhije"`
}

// mozda da se jos stavi na zilazu i u bazi status i comment
type HramCrtResWo struct {
	HramID     uint   `json:"hram_id"`
	NazivHrama string `json:"naziv_hrama"`
	Status     string `json:"status"`
	CreatedAt  string `json:"created_at"`
}

type EparhijaCrtResWo struct {
	EparhijaID    uint   `json:"eparhija_id"`
	NazivEparhije string `json:"naziv_eparhije"`
	Status        string `json:"status"`
	CreatedAt     string `json:"created_at"`
}

func makeResultSysApplication(sysApp *dao.HramDo) *HramCrtResWo {
	result := &HramCrtResWo{
		HramID:     sysApp.HramID,
		NazivHrama: sysApp.HramName,
		Status:     sysApp.Status,
		CreatedAt:  sysApp.CreatedAt.String(),
	}

	return result
}
