package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

//Cron cron detail
type Cron struct {
	ID       string   `json:"id"`
	CMD      string   `json:"cmd"`
	Args     []string `json:"args"`
	Interval int64    `json:"interval"`
	tk       *time.Ticker
}

//ServerResponse http respone
type ServerResponse struct {
	OK    bool   `json:"ok"`
	ID    string `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

//HTTPAPI http handle
func (run *Run) HTTPAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	sr := &ServerResponse{}
	if r.Body == nil {
		sr.Error = "post data error"
		run.XHTTPRespone(w, sr, http.StatusBadRequest)
		return
	}
	cron := &Cron{}
	err := json.NewDecoder(r.Body).Decode(cron)
	if err != nil {
		sr.Error = fmt.Sprintf("post data can not decode with error %v", err)
		run.XHTTPRespone(w, sr, http.StatusBadRequest)
		return
	}
	if cron.ID == "" {
		sr.Error = "post data error, miss some params"
		run.XHTTPRespone(w, sr, http.StatusBadRequest)
		return
	}
	log.Printf("add data is %v", cron)
	sr.ID = cron.ID

	var httpStatus int
	if r.Method == http.MethodPost {
		httpStatus = run.XHTTPPost(sr, cron)
	} else if r.Method == http.MethodDelete {
		httpStatus = run.XHTTPDelete(sr, cron)
	} else {
		log.Printf("can not accept method %s", r.Method)
		return
	}

	run.XHTTPRespone(w, sr, httpStatus)
}

//StartHTTP start http
func (run *Run) StartHTTP() error {
	log.Printf("start http...")
	http.HandleFunc("/", run.HTTPAPI)

	log.Printf("listen port %d...", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

//XHTTPDelete delete task
func (run *Run) XHTTPDelete(sr *ServerResponse, cron *Cron) int {
	if run.OperateTask(cron, true) {
		sr.OK = true
		return http.StatusOK
	}
	sr.Error = fmt.Sprintf("The task %s is not found.", sr.ID)
	log.Printf("delete task success...")
	return http.StatusNotFound
}

//XHTTPPost handle http post
func (run *Run) XHTTPPost(sr *ServerResponse, cron *Cron) int {
	if run.OperateTask(cron, false) {
		sr.Error = fmt.Sprintf("The task %s already exists.", sr.ID)
		return http.StatusConflict
	}
	sr.OK = true
	log.Printf("add task success...")
	return http.StatusOK
}

//XHTTPRespone http post respone
func (run *Run) XHTTPRespone(w http.ResponseWriter, sr *ServerResponse, httpStatus int) {
	log.Printf("respone data %v", sr)
	j, _ := json.Marshal(sr)
	w.WriteHeader(httpStatus)
	w.Write(j)
}

//OperateTask check cron exist
func (run *Run) OperateTask(cron *Cron, del bool) bool {
	m := run.MuxCrons
	m.Lock()
	defer m.Unlock()
	cronNum := len(m.Crons)
	for i := 0; i < cronNum; i++ {
		c := <-m.Crons
		if c.ID == cron.ID && !del {
			m.Crons <- c
			//task exist
			return true
		} else if c.ID == cron.ID && del {
			c.tk.Stop()
			//task exist and delete
			return true
		}
	}
	if !del {
		cron.tk = time.NewTicker(time.Millisecond * time.Duration(cron.Interval))
		go run.RunTask(cron)
		run.MuxCrons.Crons <- cron
	}
	return false
}
