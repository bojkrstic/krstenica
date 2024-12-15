package api

import (
	"krstenica/pkg/apiutil"
	"log"
)

func (val *HramUpdateData) Validate(ID uint) (map[string]interface{}, error) {
	var errItems []*apiutil.ErrorItem

	update := map[string]interface{}{}

	if val.NazivHrama != nil {
		if *val.NazivHrama == "" {
			errItems = append(errItems, &apiutil.ErrorItem{
				Name:        "Name",
				Message:     "Pogresna vrednost",
				Description: "Ime nije odredjeno."})
		}
		update["naziv_hrama"] = *val.NazivHrama
	}

	if val.Status != nil {
		if *val.Status == "deleted" {
			errItems = append(errItems, &apiutil.ErrorItem{
				Name:        "Deleted",
				Message:     "Ne mozete da menajte obrisan slog",
				Description: "Hram je vec obrisan",
			})
		}
		update["status"] = *val.Status
	}

	if len(errItems) > 0 {
		for _, item := range errItems {
			log.Println(item.Name, item.Message, item.Description)
		}
		return nil, apiutil.NewValidationError(errItems)
	}
	log.Print("UPDATE", update)
	return update, nil

}
