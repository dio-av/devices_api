package server

import (
	"devices_api/internal/devices"
	"encoding/json"
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
		httpSwagger.URL("http://localhost:1323/swagger/doc.json"),
	))

	// api router
	apiRouter := chi.NewRouter()
	r.Mount("/api/v1", apiRouter)

	apiRouter.Post("/devices/new", s.NewDevice)

	apiRouter.Put("/devices/update", s.UpdateDevice)

	apiRouter.Get("/devices/{id}", s.DeviceById)

	apiRouter.Get("/devices/brand/{brand}", s.DevicesByBrand)

	apiRouter.Get("/devices/all", s.AllDevices)

	apiRouter.Get("devices/state/{state}", s.DevicesByState)
	// end of api router

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

// NewDevice godoc
//
//	@Summary		Add a device
//	@Description	add by json create device
//	@Tags			devices
//	@Accept			json
//	@Produce		json
//	@Param			create_device	body		model.CreateDevice	true	"Add device"
//	@Success		200		{object}	model.Device
//	@Failure		400		{object}	httputil.HTTPError
//	@Failure		404		{object}	httputil.HTTPError
//	@Failure		500		{object}	httputil.HTTPError
//	@Router			/devices/new [post]
func (s *Server) NewDevice(w http.ResponseWriter, r *http.Request) {
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

func (s *Server) AllDevices(w http.ResponseWriter, r *http.Request) {
	dd, err := s.db.All(r.Context())
	if err != nil {
		log.Println(w, r, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(dd)
}

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

func (s *Server) DeleteDevice(w http.ResponseWriter, r *http.Request) {
	var device devices.Device
	json.NewDecoder(r.Body).Decode(&device)

	result, err := s.db.Delete(r.Context(), device)
	if err != nil {
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
