package seriliaze

//seriliaze obj to byte arr

type Persistence interface {
	Seriliazer
	Deriliazer
}

type Seriliazer interface {
	Seriliaze(obj interface{}) ([]byte, error)
}

type Deriliazer interface {
	DeSeriliaze(bytesArr []byte, vPointer interface{}) error
}

func JsonSeriliaze(obj interface{}) ([]byte, error) {
	j := JsonPersistence{}
	data, err := j.Seriliaze(obj)
	return data, err
}

func JsonDeriliaze(byteArr []byte, vPointer interface{}) error {
	j := JsonPersistence{}
	err := j.Deriliaze(byteArr, vPointer)
	return err
}
