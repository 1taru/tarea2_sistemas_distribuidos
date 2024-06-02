package model

import (
	"time"
)

// Solicitud defines the structure of a request
type Solicitud struct {
	ID                         string    `json:"id"`
	Fecha                      string    `json:"fecha"`
	LocalID                    string    `json:"local_id"`
	LocalNombre                string    `json:"local_nombre"`
	ComunaNombre               string    `json:"comuna_nombre"`
	LocalidadNombre            string    `json:"localidad_nombre"`
	LocalDireccion             string    `json:"local_direccion"`
	FuncionamientoHoraApertura string    `json:"funcionamiento_hora_apertura"`
	FuncionamientoHoraCierre   string    `json:"funcionamiento_hora_cierre"`
	LocalTelefono              string    `json:"local_telefono"`
	LocalLat                   string    `json:"local_lat"`
	LocalLng                   string    `json:"local_lng"`
	FuncionamientoDia          string    `json:"funcionamiento_dia"`
	FkRegion                   string    `json:"fk_region"`
	FkComuna                   string    `json:"fk_comuna"`
	FkLocalidad                string    `json:"fk_localidad"`
	Correo                     string    `json:"correo"`
	Estado                     string    `json:"estado"`
	Time                       time.Time `json:"time"`
}
