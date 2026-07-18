package usecase

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/caio-bernardo/dragonite/internal/domain"
	"github.com/caio-bernardo/dragonite/internal/domain/types"
	"github.com/caio-bernardo/dragonite/internal/util"
)

// KeysService contém a lógica de gerenciamento de chaves de dispositivo E2EE (device keys,
// one-time keys, fallback keys)
type KeysService struct {
	uow        WorkUnit
	keysStore  KeysStorage
	federation *FederationService
	serverName string
}

func NewKeysService(uow WorkUnit, keysStore KeysStorage, federation *FederationService, serverName string) *KeysService {
	return &KeysService{uow: uow, keysStore: keysStore, federation: federation, serverName: serverName}
}

// UploadKeysParams agrupa os dados recebidos em POST /_matrix/client/v3/keys/upload
type UploadKeysParams struct {
	UserID       string
	DeviceID     string
	Algorithms   []string
	IdentityKeys json.RawMessage
	Signatures   json.RawMessage
	OneTimeKeys  map[string]json.RawMessage
	FallbackKeys map[string]json.RawMessage
}

// UploadKeys persiste as chaves de identidade, one-time keys e fallback keys de um dispositivo
// numa única transação, retornando a contagem de one-time keys remanescentes por algoritmo
func (s *KeysService) UploadKeys(ctx context.Context, params UploadKeysParams) (map[string]int, error) {
	var counts map[string]int

	err := s.uow.Execute(ctx, func(txCtx context.Context) error {
		if len(params.Algorithms) > 0 {
			err := s.keysStore.UpsertDeviceKeys(txCtx, domain.ChavesDispositivo{
				DispositivoID: params.DeviceID,
				UsuarioID:     params.UserID,
				Algorithms:    params.Algorithms,
				Keys:          params.IdentityKeys,
				Signatures:    params.Signatures,
			})
			if err != nil {
				return err
			}
		}

		if len(params.OneTimeKeys) > 0 {
			otks := make([]domain.ChaveUsoUnico, 0, len(params.OneTimeKeys))
			for keyID, data := range params.OneTimeKeys {
				algorithm, _, _ := strings.Cut(keyID, ":")
				otks = append(otks, domain.ChaveUsoUnico{
					DispositivoID: params.DeviceID,
					KeyID:         keyID,
					Algorithm:     algorithm,
					KeyData:       data,
				})
			}
			if err := s.keysStore.UpsertOneTimeKeys(txCtx, params.DeviceID, otks); err != nil {
				return err
			}
		}

		for keyID, data := range params.FallbackKeys {
			algorithm, _, _ := strings.Cut(keyID, ":")
			err := s.keysStore.UpsertFallbackKey(txCtx, domain.ChaveFallback{
				DispositivoID: params.DeviceID,
				Algorithm:     algorithm,
				KeyID:         keyID,
				KeyData:       data,
				Usada:         false,
			})
			if err != nil {
				return err
			}
		}

		var err error
		counts, err = s.keysStore.CountOneTimeKeys(txCtx, params.DeviceID)
		return err
	})
	if err != nil {
		return nil, types.InternalError(err)
	}

	return counts, nil
}

// QueryKeysResult é o resultado agregado (local + federado) de uma consulta de chaves
type QueryKeysResult struct {
	DeviceKeys map[string]map[string]domain.ChavesDispositivo // userID -> deviceID -> chaves
	Failures   map[string]any                                 // servidor -> motivo
}

