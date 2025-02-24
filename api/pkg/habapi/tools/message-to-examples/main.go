package main

import (
	"context"
	"fmt"
	"github.com/habiliai/alice/api/pkg/config"
	"github.com/habiliai/alice/api/pkg/digo"
	"github.com/habiliai/alice/api/pkg/habapi"
	"github.com/habiliai/alice/api/pkg/helpers"
	"github.com/habiliai/alice/api/pkg/services"
	"gorm.io/gorm"
	"strings"
)

func main() {
	ctx := context.TODO()

	conf := config.ReadHabApiConfig(".env")
	container := digo.NewContainer(ctx, digo.EnvProd, &conf)

	server, err := digo.Get[habapi.HabiliApiServer](container, habapi.ServerKey)
	if err != nil {
		panic(err)
	}
	db, err := digo.Get[*gorm.DB](container, services.ServiceKeyDB)
	if err != nil {
		panic(err)
	}

	// Get thread
	ctx = helpers.WithTx(ctx, db)
	thread, err := server.GetThread(ctx, &habapi.ThreadId{Id: 26})
	if err != nil {
		panic(err)
	}

	// Print messages
	println("-- Thread messages:\n")
	for _, msg := range thread.Messages {
		text := strings.ReplaceAll(msg.Text, "\n", "\\n")
		text = strings.ReplaceAll(text, "\r", "\\r")
		text = strings.ReplaceAll(text, "\t", "\\t")
		if strings.ToLower(msg.Role.String()) == "assistant" {
			fmt.Printf("- %s(@%s): \"%s\"\n", msg.Role, msg.Agent.Name, text)
		} else {
			fmt.Printf("- %s: \"%s\"\n", msg.Role, text)
		}
	}
	println("--\n")
}
