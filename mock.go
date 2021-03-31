package configuration

type mockClient struct {
	store map[string]string
}

func NewMockClient(store map[string]string) (*mockClient, error) {
	return &mockClient{store: store}, nil
}

func (c *mockClient) GetValues(keys []string) (vls map[string]string, err error) {
	vls = make(map[string]string, len(keys))
	for _, v := range keys {
		vls[v] = c.store[v]
	}
	return
}
func (c *mockClient) WatchPrefix(keys []string, waitIndex uint64) (index uint64, err error) {
	return
}
