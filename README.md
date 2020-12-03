# Squads (Project Never)

This is an experiment/game involving turns-based strategy, and the Entity
Component System paradigm.

![20200121](https://user-images.githubusercontent.com/11085049/100997614-92e51380-35ae-11eb-961b-c1efc6c6f9fa.png "Combat")

## Install

```bash
go install github.com/griffithsh/squads/...
$GOPATH/bin/squads
```

## Playing with it

Pan the camera with the arrow keys. Zoom the camera with Z and X.

![20200121](https://user-images.githubusercontent.com/11085049/72780648-78929780-3c73-11ea-98a2-2d5f1f990625.png "Embarking")
![20200121](https://user-images.githubusercontent.com/11085049/72780650-792b2e00-3c73-11ea-9ca1-f11c06d08aa2.png "The Overworld")

## Profiling

```bash
$ go run . -cpuprofile pprof/cpu.prof
$ go tool pprof -http :8080 pprof/cpu.prof
```
