package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
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
func (x Res) encryptRes() Res {
	x.First_name = Encrypt(x.First_name)
	x.Last_name = Encrypt((x.Last_name))
	x.Maide_name = Encrypt((x.Maide_name))
	x.Relation_partner.Name = Encrypt((x.Relation_partner.Name))
	x.Relatives[0].Name = Encrypt((x.Relatives[0].Name))
	return x
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


func (res Response) encryptResponse() Response {
	for i, x := range res.Response {
		res.Response[i] = x.encryptRes()
	}
	return res
}

func main() {

	input := os.Args[1]
	//input := "input.json"
    fmt.Println("Hello, little hacker!")

	file, e := os.Open(input)
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}
	defer file.Close()

	byteValue, _ := ioutil.ReadAll(file)

	var res *Response

	json.Unmarshal([]byte(byteValue), &res)
	res.encryptResponse()
	b, _ := json.MarshalIndent(res, "", "    ")

	//output := os.Args[2]
	output := "output.json"
	jsonFile, err := os.Create(output)

	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	jsonFile.Write(b)
	jsonFile.Close()
}
