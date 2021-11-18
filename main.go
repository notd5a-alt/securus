package main

import (
	"bufio"
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"

	secui "github.com/notd5a-alt/securus/ui"
)

type User struct {
	id     uint64
	handle string
}

type Message struct {
	id        uint32
	user      User
	timestamp string
	contents  string
}

func (u *User) initID() (err error) {
	// md5 hash
	if u.handle == "" {
		// empty handle so return an error
		return err
	}

	h := md5.New()
	// md5 hash of user handle
	io.WriteString(h, u.handle)
	// generate seed according to md5 sum
	var seed uint64 = binary.BigEndian.Uint64(h.Sum(nil))
	rand.Seed(int64(seed))
	// assign id to user
	u.id = rand.Uint64()
	return nil
}

func handleInput(prefix string) (str string, err error) {
	// create reader and read input into line
	fmt.Print(prefix)
	reader := bufio.NewReader(os.Stdin)
	str, err = reader.ReadString('\n')
	if err != nil {
		// return an error
		return "", err
	}
	fmt.Println()

	str = strings.TrimSuffix(str, "\n")
	return str, nil

}

func main() {
	// Securus interface.
	secui.Title()
	var user User

	// create reader to handle input
	fmt.Println("username:")
	handle, err := handleInput("# ")
	if err != nil || handle == "" {
		return
	}

	user.handle = handle
	err = user.initID()
	if err != nil {
		// terminate program
		log.Println(err)
		os.Exit(1)
	}

	fmt.Println("ID Generated: ", user.id)
}
