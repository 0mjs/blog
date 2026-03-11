---
{
    "title": "Why Go?",
    "date": "2025-10-28T23:13:06Z",
    "read_time": 4,
    "draft": false,
    "tags": ["go", "typescript", "dx", "programming languages"]
}
---

In my day-to-day, TypeScript pays the bills. It keeps the lights on, as the saying goes.

Having worked with Node.js for so much of my career, the Nest.js framework gives things structure, polish, and a level of sanity you start to take for granted. The ecosystem is enormous, the tooling is solid, and everything feels battle-tested. It does exactly what you'd expect from a mature stack.

But when I'm building something for myself? I keep choosing to work in Go.

It took me a while to properly understand why.

## The Thing About TypeScript

TypeScript solved real problems. Type safety in JavaScript wasn't just nice to have, it was transformative. I remember the before-times: runtime errors that should've been caught at compile time, refactors powered by `grep`, and the constant uncertainty of what a function actually returned.

But there's something else that rarely gets said out loud: TypeScript is still a layer on top of JavaScript. A very good layer, but still a layer.

Every new feature adds another abstraction trying to shape JavaScript into something more structured. Generics that disappear at runtime. Decorators that stayed "experimental" for years on end. Async/await wrapping promises, wrapping callbacks. Enums that exist, but that you're often advised not to use.

None of this makes it bad. It just makes it heavier than it sometimes needs to be.

## When I Found Go

I picked up *The Go Programming Language* fully expecting to dislike it. Everyone said it was too simple. Too opinionated. Missing things I'd come to consider essential.

Then I spent a weekend building a small API, and something clicked.

The language just got out of my way. No build pipeline gymnastics. No `tsconfig.json`. No Babel. No debates about CJS versus ESM. Just code that compiled to a single binary and ran with one command. Fast.

Later, watching Rob Pike's talks helped it all make sense — especially [Simplicity is Complicated](https://www.youtube.com/watch?v=rFejpH_tAHM). The idea that constraints can be freeing. That removing features can sometimes create better systems.

Go doesn't have inheritance because inheritance introduces more problems than it solves. It has generics now, but most of the time you don't actually need them. Interfaces and concrete types are usually clearer and easier to reason about.

The standard library gives you what you need without forcing you into a framework treadmill. And if anything, there is actually a massive contingent of Go experts that swear against libraries or frameworks of any sort, especially for APIs, simply because the standard library is so powerful and feature-rich.

<div class="video-embed">
  <iframe
    src="https://www.youtube.com/embed/rFejpH_tAHM"
    title="Simplicity is Complicated — Rob Pike"
    frameborder="0"
    allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share"
    allowfullscreen
  ></iframe>
</div>

> "Go doesn't have type hierarchy because hierarchies are brittle. Composition is more flexible." — Rob Pike (Google I/O 2012)

## What I Actually Like

Goroutines make concurrency feel natural. Not easy, concurrency is never truly easy — but natural. No callback hell. No async/await acrobatics.

Error handling is verbose, but I've grown to really appreciate that. `if err != nil` everywhere means I always know where things can fail, and what happens when they do. No exceptions bubbling up from deep inside a dependency. No unhandled promises surfacing hours later.

Compilation speed matters more than I realised it would. Most of my Go projects build in seconds. The binary is self-contained and "just runs". No `node_modules`. No container unless I actually need one. No worrying about Node versions or navigating the package-manager hellscape.

It's a calmer way to build software.

## The Honest Part

I'm not ditching TypeScript. Realistically, I can't. Too much of the professional world runs on it. Nest.js is genuinely excellent at what it does. I'm not trying to fight that reality.

But when I'm prototyping something, building internal tools, or just want to focus on the problem instead of the tooling, I reach for Go almost every time.

It's not about one language being objectively better.  
It's about which one keeps me focused on building.

For me, that's Go.

## If You're Curious

- *The Go Programming Language* (Donovan & Kernighan) — genuinely worth reading cover to cover  
- [Effective Go](https://go.dev/doc/effective_go) — free, official, and teaches idiomatic Go  
- Rob Pike's [Concurrency is Not Parallelism](https://www.youtube.com/watch?v=oV9rvDllKEg) — ~30 minutes that will change how you think about systems

You don't need a 10-hour course.

Read the docs. 
**Build something real**, that you're passionate about. 
See if it clicks.
