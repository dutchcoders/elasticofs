# ElasticoFS
ElasticoFS is a fuse driver for Elasticsearch clusters. Currently the most basic options are working, retrieving cluster and node status.

Navigating through Elasticsearch is now as easy as through your /proc file system. Tab completion, grepping and lessing are all completely working.

## Samples

```
$ cat /elasticsearch/_cluster/health 
{"cluster_name":"elasticsearch_remco","status":"yellow","timed_out":false,"number_of_nodes":1,"number_of_data_nodes":1,"active_primary_shards":41,"active_shards":41,"relocating_shards":0,"initializing_shards":0,"unassigned_shards":41,"delayed_unassigned_shards":0,"number_of_pending_tasks":0,"number_of_in_flight_fetch":0,"task_max_waiting_in_queue_millis":0,"active_shards_percent_as_number":50.0}cat: e: No such file or directory
```

```
$ cat /elasticsearch/_cat/nodes
172.16.84.1 172.16.84.1 10 99 2.41 d * Dead Girl
```

```
$ cat /elasticsearch/darksearch/_mapping |json_pp |less
```

```
$ cat /elasticsearch/_stats |grep -o -E \"total\":[0-9]+
```

```
$ find /elasticsearch/
/elasticsearch/
/elasticsearch/_stats
/elasticsearch/_field_stats
/elasticsearch/_cluster
/elasticsearch/_cluster/health
/elasticsearch/_cluster/stats
/elasticsearch/_cluster/settings
/elasticsearch/_cluster/pending_tasks
/elasticsearch/_cat
/elasticsearch/_cat/count
/elasticsearch/_cat/master
/elasticsearch/_cat/nodes
/elasticsearch/_cat/plugins
/elasticsearch/_cat/repositories
/elasticsearch/_cat/tasks
/elasticsearch/_cat/aliases
/elasticsearch/_cat/allocation
/elasticsearch/_cat/health
/elasticsearch/_nodes
/elasticsearch/_nodes/stats
/elasticsearch/_nodes/dr-bSDj2Qgq82-9-wXyFSw
/elasticsearch/_nodes/dr-bSDj2Qgq82-9-wXyFSw/stats
/elasticsearch/index1
/elasticsearch/index1/_field_stats
/elasticsearch/index1/_mapping
/elasticsearch/index1/_stats
/elasticsearch/index1/type1
/elasticsearch/index1/type1/_mapping
```

## Build

```
$ go get github.com/dutchcoders/elasticofs
$ go build -o /usr/bin/elasticofs
```

## Installation

```
$ ln -s /usr/bin/elasticofs /usr/sbin/mount.elasticofs
```

## Mount

```
$ mount -t elasticofs -o gid=0,uid=0 http://127.0.0.1:9000/ /elasticsearch
```

## Unmount

```
$ umount /elasticsearch
```

## Options

* **GID**: The default gid to assign for files from storage.
* **UID**: The default gid to assign for files from storage.

