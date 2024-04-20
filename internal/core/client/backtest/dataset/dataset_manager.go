package dataset

import "math/rand"

// type DatasetManager interface {
// 	AddDataset(dataset *Dataset) uint32
// 	GetDataset(id uint32) (*Dataset, error)
// }

type DatasetManager struct {
	datasets map[uint32]*Dataset
}

func (t *DatasetManager) generateDatasetId() uint32 {
	id := rand.Uint32()
	for _, ok := t.datasets[id]; ok; {
		id = rand.Uint32()
	}
	return id
}

func (t *DatasetManager) AddDataset(dataset *Dataset) uint32 {
	id := t.generateDatasetId()
	t.datasets[id] = dataset
	return id
}

func (t *DatasetManager) GetDataset(id uint32) (*Dataset, error) {
	dataset, ok := t.datasets[id]
	if !ok {
		return nil, IdNotFoundError{}
	}
	return dataset, nil
}

type DatasetManagerOpts struct {
}

func NewDatasetManager(opts DatasetManagerOpts) *DatasetManager {
	return &DatasetManager{
		datasets: make(map[uint32]*Dataset),
	}
}
