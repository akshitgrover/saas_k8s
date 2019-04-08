package controllers

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"

	"github.com/akshitgrover/saas_k8s/models"
	"github.com/globalsign/mgo"
	yaml "gopkg.in/yaml.v2"
)

var portCounter = 30000

func createService(r models.Pod) string {
	f, _ := os.OpenFile("yamls/service.yml", os.O_RDWR, 0755)
	defer f.Close()
	stat, _ := f.Stat()
	b := make([]byte, stat.Size())
	f.Read(b)
	fmt.Println(string(b))
	var r_ models.Service
	yaml.Unmarshal(b, &r_)
	r_.Metadata["name"] = r.Metadata.Name
	r_.Metadata["namespace"] = r.Metadata.Namespace
	r_.Spec.Selector["tenant"] = r.Metadata.Labels["tenant"]
	r_.Spec.Selector["user"] = r.Metadata.Labels["user"]
	r_.Spec.Ports[0]["targetPort"] = 80
	r_.Spec.Ports[0]["nodePort"] = portCounter
	portCounter += 1
	b_, _ := yaml.Marshal(r_)
	f.WriteAt(b_, 0)
	cmd := exec.Command("kubectl", "apply", "-f", "yamls/service.yml")
	registerStreams(cmd)
	cmd.Run()
	return strconv.Itoa(r_.Spec.Ports[0]["nodePort"])
}

func Create(db *mgo.Database) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		_ = req.ParseForm()
		f, _ := os.OpenFile("yamls/pod.yml", os.O_RDWR, 0755)
		defer f.Close()
		stat, _ := f.Stat()
		b := make([]byte, stat.Size())
		f.Read(b)
		var r models.Pod
		err := yaml.Unmarshal(b, &r)
		if err != nil {
			fmt.Println(err)
			res.WriteHeader(500)
			res.Write([]byte(""))
			return
		}
		user := req.PostFormValue("user")
		tenant := req.PostFormValue("tenantid")
		memoryl := req.PostFormValue("memoryl")
		memory := req.PostFormValue("memory")
		cpul := req.PostFormValue("cpul")
		cpu := req.PostFormValue("cpu")

		r.Metadata.Name = user + "-" + tenant
		r.Metadata.Labels["tenant"] = tenant
		r.Metadata.Labels["user"] = user
		r.Spec.Containers[0].Name = user + "-" + tenant + "-container"
		r.Spec.Containers[0].Resources["limits"]["memory"] = memoryl
		r.Spec.Containers[0].Resources["limits"]["cpu"] = cpul
		r.Spec.Containers[0].Resources["requests"]["memory"] = memory
		r.Spec.Containers[0].Resources["requests"]["cpu"] = cpu

		b_, _ := yaml.Marshal(r)
		f.WriteAt(b_, 0)
		cmd := exec.Command("kubectl", "apply", "-f", "yamls/pod.yml")
		registerStreams(cmd)
		cmd.Run()
		res.WriteHeader(200)
		port := createService(r)
		res.Write([]byte(port))
	}
}
