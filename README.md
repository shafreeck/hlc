# Hybrid Logical Clock

# Usage

## Get time now
```go
ts := hlc.Now()
```

## Tick when send or has local event
```go
ts := hlc.Tick()
```

## Tick when receives an event
```go
ts := hlc.Tick(hlc.Sync(receivedTS))
```