# static-cling
A no nonsense static static site generator

[![Go Report Card](https://goreportcard.com/badge/github.com/DataDrake/static-cling)](https://goreportcard.com/report/github.com/DataDrake/static-cling) [![license](https://img.shields.io/github/license/DataDrake/static-cling.svg)]()

## Motivation

There are many existing static site generators and all of them have one issue or another that makes them not something that I want to use. Some of these things include:

* Not supporting HAML as a template language
* Forcing the use of template languages that are hard to write
* Enforcing a complicated or unintuitive project structure
* Having so many features that it's possible to achieve the same result, with various solutions of differing complexity
* Requiring many additional dependencies to handle various features that you may or may not use (ie. scripting languages, libraries)

## Goals

* Easy to use
* As few features as necessary
* Support multiple templating languages:
  - [ ] Go's `html/template`
  - [ ] HAML (WIP)
* Support multiple content languages:
  - [ ] HAML (WIP)
  - [ ] HTML
  - [ ] Markdown (blackfriday)
  - [ ] TimberText (my own language)
* Documentation will be self-hosted as a static-cling project
* Stretch Goals
  - [ ] CSS Preprocessor Support
  - [ ] RSS Feeds
    - [ ] Podcast Support
* A+ Rating on [Report Card](https://goreportcard.com/report/github.com/DataDrake/static-cling)

 
## Requirements

#### Compile-Time
* Go 1.16 (tested)
* Make

## Installation

1. Clone repo and enter its directory
2. `make`
3. `sudo make install`

## Usage

> Under construction

## License
 
Copyright 2021 Bryan T. Meyers <root@datadrake.com>
 
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
 
http://www.apache.org/licenses/LICENSE-2.0
 
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
