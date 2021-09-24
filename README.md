# Hybrid Logical Clock

# Usage

## Get time now
```go
ts := hlc.Now()
```

## Tick when sending or having a local event
```go
ts := hlc.Tick()
```

## Tick when receiving an event
```go
ts := hlc.Tick(hlc.Sync(receivedTS))
```