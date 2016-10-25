package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/josephroberts/edge-node-manager/application"
	"github.com/josephroberts/edge-node-manager/database"

	log "github.com/Sirupsen/logrus"
)

func DependantDeviceUpdate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deviceUUID := vars["uuid"]

	type dependantDeviceUpdate struct {
		Commit      string      `json:"commit"`
		Environment interface{} `json:"environment"`
	}

	decoder := json.NewDecoder(r.Body)
	var content dependantDeviceUpdate
	if err := decoder.Decode(&content); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Unable to decode Dependant device update hook")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	applicationUUID, localUUID, err := database.GetDeviceMapping(deviceUUID)
	if err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Unable to get device mapping")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	log.WithFields(log.Fields{
		"ApplicationUUID": applicationUUID,
		"DeviceUUID":      deviceUUID,
		"LocalUUID":       localUUID,
		"Target commit":   content.Commit,
	}).Debug("Dependant device update hook")

	application.List[applicationUUID].TargetCommit = content.Commit
	application.List[applicationUUID].Devices[localUUID].TargetCommit = content.Commit

	w.WriteHeader(http.StatusOK)
}

func DependantDeviceDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deviceUUID := vars["uuid"]

	applicationUUID, localUUID, err := database.GetDeviceMapping(deviceUUID)
	if err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Unable to get device mapping")
		// Send back 200 as the devce must of already been deleted if we can't find it in the DB
		w.WriteHeader(http.StatusOK)
		return
	}

	log.WithFields(log.Fields{
		"ApplicationUUID": applicationUUID,
		"DeviceUUID":      deviceUUID,
		"LocalUUID":       localUUID,
	}).Debug("Dependant device delete hook")

	delete(application.List[applicationUUID].Devices, localUUID)
	if err := database.DeleteDevice(applicationUUID, deviceUUID); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Unable to delete device")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DependantDeviceRestart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deviceUUID := vars["uuid"]

	applicationUUID, localUUID, err := database.GetDeviceMapping(deviceUUID)
	if err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Unable to get device mapping")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	log.WithFields(log.Fields{
		"UUID": deviceUUID,
	}).Debug("Dependant device restart hook")

	application.List[applicationUUID].Devices[localUUID].RestartFlag = true

	w.WriteHeader(http.StatusOK)
}

func SetState(w http.ResponseWriter, r *http.Request) {
	// decoder := json.NewDecoder(r.Body)
	// if err := decoder.Decode(&process.State); err != nil {
	// 	log.WithFields(log.Fields{
	// 		"Error": err,
	// 	}).Error("Unable to decode State hook")
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	// log.WithFields(log.Fields{
	// 	"State": process.State,
	// }).Debug("Set state")

	// w.WriteHeader(http.StatusOK)
}

func GetState(w http.ResponseWriter, r *http.Request) {
	// bytes, err := json.Marshal(process.State)
	// if err != nil {
	// 	log.WithFields(log.Fields{
	// 		"Error": err,
	// 	}).Error("Unable to encode State hook")
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	// log.WithFields(log.Fields{
	// 	"State": process.State,
	// }).Debug("Get state")

	// w.Header().Set("Content-Type", "application/json")
	// w.Write(bytes)
	// w.WriteHeader(http.StatusOK)
}
