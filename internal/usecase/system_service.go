package usecase

type SystemService struct {
	serverName       string
	serverVersion    string
	serverPublicKey  []byte
	serverPrivateKey []byte
	serverKeyID      string
	storage          SystemStorage
}

func NewSystemService(serverName, serverVersion string, serverPublicKey, serverPrivateKey []byte, serverKeyID string, storage SystemStorage) *SystemService {
	return &SystemService{
		serverName:       serverName,
		serverVersion:    serverVersion,
		serverPublicKey:  serverPublicKey,
		serverPrivateKey: serverPrivateKey,
		serverKeyID:      serverKeyID,
		storage:          storage,
	}
}

func (s *SystemService) PingDB() map[string]string {
	return s.storage.PingDB()
}

func (s *SystemService) GetServerName() string {
	return s.serverName
}

func (s *SystemService) GetServerVersion() string {
	return s.serverVersion
}

func (s *SystemService) GetPublicKey() []byte {
	return s.serverPublicKey
}

func (s *SystemService) GetPrivateKey() []byte {
	return s.serverPrivateKey
}

func (s *SystemService) GetServerKeyID() string {
	return s.serverKeyID
}

func (s *SystemService) GetServerPrivateKey() []byte {
	return s.serverPrivateKey
}

func (s *SystemService) GetServerPublicKey() []byte {
	return s.serverPublicKey
}
