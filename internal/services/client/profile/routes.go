package profile

import (
	"net/http"

	"github.com/caio-bernardo/dragonite/internal/model"
	"github.com/caio-bernardo/dragonite/internal/repository"
	"github.com/caio-bernardo/dragonite/internal/types"
	"github.com/caio-bernardo/dragonite/internal/util"
)

type Handler struct {
	userStore repository.UserStore
}

func NewHandler(userStore repository.UserStore) *Handler {
	return &Handler{userStore: userStore}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux, authMiddleware types.Middleware) {

	// Chave do perfil do usuário
	mux.HandleFunc("GET /_matrix/client/v3/profile/{userId}/keys", h.getProfileKey)

	// Alterar chave do perfil do usuário
	mux.Handle("PUT /_matrix/client/v3/profile/{userId}/keys", authMiddleware(http.HandlerFunc(h.putProfileKey)))

	// Remover chave do perfil do usuário
	mux.Handle("DELETE /_matrix/client/v3/profile/{userId}/keys", authMiddleware(http.HandlerFunc(h.deleteProfileKey)))

}

// GET /_matrix/client/v3/profile/{userId}
func (h *Handler) getProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userId")

	usuario, err := h.userStore.GetNameAndPhotoByID(r.Context(), userID)

	if err != nil {
		if err == types.ErrNotFound {
			util.WriteError(w, http.StatusNotFound, types.ErrorResponse{
				ErrCode: "M_NOT_FOUND",
				Message: "Profile not found",
			})
			return
		}
		util.WriteError(w, http.StatusInternalServerError, types.ErrorResponse{
			ErrCode: "M_UNKNOWN",
			Message: "Internal server error",
		})
		return
	}

	response := model.ProfileResponse{
		DisplayName: usuario.Nome,
		AvatarURL:   usuario.Foto,
	}

	util.WriteJSON(w, http.StatusOK, response)
}

// GET /_matrix/client/v3/profile/{userId}/{keyName}
func (h *Handler) getProfileKey(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userId")
	keyName := r.PathValue("keyName") // Vai ser "displayname" ou "avatar_url"

	usuario, err := h.userStore.GetNameAndPhotoByID(r.Context(), userID)

	if err != nil {
		if err == types.ErrNotFound {
			util.WriteError(w, http.StatusNotFound, types.ErrorResponse{
				ErrCode: "M_NOT_FOUND",
				Message: "User not found",
			})
			return
		}
		util.WriteError(w, http.StatusInternalServerError, types.ErrorResponse{
			ErrCode: "M_UNKNOWN",
			Message: "Internal server error",
		})
		return
	}

	var valor string
	if keyName == "displayname" {
		valor = usuario.Nome
	} else if keyName == "avatar_url" {
		valor = usuario.Foto
	} else {
		// Se pediu uma chave que não existe no Matrix
		util.WriteError(w, http.StatusBadRequest, types.ErrorResponse{
			ErrCode: "M_BAD_JSON",
			Message: "Invalid profile key",
		})
		return
	}

	response := map[string]string{
		keyName: valor,
	}

	util.WriteJSON(w, http.StatusOK, response)
}

// PUT /_matrix/client/v3/profile/{userId}/{keyName}
func (h *Handler) putProfileKey(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userId")
	keyName := r.PathValue("keyName") // "displayname" ou "avatar_url"

	colunaDB := mapMatrixKeyToDB(keyName)
	if colunaDB == "" {
		util.WriteError(w, http.StatusBadRequest, types.ErrorResponse{
			ErrCode: "M_BAD_JSON",
			Message: "Invalid profile key",
		})
		return
	}

	var reqBody map[string]string
	if err := util.ParseBody(r, &reqBody); err != nil {
		util.WriteError(w, http.StatusBadRequest, types.ErrorResponse{
			ErrCode: "M_NOT_JSON",
			Message: "Request body must be JSON",
		})
		return
	}

	novoValor, existe := reqBody[keyName]
	if !existe {
		util.WriteError(w, http.StatusBadRequest, types.ErrorResponse{
			ErrCode: "M_BAD_JSON",
			Message: "Missing key in request body",
		})
		return
	}

	err := h.userStore.UpdateProfileKey(r.Context(), userID, colunaDB, novoValor)
	if err != nil {
		if err == types.ErrNotFound {
			util.WriteError(w, http.StatusNotFound, types.ErrorResponse{
				ErrCode: "M_NOT_FOUND",
				Message: "User not found",
			})
			return
		}

		util.WriteError(w, http.StatusInternalServerError, types.ErrorResponse{
			ErrCode: "M_UNKNOWN",
			Message: "Database error",
		})
		return
	}

	util.WriteJSON(w, http.StatusOK, map[string]string{})
}

// DELETE /_matrix/client/v3/profile/{userId}/{keyName}
func (h *Handler) deleteProfileKey(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userId")
	keyName := r.PathValue("keyName")

	colunaDB := mapMatrixKeyToDB(keyName)
	if colunaDB == "" {
		util.WriteError(w, http.StatusBadRequest, types.ErrorResponse{
			ErrCode: "M_BAD_JSON",
			Message: "Invalid profile key",
		})
		return
	}

	err := h.userStore.ClearProfileKey(r.Context(), userID, colunaDB)

	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, types.ErrorResponse{
			ErrCode: "M_UNKNOWN",
			Message: "Database error",
		})
		return
	}

	util.WriteJSON(w, http.StatusOK, map[string]string{})
}

// mapMatrixKeyToDB converte a chave do Matrix para a coluna do bd.
// Isso evita SQL Injection, pois não usamos a string do usuário direto na query.
func mapMatrixKeyToDB(keyName string) string {
	switch keyName {
	case "displayname":
		return "nome_usuario"
	case "avatar_url":
		return "foto_usuario"
	default:
		return ""
	}
}
