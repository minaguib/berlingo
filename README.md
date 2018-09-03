abandonware notice
==================

With berlin-ai.com itself falling off the map, this client is now abandonware and exists only to serve historical curiosities.

berlingo
========

A Go framework for writing AIs for http://berlin-ai.com/

The nomenclature and object names/attributes in this framework tries to mirror as much as possible the ones in the official berlin-ai ruby client https://github.com/thirdside/berlin-ai


Usage
=====
Implement your AI algorithm by creating a new struct that satisfies the AI interface:
```go
type AwesomeAi struct{}

func (ai *AwesomeAi) GameStart(game *berlingo.Game) {
}

func (ai *AwesomeAi) Turn(game *berlingo.Game) {
}

func (ai *AwesomeAi) GameOver(game *berlingo.Game) {
}

func (ai *AwesomeAi) Ping(game *berlingo.Game) {
}

```

Invoke the AI by passing it to _berlingo.Serve_

For a working example which moves soldiers randomly, see cmd/berlingo-bot-random/main.go example

API documentation at: http://godoc.org/github.com/minaguib/berlingo

Invocation
==========
The example cmd/berlingo-bot-random/main.go and any similar AIs you write can be invoked in any of the following manners:
 * ./botname            - this will start the AI in single-request local mode - the request JSON will be read from STDIN
 * ./botname filename   - this will start the AI in single-request local mode - the request JSON will be read from the given filename
 * ./botname port       - this will start the AI in multi-request web mode    - the request JSON will be received over HTTP

In web mode, you may POST request JSONs to your bot, or sign your AI up at http://berlin-ai.com to duel other AIs
