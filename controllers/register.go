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

func registerStreams(cmd *exec.Cmd) {
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
}

func Register(db *mgo.Database) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		_ = req.ParseForm()
		f, _ := os.OpenFile("yamls/resourcequota.yml", os.O_RDWR, 0755)
		defer f.Close()
		stat, _ := f.Stat()
		b := make([]byte, stat.Size())
		f.Read(b)
		var r models.ResourceQuota
		err := yaml.Unmarshal(b, &r)
		if err != nil {
			fmt.Println(err)
		} else {
			t := models.Tenant{Username: req.PostFormValue("tenantid")}
			r.Metadata["namespace"] = t.Username
			b_, _ := yaml.Marshal(r)
			f.WriteAt(b_, 0)
			err = db.C("tenants").Insert(t)
			if err != nil {
				res.WriteHeader(500)
				res.Write([]byte(""))
				return
			}
			cmd := exec.Command("kubectl", "create", "namespace", t.Username)
			registerStreams(cmd)
			cmd.Run()
			cmd = exec.Command("kubectl", "apply", "-f", "yamls/resourcequota.yml")
			registerStreams(cmd)
			cmd.Run()
		}
	}
}
