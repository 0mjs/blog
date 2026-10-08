---
{
  "title": "Gopher Nation",
  "subtitle": "Becoming a Gopher.",
  "date": "2025-10-28T23:13:06Z",
  "read_time": 7,
  "draft": false,
  "tags": ["golang"],
}
---

In my day-to-day job, TypeScript pays the bills. It keeps the lights on, if you will.

I've been working with Node.js for years now, and honestly, it's good. Using something like Nest.js every day gives a backend the kind of structure that you don't really appreciate until you go back and raw-dog an Express app and suddenly you're deciding where `user.service.ts` should live like it's 2018 again.

The ecosystem is enormous, everything has a package, and most of the boring stuff you actually need in a production backend - Kafka, CRON jobs, validation, queues, whatever - is basically a solved problem.

I like TypeScript.

But lately, for anything that's just mine, I keep reaching for something else.

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

<figure class="fig" data-reveal>
<div class="stacks"><div class="stack"><b>typescript on node</b><span style="--n:0">your code</span><i style="--n:1">↓</i><span style="--n:2">tsconfig.json</span><i style="--n:3">↓</i><span style="--n:4">tsc, or a bundler</span><i style="--n:5">↓</i><span style="--n:6">node_modules</span><i style="--n:7">↓</i><span style="--n:8">the right Node version</span><i style="--n:9">↓</i><span class="end" style="--n:10">your program</span></div><div class="stack"><b>go</b><span style="--n:0">your code</span><i style="--n:1">↓</i><span style="--n:2">go build</span><i style="--n:3">↓</i><span class="end" style="--n:4">your program</span></div></div>
<figcaption>Everything that sits between writing the thing and running it.</figcaption>
</figure>

It felt weirdly refreshing.

