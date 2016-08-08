package rbd

import (
	"app/backend/common/util/placeholder"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

const (
	RBD_CREATE         = "rbd create <image> -s <size> -p <pool>"
	RBD_SHOWMAPPED     = "rbd showmapped"
	RBD_UNMAP          = "rbd unmap <device>"
	RBD_MAP            = "rbd map <image>"
	MAKEFS             = "mkfs.<fs> <device>"
	RBD_REMOVE         = "rbd rm <image>"
	DEFAULT_FILESYSTEM = "ext4"
)

const (
	ID = iota
	POOL
	IMAGE
	SNAP
	DEVICE
)

type RbdBlock struct {
	Image      string
	Pool       string
	FileSystem string
	Size       int32
	Snap       string
	Device     string
}

func NewRbdBlock(image, pool, filesystem string, size int32) *RbdBlock {
	return &RbdBlock{
		Image:      image,
		Pool:       pool,
		FileSystem: filesystem,
		Size:       size,
		Snap:       "-",
		Device:     "",
	}
}

func (rb *RbdBlock) Create() error {

	// Makeup shell command
	ph := placeholder.NewPlaceHolder(RBD_CREATE)

	size := strconv.Itoa(int(rb.Size))

	cmd := ph.Replace("<image>", rb.Image, "<size>", size, "<pool>", rb.Pool)

	parts := strings.Fields(cmd)
	head := parts[0]
	parts = parts[1:len(parts)]

	// Exec shell command
	out, err := exec.Command(head, parts...).Output()

	if err != nil {
		log.Fatal(err)
		return err

	}

	fmt.Printf("%s\n", out)
	return err

}

func (rb *RbdBlock) GetMappedDevice(output []byte, image string) (pool, device string) {

	readBuf := bytes.NewBuffer(output)

	reader := bufio.NewReader(readBuf)

	// Ignor the head line: id pool image         snap device
	reader.ReadString('\n')

	for {
		str, err := reader.ReadString('\n')
		if err != nil {
			// log.Fatal(err)
			break
		}

		str = strings.Replace(str, "\n", "", -1)

		ss := strings.Split(str, " ")

		// Read output line into a slice, which contruct the rbd block
		slice := make([]string, 0)
		for _, s := range ss {
			s = strings.TrimSpace(s)
			if !strings.EqualFold(s, "") {
				slice = append(slice, s)
			}
		}

		if strings.EqualFold(slice[IMAGE], image) {
			return slice[POOL], slice[DEVICE]
		}
	}

	return "", ""

}

func (rb *RbdBlock) ShowMapped() error {

	parts := strings.Fields(RBD_SHOWMAPPED)
	head := parts[0]
	parts = parts[1:len(parts)]

	out, err := exec.Command(head, parts...).Output()

	if err != nil {
		log.Fatal(err)
		return err
	}

	pool, device := rb.GetMappedDevice(out, rb.Image)

	// if pool != rb.pool, big problem!
	if !strings.EqualFold(pool, rb.Pool) {
		log.Fatalf("RBD Problem: pool name mismatch! GetMappedDevice.pool=%s, rb.pool=%s\n", pool, rb.Pool)
		panic("RBD Problem: pool name mismatch! GetMappedDevice.pool=%s, rb.pool=%s\n")
	}

	rb.Device = device

	return nil
}

func (rb *RbdBlock) Map() error {

	// Makeup shell command
	ph := placeholder.NewPlaceHolder(RBD_MAP)

	cmd := ph.Replace("<image>", rb.Image)

	fmt.Printf("CMD: %s\n", cmd)

	// Exec unmap command
	parts := strings.Fields(cmd)
	head := parts[0]
	parts = parts[1:len(parts)]

	out, err := exec.Command(head, parts...).Output()
	if err != nil {
		log.Fatal(err)
		return err
	}

	log.Println(string(out))
	fmt.Println(string(out))

	return nil
}

func (rb *RbdBlock) UnMap() error {

	// Makeup shell command
	ph := placeholder.NewPlaceHolder(RBD_UNMAP)

	cmd := ph.Replace("<device>", rb.Device)

	// Exec unmap command
	parts := strings.Fields(cmd)
	head := parts[0]
	parts = parts[1:len(parts)]

	out, err := exec.Command(head, parts...).Output()
	if err != nil {
		log.Fatal(err)
		return err
	}

	log.Println(string(out))

	return nil
}

func (rb *RbdBlock) MakeFileSystem() error {

	// Makeup mkfs command
	ph := placeholder.NewPlaceHolder(MAKEFS)

	cmd := ph.Replace("<fs>", rb.FileSystem, "<device>", rb.Device)

	// Exec unmap command
	parts := strings.Fields(cmd)
	head := parts[0]
	parts = parts[1:len(parts)]

	out, err := exec.Command(head, parts...).Output()
	if err != nil {
		log.Fatal(err)
		return err
	}

	log.Println(string(out))

	return nil
}

func (rb *RbdBlock) Remove() error {

	// Makeup mkfs command
	ph := placeholder.NewPlaceHolder(RBD_REMOVE)

	cmd := ph.Replace("<image>", rb.Image)

	// Exec unmap command
	parts := strings.Fields(cmd)
	head := parts[0]
	parts = parts[1:len(parts)]

	out, err := exec.Command(head, parts...).Output()
	if err != nil {
		log.Fatal(err)
		return err
	}

	log.Println(string(out))
	return nil
}

func (rb *RbdBlock) DecodeJson(data string) {
	err := json.Unmarshal([]byte(data), rb)

	if err != nil {
		log.Fatal(err)
		panic(err.Error())
	}
}

func (rb *RbdBlock) EncodeJson() string {
	data, err := json.MarshalIndent(rb, "", " ")
	if err != nil {
		log.Fatal(err)
		panic(err.Error())
	}
	return string(data)
}
