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
 * date  : 2016-05-30
 */

package main

import (
    "log"
    "fmt"
    "strconv"
    "strings"
)

type Data struct {
    Id string
    Name string
}

type Group struct {
    Length int
    Data Data
}

type Icon struct {
    Id string
    Name string
    GroupId string `json:"group_id"`
    Keywords []string
    Ligature string
    Codepoint string
    IsNew bool `json:"is_new"`
}

type Root struct {
    Groups []Group
    BasePath string `json:"base_path"`
    Icons []Icon
}

type Search struct {
    Json Root
    Query string
}

func (s *Search) byKeywords() []Icon {
    icons := s.Json.Icons
    res := make([]Icon, 0)
    for _, icon := range icons {
        keywords := icon.Keywords
        hit := false
        for _, keyword := range keywords {
            if strings.Contains(keyword, s.Query) {
                hit = true
            }
        }
        if hit {
            res = append(res, icon)
        }
    }
    return res
}

func (s *Search) byName() *Icon {
    icons := s.Json.Icons
    for _, icon := range icons {
        if s.Query == icon.Name {
            return &icon
        }
    }
    return nil
}

func (s *Search) Search(st, color string, size int) *Icon {
    switch(st) {
    case "name":
        return s.byName()
    case "keyword":
        res := s.byKeywords()
        for i, r := range res {
            fmt.Println(strconv.Itoa(i) + ": " + r.Name)
        }
        idx := -1
        for ;; {
            str := readFromStdin()
            i, err := strconv.Atoi(strings.TrimSpace(str))
            if err != nil {
                log.Println(err)
                continue
            }
            if -1 < i && i < len(res) {
                idx = i
                break
            } else {
                continue
            }
        }
        if idx != -1 {
            return &res[idx]
        }
    }
    return nil
}
