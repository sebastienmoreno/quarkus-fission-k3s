package main

import (
  "encoding/json"
  "github.com/pkg/errors"
  "go.uber.org/zap"
  "io/ioutil"
  "log"
  "net/http"
  "net/http/httputil"
  "net/url"
  "os"
  "os/exec"
)

const (
  CODE_PATH          = "/userfunc/user"
)

var specialized bool

type (
  FunctionLoadRequest struct {
    // FilePath is an absolute filesystem path to the
    // function. What exactly is stored here is
    // env-specific. Optional.
    FilePath string `json:"filepath"`

    // FunctionName has an environment-specific meaning;
    // usually, it defines a function within a module
    // containing multiple functions. Optional; default is
    // environment-specific.
    FunctionName string `json:"functionName"`

    // URL to expose this function at. Optional; defaults
    // to "/".
    URL string `json:"url"`
  }
)

func specializeHandler(logger *zap.Logger) func(http.ResponseWriter, *http.Request) {
  return func(w http.ResponseWriter, r *http.Request) {
    if specialized {
      w.WriteHeader(http.StatusBadRequest)
      w.Write([]byte("Not a generic container"))
      return
    }

    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
      logger.Error("error reading request body", zap.Error(err))
      w.WriteHeader(http.StatusInternalServerError)
      return
    }
    var loadreq FunctionLoadRequest
    err = json.Unmarshal(body, &loadreq)
    if err != nil {
      w.WriteHeader(http.StatusBadRequest)
      return
    }

    _, err = os.Stat(loadreq.FilePath)
    if err != nil {
      if os.IsNotExist(err) {
        logger.Error("code path does not exist",
          zap.Error(err),
          zap.String("code_path", loadreq.FilePath))
        w.WriteHeader(http.StatusNotFound)
        w.Write([]byte(loadreq.FilePath + ": not found"))
        return
      } else {
        logger.Error("unknown error looking for code path",
          zap.Error(err),
          zap.String("code_path", loadreq.FilePath))
        err = errors.Wrap(err, "unknown error")
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte(err.Error()))
        return
      }
    }

    logger.Info("Changing execution rights on "+loadreq.FilePath+"...")
    chmodErr := os.Chmod(loadreq.FilePath, 0755)
    if chmodErr != nil {
      panic(chmodErr)
    }
    logger.Info("Executing in background "+loadreq.FilePath+"...")
    cmd := exec.Command(loadreq.FilePath, "-Dquarkus.http.host=0.0.0.0")
    execErr := cmd.Start()
    if execErr != nil {
      panic(execErr)
    }

    logger.Info("Specializing ...")
    specialized = true
    logger.Info("done")
  }
}

func readinessProbeHandler(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusOK)
}

func main() {
  logger, err := zap.NewProduction()
  if err != nil {
    log.Fatalf("can't initialize zap logger: %v", err)
  }
  defer logger.Sync()

  http.HandleFunc("/healthz", readinessProbeHandler)
  http.HandleFunc("/specialize", specializeHandler(logger.Named("specialize_handler")))
  http.HandleFunc("/v2/specialize", specializeHandler(logger.Named("specialize_v2_handler")))

  origin, _ := url.Parse("http://localhost:8080/")

  director := func(req *http.Request) {
    req.Header.Add("X-Forwarded-Host", req.Host)
    req.Header.Add("X-Origin-Host", origin.Host)
    req.URL.Scheme = "http"
    req.URL.Host = origin.Host
  }

  proxy := &httputil.ReverseProxy{Director: director}

  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    log.Println("Calling URL ", r.URL)
    proxy.ServeHTTP(w, r)
  })

  logger.Info("listening on 8888 ...")
  err = http.ListenAndServe(":8888", nil)
  if err != nil {
    panic(err)
  }
}
