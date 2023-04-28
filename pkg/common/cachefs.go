/*
Copyright 2021 The Fluid Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package common

// Runtime for CacheFS
const (
	CacheFSRuntime               = "cachefs"
	CacheFSMountType             = "JuiceFS"
	CacheFSChart                 = CacheFSRuntime
	CacheFSFuseContainer         = "cachefs-fuse"
	CacheFSDefaultCacheDir       = "/var/jfsCache"
	CacheFSFuseImageEnv          = "CACHEFS_FUSE_IMAGE_ENV"
	CacheFSWorkerImageEnv        = "CACHEFS_WORKER_IMAGE_ENV"
	CacheFSFuseDefaultImage      = "juicedata/juicefs-fuse"
	CacheFSFuseDefaultImageTag   = "v1.0.0"
	CacheFSWorkerDefaultImage    = "juicedata/juicefs-fuse"
	CacheFSWorkerDefaultImageTag = "v1.0.0"
	CacheFSWorkerContainer       = "cachefs-worker"
	CacheFSCliPath               = "/usr/local/bin/juicefs"
	CacheFSMountPath             = "/bin/mount.juicefs"
)
