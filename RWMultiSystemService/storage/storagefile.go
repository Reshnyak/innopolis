package storage

//Instance of file storage
type StorageFile struct {
	config   *Config
	fileName string
}

func (sf *StorageFile) WriteData(fName string, data []byte) error {
	return nil
}
func (sf *StorageFile) ReadData(fName string) ([]byte, error) {
	var data []byte
	return data, nil
}
