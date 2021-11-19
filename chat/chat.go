package chat

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	peerstore "github.com/libp2p/go-libp2p-core/peerstore"
	"github.com/multiformats/go-multiaddr"
	secui "github.com/notd5a-alt/securus/ui"
)

func HandleStream(s network.Stream) {
	secui.PrintInfo("Got a new Stream")

	// create a buffer stream for nn blocking read and write
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	go ReadData(rw)
	go WriteData(rw)

	// stream s stays open till either client closes it
}

// reads a string froom the network stream and prints it out to console
func ReadData(rw *bufio.ReadWriter) {
	for {
		str, err := rw.ReadString('\n')
		if err != nil {
			fmt.Println(err)
		}

		if str == "" {
			return
		}

		if str != "\n" {
			fmt.Printf("\x1b[32m%s\x1b[0m", str) // TODO: Change to use updated UI
		}

	}
}

// reads data from stdin and writes it to network stream
func WriteData(rw *bufio.ReadWriter) {
	// create a stdin reader
	stdinReader := bufio.NewReader(os.Stdin)

	secui.PrintInputPrefix(fmt.Sprintf("[#] "))
	for {
		sendData, err := stdinReader.ReadString('\n')
		if err != nil {
			secui.PrintError(err)
			return
		}

		secui.PrintInputPrefix("[#] ")
		rw.WriteString(fmt.Sprintf("%s\n", sendData))
		rw.Flush()
	}

}

// constructs a new libp2p host with a new RSA key pair.
func MakeHost(port int, ctx context.Context, randomness io.Reader) (host.Host, error) {
	privateKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, randomness)
	if err != nil {
		secui.PrintError(err)
		return nil, err
	}

	// 0.0.0 will listen on any interface device
	sourceMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port))

	// make a new libp2p host
	return libp2p.New(ctx,
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.Identity(privateKey),
	)
}

func StartPeer(ctx context.Context, h host.Host, streamHandler network.StreamHandler) {
	// Set a function as stream handler.
	// This function is called when a peer connects, and starts a stream with this protocol.
	// Only applies on the receiving side.
	h.SetStreamHandler("/chat/1.0.0", streamHandler)

	// get tcp port from our listen multiaddr
	var port string
	for _, la := range h.Network().ListenAddresses() {
		if p, err := la.ValueForProtocol(multiaddr.P_TCP); err == nil {
			port = p
			break
		}
	}

	if port == "" {
		secui.PrintWarn("Was not able to find the Actual Port")
		return
	}

	secui.PrintSpinnerSuccess("Generating Multi Address.")
	secui.PrintSuccess(fmt.Sprintf("MultiAddr: /ip4/127.0.0.1/tcp/%v/p2p/%s", port, h.ID().Pretty()))
	secui.PrintInfo("You can also replace 127.0.0.1 with a public IP address")

	// TODO: custom pterm spinner waitng for user to connect to host.
	secui.PrintSuccess("Waiting for incoming connection...") // << placeholder
	fmt.Println()

}

func StartPeerAndConnect(ctx context.Context, h host.Host, destination string) (*bufio.ReadWriter, error) {
	maddr, err := multiaddr.NewMultiaddr(destination)
	if err != nil {
		secui.PrintError(err)
		return nil, err
	}

	info, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		secui.PrintError(err)
		return nil, err
	}

	h.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)

	s, err := h.NewStream(context.Background(), info.ID, "/chat/1.0.0")
	if err != nil {
		secui.PrintError(err)
		return nil, err
	}

	secui.PrintInfo("Established connection to destination")
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	return rw, nil
}
