package utils

import (
	"fmt"

	"github.com/goccy/go-json"
	"github.io/victor99z/rinha/models"
)

func JsonPrettyPrint(newTx *models.Transaction) string {
	b, err := json.MarshalIndent(newTx, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	return string(b)
}
