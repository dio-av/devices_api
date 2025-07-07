package server

import (
	"devices_api/internal/devices"
	"devices_api/internal/server/rest"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:1323/swagger.json"),
	))

	// REST api routes begin
	// TODO: REST api router should be moved to the rest package
	apiRouter := chi.NewRouter()
	r.Mount("/api/v1", apiRouter)

	apiRouter.Post("/devices", s.createDevice)

	apiRouter.Put("/devices/{deviceID:^[0-9]}", s.updateDevice)

	apiRouter.Get("/devices/{deviceID:^[0-9]}", s.deviceById)

	apiRouter.Get("/devices/brands/{brand}", s.devicesByBrand)

	apiRouter.Get("devices/states/{state}", s.devicesByState)

	apiRouter.Get("/devices", s.AllDevices)

	apiRouter.Delete("devices/{deviceID}", s.DeleteDevice)
	// end of REST api routes

	r.Get("/", s.HelloWorldHandler)

	r.Get("/health", s.healthHandler)

	return r
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Devices API"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
		return
	}

	_, err = w.Write(jsonResp)
	if err != nil {
		log.Println(w, r, err.Error())
	}
}

// CreateDevice swagger:route POST /devices device devices.CreateDevice
//
// Creates a new device.
//
// Responses:
//
//		201: statusCreated
//	 	404: statusNotFound
func (s *Server) createDevice(w http.ResponseWriter, r *http.Request) {

	var device devices.CreateDevice
	if r.Body == http.NoBody {
		err := errors.New("empty body request")
		log.Println(w, r, err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&device); err != nil {
		log.Println(w, r, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	d, err := s.db.Create(r.Context(), device)
	if err != nil {
		log.Println(w, r, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type success struct {
		Device     *devices.Device
		StatusCode int
	}
	successResponse := success{
		Device:     d,
		StatusCode: http.StatusCreated,
	}
	if err := json.NewEncoder(w).Encode(successResponse); err != nil {
		log.Println(w, r, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// deviceById swagger:route GET /devices/{id}
//
// Get a device by its ID.
//
// Responses:
//
//		200: device
//		404: statusNotFound
//	 	500: internalServerError
func (s *Server) deviceById(w http.ResponseWriter, r *http.Request) {
	idUrl := chi.URLParam(r, "deviceID")

	id, err := strconv.ParseInt(idUrl, 10, 64)
	if err != nil {
		log.Println(w, r, err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	d, err := s.db.GetById(r.Context(), id)
	if err != nil {
		log.Println(w, r, err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(d); err != nil {
		log.Println(w, r, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// AllDevices swagger:route GET /devices
//
// Get all devices.
//
// Responses:
//
//		200: []device
//		404: statusNotFound
//	 	500: internalServerError
func (s *Server) AllDevices(w http.ResponseWriter, r *http.Request) {
	dd, err := s.db.All(r.Context())
	if err != nil {
		log.Println(w, r, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(dd) <= 0 {
		log.Println(w, r, "no devices")
		http.Error(w, "no devices found", http.StatusNotFound)
		return
	}

	type success struct {
		Device     []devices.Device
		StatusCode int
	}
	successResponse := success{
		Device:     dd,
		StatusCode: http.StatusOK,
	}
	if err := json.NewEncoder(w).Encode(successResponse); err != nil {
		log.Println(w, r, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// devicesByBrand swagger:route GET /devices/brands/{deviceID}
//
// Get Devices By Brand.
//
// Responses:
//
//		200: []device
//	 	500: internalServerError
func (s *Server) devicesByBrand(w http.ResponseWriter, r *http.Request) {
	brand := chi.URLParam(r, "brand")

	dd, err := s.db.GetByBrand(r.Context(), brand)
	if err != nil {
		log.Println(w, r, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type success struct {
		Device     []devices.Device
		StatusCode int
	}
	successResponse := success{
		Device:     dd,
		StatusCode: http.StatusOK,
	}
	if err := json.NewEncoder(w).Encode(successResponse); err != nil {
		log.Println(w, r, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// UpdateDevcie swagger:route PUT /devices/{deviceID} devices updateDevice
//
// Updates the parameters for a device.
//
// Responses:
//
//	    200: device
//		400: statusBadRequest
//	    500: internalServerError
func (s *Server) updateDevice(w http.ResponseWriter, r *http.Request) {
	var device devices.Device
	// idUrl := chi.URLParam(r, "deviceID")

	// id, err := strconv.ParseInt(idUrl, 10, 64)
	// if err != nil {
	// 	log.Println(w, r, err.Error())
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }

	if err := json.NewDecoder(r.Body).Decode(&device); err != nil {
		log.Println(w, r, err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	d, err := s.db.GetById(r.Context(), device.Id)
	if err != nil {
		log.Println(w, r, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = s.db.Update(r.Context(), *d)
	if err != nil {
		log.Println(w, r, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: Add the device object in the response
	if err := json.NewEncoder(w).Encode(http.StatusOK); err != nil {
		log.Println(w, r, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// devicesByState swagger:route PUT /devices/states/{state} devices devicesByState
//
// Get devices in the parameter state.
//
// Responses:
//
//	200: []device
//	400: statusBadRequest
//	404: statusNotFound
//	500: internalServerError
func (s *Server) devicesByState(w http.ResponseWriter, r *http.Request) {
	state := chi.URLParam(r, "state")

	st, err := strconv.ParseInt(state, 10, 64)
	if err != nil {
		log.Println(w, r, err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dd, err := s.db.GetByState(r.Context(), devices.DeviceState(st))
	if err != nil {
		log.Println(w, r, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if len(dd) <= 0 {
		log.Println(w, r, "devices not found")
		http.Error(w, "devices not found", http.StatusNotFound)
		return
	}

	successResponse := rest.DevicesResponse{}
	successResponse.Body.Devices = dd
	successResponse.Body.StatusCode = http.StatusOK

	if err := json.NewEncoder(w).Encode(&successResponse); err != nil {
		log.Println(w, r, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// DeleteDevice swagger:route DELETE /devices/{id} device deleteDevice
//
// Deletes a device.
//
// Responses:
//
//		204: statusAccepted
//		404: statusNotFound
//	    500: internalServerError
func (s *Server) DeleteDevice(w http.ResponseWriter, r *http.Request) {
	idUrl := chi.URLParam(r, "deviceID")
	id, err := strconv.ParseInt(idUrl, 10, 64)
	if err != nil {
		log.Println(w, r, err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	d := devices.Device{Id: id}
	result, err := s.db.Delete(r.Context(), d)
	if err != nil {
		if errors.Is(err, devices.ErrDeviceInUse) {
			log.Println(w, r, err.Error())
			http.Error(w, err.Error(), http.StatusAccepted)
			return
		} else if n, err := result.RowsAffected(); n <= 0 {
			if err != nil {
				log.Println(w, r, err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			log.Println(w, r, "device not found")
			http.Error(w, "device not found", http.StatusNotFound)
			return
		}
		log.Println(w, r, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(http.StatusNoContent); err != nil {
		log.Println(w, r, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, _ := json.Marshal(s.db.Health())
	_, err := w.Write(jsonResp)
	if err != nil {
		log.Println(w, r, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
