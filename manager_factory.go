package bgc

import (
	"encoding/json"
	"fmt"
	"github.com/viant/dsc"
	"github.com/viant/toolbox"
)

type managerFactory struct{}

func (f *managerFactory) Create(config *dsc.Config) (dsc.Manager, error) {
	var connectionProvider = newConnectionProvider(config)
	manager := &manager{}
	var self dsc.Manager = manager
	super := dsc.NewAbstractManager(config, connectionProvider, self)
	manager.AbstractManager = super

	for _, key := range []string{ProjectIDKey, DataSetIDKey} {
		if !config.Has(key) {
			return nil, fmt.Errorf("config.parameters.%v was missing", key)
		}
	}

	manager.bigQueryConfig = newConfig(config)
	return self, nil
}

func (f managerFactory) CreateFromURL(url string) (dsc.Manager, error) {
	reader, _, err := toolbox.OpenReaderFromURL(url)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	config := &dsc.Config{}
	err = json.NewDecoder(reader).Decode(config)
	if err != nil {
		return nil, err
	}
	config.Init()
	return f.Create(config)
}

func newManagerFactory() dsc.ManagerFactory {
	var result dsc.ManagerFactory = &managerFactory{}
	return result
}
