package room

import (
	"sync"

	"github.com/cloudwebrtc/go-protoo/logger"
	"github.com/cloudwebrtc/go-protoo/peer"
	"github.com/cloudwebrtc/go-protoo/transport"
)

type Room struct {
	*sync.Mutex
	peers  map[string]*peer.Peer
	closed bool
	id     string
}

func NewRoom(roomId string) *Room {
	room := &Room{
		peers:  make(map[string]*peer.Peer),
		closed: false,
		id:     roomId,
	}
	room.Mutex = new(sync.Mutex)
	return room
}

func (room *Room) CreatePeer(peerId string, transport *transport.WebSocketTransport) *peer.Peer {
	newPeer := peer.NewPeer(peerId, transport)
	newPeer.On("close", func() {
		room.Lock()
		defer room.Unlock()
		delete(room.peers, peerId)
	})
	room.Lock()
	defer room.Unlock()
	room.peers[peerId] = newPeer
	return newPeer
}

func (room *Room) GetPeer(peerId string) *peer.Peer {
	room.Lock()
	defer room.Unlock()
	if peer, ok := room.peers[peerId]; ok {
		return peer
	}
	return nil
}

func (room *Room) GetPeers() map[string]*peer.Peer {
	return room.peers
}

func (room *Room) ID() string {
	return room.id
}

func (room *Room) HasPeer(peerId string) bool {
	room.Lock()
	defer room.Unlock()
	_, ok := room.peers[peerId]
	return ok
}

func (room *Room) Notify(from *peer.Peer, method string, data map[string]interface{}) {
	room.Lock()
	defer room.Unlock()
	for id, peer := range room.peers {
		//send to other peers
		if id != from.ID() {
			peer.Notify(method, data)
		}
	}
}

func (room *Room) Close() {
	logger.Warnf("Close all peers !")
	room.Lock()
	defer room.Unlock()
	for id, peer := range room.peers {
		logger.Warnf("Close => peer(%s).", id)
		peer.Close()
	}
	room.closed = true
}
