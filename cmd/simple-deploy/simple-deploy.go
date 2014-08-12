package main

import (
	"net/url"
	"os"
	"strconv"
	"strings"

	"code.google.com/p/go-uuid/uuid"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	kube_client "github.com/GoogleCloudPlatform/kubernetes/pkg/client"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/labels"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/util"
	"github.com/golang/glog"
	"gopkg.in/v1/yaml"
)

func main() {
	util.InitLogs()
	defer util.FlushLogs()

	var masterServer string
	if len(os.Getenv("KUBERNETES_MASTER")) > 0 {
		masterServer = os.Getenv("KUBERNETES_MASTER")
	} else {
		masterServer = "http://localhost:8080"
	}
	_, err := url.Parse(masterServer)
	if err != nil {
		glog.Fatalf("Unable to parse %v as a URL\n", err)
	}
	client := kube_client.New(masterServer, nil)
	deployTarget(client)
}

func deployTarget(client *kube_client.Client) {
	var selector labels.Selector
	var targetImage, targetReplicas, portMapping, deploymentName string
	var replicas int
	if deploymentName = os.Getenv("NAME"); len(deploymentName) == 0 {
		glog.Fatal("No name specified. Expected NAME environment variable")
		return
	}
	if targetImage = os.Getenv("TARGET_IMAGE"); len(targetImage) == 0 {
		glog.Fatal("No target image specified. Expected TARGET_IMAGE environment variable")
		return
	}
	if targetReplicas = os.Getenv("TARGET_REPLICAS"); len(targetReplicas) == 0 {
		glog.Fatal("No target replicas specified. Expected TARGET_REPLICAS environment variable")
		return
	}
	if portMapping = os.Getenv("PORT_MAPPING"); len(portMapping) == 0 {
		glog.Fatal("No port mapping specified. Expected PORT_MAPPING environment variable")
		return
	}

	replicas, _ = strconv.Atoi(targetReplicas)

	selector, _ = labels.ParseSelector("deployment=" + deploymentName)
	replicationControllers, err := client.ListReplicationControllers(selector)
	if err != nil {
		glog.Fatalf("Unable to get list of replication controllers %v\n", err)
		return
	}

	controllerName := uuid.NewUUID().String()

	controller := api.ReplicationController{
		JSONBase: api.JSONBase{
			ID: controllerName,
		},
		DesiredState: api.ReplicationControllerState{
			Replicas: replicas,
			ReplicaSelector: map[string]string{
				"name": controllerName,
			},
			PodTemplate: api.PodTemplate{
				DesiredState: api.PodState{
					Manifest: api.ContainerManifest{
						Version: "v1beta2",
						Containers: []api.Container{
							{
								Name:  controllerName,
								Image: targetImage,
								Ports: makePorts(portMapping),
							},
						},
					},
				},
				Labels: map[string]string{
					"name":       controllerName,
					"deployment": deploymentName,
				},
			},
		},
		Labels: map[string]string{
			"name":       controllerName,
			"deployment": deploymentName,
		},
	}

	glog.Info("Creating replication controller: ")
	obj, _ := yaml.Marshal(controller)
	glog.Info(string(obj))

	if _, err := client.CreateReplicationController(controller); err != nil {
		glog.Fatalf("An error occurred creating the replication controller: %v", err)
		return
	}

	// For this simple deploy, remove previous replication controllers
	for _, rc := range replicationControllers.Items {
		glog.Info("Stopping replication controller: ")
		obj, _ := yaml.Marshal(rc)
		glog.Info(string(obj))
		rcObj, err1 := client.GetReplicationController(rc.ID)
		if err1 != nil {
			glog.Fatalf("Unable to get replication controller %s - error: %#v\n", rc.ID, err1)
		}
		rcObj.DesiredState.Replicas = 0
		_, err := client.UpdateReplicationController(rcObj)
		if err != nil {
			glog.Fatalf("Unable to stop replication controller %s - error: %#v\n", rc.ID, err)
		}
	}

	for _, rc := range replicationControllers.Items {
		glog.Infof("Deleting replication controller %s", rc.ID)
		err := client.DeleteReplicationController(rc.ID)
		if err != nil {
			glog.Fatalf("Unable to remove replication controller %s - error: %#v\n", rc.ID, err)
		}
	}

}

func makePorts(spec string) []api.Port {
	parts := strings.Split(spec, ",")
	var result []api.Port
	for _, part := range parts {
		pieces := strings.Split(part, ":")
		if len(pieces) != 2 {
			glog.Infof("Bad port spec: %s", part)
			continue
		}
		host, err := strconv.Atoi(pieces[0])
		if err != nil {
			glog.Errorf("Host part is not integer: %s %v", pieces[0], err)
			continue
		}
		container, err := strconv.Atoi(pieces[1])
		if err != nil {
			glog.Errorf("Container part is not integer: %s %v", pieces[1], err)
			continue
		}
		result = append(result, api.Port{ContainerPort: container, HostPort: host})
	}
	return result

}
