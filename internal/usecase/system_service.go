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

func (h *SystemService) PingDB() map[string]string {
	return h.storage.PingDB()
}

func (h *SystemService) GetServerName() string {
	return h.serverName
}

func (h *SystemService) GetServerVersion() string {
	return h.serverVersion
}

func (h *SystemService) GetPublicKey() []byte {
	return h.serverPublicKey
}

func (h *SystemService) GetPrivateKey() []byte {
	return h.serverPrivateKey
}

func (h *SystemService) GetServerKeyID() string {
	return h.serverKeyID
}

func (h *SystemService) GetServerPrivateKey() []byte {
	return h.serverPrivateKey
}

func (h *SystemService) GetServerPublicKey() []byte {
	return h.serverPublicKey
}
