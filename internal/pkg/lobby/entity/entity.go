package entity

import (
	"github.com/google/uuid"
	"strings"
)

type entity struct {
	ID ID
}

func (e *entity) GetID() ID {
	return e.ID
}

func New() Entity {
	return &entity{
		ID: ID(strings.Split(uuid.NewString(), "-")[0]),
	}
}
