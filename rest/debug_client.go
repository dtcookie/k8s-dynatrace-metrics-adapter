package rest

type debugClient struct {
	client Client
}

func (c *debugClient) GET(path string, expectedStatusCode int) ([]byte, error) {
	// if path == "" {
	// klog.Info("GET ", path)
	// }
	data, err := c.client.GET(path, expectedStatusCode)
	// fmt.Println("  ", string(data))
	return data, err
}
