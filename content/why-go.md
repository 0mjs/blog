---
{
    "title": "Why Go?",
    "date": "2025-10-28T23:13:06Z",
    "read_time": 4,
    "draft": false
}
---
I still write TypeScript every day. Nest.js pays the bills. The ecosystem is massive, the tooling is solid, and honestly, it gets the job done. But when I'm building something for myself? I keep choosing Go.

```go
package main

import "net/http"

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello World!"))
    })

    http.ListenAndServe(":8080", nil)
}
```


It took me a while to figure out why.

## The Thing About TypeScript

TypeScript solved real problems. Type safety in JavaScript was genuinely revolutionary. I remember the before times - runtime errors that could've been caught at compile time, refactoring that required searching the entire codebase with grep, the constant uncertainty about what a function actually returned.

But here's what nobody talks about: TypeScript is a bandaid on JavaScript. A really good bandaid, but still a bandaid. Every new feature they add is another layer of abstraction trying to make JavaScript behave like something it's not. Generics that disappear at runtime. Decorators that have been "experimental" for years. Async/await wrapping promises wrapping callbacks. Enums we're eventually told to avoid.

I'm not saying it's bad. I'm saying it's complicated in ways that aren't always necessary.

## When I Found Go

I picked up "The Go Programming Language" expecting to hate it. Everyone said it was too simple, too opinionated, missing features. Then I spent a weekend building an API and something clicked.

The language got out of my way. No [Webpack](https://webpack.js.org/) gymnastics, no `tsconfig.json`, no [Babel](https://babeljs.io/), no fighting with module resolution between CJS and ESM. Just code that compiled to a binary, and ran with a single command. Fast.

Rob Pike's talks helped it make sense. Especially [Simplicity is Complicated](https://www.youtube.com/watch?v=rFejpH_tAHM) - the idea that constraints can be freeing. Go doesn't have inheritance because inheritance causes more problems than it solves. It has generics now, but you rarely reach for them - interfaces and concrete types are usually clearer. The standard library gives you what you need without the framework churn.

## What I Actually Like

Goroutines make concurrency feel natural. Not easy - concurrency is never easy - but at least it's not callback hell or trying to reason about async/await chains. CSP got it right decades ago and Go actually implements it properly.

Error handling is verbose but I prefer it now. `if err != nil` everywhere means I know exactly where things can fail and what happens when they do. No exceptions bubbling up from somewhere deep in a library. No unhandled promise rejections in production at 2am.

Compilation speed matters more than I thought. My Go projects build in seconds. The resulting binary is self-contained, small enough, and just runs. No `node_modules`, no containers unless I actually need them, no worrying about Node versions or package manager hellscape.

## The Honest Part

I'm not ditching TypeScript. I can't. Too much of the professional world runs on it. Nest.js is genuinely good at what it does. The React ecosystem isn't going anywhere. I'm not fighting that fight.

But when I'm prototyping something, building tools for myself, or just want to think about the problem instead of the build pipeline? Go wins every time.

It's not about one being "better" in some abstract sense. It's about what lets me focus on actually solving problems instead of wrestling with tooling. For me, that's increasingly Go.

## If You're Curious

- The Go Programming Language (Donovan & Kernighan) - Actually worth reading cover to cover
- [Effective Go](https://go.dev/doc/effective_go) - Free, official, teaches you to write idiomatic code
- Rob Pike's [Concurrency is Not Parallelism](https://www.youtube.com/watch?v=oV9rvDllKEg) - ~30 minutes that'll change how you think about systems

You don't need a course. Read the docs, build something real, see if it clicks for you, like it did for me.
