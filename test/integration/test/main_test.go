package test

import (
	"context"
	"log"

	"github.com/docker/compose/v2/pkg/api"
	"github.com/testcontainers/testcontainers-go/modules/compose"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setup(ctx context.Context) compose.ComposeStack {
	tc, err := compose.NewDockerComposeWith(compose.WithStackFiles("../testdata/docker-compose.yml"), compose.WithLogger(log.Default()))
	if err != nil {
		log.Fatal(err)
	}

	//Start docker compose
	tc.WaitForService("app-consumer-test", wait.ForAll(
		wait.ForLog("SERVER is UP"),
		wait.ForLog("Потребитель Kafka начал работу...")))
	err = tc.Up(ctx, compose.WithRecreate(api.RecreateDiverged), compose.WithRecreateDependencies(api.RecreateDiverged), compose.Wait(true))
	if err != nil {
		cleanup(ctx, tc)
		log.Fatal(err)
	}
	return tc
}

func cleanup(ctx context.Context, tc compose.ComposeStack) {
	err := tc.Down(ctx, compose.RemoveOrphans(true), compose.RemoveImagesLocal, compose.RemoveVolumes(true))
	if err != nil {
		log.Println(err)
	}
}
