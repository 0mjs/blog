---
{
    "title": "Embracing Intangibility",
    "date": "2025-04-23T10:35:54Z"
}
---

In my early career, I was **enthralled** by the immediate feedback loop of frontend development. There's something uniquely satisfying about tweaking something and seeing instant feedback, often in realtime. In a very specific way, whatever way my passion for the subject began was more in how a "Login" page was created visually rather than how the server processed a `POST /auth/signin` request. TypeScript was certainly gaining traction at the time, but hadn't quite caught on yet, so plain old JavaScript, quirks and all, was the only learning curve to conquer, with the goal of building user interfaces on various teams. I perched happily on the summit of Mount Stupid, convinced quite quickly after learning the very basics of JavaScript and React that `.map()` and `.filter()` and knowing where and when to use `useState` or `useEffect` were peak computer-science. Plus, the developer experience I had initially was relatively primitive and pretty tangible... Need to debug? Fire in a `console.log()` and carry on. Unsure about your layout? Throw on a `border: 1px solid red` and adjust the CSS until it all lines up, as one does!

## The Tangibility of Frontend

Frontend work is gloriously tangible. The code you write turns into something you can see, click, and instantly break if you so choose. Change a colour, hit refresh, and there's your handiwork staring back at you. That immediacy gave me confidence, the browser becomes a canvas in a way, and every change leaves a certain brushstroke.

Now, lest the tooling brigade come for me: yes, there are plenty of **intangible** concepts even on the "front of the front-end". Bundlers, build pipelines, runtimes, all that oldschool webpack-shaped sorcery, absolutely. But for the purposes of this tale I'm talking about the bit the user actually interacts with in the browser, not the machinery that wrangles our code into shape.

## MAMP, End-to-End

In around 2016 when I was beginning my journey in professional Software Engineering, before meta-frameworks like Next and Nuxt had taken shape, my journey into full-stack began on a classic MAMP stack, where I was exposed to full-stack development in a professional setting for this first time.

I vividly recall crafting simple login pages with raw HTML and CSS, sprinkled with JavaScript for client-side validation where and when I could understand it enough to do so, and tying it all together with PHP to handle `POST` requests, sessions, and authentication logic. Seeing my CSS-styled form submit to a PHP script, verifying credentials against a MySQL database, and redirecting users upon success was my first taste of end-to-end development. It was here I learned the fundamentals of HTTP, form handling, and server-side sessions.

## The Dawn of React.js + Node.js

A quick tip of the hat to frameworks like Next.js. Those early versions around v9 and v10, were my first real glimpse of backend thinking inside a front-end comfort zone. Pages Router, file-based routing, and the wonderfully simple API routes felt like guard-rails for Express and Node.js without the boilerplate. Helpers like `Link`, `Image`, and `Head` took care of the messy bits so I could focus on wiring data together.

Dropping an ORM like Prisma or even the native MongoDB SDK into a `/api/*` route was when the penny finally dropped: "Ah, this is just a server under the hood." Fetching from a database, shaping the response, and watching it hydrate a React page all in one repo bridged a gap I hadn't fully understood. Next.js didn't just abstract the complexity, it introduced it in a way that made complete sense - all in one language. This gives developers space to explore common backend patterns while still shipping features that look great in the browser.

That gentle on-ramp turned out to be invaluable once I started designing standalone services and databases from scratch.

## REST API's & The Gateway Drug

Things changed when I discovered HTTP clients for the first time, such as Postman and later Insomnia which largely to this day remain the two most popular (AI-slop HTTP client anyone?). Suddenly I was captivated by the silent dance of requests and responses, back-and-forth, and how you could construct distinct API's across various routes. Building even basic CRUD APIs at the time was really interesting, I found it difficult to not be in that part of the code on any app or system I was working on at the time both professionally and in side projects.

This was my bridge to full-stack land with isomorphic JavaScript. I was still living in React, but my reach now extended to the server and databases. The feedback loop lengthened, but the results became far more powerful.

## Fading Tangibility & Broadening Horizons

Venturing deeper, first with vanilla Node.js, Express.js, and onto frameworks like Feathers.js and Nest.js, then Python's FastAPI, and eventually Go, 1.22's `net/http` and to Fiber and Echo. All of this coupled with the dawning of infrastructure as a concept to me - and something I felt almost obligated to do and learn. Docker needs to be configured, servers need to be hosted, and developers have to have some way of deploying new versions. I felt the tangibility fade. Debugging was no longer as simple as a quick page refresh and development no longer included whimsically altering the height of a delightfully over-engineered `TopNav.tsx` component. Instead it meant considering how systems interact, trawling logs, tracing request lifecycles, and following execution paths that quite often refused to show their faces.

Moving off NoSQL/Document databases for development also helped to further understand important relational concepts. Schema design, normalisation, query optimisation, these ideas live in the mind long before they emerge in running code. You can't just slap a red border on an index in order to remain steadfast that it's working (though, that could well be handy). Delbert, note that down.

![At once, sire](/static/img/delbert.jpg)

## The Invisible Architecture

My thinking shifted from concrete to conceptual. Frontend work is about visualising data, components, animations, and interactions you can *see*. Backend work is about structuring data, flows, how its transported, morphed, and persisted across an often unseen landscape. This concept only balloons in complexity further as other parts of a system are included, or are needed to be plugged in and talk to the rest of the system.

Instead of visualising a fancy UI, a form, or how I might validate it clientside, I was picturing endpoints, data structures and service diagrams. The payoff wasn't instant, but there's a sublime thrill in watching a well-designed system shoulder real-world load without breaking a sweat.

## Shiny Object Syndrome

These days the hype-cycle spins fast enough to make your head spin, frameworks are released and updated with breaking changes seemingly on a conveyor belt. Abstractions are marvellous, but they can hide so much that newcomers or what is now known as "vibe coders" miss the bigger picture, if the care for any "picture" at all other than the final one. It's a little worrying, but also a reminder that staying curious about what's under the hood is a really integral part of the job.

I mean no disregard to product-focused builders that use AI tools and are blissfully unaware of how the software their producing actually works, or who don't even care. In addition, I'm certainly not a wise old ex-Microsoft DOS engineer that can rightfully throw around my weight as an intellect in my field in absolution, nor would I. I still Google (or ask an LLM) the shape of that cool `curl` command variant you can do for only posting a JSON blob with a `POST` (it's `--json`, sigh) more often than I'd care to admit! The point is to keep lifting the curtain whenever you can and _learn_ about why things occur, to be a better engineer and service a product you are working on even better for your clients and users.

## A New Frontier

My next planned adventure is to *dabble* in DevOps, AWS, Terraform, Docker, Kubernetes, the corners of tech where even the logs have logs. These are fresh interests, not a full-blown career pivot. I'm simply curious to understand how the boxes fit together at scale. The feedback loop will likely stretch into *days* or *weeks*, and I might not know if a design decision was misguided until nasty *bleep* on the production environment taps me on the shoulder.

## Embracing the Intangible

This journey has taught me that growth is often an exercise in embracing ever-increasing abstraction. The skills honed on the frontend, attention to detail, empathy for the user, a good eye for polish, are still priceless. They're just joined by systems thinking, architectural design, and performance tuning.

Both worlds have their own charm. Frontend offers instant gratification and direct user impact. Backend and infrastructure offer scale, depth, and the quiet pride of elegant engineering.

So if you're eyeing a similar path, lean into the abstraction. Learn to enjoy the graceful flow of data as much as the perfect pixel. Build mental models in place for things you can't physically see or interact with. Embrace intangibility.