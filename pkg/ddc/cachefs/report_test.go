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

package cachefs

import (
	"reflect"
	"testing"

	"github.com/brahma-adshonor/gohook"
	"github.com/fluid-cloudnative/fluid/pkg/ddc/cachefs/operations"
	"github.com/fluid-cloudnative/fluid/pkg/utils/fake"
)

func mockCacheFSMetric() string {
	return `juicefs_blockcache.blocks: 9708
juicefs_blockcache.bytes: 40757435762
juicefs_blockcache.evict: 0
juicefs_blockcache.evictBytes: 0
juicefs_blockcache.evictDur: 0
juicefs_blockcache.hitBytes: 40717671794
juicefs_blockcache.hits: 9708
juicefs_blockcache.miss: 0
juicefs_blockcache.missBytes: 0
juicefs_blockcache.readDuration: 2278386748
juicefs_blockcache.transfer: 0
juicefs_blockcache.transferBytes: 0
juicefs_blockcache.transferDur: 0
juicefs_blockcache.transferScanDur: 0
juicefs_blockcache.write: 0
juicefs_blockcache.writeBytes: 0
juicefs_blockcache.writeDuration: 0
juicefs_cpuusage: 90497392
juicefs_fuse_ops.access: 0
juicefs_fuse_ops.copy_file_range: 0
juicefs_fuse_ops.create: 0
juicefs_fuse_ops.fallocate: 0
juicefs_fuse_ops.flock: 0
juicefs_fuse_ops.flush: 1
juicefs_fuse_ops.fsync: 0
juicefs_fuse_ops.getattr: 163391
juicefs_fuse_ops.getlk: 0
juicefs_fuse_ops.getxattr: 0
juicefs_fuse_ops.link: 0
juicefs_fuse_ops.listxattr: 0
juicefs_fuse_ops.lookup.cache: 0
juicefs_fuse_ops.lookup: 2
juicefs_fuse_ops.mkdir: 0
juicefs_fuse_ops.mknod: 0
juicefs_fuse_ops.open: 2
juicefs_fuse_ops.opendir: 3
juicefs_fuse_ops.read: 310652
juicefs_fuse_ops.readdir: 6
juicefs_fuse_ops.readlink: 0
juicefs_fuse_ops.release: 1
juicefs_fuse_ops.releasedir: 3
juicefs_fuse_ops.removexattr: 0
juicefs_fuse_ops.rename: 0
juicefs_fuse_ops.resolve: 0
juicefs_fuse_ops.rmdir: 0
juicefs_fuse_ops.setattr: 0
juicefs_fuse_ops.setlk: 0
juicefs_fuse_ops.setxattr: 0
juicefs_fuse_ops.statfs: 97
juicefs_fuse_ops.summary: 0
juicefs_fuse_ops.symlink: 0
juicefs_fuse_ops.truncate: 0
juicefs_fuse_ops.unlink: 0
juicefs_fuse_ops.write: 0
juicefs_fuse_ops: 474158
juicefs_gcPause: 5553281
juicefs_get_bytes: 0
juicefs_goroutines: 50
juicefs_handles: 1
juicefs_heapCacheUsed: 0
juicefs_heapInuse: 203571200
juicefs_heapSys: 360772680
juicefs_memusage: 335941632
juicefs_meta.bytes_received: 65380
juicefs_meta.bytes_sent: 73711
juicefs_meta.dircache.access: 0
juicefs_meta.dircache.add: 2
juicefs_meta.dircache.addEntry: 0
juicefs_meta.dircache.getattr: 163280
juicefs_meta.dircache.lookup: 1
juicefs_meta.dircache.newDir: 0
juicefs_meta.dircache.open: 0
juicefs_meta.dircache.readdir: 1
juicefs_meta.dircache.remove: 0
juicefs_meta.dircache.removeEntry: 0
juicefs_meta.dircache.setattr: 0
juicefs_meta.dircache0.dirs: 1
juicefs_meta.dircache0.inodes: 6
juicefs_meta.dircache0: 0
juicefs_meta.dircache: 163284
juicefs_meta.packets_received: 1293
juicefs_meta.packets_sent: 1357
juicefs_meta.reconnects: 0
juicefs_meta.usec_ping: [1799]
juicefs_meta.usec_timediff: [39520]
juicefs_meta: 305025
juicefs_metaDuration: 2306443
juicefs_metaRequest: 1258
juicefs_object.copy: 0
juicefs_object.delete: 0
juicefs_object.error: 0
juicefs_object.get: 0
juicefs_object.head: 0
juicefs_object.list: 0
juicefs_object.put: 0
juicefs_object: 0
juicefs_objectDuration.delete: 0
juicefs_objectDuration.get: 0
juicefs_objectDuration.head: 0
juicefs_objectDuration.list: 0
juicefs_objectDuration.put: 0
juicefs_objectDuration: 0
juicefs_offHeapCacheUsed: 0
juicefs_openfiles: 14
juicefs_operationDuration: 514269353
juicefs_operations: 474157
juicefs_put_bytes: 0
juicefs_readBufferUsed: 0
juicefs_read_bytes: 40717671794
juicefs_remotecache.errors: 0
juicefs_remotecache.get: 2
juicefs_remotecache.getBytes: 8
juicefs_remotecache.getDuration: 1575
juicefs_remotecache.put: 0
juicefs_remotecache.putBytes: 0
juicefs_remotecache.putDuration: 0
juicefs_remotecache.receive: 0
juicefs_remotecache.receiveBytes: 0
juicefs_remotecache.recvDuration: 0
juicefs_remotecache.send: 0
juicefs_remotecache.sendBytes: 0
juicefs_remotecache.sendDuration: 0
juicefs_symlink_cache.inserts: 0
juicefs_symlink_cache.search_hits: 0
juicefs_symlink_cache.search_misses: 0
juicefs_symlink_cache: 0
juicefs_threads: 56
juicefs_totalBufferUsed: 0
juicefs_uptime: 487
juicefs_write_bytes: 0`
}

