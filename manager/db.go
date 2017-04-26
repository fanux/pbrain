package manager

import (
	"encoding/json"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"
)

//PluginDB is
type PluginDB struct {
	*leveldb.DB
}

//Dber is
type Dber interface {
	Save(*PluginStrategy)
	Load(string) *PluginStrategy
	LoadAll() map[string]*PluginStrategy
	Delete(string)
}

//Save is
func (db *PluginDB) Save(p *PluginStrategy) {
	key := []byte(p.Name)
	value, err := json.Marshal(p)
	if err != nil {
		logrus.Errorf("Plugin json marshal error: %s", err)
	}

	err = db.Put(key, value, nil)
	if err != nil {
		fmt.Println(err)
	}
}

//Load is
func (db *PluginDB) Load(name string) *PluginStrategy {
	key := []byte(name)
	data, err := db.Get(key, nil)
	if err != nil {
		logrus.Errorf("Load error: %s", err)
	}
	p := &PluginStrategy{}
	err = json.Unmarshal(data, p)
	if err != nil {
		logrus.Errorf("Load plugin json error: %s", err)
		return nil
	}

	return p
}

//LoadAll is
func (db *PluginDB) LoadAll() (m map[string]*PluginStrategy) {
	m = make(map[string]*PluginStrategy)

	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		// Remember that the contents of the returned slice should not be modified, and
		// only valid until the next call to Next.
		key := string(iter.Key())
		value := iter.Value()

		p := &PluginStrategy{}
		err := json.Unmarshal(value, p)
		if err != nil {
			logrus.Errorf("unmarsal plugin json faied! %s", err)
			continue
		}

		m[key] = p
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		logrus.Errorf("%s", err)
	}

	return
}

//Delete is
func (db *PluginDB) Delete(name string) {
	key := []byte(name)
	err := db.DB.Delete(key, nil)
	if err != nil {
		logrus.Errorf("%s", err)
	}
}

//NewDb is
func NewDb(name string, path string) Dber {
	if path == "" {
		path = "./db"
	}
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		logrus.Errorf("%s", err)
	}

	return &PluginDB{db}
}
