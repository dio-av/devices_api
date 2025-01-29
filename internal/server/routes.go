package server

import (
	"devices_api/internal/devices"
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

	apiRouter.Post("/devices/new", s.CreateDevice)

	apiRouter.Put("/devices/{id}", s.UpdateDevice)

	apiRouter.Get("/devices/{id}", s.DeviceById)

	apiRouter.Get("/devices/brand/{brand}", s.DevicesByBrand)

	apiRouter.Get("devices/state/{state}", s.DevicesByState)

	apiRouter.Get("/devices/all", s.AllDevices)

	apiRouter.Delete("devices/delete/", s.DeleteDevice)
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
	}

	_, _ = w.Write(jsonResp)
}

// CreateDevice swagger:route POST /devices/new device devices.CreateDevice
//
// Creates a new device.
//
// Responses:
//
//		default: genericError
//		200: device
//	 	500: internalServerError
func (s *Server) CreateDevice(w http.ResponseWriter, r *http.Request) {
	var device devices.CreateDevice

	json.NewDecoder(r.Body).Decode(&device)

	d, err := s.db.Create(r.Context(), device)
	if err != nil {
		log.Println(w, r, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(d)
}

// DeviceById swagger:route GET /devices/{id}
//
// Get a device by its ID.
//
// Responses:
//
//		default: genericError
//		200: device
//	 	500: internalServerError
func (s *Server) DeviceById(w http.ResponseWriter, r *http.Request) {
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

	json.NewEncoder(w).Encode(d)
}

// AllDevices swagger:route GET /devices/all
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
	json.NewEncoder(w).Encode(dd)
}

// DevicesByBrand swagger:route GET /devices/brand/{brand}
//
// Get Devices By Brand.
//
// Responses:
//
//		default: genericError
//		200: []device
//	 	500: internalServerError
func (s *Server) DevicesByBrand(w http.ResponseWriter, r *http.Request) {
	brand := chi.URLParam(r, "brand")

	dd, err := s.db.GetByBrand(r.Context(), brand)
	if err != nil {
		log.Println(w, r, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(dd)
}

// UpdateDevcie swagger:route PUT /devices/{id} devices updateDevice
//
// Updates the parameters for a device.
//
// Responses:
//
//	default: genericError
//	    200: device
//	    500: internalServerError
func (s *Server) UpdateDevice(w http.ResponseWriter, r *http.Request) {
	var device devices.Device
	json.NewDecoder(r.Body).Decode(&device)

	d, err := s.db.GetById(r.Context(), device.Id)
	if err != nil {
		log.Println(w, r, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	du, err := s.db.Update(r.Context(), *d)
	if err != nil {
		log.Println(w, r, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(du)
}

// DevicesByState swagger:route PUT /devices/state/{state} devices DevicesByState
//
// Get devices in the parameter state.
//
// Responses:
//
//	default: genericError
//	    200: []device
//	    500: internalServerError
func (s *Server) DevicesByState(w http.ResponseWriter, r *http.Request) {
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
	json.NewEncoder(w).Encode(dd)
}

// DeleteDevice swagger:route DELETE /devices/{id} device deleteDevice
//
// Deletes a device.
//
// Responses:
//
//		default: genericError
//		204:
//	    500: internalServerError
func (s *Server) DeleteDevice(w http.ResponseWriter, r *http.Request) {
	var device devices.Device
	json.NewDecoder(r.Body).Decode(&device)

	result, err := s.db.Delete(r.Context(), device)
	if err != nil {
		if errors.Is(err, devices.ErrDeviceInUse) {
			http.Error(w, err.Error(), http.StatusAccepted)
		}
		log.Println(w, r, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	i, err := result.RowsAffected()
	json.NewEncoder(w).Encode(i)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, _ := json.Marshal(s.db.Health())
	_, _ = w.Write(jsonResp)
}
