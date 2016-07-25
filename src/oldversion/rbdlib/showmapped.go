package main

import (
	"fmt"
	// 	"io/ioutil"
	"encoding/json"
	"os/exec"
	"strings"
)

type RbdImage struct {
	Id     string
	Pool   string
	Image  string
	Snap   string
	Device string
}

func main() {

	rbd, err := exec.Command("rbd", "showmapped").Output()

	if err != nil {
		panic(err)
		return
	}

	// fmt.Println(">rbd showmapped")
	// fmt.Println(string(images))

	outputs := strings.Split(string(rbd), "\n")[1:]
	// fmt.Println(outputs)

	var images []RbdImage
	for _, value := range outputs {
		v := strings.Split(value, "\t")
		fmt.Printf("---> %v, %s, %s, %s, %s\n", v, v[0], v[1], v[2], v[3], v[4])
		if !strings.EqualFold(v[0], "") {
			image := RbdImage{Id: v[0], Pool: v[1], Image: v[2], Snap: v[3], Device: v[4]}
			images = append(images, image)
		}
	}

	res, err := json.Marshal(images)
	if err != nil {
		panic(err)
		return
	}

	fmt.Println(string(res))
	// outputs = outputs[1:]
	// fmt.Println(outputs)
}
