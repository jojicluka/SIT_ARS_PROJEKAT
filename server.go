package main

import (
	"errors"
	"github.com/gorilla/mux"
	"mime"
	"net/http"
)

type Service struct {
	store *ConfigStore
}

func (ts *Service) createConfigHandler(w http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if mediatype != "application/json" {
		err := errors.New("Expect application/json Content-Type")
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}

	rt, err := decodeConfigBody(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := createId()
	ts.store.Post(rt)
	renderJSON(w, rt)
	w.Write([]byte(id))
}

func (ts *Service) getConfigHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	//fmt.Println(id)
	version := mux.Vars(req)["version"]
	//fmt.Println(version)
	task, ok := ts.store.Get(id, version)
	//fmt.Println("task: ", task, " | ok: ", ok)
	if ok != nil {
		err := errors.New("key not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	renderJSON(w, task)
}

func (ts *Service) delConfigHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]

	conf, err := ts.store.Delete(id, version)

	if err != nil {
		http.Error(w, "config does not exist", http.StatusBadRequest)
		return
	}
	renderJSON(w, conf)

}

func (ts *Service) createGroupHandler(w http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if mediatype != "application/json" {
		err := errors.New("Expect application/json Content-Type")
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}

	rt, err := decodeGroupBody(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	group, err := ts.store.PostGroup(rt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, group)
}

func (ts *Service) getGroupHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	task, ok := ts.store.GetGroup(id, version)
	if ok != nil {
		err := errors.New("key not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	renderJSON(w, task)
}

func (ts *Service) delGroupHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]

	_, err := ts.store.DeleteGroup(id, version)

	if err != nil {
		http.Error(w, "Could not delete group", http.StatusBadRequest)
	} else {
		http.Error(w, "Group is deleted", http.StatusOK)
	}

}

//func (ts *Service) filterConfigHandler(writer http.ResponseWriter, request *http.Request) {
//	label := mux.Vars(request)["label"]
//
//	task, ok := ts.store.FilterLabel(label)
//	if ok != nil {
//		err := errors.New("key not found")
//		http.Error(writer, err.Error(), http.StatusNotFound)
//		return
//	}
//	if *task == nil {
//		err := errors.New("Config with this label does not exist")
//		http.Error(writer, err.Error(), http.StatusNotFound)
//		return
//	}
//	renderJSON(writer, task)
//}

func (ts *Service) getGroupLabelHandler(w http.ResponseWriter, req *http.Request) {

	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	label := mux.Vars(req)["label"]
	//list := strings.Split(label, ";")
	//sort.Strings(list)
	//sortedLabel := ""
	//for _, v := range list {
	//	sortedLabel += v + ";"
	//}
	//sortedLabel = sortedLabel[:len(sortedLabel)-1]
	returnConfigs, error := ts.store.GetGroupByLabel(id, version, label)

	if error != nil {
		renderJSON(w, "Not found")
	}
	renderJSON(w, returnConfigs)
}

//func (ts *Service) filterGroupHandler(writer http.ResponseWriter, request *http.Request) {
//	label := mux.Vars(request)["label"]
//
//	task, ok := ts.store.FilterLabel(label)
//	if ok != nil {
//		err := errors.New("key not found")
//		http.Error(writer, err.Error(), http.StatusNotFound)
//		return
//	}
//	if *task == nil {
//		err := errors.New("Config with this label does not exist")
//		http.Error(writer, err.Error(), http.StatusNotFound)
//		return
//	}
//	renderJSON(writer, task)
//}

func (ts *Service) updateGroupHandler(w http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if mediatype != "application/json" {
		err := errors.New("Expect application/json Content-Type")
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}

	rt, err := decodeGroupBody(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	v := mux.Vars(req)
	id := v["id"]
	version := v["version"]
	rt.Id = id
	rt.Version = version

	configgroup, err := ts.store.UpdateConfigGroup(rt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	renderJSON(w, configgroup)
}
