package randutil

import (
	"encoding/hex"
	"github.com/satori/go.uuid"
)

func GenerateID(name string) string {
	gen := uuid.NewV5(uuid.NewV4(), name)
	return hex.EncodeToString(gen.Bytes())
}
