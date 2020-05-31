package ring

import (
	"fmt"
	"hash/adler32"
	"strings"
)

type Node struct {
	Host string
	Port string
	Disk string
}

func (n *Node) String() string {
	return n.Host + ":" + n.Port + ":" + n.Disk
}

func (n *Node) Url() string {
	return "http://" + n.Host + ":" + n.Port + "/obj" + "/" + fmt.Sprintf("%d", adler32.Checksum([]byte(n.Disk))) + "/"
}

func (n *Node) Unmarshal(nodeS string) {
	part := strings.Split(nodeS, ":")

	n.Host = part[0]
	n.Port = part[1]
	n.Disk = part[2]
}
