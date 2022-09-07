package domain

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/rs/zerolog/log"
)

type Menu struct {
	FoodsCount int
	Foods      []Food
}

func GetMenu() Menu {
	file, err := os.Open("config/menu.json")
	if err != nil {
		log.Fatal().Err(err).Msg("Error opening menu.json")
	}

	byteValue, _ := ioutil.ReadAll(file)
	var menu Menu
	json.Unmarshal(byteValue, &menu)

	menu.FoodsCount = len(menu.Foods)
	return menu
}
