# Guide

This guide shows how you can modify the topology.

Topologies are always specified in `GenerateAll()` of `topologies/topology_util/util.go`. Just follow the format of the
existing examples; you can set the name, number of clients, validators, shards along with some configurations. The
topologies are generated into JSON files into `topologies/` every time `GenerateAll()` is run. Therefore: do NOT modify
the JSON files, because they are overwritten everytime `GenerateAll()` is run! Instead, specify everything in
`GenerateAll()` itself.

For AWS, use `merkleAWS()` as a starting example if you want to create a new topology.

If you want to benchmark a topology with `launch.go`, do the following:
1. Modify `Topology` (`launch.go`)
1. Set `NumShards` (`launch.go`) to the same value as in the selected topology
1. Important: set the number of agents (clients) in `NumAgentsInstances` (`launch.go`).
Usually I set the number to 3 times the number of validator servers. So for example if we have 4 validators, we set
`NumAgentsInstances` to `3 * 4 * NumShards`
1. Launch `main()` in launch.go 

`launch.go` must be launched from the path of the repo.

If you have any other questions, just let me know.