package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math"

	"gopkg.in/yaml.v3"
)

//bp parsing
type Material struct {
	Quantity int `yaml:"quantity"`
	TypeID   int `yaml:"typeID"`
}

type Skills struct {
	Level  int `yaml:"level"`
	TypeID int `yaml:"typeID"`
}

type Products struct {
	Probability float32 `yaml:"probability"`
	Quantity int `yaml:"quantity"`
	TypeID   int `yaml:"typeID"`
}

type Copying struct {
	Materials []Material `yaml:"materials"`
	Skills    []Skills `yaml:"skills"`
}

type Invention struct {
	Materials []Material `yaml:"materials"`
	Products []Products `yaml:"products"`
	Skills    []Skills `yaml:"skills"`
	Time int `yaml:"time"`
}

type Manufacturing struct {
	Materials []Material `yaml:"materials"`
	Products  []Products `yaml:"products"`
	Skills    []Skills `yaml:"skills"`
	Time      int `yaml:"time"`
}

type Reaction struct {
	Materials []Material `yaml:"materials"`
	Products  []Products `yaml:"products"`
	Skills    []Skills `yaml:"skills"`
	Time      int `yaml:"time"`
}

type Research_material struct {
	Materials []Material `yaml:"materials"`
	Skills    []Skills `yaml:"skills"`
	Time      int `yaml:"time"`
}

type Research_time struct {
	Materials []Material `yaml:"materials"`
	Skills    []Skills `yaml:"skills"`
	Time      int `yaml:"time"`
}

type Activities struct {
	Invention         Invention `yaml:"invention"`
	Copying           Copying `yaml:"copying"`
	Manufacturing     Manufacturing `yaml:"manufacturing"`
	Reaction          Reaction `yaml:"reaction"`
	Research_material Research_material `yaml:"research_material"`
	Research_time     Research_time `yaml:"research_time"`
}

type bp struct {
	Activities Activities `yaml:"activities"`
	BlueprintTypeID int `yaml:"blueprintTypeID"`
	MaxProductionLimit int `yaml:"maxProductionLimit"`
}

//types parsing
type Description struct {
	De string `yaml:"de"`
	En string `yaml:"en"`
	Fr string `yaml:"fr"`
	Ja string `yaml:"ja"`
	Ru string `yaml:"ru"`
	Zh string `yaml:"zh"`
}

type Name struct {
	De string `yaml:"de"`
	En string `yaml:"en"`
	Fr string `yaml:"fr"`
	Ja string `yaml:"ja"`
	Ru string `yaml:"ru"`
	Zh string `yaml:"zh"`
}

type entity struct {
	BasePrice float32 `yaml:"basePrice"`
	Capacity float32 `yaml:"capacity"`
	Description Description `yaml:"description"`
	FactionID int `yaml:"factionID"`
	GraphicID int `yaml:"graphicID"`
	GroupID int `yaml:"groupID"`
	MarketGroupID int `yaml:"marketGroupID"`
	Mass float32 `yaml:"mass"`
	//Masteries int `yaml:"masteries"`
	Name Name `yaml:"name"`
	Volume float32 `yaml:"volume"`
}

func compositeCheck(listOfComponents []Material, bpArray *map[int]bp) ([]Material, bool, error){
	var checker, flag bool
	var bpID int
	var err error

	var newListOfComponents []Material

	flag = false

	for _, v := range listOfComponents {
		bpID, checker, err = innerCheck(v, bpArray)
		if err != nil {
			fmt.Println("some error with searching:", v)
		}

		if checker {
			flag = true
			var newBP []Material
			var newProductQuantity int

			if (*bpArray)[bpID].Activities.Manufacturing.Materials != nil {
				newBP = (*bpArray)[bpID].Activities.Manufacturing.Materials
				newProductQuantity = (*bpArray)[bpID].Activities.Manufacturing.Products[0].Quantity
			}
			if (*bpArray)[bpID].Activities.Reaction.Materials != nil {
				newBP = (*bpArray)[bpID].Activities.Reaction.Materials
				newProductQuantity = (*bpArray)[bpID].Activities.Reaction.Products[0].Quantity
			}

			for i := 0; i < len(newBP); i++ {
				newBP[i].Quantity = int(float64(newBP[i].Quantity) * math.Ceil(float64(v.Quantity) / float64(newProductQuantity)))
			}

			newListOfComponents = append(newListOfComponents, newBP...)
		} else {
			newListOfComponents = append(newListOfComponents, v)
		}
	}

	return newListOfComponents, flag, err
}

