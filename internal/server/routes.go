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

	apiRouter.Patch("/devices/{id}", s.updateDevice)

	apiRouter.Get("/devices/{id}", s.deviceById)

	apiRouter.Get("/devices/brand/{brand}", s.devicesByBrand)

	apiRouter.Get("devices/state/{state}", s.devicesByState)

	apiRouter.Get("/devices", s.AllDevices)

	apiRouter.Delete("devices/{id}", s.DeleteDevice)
	// end of REST api routes

	r.Get("/", s.HelloWorldHandler)

	r.Get("/health", s.healthHandler)

	return r
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

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
//		default: genericError
//		200: device
//	 	500: internalServerError
func (s *Server) createDevice(w http.ResponseWriter, r *http.Request) {

	var device devices.CreateDevice
	if err := json.NewDecoder(r.Body).Decode(&device); err != nil {
		log.Println(w, r, err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	d, err := s.db.Create(r.Context(), device)
	if err != nil {
		log.Println(w, r, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type success struct {
		device     *devices.Device
		statusCode int
	}
	successResponse := success{
		device:     d,
		statusCode: http.StatusOK,
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
//		default: genericError
//		200: device
//	 	500: internalServerError
func (s *Server) deviceById(w http.ResponseWriter, r *http.Request) {
	idUrl := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(idUrl, 10, 64)
	if err != nil {
		log.Println(w, r, err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	d, err := s.db.GetById(r.Context(), id)
	if err != nil {
		log.Println(w, r, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(d); err != nil {

	}
}

// AllDevices swagger:route GET /devices
//
// Get all devices.
//
// Responses:
//
//		default: genericError
//		200: []device
//	 	500: internalServerError
func (s *Server) AllDevices(w http.ResponseWriter, r *http.Request) {
	dd, err := s.db.All(r.Context())
	if err != nil {
		log.Println(w, r, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type success struct {
		device     []devices.Device
		statusCode int
	}
	successResponse := success{
		device:     dd,
		statusCode: http.StatusOK,
	}
	if err := json.NewEncoder(w).Encode(successResponse); err != nil {
		log.Println(w, r, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// devicesByBrand swagger:route GET /devices/brand/{brand}
//
// Get Devices By Brand.
//
// Responses:
//
//		default: genericError
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
		device     []devices.Device
		statusCode int
	}
	successResponse := success{
		device:     dd,
		statusCode: http.StatusOK,
	}
	if err := json.NewEncoder(w).Encode(successResponse); err != nil {
		log.Println(w, r, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// UpdateDevcie swagger:route PUT /devices/{id} devices updateDevice
//
// Updates the parameters for a device.
//
// Responses:
//
//	default: genericError
//	    200: device
//		400: statusBadRequest
//	    500: internalServerError
func (s *Server) updateDevice(w http.ResponseWriter, r *http.Request) {
	var device devices.Device
	if err := json.NewDecoder(r.Body).Decode(&device); err != nil {
		log.Println(w, r, err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

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

// devicesByState swagger:route PUT /devices/state/{state} devices devicesByState
//
// Get devices in the parameter state.
//
// Responses:
//
//	default: genericError
//	    200: []device
//	    500: internalServerError
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
	}

	successResponse := rest.DevicesResponse{
		Devices:    dd,
		StatusCode: http.StatusOK,
	}
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
//		default: genericError
//		204: statusNoContent
//	    500: internalServerError
func (s *Server) DeleteDevice(w http.ResponseWriter, r *http.Request) {
	idUrl := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idUrl, 10, 64)
	if err != nil {
		log.Println(w, r, err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	d := devices.Device{Id: id}
	_, err = s.db.Delete(r.Context(), d)
	if err != nil {
		if errors.Is(err, devices.ErrDeviceInUse) {
			log.Println(w, r, err.Error())
			http.Error(w, err.Error(), http.StatusAccepted)
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
