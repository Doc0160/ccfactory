package storage

import (
	"ccfactory/server/factory"
	"encoding/json"
	"fmt"
)

type ChestConfig struct {
	Client  string
	InvAddr string
	BusAddr string
}

func (s *ChestConfig) Build(f *factory.Factory) factory.Storage {
	return &Chest{
		config:  s,
		factory: f,
	}
}

type Chest struct {
	config  *ChestConfig
	factory *factory.Factory
}

func (s *Chest) Size() (float64, error) {
	call := s.factory.Call(s.config.Client,
		factory.ActionPeripheralCall,
		s.config.InvAddr,
		"size")
	if call.Error != nil {
		return 0, fmt.Errorf("Error size: %#+v", call.Error)
	}
	return call.Result.([]any)[0].(float64), nil
}

func (s *Chest) GetItemDetail(slot int) (*ItemDetails, error) {
	id := ItemDetails{}
	call := s.factory.Call(s.config.Client,
		factory.ActionPeripheralCall,
		s.config.InvAddr,
		"getItemDetail",
		slot)
	if call.Error != nil {
		return nil, fmt.Errorf("Error size: %#+v", call.Error)
	}
	log.Info(call)
	_, ok := call.Result.([]any)
	if !ok {
		return nil, nil
	}
	bytes, err := json.Marshal(call.Result.([]any)[0])
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bytes, &id)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

type ItemDetails struct {
	Count       int    `json:"count"`
	DisplayName string `json:"displayName"`
	ItemGroups  []struct {
		DisplayName string `json:"displayName"`
		ID          string `json:"id"`
	} `json:"itemGroups"`
	MaxCount int             `json:"maxCount"`
	Name     string          `json:"name"`
	Tags     map[string]bool `json:"tags"`
}

func (s *Chest) Update() {
	if !s.factory.ClientConnected(s.config.Client) {
		log.Warn("Client missing", "name", s.config.Client)
		return
	}

	/*size, err := s.Size()
	if err != nil {
		log.Error(err)
	}
	log.Debug("", "size", size)

	id, err := s.GetItemDetail(1)
	if err != nil {
		log.Error(err)
	}
	log.Debug("", "item", id)*/

	call := s.factory.Call(s.config.Client,
		factory.ActionPeripheralCall,
		s.config.InvAddr,
		"list")
	log.Debug("", "item", call)

}
