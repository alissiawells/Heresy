package main

import (
	"os"
	"fmt"
	"io/ioutil"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

type Response struct {
	Response []Res
}

type Res struct {
	Id         int `json:"-"`
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
		Name string
	}
	Relatives    []Rel
	Verified     int
	Home_phone   int `json:"-"`
	Mobile_phone int `json:"-"`
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
	Name string
}

func Encrypt(data string) string {
	h := sha256.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func main() {

	input := os.Args[1]
	file, e := os.Open(input)
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}
	defer file.Close()

	byteValue, _ := ioutil.ReadAll(file)

	var res *Response

	json.Unmarshal([]byte(byteValue), &res)

	res.Response[0].First_name = Encrypt(res.Response[0].First_name)
	res.Response[0].Last_name = Encrypt((res.Response[0].Last_name))
	res.Response[0].Maide_name = Encrypt((res.Response[0].Maide_name))
	res.Response[0].Relation_partner.Name = Encrypt((res.Response[0].Relation_partner.Name))
	res.Response[0].Relatives[0].Name = Encrypt((res.Response[0].Relatives[0].Name))

	b, _ := json.Marshal(res)

	output := os.Args[2]
	jsonFile, err := os.Create(output)

	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	jsonFile.Write(b)
	jsonFile.Close()
}
