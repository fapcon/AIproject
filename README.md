# Targets
    - make run-okx
    - make run-local-piston

    - make generate
    - make lint


### How to collect cpu and memory profile
    
#### CPU
    curl -o cpu.pprof http://localhost:9000/debug/pprof/profile?seconds=30

    go tool pprof cpu.pprof

    go tool pprof -http=:9090 cpu.pprof

#### Memory
    curl -o mem.pprof http://localhost:9000/debug/pprof/heap

    go tool pprof -http=:9090 mem.pprof

    


