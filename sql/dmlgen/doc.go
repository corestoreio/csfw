// Copyright 2015-2017, Cyrill @ Schumacher.fm and the CoreStore contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package dmlgen generates Go code and protocol buffer files from database tables.
//
// TODO check for https://github.com/improbable-eng/ts-protoc-gen and https://github.com/improbable-eng/grpc-web
//
// Join Preload
// Preload loads the association data in a separate query, Join Preload will
// loads association data using inner join It can handle null association
package dmlgen
