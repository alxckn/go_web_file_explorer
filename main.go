package main

import (
  "crypto/md5"
  "encoding/hex"
  "fmt"
  "html/template"
  "net/http"
  "log"
  "io/ioutil"
  "path/filepath"
  "os"
  "net/url"
  "path"
  "regexp"
)

const port = "8080"
const default_rel_path = "/"

var root_dir = "test"
var nginx_url = "/"
var nginx_passphrase = "titi"

func md5HexEncode(str string) string {
  re := regexp.MustCompile(`^/`)
  hasher := md5.New()
  hasher.Write([]byte(re.ReplaceAllString(str + nginx_passphrase, "")))
  return hex.EncodeToString(hasher.Sum(nil))
}

func secureNginxUrl(rel_path string) string {
  hash := md5HexEncode(rel_path)
  u, err := url.Parse(nginx_url)
  if err != nil {
    log.Println(err)
    return ""
  }
  u.Path = path.Join(u.Path, "file", hash, rel_path)
  return u.String()
}

func previousPath(path string) (prev_path string) {
  prev_path = filepath.Dir(path)
  if prev_path == "." {
    prev_path = default_rel_path
  }
  return
}

func listFiles(dir string) (files []map[string]string, dirs []map[string]string) {
  infos, err := ioutil.ReadDir(filepath.Join(root_dir, dir))
  if err != nil {
    log.Println(err)
    return
  }

  for _, info := range infos {
    if info.IsDir() {
      dirs = append(dirs, map[string]string{
        "name": info.Name(),
        "path": filepath.Join(dir, info.Name()),
      })
    } else {
      files = append(files, map[string]string{
        "name": info.Name(),
        "url": secureNginxUrl(filepath.Join(dir, info.Name())),
      })
    }
  }

  return
}

type PageData struct {
  Previous string
  Current string
  Files []map[string]string
  Dirs []map[string]string
}

func main() {
  if os.Getenv("ROOT_DIR") != "" {
    root_dir = os.Getenv("ROOT_DIR")
  }
  if os.Getenv("NGINX_PASSPHRASE") != "" {
    nginx_passphrase = os.Getenv("NGINX_PASSPHRASE")
  }
  if os.Getenv("NGINX_URL") != "" {
    nginx_url = os.Getenv("NGINX_URL")
  }

  log.Println(fmt.Sprintf("Listing files in %s", root_dir))

  tmpl := template.Must(template.ParseFiles("templates/index.tmpl"))
  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    var curr_path string

    keys, ok := r.URL.Query()["dir"]
    if !ok || len(keys[0]) < 1 {
      curr_path = default_rel_path
    } else {
      curr_path = keys[0]
    }

    log.Println(fmt.Sprintf("Received request for %s", curr_path))

    files, dirs := listFiles(curr_path)
    data := PageData{
      Previous: previousPath(curr_path),
      Current: curr_path,
      Files: files,
      Dirs: dirs,
    }

    tmpl.Execute(w, data)
  })

  log.Println(fmt.Sprintf("Serving on %s", port))
  http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