Later I started watching Rob Pike talk about the philosophy behind the language, particularly [Simplicity is Complicated](https://www.youtube.com/watch?v=rFejpH_tAHM), and that was probably the point where Go stopped feeling "primitive" and started feeling deliberately small.

<figure class="fig video" data-reveal>
<iframe src="https://www.youtube.com/embed/rFejpH_tAHM" title="Simplicity is Complicated - Rob Pike" loading="lazy" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" allowfullscreen></iframe>
<figcaption>Simplicity is Complicated · Rob Pike, dotGo 2015</figcaption>
</figure>

That's the bit I hadn't really understood.

Go isn't missing a load of features because nobody thought of them.

A lot of them were left out on purpose.

<aside class="lesson" data-label="// the gopher way"><strong>Less is a feature.</strong> Every feature a language adds is one more thing every reader of your code has to know. Go's bet is that a small language read by lots of people beats a big one only a few people fully understand.</aside>

There's no traditional class inheritance hierarchy. Interfaces are implicit. Composition is everywhere. Generics eventually arrived, but you can write a surprising amount of useful software without touching them.

<figure class="fig" data-reveal>
<div class="split stacked"><div class="box"><span>your type</span>type Shout struct{}
func (Shout) <em>Write</em>(p []byte) (int, error)</div><i>↓ fits ↓</i><div class="box"><span>io.Writer</span>type Writer interface {
    <em>Write</em>(p []byte) (n int, err error)
}</div></div>
<figcaption>No implements keyword anywhere. If the methods match, it's an io.Writer, and anything that writes to one will happily write to it.</figcaption>
</figure>

> "Go doesn't have type hierarchy because hierarchies are brittle. Composition is more flexible." - Rob Pike (Google I/O 2012)

<aside class="lesson" data-label="// the gopher way"><strong>Small interfaces, defined where they're used.</strong> Because a type never has to declare what it implements, you write the interface next to the code that needs it, and it tends to stay tiny. <code>io.Writer</code> is one method, and half the standard library speaks it.</aside>

And the standard library is ridiculous.

The first time you realise you can build a completely respectable HTTP API using `net/http` without immediately reaching for a framework feels almost wrong if you've spent years in Node.

```go
http.HandleFunc("GET /hello/{name}", func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello, %s\n", r.PathValue("name"))
})
log.Fatal(http.ListenAndServe(":8080", nil))
```

That's a web server. Routing, path parameters, the lot, straight out of the standard library.

There's an entire species of Go developer who will practically appear behind you and slap Gin out of your hand if you suggest installing it.

And I'm starting to understand them.

## What I Actually Like

Goroutines are probably the obvious one.

Concurrency is still concurrency. You can absolutely make a complete balls of it.

But the basic model feels natural.

```go
go doSomething()
```

There you go. It's doing something.

Channels took me slightly longer to get my head around, but once they click, you start seeing why Go ended up absolutely everywhere in infrastructure and distributed systems.

<figure class="fig" data-reveal>
<div class="chan-wrap"><div class="chan"><div class="side"><span>go worker(1)</span><span>go worker(2)</span><span>go worker(3)</span></div><div class="pipe"><b style="animation-delay: 0s">1</b><b style="animation-delay: 1.2s">2</b><b style="animation-delay: 2.4s">3</b></div><div class="side"><span class="recv">v := &lt;-ch</span></div></div></div>
<figcaption>Three goroutines, one channel. Whoever's receiving just takes the next value off the pipe.</figcaption>
</figure>

<aside class="lesson" data-label="// the gopher way"><strong>Share memory by communicating.</strong> One of the Go proverbs is "don't communicate by sharing memory; share memory by communicating". Instead of several threads locking the same value and fighting over it, you pass the value down a channel, and only one goroutine holds it at a time.</aside>

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

<figure class="fig" data-reveal>
<div class="errs"><div class="col"><b>an exception</b><div class="layers"><span class="hot">handler<em>catch</em></span><span>service</span><span>repository</span><span>client</span><span>library</span><span class="hot">driver<em>throw</em></span><i class="throw" aria-hidden="true"></i></div></div><div class="col"><b>a go error</b><div class="layers"><span class="hot">handler<em>handle it</em></span><span>service<em>return err</em></span><span>repository<em>return err</em></span><span>client<em>return err</em></span><span>library<em>return err</em></span><span>driver<em>return err</em></span></div></div></div>
<figcaption>Same failure, two journeys. One skips every layer in between; the other walks past each one in plain sight.</figcaption>
</figure>

<aside class="lesson" data-label="// the gopher way"><strong>Errors are values.</strong> They're just things you return, check, wrap and pass on, like any other value. Rob Pike wrote a <a href="https://go.dev/blog/errors-are-values">whole post</a> with exactly that title, and once it sinks in, <code>if err != nil</code> stops looking like noise and starts looking like the actual control flow.</aside>

The tooling is another huge part of it.

<figure class="fig" data-reveal>
<div class="chain-notes pairs"><p><b>go fmt</b>One format for everyone. Nothing to configure.</p><p><b>go test</b>Tests are just functions in a <code>_test.go</code> file.</p><p><b>go build</b>A single binary for whatever OS you ask for.</p><p><b>go mod</b>Dependencies, versions and checksums, built in.</p></div>
<figcaption>All of it ships with the language. There's nothing to choose.</figcaption>
</figure>

You don't spend an afternoon deciding which formatter your team should use.

Go already decided.

You don't spend another afternoon arguing about the formatter configuration.

Go doesn't care.

Shut up and write the code.

And then there's deployment.

Build binary. Copy binary. Run binary.

<figure class="fig" data-reveal>
<div class="chain"><span class="key" style="--n:0">go</span><i style="--n:1">·</i><span style="--n:2">build</span><i style="--n:3">→</i><span style="--n:4">copy</span><i style="--n:5">→</i><span class="out" style="--n:6">run</span></div>
<div class="chain long" style="margin-top: .9rem"><span class="key" style="--n:7">node</span><i style="--n:8">·</i><span style="--n:9">install Node</span><i style="--n:10">→</i><span style="--n:11">npm ci</span><i style="--n:12">→</i><span style="--n:13">build</span><i style="--n:14">→</i><span style="--n:15">copy + node_modules</span><i style="--n:16">→</i><span class="out" style="--n:17">run</span></div>
<figcaption>Shipping the same small service, both ways.</figcaption>
</figure>

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

## What Came of It

_Update, one year on._ I wrote everything above in October 2025. It's October 2026 now, so here's what came of it.

That weekend API turned into something bigger than I expected.

I liked `net/http` so much that I started building a thin layer on top of it. Routing, binding, validation, error handling, OpenAPI docs generated straight from your types: the bits I kept rewriting for every project. A Zinc app is still a plain `http.Handler`, so nothing you already have stops working.

It's called [Zinc](https://zinc.carbonsoft.sh), and it's open source.

You're reading this on it. This blog is a Zinc app, exported to static files when it ships.

And it's now on [HttpArena](https://www.http-arena.com/), an independent benchmark that runs dozens of frameworks, across loads of languages, through the same tests. It's not at the top. But it's on the same board as frameworks I've been reading about for years, and that still feels a bit mad.

Not bad for a language I once looked at and thought, "that's it?"

## If You're Curious

- _The Go Programming Language_ (Donovan & Kernighan) - probably still the best proper introduction I've found
- [Effective Go](https://go.dev/doc/effective_go) - free, official and worth reading even if parts of it are showing their age
- Rob Pike's [Concurrency is Not Parallelism](https://www.youtube.com/watch?v=oV9rvDllKEg) - one of those talks that makes a concept you've heard 400 times suddenly make sense

You don't need a 10-hour YouTube course.

Read enough to stop being completely lost.

Then build something you actually care about.

That's what made Go click for me.
