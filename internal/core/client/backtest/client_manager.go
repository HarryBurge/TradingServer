package backtest

import "math/rand"

type BacktestClientManager struct {
	clients map[uint32]*BacktestClient
}

func (t *BacktestClientManager) generateClientId() uint32 {
	id := rand.Uint32()
	for _, ok := t.clients[id]; ok; {
		id = rand.Uint32()
	}
	return id
}

func (t *BacktestClientManager) AddClient(client *BacktestClient) uint32 {
	id := t.generateClientId()
	t.clients[id] = client
	return id
}

func (t *BacktestClientManager) GetClient(id uint32) (*BacktestClient, error) {
	client, ok := t.clients[id]
	if !ok {
		return nil, IdNotFoundError{}
	}
	return client, nil
}

type ClientManagerOpts struct {
}

func NewClientManager(opts ClientManagerOpts) *BacktestClientManager {
	return &BacktestClientManager{
		clients: make(map[uint32]*BacktestClient),
	}
}
