package game

import (
	"math/rand"
	"io/ioutil"
	"os"
	"strconv"
	"encoding/json"
)

// define the NewField function in game package
func NewField (width, height float64) *Size {
	field := Size {
		W: width,
		H: height,
	}
	return &field
}

// define the NewStuff function in game package
func NewStuff() *Stuff {
	jsonFile, err := os.Open("../utils/stuffType.json")
	if (err != nil) {
		return nil
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var stuffType map[string]interface{}
  json.Unmarshal([]byte(byteValue), &stuffType)

	var type_num = int(rand.Float64() * 4 + 1)
	var type_attr = stuffType[strconv.Itoa(type_num)].(StuffAttribute)
	var new_stuff = Stuff {
		Type: type_num,
		Position: Point {
			X: rand.Float64() * 1023,
			Y: rand.Float64() * 1023,
		},
		Attr: StuffAttribute {
			HP: type_attr.HP,
			EXP: type_attr.EXP,
			BodyDamage: type_attr.BodyDamage,
		},
	}
	return &new_stuff
}
