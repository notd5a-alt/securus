package main

import (
	"bufio"
	"context"
	"crypto/md5"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	mrand "math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/notd5a-alt/securus/chat"
	secui "github.com/notd5a-alt/securus/ui"
)

type User struct {
	id     uint64
	handle string
}

type Connection struct {
	Port        int
	Destination string
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
	mrand.Seed(int64(seed))
	// assign id to user
	u.id = mrand.Uint64()
	return nil
}

func handleInput(prefix string) string {
	// create reader and read input into line
	secui.PrintInputPrefix(prefix)
	reader := bufio.NewReader(os.Stdin)
	str, err := reader.ReadString('\n')
	if err != nil {
		// return an error
		return ""
	}
	fmt.Println()

	str = strings.TrimSuffix(str, "\n")
	return str

}

// anything we need to initalize before the main function goes here
func init() {
	secui.Title()

}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var user User

	// create reader to handle input
	secui.PrintSection("User Details")
	handle := handleInput("Handle [#] ")
	if handle == "" {
		secui.PrintError(errors.New("invalid input"))
		return
	}

	user.handle = handle
	err := user.initID()
	if err != nil {
		// terminate program
		secui.PrintError(err)
		os.Exit(1)
	}
	secui.PrintSpinnerSuccess("Generating ID")
	secui.PrintSuccess(fmt.Sprintf("ID: %d", user.id))

	// User details stored, TODO: may want to extend this so that it stores the details into
	// a file for auto login.

	// New section for connection details
	var connection Connection

	secui.PrintSection("Connection")
	secui.PrintInfo("1. Host")
	secui.PrintInfo("2. Connect")
	choice, err := strconv.Atoi(handleInput("Choice [#] "))
	if choice == 0 {
		secui.PrintWarn("Invalid Input")
		secui.PrintError(err)
		return
	}

	// iniialize peer id
	var reader io.Reader = rand.Reader

	connection.Port = 0 // just make sure its set to its 0 value
	host, err := chat.MakeHost(connection.Port, ctx, reader)
	if err != nil {
		secui.PrintError(err)
		return
	}

	switch choice {
	case 1:
		// Host
		secui.PrintSection("Host Information")
		connection.Port, err = strconv.Atoi(handleInput("Source Port [#] "))
		if err != nil {
			secui.PrintWarn(fmt.Sprintf("%d", connection.Port))
			secui.PrintError(err)
		}

		chat.StartPeer(ctx, host, chat.HandleStream)

	case 2:
		// Connect
		secui.PrintSection("Connection 	Information")
		connection.Destination = handleInput("MultiAddr Destination [#] ")
		rw, err := chat.StartPeerAndConnect(ctx, host, connection.Destination)
		if err != nil {
			secui.PrintError(err)
			return
		}

		go chat.WriteData(rw)
		go chat.ReadData(rw)

	default:
		// Host
		secui.PrintSection("Host Information")
		connection.Port, err = strconv.Atoi(handleInput("Source Port [#] "))
		if err != nil {
			secui.PrintWarn(fmt.Sprintf("%d", connection.Port))
			secui.PrintError(err)
		}

		chat.StartPeer(ctx, host, chat.HandleStream)
	}

	// wait forever
	select {}

}
