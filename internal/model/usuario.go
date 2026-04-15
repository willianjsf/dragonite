package model

import "time"

type Usuario struct {
	ID          string    `json:"id_usuario"`
	Localpart   string    `json:"localpart_usuario"`
	Nome        string    `json:"nome_usuario"`
	Senha       string    `json:"senha_usuario"`
	Foto        string    `json:"foto_usuario"`
	DataCriacao time.Time `json:"data_criacao_usuario"`
}

type UsuarioCreate struct {
	Localpart   string    `json:"localpart_usuario"`
	Nome        string    `json:"nome_usuario"`
	Senha       string    `json:"senha_usuario"`
	Foto        string    `json:"foto_usuario"`
	DataCriacao time.Time `json:"data_criacao_usuario"`
}

func (uc UsuarioCreate) ToUsuario() Usuario {
	return Usuario{
		Localpart:   uc.Localpart,
		Nome:        uc.Nome,
		Senha:       uc.Senha,
		Foto:        uc.Foto,
		DataCriacao: uc.DataCriacao,
	}
}