func TestCacheFSEngine_parseMetric(t *testing.T) {
	type args struct {
		metrics string
	}
	tests := []struct {
		name          string
		args          args
		wantPodMetric fuseMetrics
	}{
		{
			name: "test",
			args: args{
				metrics: mockCacheFSMetric(),
			},
			wantPodMetric: fuseMetrics{
				blockCacheBytes:     40757435762,
				blockCacheHits:      9708,
				blockCacheMiss:      0,
				blockCacheHitsBytes: 40717671794,
				blockCacheMissBytes: 0,
				usedSpace:           0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := CacheFSEngine{}
			if gotPodMetric := j.parseMetric(tt.args.metrics); !reflect.DeepEqual(gotPodMetric, tt.wantPodMetric) {
				t.Errorf("parseMetric() = %v, want %v", gotPodMetric, tt.wantPodMetric)
			}
		})
	}
}

func TestCacheFSEngine_getPodMetrics(t *testing.T) {
	GetMetricCommon := func(a *operations.CacheFSFileUtils, cachefsPath string) (metric string, err error) {
		return mockCacheFSMetric(), nil
	}
	err := gohook.Hook((*operations.CacheFSFileUtils).GetMetric, GetMetricCommon, nil)
	if err != nil {
		t.Fatal(err.Error())
	}
	j := CacheFSEngine{
		Log: fake.NullLogger(),
	}

	gotMetrics, err := j.GetPodMetrics("test", "test")
	if err != nil {
		t.Errorf("getPodMetrics() error = %v", err)
		return
	}
	if gotMetrics != mockCacheFSMetric() {
		t.Errorf("getPodMetrics() gotMetrics = %v, want %v", gotMetrics, mockCacheFSMetric())
	}
}
