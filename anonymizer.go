package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

type User struct {
	Response []Res
}

type Res struct {
	Id         int `json:"-"`
	HiddenId string
	First_name string
	Last_name  string
	City       []struct {
		Id    int
		Title string
	}
	Maide_name       string
	Relation         int
	Relation_partner struct {
		Id   int `json:"-"`
		HiddenId string
		Name string
	}
	Relatives    []Rel
	Verified     int
	Home_phone   int `json:"-"`
	Hidden_hphone string
	Mobile_phone int `json:"-"`
	Hidden_mphone string
	Personal     struct {
		Religion string
	}

	Work struct {
		Id   int
		Name string
	}

	Universities struct {
		Id   int
		Name string
	}
}

type Rel struct {
	Type string
	Id   int `json:"-"`
	HiddenId string
	Name string
}

func Encrypt(data string) (hash string) {
	h := sha256.New()
	h.Write([]byte(data))
	hash = hex.EncodeToString(h.Sum(nil))
	return hash
}

func main() {

	input := os.Args[1] // private.json
	file, e := os.Open(input)
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}
	defer file.Close()

	byteValue, _ := ioutil.ReadAll(file)

	var res *User

	json.Unmarshal([]byte(byteValue), &res)

	res.Response[0].First_name = Encrypt((res.Response[0].First_name))
	res.Response[0].Last_name = Encrypt((res.Response[0].Last_name))
	res.Response[0].Maide_name = Encrypt((res.Response[0].Maide_name))
	res.Response[0].Relation_partner.Name = Encrypt((res.Response[0].Relation_partner.Name))
	res.Response[0].Relatives[0].Name = Encrypt((res.Response[0].Relatives[0].Name))
		
	res.Response[0].HiddenId = Encrypt((strconv.Itoa(res.Response[0].Id)))
	res.Response[0].Relation_partner.HiddenId = Encrypt((strconv.Itoa(res.Response[0].Relation_partner.Id)))
	res.Response[0].Relatives[0].HiddenId = Encrypt((strconv.Itoa(res.Response[0].Relatives[0].Id)))
	res.Response[0].Hidden_hphone = Encrypt((strconv.Itoa(res.Response[0].Home_phone)))
	res.Response[0].Hidden_mphone = Encrypt((strconv.Itoa(res.Response[0].Mobile_phone)))

	b, _ := json.Marshal(res)

	output := os.Args[2] // depersonalized.json
	jsonFile, err := os.Create(output)

         if err != nil {
                 panic(err)
         }
         defer jsonFile.Close()

         jsonFile.Write(b)
         jsonFile.Close()
}
