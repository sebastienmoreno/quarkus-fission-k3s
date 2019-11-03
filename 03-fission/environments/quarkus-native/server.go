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
  "syscall"
)

const (
	CODE_PATH          = "/userfunc/user"
)

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

// var userFunc http.HandlerFunc
var specialized bool

// func (bs *BinaryServer) specializeHandler(logger *zap.Logger) func(w http.ResponseWriter, r *http.Request)
// 	if specialized {
// 		w.WriteHeader(http.StatusBadRequest)
// 		w.Write([]byte("Not a generic container"))
// 		return
// 	}

// 	request := FunctionLoadRequest{}

// 	codePath := bs.fetchedCodePath
// 	err := json.NewDecoder(r.Body).Decode(&request)
// 	switch {
// 	case err == io.EOF:
// 	case err != nil:
// 		panic(err)
// 	}

// 	if request.FilePath != "" {
// 		fileStat, err := os.Stat(request.FilePath)
// 		if err != nil {
// 			panic(err)
// 		}

// 		codePath = request.FilePath
// 		switch mode := fileStat.Mode(); {
// 		case mode.IsDir():
// 			codePath = filepath.Join(request.FilePath, request.FunctionName)
// 		}
// 	}

// 	_, err = os.Stat(codePath)
// 	if err != nil {
// 		if os.IsNotExist(err) {
// 			w.WriteHeader(http.StatusNotFound)
// 			w.Write([]byte(codePath + ": not found"))
// 			return
// 		} else {
// 			panic(err)
// 		}
// 	}

// 	// Future: Check if executable is correct architecture/executable.

// 	// Copy the executable to ensure that file is executable and immutable.
// 	userFunc, err := ioutil.ReadFile(codePath)
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		w.Write([]byte("Failed to read executable."))
// 		return
// 	}
// 	err = ioutil.WriteFile(bs.internalCodePath, userFunc, 0555)
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		w.Write([]byte("Failed to write executable to target location."))
// 		return
// 	}

// 	logger.Info("Specializing ...")
//   specialized = true
  
//   logger.Info("Executing function process ...", bs.internalCodePath)
//   args := []string{"-Dquarkus.http.host=0.0.0.0"}
//   env := os.Environ()
//   execErr := syscall.Exec(bs.internalCodePath, args, env)
//   if execErr != nil {
//       panic(execErr)
//   }
// 	logger.Info("Done")
// }


// func specializeHandler(logger *zap.Logger) func(http.ResponseWriter, *http.Request) {
//   return func(w http.ResponseWriter, r *http.Request) {
//     if specialized {
//       w.WriteHeader(http.StatusBadRequest)
//       w.Write([]byte("Not a generic container"))
//       return
//     }

//     _, err := os.Stat(CODE_PATH)
//     if err != nil {
//       if os.IsNotExist(err) {
//         w.WriteHeader(http.StatusNotFound)
//         logger.Error("code path does not exist",
//           zap.Error(err),
//           zap.String("code_path", CODE_PATH))
//         w.Write([]byte(CODE_PATH + ": not found"))
//         return
//       } else {
//         logger.Error("unknown error looking for code path",
//           zap.Error(err),
//           zap.String("code_path", CODE_PATH))
//         err = errors.Wrap(err, "unknown error")
//         w.WriteHeader(http.StatusInternalServerError)
//         w.Write([]byte(err.Error()))
//         return
//       }
//     }

//   // TODO

//     // userFunc, err = loadPlugin(logger, CODE_PATH, "Handler")
//     // if err != nil {
//     //   e := "error specializing function"
//     //   logger.Error(e, zap.Error(err))
//     //   w.WriteHeader(http.StatusInternalServerError)
//     //   w.Write([]byte(errors.Wrap(err, e).Error()))
//     //   return
//     // }
// ///////////////////
//     //args := []string{"ls", "-a", "-l", "-h"}
//     args := []string{}

//     env := os.Environ()

//     logger.Info("Executing function process ...")
//     execErr := syscall.Exec(CODE_PATH, args, env)
//     if execErr != nil {
//         panic(execErr)
//     }
//     logger.Info("Specializing ...")
//     specialized = true
// //////////////////
//     // package main

//     // import (
//     //     "os"
//     //     "os/exec"
//     //     "syscall"
//     // )
    
//     // func main() {
    
//     //     binary, lookErr := exec.LookPath("ls")
//     //     if lookErr != nil {
//     //         panic(lookErr)
//     //     }
    
//     //     args := []string{"ls", "-a", "-l", "-h"}
    
//     //     env := os.Environ()
    
//     //     execErr := syscall.Exec(binary, args, env)
//     //     if execErr != nil {
//     //         panic(execErr)
//     //     }
//     // }
// ///////////////////
//     logger.Info("done")
//   }
// }

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

    logger.Info("CODE_PATH", zap.String("CODE_PATH", CODE_PATH))
    logger.Info("loadreq.FilePath", zap.String("loadreq.FilePath", loadreq.FilePath))
    logger.Info("loadreq.FunctionName", zap.String("loadreq.FunctionName", loadreq.FunctionName))

    args := []string{}
    env := os.Environ()
    logger.Info("Executing function process ...")
    execErr := syscall.Exec(loadreq.FilePath, args, env)
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

  // Generic route -- all http requests go to the user function.
  // http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
  // 	if userFunc == nil {
  // 		w.WriteHeader(http.StatusInternalServerError)
  // 		w.Write([]byte("Generic container: no requests supported"))
  // 		return
  // 	}
  // 	userFunc(w, r)
  // })
  userfuncUrl, err := url.Parse("http://localhost:8080")
  if err != nil {
    panic(err)
  }
  proxy := httputil.NewSingleHostReverseProxy(userfuncUrl)
  http.HandleFunc("/", handler(proxy))

  logger.Info("listening on 8888 ...")
  err = http.ListenAndServe(":8888", nil)
  if err != nil {
    panic(err)
  }
}

func handler(p *httputil.ReverseProxy) func(w http.ResponseWriter, r *http.Request) {
  if !specialized {
    return func(w http.ResponseWriter, r *http.Request) {
      w.WriteHeader(http.StatusInternalServerError)
      w.Write([]byte("Generic container: no requests supported"))
    }
  } else {
    return func(w http.ResponseWriter, r *http.Request) {
            log.Println(r.URL)
            w.Header().Set("X-Ben", "Rad")
            p.ServeHTTP(w, r)
    }
  }
}
