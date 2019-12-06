package services

import "github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/models"

type DatastoreScope interface {
	GetScopes([]int) (models.Scope, error)
	UpdateScope(int, models.Scope) (models.Scope, error)
	DeleteScope(int) error
	CreateScope(int, models.Scope) (models.Scope, error)
}
