package controllers

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"

	"github.com/akshitgrover/saas_k8s/models"
	"github.com/globalsign/mgo"
	yaml "gopkg.in/yaml.v2"
)

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
	}
}