// QueryKeys busca as chaves de identidade dos dispositivos pedidos, roteando para o storage
// local ou pra federação de acordo com o domínio de cada userID
func (s *KeysService) QueryKeys(ctx context.Context, requested map[string][]string) QueryKeysResult {
	result := QueryKeysResult{
		DeviceKeys: make(map[string]map[string]domain.ChavesDispositivo),
		Failures:   make(map[string]any),
	}

	remoteByServer := make(map[string]map[string][]string)

	for userID, deviceIDs := range requested {
		domainName := util.ExtractDomainFromUserID(userID)
		if domainName == s.serverName {
			keys, err := s.keysStore.GetDeviceKeys(ctx, userID, deviceIDs)
			if err != nil {
				log.Printf("[ERROR] QueryKeys (local user=%s): %v", userID, err)
				continue
			}
			if len(keys) == 0 {
				continue
			}
			devMap := make(map[string]domain.ChavesDispositivo, len(keys))
			for _, k := range keys {
				devMap[k.DispositivoID] = k
			}
			result.DeviceKeys[userID] = devMap
			continue
		}

		if remoteByServer[domainName] == nil {
			remoteByServer[domainName] = make(map[string][]string)
		}
		remoteByServer[domainName][userID] = deviceIDs
	}

	// Uma requisição de federação por servidor remoto envolvido
	for server, deviceKeysReq := range remoteByServer {
		remoteResult, err := s.federation.QueryKeysCall(ctx, server, deviceKeysReq)
		if err != nil {
			log.Printf("[ERROR] QueryKeys (remote server=%s): %v", server, err)
			result.Failures[server] = map[string]any{}
			continue
		}
		for userID, devices := range remoteResult {
			devMap := make(map[string]domain.ChavesDispositivo, len(devices))
			for deviceID, raw := range devices {
				var parsed struct {
					Algorithms []string        `json:"algorithms"`
					Keys       json.RawMessage `json:"keys"`
					Signatures json.RawMessage `json:"signatures"`
				}
				if err := json.Unmarshal(raw, &parsed); err != nil {
					continue
				}
				devMap[deviceID] = domain.ChavesDispositivo{
					DispositivoID: deviceID,
					UsuarioID:     userID,
					Algorithms:    parsed.Algorithms,
					Keys:          parsed.Keys,
					Signatures:    parsed.Signatures,
				}
			}
			result.DeviceKeys[userID] = devMap
		}
	}

	return result
}

// ClaimKeysResult é o resultado agregado (local + federado) de uma reivindicação de one-time keys
type ClaimKeysResult struct {
	OneTimeKeys map[string]map[string]map[string]json.RawMessage // userID -> deviceID -> keyID -> KeyObject
	Failures    map[string]any
}

// ClaimKeys reivindica one-time keys (ou fallback, se as OTKs se esgotaram) dos dispositivos pedidos
func (s *KeysService) ClaimKeys(ctx context.Context, requested map[string]map[string]string) ClaimKeysResult {
	result := ClaimKeysResult{
		OneTimeKeys: make(map[string]map[string]map[string]json.RawMessage),
		Failures:    make(map[string]any),
	}

	remoteByServer := make(map[string]map[string]map[string]string)

	for userID, devices := range requested {
		domainName := util.ExtractDomainFromUserID(userID)
		if domainName == s.serverName {
			for deviceID, algorithm := range devices {
				claimed, err := s.keysStore.ClaimOneTimeKey(ctx, deviceID, algorithm)
				if err != nil {
					log.Printf("[ERROR] ClaimKeys (local device=%s): %v", deviceID, err)
					continue
				}
				if claimed == nil {
					// Sem OTKs sobrando, tenta a fallback key
					fallback, err := s.keysStore.ClaimFallbackKey(ctx, deviceID, algorithm)
					if err != nil || fallback == nil {
						continue
					}
					s.addClaimedKey(result.OneTimeKeys, userID, deviceID, fallback.KeyID, fallback.KeyData)
					continue
				}
				s.addClaimedKey(result.OneTimeKeys, userID, deviceID, claimed.KeyID, claimed.KeyData)
			}
			continue
		}

		if remoteByServer[domainName] == nil {
			remoteByServer[domainName] = make(map[string]map[string]string)
		}
		remoteByServer[domainName][userID] = devices
	}

	for server, oneTimeKeysReq := range remoteByServer {
		remoteResult, err := s.federation.ClaimKeysCall(ctx, server, oneTimeKeysReq)
		if err != nil {
			log.Printf("[ERROR] ClaimKeys (remote server=%s): %v", server, err)
			result.Failures[server] = map[string]any{}
			continue
		}
		for userID, devices := range remoteResult {
			if result.OneTimeKeys[userID] == nil {
				result.OneTimeKeys[userID] = make(map[string]map[string]json.RawMessage)
			}
			for deviceID, keys := range devices {
				result.OneTimeKeys[userID][deviceID] = keys
			}
		}
	}

	return result
}

func (s *KeysService) addClaimedKey(dst map[string]map[string]map[string]json.RawMessage, userID, deviceID, keyID string, data json.RawMessage) {
	if dst[userID] == nil {
		dst[userID] = make(map[string]map[string]json.RawMessage)
	}
	if dst[userID][deviceID] == nil {
		dst[userID][deviceID] = make(map[string]json.RawMessage)
	}
	dst[userID][deviceID][keyID] = data
}