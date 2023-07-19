package app

import (
	"github.com/google/wire"
)

var InjectorSet = wire.NewSet(wire.Struct(new(Injector), "*"))

type Injector struct {
	// Engine *gin.Engine
	DB *DB
}
