package rbd

import (
	"fmt"
	"strconv"
	"app/backend/common/util/placeholder"
	"strings"
	"log"
	"os"
	"os/exec"
)

const (
	RBD_CREATE = "rbd create <image> -s <size> -p <pool>"

	RBD_SHOWMAPPED = "rbd show mapped"

	RBD_UNMAP = "rbd unmap <path>"

	MAKEFS = "mkfs.<fs> <path>"

	DEFAULT_FILESYSTEM = "ext4"

)

type RbdBlock struct {
	Name string
	Pool string
	FileSystem string
	Size int32
}

func NewRbdBlock(name, pool, filesystem string, size int32) *RbdBlock {
	return &RbdBlock{
		Name: name,
		Pool: pool,
		FileSystem: filesystem,
		Size: size,
	}
}

func (rb *RbdBlock) Create() error {

	// Makeup shell command
	ph := placeholder.NewPlaceHolder(RBD_CREATE)

	size := strconv.Itoa(rb.Size)

	cmd := ph.Replace("<image>", rb.Name, "<size>", size, "<pool>", rb.Pool)

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

func (rb *RbdBlock) ShowMapped(rbd string) (string, error) {

	parts := strings.Fields(RBD_SHOWMAPPED)
	head := parts[0]
	parts = parts[1:len(parts)]

	out, err := exec.Command(head, parts...)

	if err != nil {
		log.Fatal(err)
		return "", err
	}
	 str := string(out.Output())
	// ToDo: grep the mapped path
}

