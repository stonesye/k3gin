package uuid

import "github.com/google/uuid"

func NewUUID() (uuid.UUID, error) {
	return uuid.NewRandom()
}

func MustUUID() uuid.UUID {
	uuid, err := uuid.NewUUID()
	if err != nil {
		panic(err)
	}

	return uuid
}

func MustStrig() string {
	return MustUUID().String()
}
