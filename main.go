/*
 * Copyright 2016 Yoshihiro Tanaka
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * Author: Yoshihiro Tanaka <contact@cordea.jp>
 * date  : 2016-05-29
 */

package main

import (
    "flag"
    "log"
    "strconv"
    "os"
    "io"
    "path"
    "bufio"
    "net/url"
    "io/ioutil"
    "encoding/json"
    "net/http"
)

var (
    color = flag.String("c", "black", "")
    size = flag.Int("s", 24, "")
    fileType = flag.String("f", "zip", "")
    searchType = flag.String("t", "name", "")
    out = flag.String("o", ".", "")
)

const (
    JsonUrl = "https://design.google.com/icons/data/grid.json"
    DownloadPath = "icons/"
)

func download(basePath string, icon Icon, color string, size int, out, fileType string) {
    name := icon.Id + "_" + color + "_" + strconv.Itoa(size) + "dp." + fileType
    uri, err := url.Parse(basePath)
    if err != nil {
        log.Fatalln(err)
    }
    uri.Path = path.Join(uri.Path, DownloadPath, fileType, name)
    saveFile(uri.String(), path.Join(out, name))
}

func saveFile(url, out string) {
    outFile, err := os.Create(out)
    if err != nil {
        log.Fatalln(err)
    }
    defer outFile.Close()
    log.Println("download " + url)
    resp, err := http.Get(url)
    if err != nil {
        os.Remove(out)
        log.Fatalln(err)
    }
    defer resp.Body.Close()
    if _, err := io.Copy(outFile, resp.Body); err != nil {
        os.Remove(out)
        log.Fatalln(err)
    }
}

func readFromStdin() string {
    r := bufio.NewReader(os.Stdin)
    str, err := r.ReadString('\n')
    if err != nil {
        log.Fatalln(err)
    }
    return str
}

func checkTypes(st, ft, cl string, sz int) string {
    if st != "name" && st != "keyword" {
        return "Search type (-t) is wrong: Accept name or keyword."
    }

    if ft != "zip" && ft != "svg" {
        return "File type (-f) is wrong: Accept zip or svg."
    }

    if cl != "white" && cl != "black" {
        return "Color (-c) is wrong: Accept white or black."
    }

    if sz != 18 && sz != 24 && sz != 36 && sz != 48 {
        return "Size (-s) is wrong: Accept 18 or 24, 36, 48."
    }
    return ""
}

func main() {
    flag.Parse()
    if flag.NArg() == 0 {
        log.Fatalln("Required argument is missing.")
    }
    q := flag.Arg(0)

    if errmsg := checkTypes(*searchType, *fileType, *color, *size); errmsg != "" {
        log.Fatalln(errmsg)
    }

    resp, err := http.Get(JsonUrl)
    if err != nil {
        log.Fatalln(err)
    }

    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Fatalln(err)
    }

    var root Root
    if err := json.Unmarshal(body, &root); err != nil {
        log.Fatalln(err)
    }

    s := Search{root, q}

    icon := s.Search(*searchType, *color, *size)
    if icon != nil {
        download(root.BasePath, *icon, *color, *size, *out, *fileType)
    }
}