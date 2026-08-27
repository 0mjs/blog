---
{
  "title": "Gopher Nation",
  "subtitle": "Becoming a Gopher.",
  "date": "2025-10-28T23:13:06Z",
  "read_time": 4,
  "draft": false,
  "tags": ["golang"],
}
---

In my day-to-day job, TypeScript pays the bills. It keeps the lights on, if you will.

I've been working with Node.js for years now, and honestly, it's good. Using something like Nest.js every day gives a backend the kind of structure that you don't really appreciate until you go back and raw-dog an Express app and suddenly you're deciding where `user.service.ts` should live like it's 2018 again.

The ecosystem is enormous, everything has a package, and most of the boring stuff you actually need in a production backend - Kafka, CRON jobs, validation, queues, whatever - is basically a solved problem.

I like TypeScript.

Unfortunately, I also increasingly feel like this:

![Go TS Meme](/assets/image/meme/go-ts-node.jpg)

And it took me a while to work out why.

## The Thing About TypeScript

TypeScript solved a massive problem.

I remember writing JavaScript before TypeScript properly took over. Refactoring things with a combination of `grep`, hope and prayer. Functions returning objects that looked vaguely like the object you expected. Finding out what a value actually was by sticking a `console.log` in front of it and running the thing.

TypeScript made JavaScript dramatically better.

But at the end of the day, you're still trying to discipline JavaScript.

And sometimes you can feel it.

You've got types that disappear entirely at runtime. Decorators. `tsconfig.json` options that apparently alter the fabric of reality. CJS versus ESM. `moduleResolution`. Five different ways to import the same package depending on what phase the moon is in.

Then there are enums, which TypeScript added, and which half the TypeScript community will immediately tell you not to use.

None of this makes TypeScript bad. I use it professionally every day and will continue to.

There's just... a lot going on.

Sometimes I want to write a little program without installing 800 packages and accidentally downloading half of GitHub into `node_modules`.

## When I Found Go

I'd looked at Go a few times over the years and basically thought:

> That's it?

The syntax looked almost suspiciously boring.

Where's all the stuff?

Then I bought _The Go Programming Language_ and actually sat down with it properly.

I built a small API one weekend and somewhere along the way it just clicked.

There was no framework. No build pipeline. No Babel. No `tsconfig.json`. No package manager discourse. No wondering whether the project was ESM, CommonJS, ESM pretending to be CommonJS, or CommonJS wearing an ESM hat.

I wrote some Go.

I ran:

```sh
go build
```

And it gave me a binary.

That was basically it.

It felt weirdly refreshing.

Later I started watching Rob Pike talk about the philosophy behind the language, particularly [Simplicity is Complicated](https://www.youtube.com/watch?v=rFejpH_tAHM), and that was probably the point where Go stopped feeling "primitive" and started feeling deliberately small.

That's the bit I hadn't really understood.

Go isn't missing a load of features because nobody thought of them.

A lot of them were left out on purpose.

There's no traditional class inheritance hierarchy. Interfaces are implicit. Composition is everywhere. Generics eventually arrived, but you can write a surprising amount of useful software without touching them.

And the standard library is ridiculous.

The first time you realise you can build a completely respectable HTTP API using `net/http` without immediately reaching for a framework feels almost wrong if you've spent years in Node.

There's an entire species of Go developer who will practically appear behind you and slap Gin out of your hand if you suggest installing it.

And I'm starting to understand them.

<div class="video-embed">
  <iframe
    src="https://www.youtube.com/embed/rFejpH_tAHM"
    title="Simplicity is Complicated - Rob Pike"
    frameborder="0"
    allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share"
    allowfullscreen
  ></iframe>
</div>

> "Go doesn't have type hierarchy because hierarchies are brittle. Composition is more flexible." - Rob Pike (Google I/O 2012)

## What I Actually Like

Goroutines are probably the obvious one.

Concurrency is still concurrency. You can absolutely make a complete balls of it.

But the basic model feels natural.

```go
go doSomething()
```

There you go. It's doing something.

Channels took me slightly longer to get my head around, but once they click, you start seeing why Go ended up absolutely everywhere in infrastructure and distributed systems.

I also thought I'd hate the error handling.

And to be fair, when you first see this repeated 900 times:

```go
if err != nil {
    return err
}
```

you do wonder if the language designers were taking the piss.

But I've grown to really like it.

Errors are just there. In front of you. You can see where something fails and you decide what happens next.

There's no exception suddenly flying out of six layers of abstraction because a library decided this particular Tuesday was a good day to throw.

The tooling is another huge part of it.

- `go fmt`
- `go test`
- `go build`
- `go mod`

They're just... there.

You don't spend an afternoon deciding which formatter your team should use.

Go already decided.

You don't spend another afternoon arguing about the formatter configuration.

Go doesn't care.

Shut up and write the code.

And then there's deployment.

Build binary. Copy binary. Run binary.

Beautiful.

No `node_modules`. No production Node version to worry about. No `npm install` on a server. Half the time I don't even bother with Docker unless I actually need Docker.

It's just a calmer way to build software.

## The Honest Part

I'm not becoming one of those people who announces they've "left TypeScript".

I haven't.

TypeScript is still what I use professionally and Nest.js is genuinely excellent. If somebody asked me to build a large product backend tomorrow with a team of TypeScript developers, I'm not going to burst through the wall dressed as the Go gopher and demand we rewrite everything.

That would be insane.

But for my own stuff?

Go keeps winning.

Small APIs. Internal tools. CLI programs. Random ideas I want to get running without spending the first hour assembling a JavaScript project like IKEA furniture.

I reach for Go more and more.

I've even built some fairly substantial things with it now, and every time I go back to it I get that same feeling: there just isn't much between me and the actual thing I'm trying to build.

That's probably what I like most about it.

Not that Go is "better" than TypeScript.

That's boring language-war drivel.

It's that Go seems almost aggressively uninterested in being clever.

And after years of working in ecosystems where cleverness has a habit of becoming somebody else's maintenance problem six months later, I've started to appreciate that quite a lot.

## If You're Curious

- _The Go Programming Language_ (Donovan & Kernighan) - probably still the best proper introduction I've found
- [Effective Go](https://go.dev/doc/effective_go) - free, official and worth reading even if parts of it are showing their age
- Rob Pike's [Concurrency is Not Parallelism](https://www.youtube.com/watch?v=oV9rvDllKEg) - one of those talks that makes a concept you've heard 400 times suddenly make sense

You don't need a 10-hour YouTube course.

Read enough to stop being completely lost.

Then build something you actually care about.

That's what made Go click for me.