func innerCheck(material Material, bpArray *map[int]bp) (int, bool, error) {
	var count, bpID int

	for _, v := range *bpArray {
		if v.Activities.Manufacturing.Products != nil && 
		   v.Activities.Manufacturing.Products[0].TypeID == material.TypeID &&
		   v.BlueprintTypeID != 4313 &&
		   v.BlueprintTypeID != 4314 &&
		   v.BlueprintTypeID != 4315 &&
		   v.BlueprintTypeID != 4316  {
			count++
			bpID = v.BlueprintTypeID
		}

		if v.Activities.Reaction.Products != nil &&
		   v.Activities.Reaction.Products[0].TypeID == material.TypeID &&
		   v.BlueprintTypeID != 4313 &&
		   v.BlueprintTypeID != 4314 &&
		   v.BlueprintTypeID != 4315 &&
		   v.BlueprintTypeID != 4316 {
			count++
			bpID = v.BlueprintTypeID
		}
	}

	if count == 0 {
		return -1, false, nil
	} else if count == 1 {
		return bpID, true, nil
	} else {
		return -1, false, errors.New("error search in bp array")
	}
}

func searchId(itemName string, idArray *map[int]entity, bpArray *map[int]bp) int{
	var id int

	for i, v := range *idArray {
		//fmt.Println("catcxxxh ", v.Name.En)
		if itemName == v.Name.En {
			fmt.Println("catch ", v.Name.En)
			for _, w := range *bpArray {
				if w.Activities.Manufacturing.Products != nil && 
				w.Activities.Manufacturing.Products[0].TypeID == i {
					fmt.Println("id:", i)
					id = w.BlueprintTypeID
				}
			}
		}
	}

	fmt.Println("id:", id)
	return id
}

func clean(material []Material) []Material {
	m := make(map[int] int)
	//j := 0

	for _, v := range material {
		if m[v.TypeID] == 0 {
			m[v.TypeID] = v.Quantity
		} else {
			m[v.TypeID] = m[v.TypeID] + v.Quantity
		}
	}

	var returnMaterials []Material
	
	for i, w := range m {
		var mat Material
		mat.Quantity = w
		mat.TypeID = i
		returnMaterials = append(returnMaterials, mat)
	}

	return returnMaterials
}

func volumeCalculation(materials []Material, idArray *map[int]entity) float32{
	var volume float32
	for _, v := range materials {
		volume += (*idArray)[v.TypeID].Volume * float32(v.Quantity)
	}
	return volume
}

func main() {
	bps, err := ioutil.ReadFile("blueprints.yaml")
	if err != nil {
		fmt.Println("error reading file")
	}

	bpArray := make(map[int]bp)

	err = yaml.Unmarshal(bps, &bpArray)
	if err != nil {
		fmt.Println(err)
	}

	ID, err := ioutil.ReadFile("typeIDs.yaml")
	if err != nil {
		fmt.Println("error reading file")
	}

	IDArray := make(map[int]entity)

	err = yaml.Unmarshal(ID, &IDArray)
	if err != nil {
		fmt.Println(err)
	}

	var itemName string
	fmt.Println("item name: ")
	fmt.Scanf("%s\n", &itemName)

	materials := bpArray[searchId(itemName, &IDArray, &bpArray)].Activities.Manufacturing.Materials
	var checker bool

	for {
		materials, checker, err = compositeCheck(materials, &bpArray)
		if err != nil {
			fmt.Println("error")
		}

		materials = clean(materials)

		print(materials, &IDArray)
		fmt.Println("volume:", volumeCalculation(materials, &IDArray))
		fmt.Println("------")

		if !checker {
			break
		}
	}

	//print(materials, &IDArray)
}