package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/api/deploymentmanager/v2"
	"io/ioutil"
	"log"
	"os"
)

func readFileContents(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", nil
	}

	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func main() {
	projectName := flag.String("project", "", "Name of project to deploy to. (Required)")
	configPath  := flag.String("config", "", "Path to deployment manager config file. (Required)")

	flag.Parse()

	if *configPath == "" || *projectName == "" {
		flag.Usage()
		os.Exit(1)
	}

	configFile, err := readFileContents(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	targetConfiguration := &deploymentmanager.TargetConfiguration{
		Config: &deploymentmanager.ConfigFile{ Content: configFile },
	}

	deployment := &deploymentmanager.Deployment{
		Name: "hello-world",
		Target: targetConfiguration,
	}

	ctx := context.Background()

	service, err := deploymentmanager.NewService(ctx)
	if err != nil {
		log.Fatal(err)
	}

	insertDeployment, err := service.Deployments.Insert(*projectName, deployment).Do()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("[*] Creating deployment %s\n", insertDeployment.Name)

	getOperation := service.Operations.Get(*projectName, insertDeployment.Name)

	operation, err := getOperation.Do()
	if err != nil || operation.Error != nil {
		log.Fatal(err)
	}

	fmt.Printf("[*] Waiting for create operation %s...", operation.Name)

	for {
		operation, err := getOperation.Do()

		if err != nil || operation.Error != nil {
			log.Fatal(err)
		}

		if operation.Status == "DONE" {
			fmt.Println("done.")
			break
		}
	}

	resource, err := service.Resources.Get(*projectName, deployment.Name, "vm").Do()
	if err != nil {
		log.Fatal(err)
	}

	json, _ := resource.MarshalJSON()
	fmt.Printf("Resource: %s\n\n", string(json))
}
