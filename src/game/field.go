package game

import (
	"math/rand"
	"io/ioutil"
	"os"
	"strconv"
	"encoding/json"
	"github.com/f26401004/Lifegamer-Diep-backend/src/util"
)

// define the NewField function in game package
func NewField (width, height float64) *util.Size {
	field := util.Size {
		W: width,
		H: height,
	}
	return &field
}

// define the NewStuff function in game package
func NewStuff() *Stuff {
	jsonFile, err := os.Open("src/config/stuffType.json")
	if (err != nil) {
		return nil
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var stuffType map[string]interface{}
  json.Unmarshal([]byte(byteValue), &stuffType)

	var type_num = int(rand.Float64() * 4 + 1)
	var type_attr = stuffType[strconv.Itoa(type_num)].(StuffAttribute)

	uuid, _ := util.NewUUID()
	var new_stuff = Stuff {
		GameObject: GameObject {
			Id: uuid,
			Position: util.Point {
				X: rand.Float64() * 1023,
				Y: rand.Float64() * 1023,
			},
			Mass: 1.0,
			Radius: 50.0,
			Velocity: util.VelocityFormat {
				X: 0.0,
				Y: 0.0,
			},
			Acceleration: util.AccelerationFormat {
				Up: 0.0,
				Down: 0.0,
				Left: 0.0,
				Right: 0.0,
			},
		},
		Type: type_num,
		Attr: StuffAttribute {
			HP: type_attr.HP,
			EXP: type_attr.EXP,
			BodyDamage: type_attr.BodyDamage,
		},
	}
	return &new_stuff
}
