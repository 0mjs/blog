---
{
  "title": "Simple Beginnings",
  "subtitle": "Super Ultra Hackerman 3000.",
  "date": "2026-08-19",
  "read_time": 10,
  "draft": false,
  "tags": ["console-hacking"],
}
---

Random memories from when I was a kid keep popping into my head lately. Must be an age thing.

This one's about messing with Xbox 360 saves. I'd have been 11 or 12. Looking back, it's also where I picked up a surprising amount of what I do for a living now, without having a clue I was learning any of it.

I loved that console. By the end of that generation my Microsoft account reckoned I'd owned six of them, which sounds mad until you remember the Red Ring of Death was a thing. Me and my dad did the towel trick more than once. If you've never heard of it, you wrapped a dead Xbox in a towel, let it basically cook itself, and hoped whatever had come loose inside would stick back down long enough to get a few more days out of it. Stupid idea. Worked sometimes, though.

<figure class="fig" data-reveal>
<div class="art"><svg class="rrod" viewBox="0 0 200 200" width="200" height="200" role="img" aria-label="A console power ring with the top-right quadrant dark and the other three lit red, around a green power button"><defs><clipPath id="rrod-tl"><rect x="0" y="0" width="94" height="94"/></clipPath><clipPath id="rrod-tr"><rect x="106" y="0" width="94" height="94"/></clipPath><clipPath id="rrod-bl"><rect x="0" y="106" width="94" height="94"/></clipPath><clipPath id="rrod-br"><rect x="106" y="106" width="94" height="94"/></clipPath></defs><circle cx="100" cy="100" r="92" fill="none" stroke="currentColor" stroke-opacity=".08" stroke-width="2"/><circle cx="100" cy="100" r="70" fill="none" stroke="currentColor" stroke-opacity=".14" stroke-width="11" clip-path="url(#rrod-tr)"/><g class="red" style="filter: drop-shadow(0 0 5px #ff3b30)"><circle cx="100" cy="100" r="70" fill="none" stroke="#ff3b30" stroke-width="11" clip-path="url(#rrod-tl)"/><circle cx="100" cy="100" r="70" fill="none" stroke="#ff3b30" stroke-width="11" clip-path="url(#rrod-bl)"/><circle cx="100" cy="100" r="70" fill="none" stroke="#ff3b30" stroke-width="11" clip-path="url(#rrod-br)"/></g><circle cx="100" cy="100" r="44" fill="currentColor" fill-opacity=".05" stroke="currentColor" stroke-opacity=".14"/><path d="M 100 80 V 100 M 88 86 A 18 18 0 1 0 112 86" fill="none" stroke="#3ddc68" stroke-width="5" stroke-linecap="round" style="filter: drop-shadow(0 0 4px #3ddc68)"/></svg></div>
<figcaption>Three red quadrants and one dark, around a perfectly happy green button. Every 360 owner knew exactly what that meant.</figcaption>
</figure>

<aside class="lesson">The Red Ring was widely put down to heat. The console warmed up and cooled down every time you played, and over time the joints under its chips cracked. The towel trick cooked it hot enough to soften those joints back into place for a while, by doing more of the exact thing that broke it. <strong>That's a hotfix that treats the symptom and feeds the root cause.</strong> I've seen plenty of those since, in code. It's always the towel trick.</aside>

That was the generation that got me properly into games.

I remember bringing _Call of Duty 4: Modern Warfare_ home after begging my parents in Abbeycentre. I definitely shouldn't have been playing it at that age, but parents were a bit more lenient back then (or just had no idea what was in it), and I turned out fine.

Mostly.

Then _Call of Duty: World at War_ came out in 2008 and that was it for me. I was obsessed with anything WWII at the time. _Saving Private Ryan_, documentaries, the lot. A Call of Duty set in WWII was more or less made with me in mind.

At school a few of us used to compare how far we'd got in the campaign. Then one day the guy furthest ahead said he'd finished it and there was a zombies mode at the end.

Nobody believed him. Turns out he was telling the truth.

<figure class="fig" data-reveal>
<div class="art"><svg class="tally" viewBox="0 0 220 150" width="220" height="150" role="img" aria-label="Five red chalk tally marks: four upright and one struck across them"><g fill="none" stroke="#c62828" stroke-width="9" stroke-linecap="round" opacity=".9"><path pathLength="1" style="--n:0" d="M 52 28 Q 49 75 54 122"/><path pathLength="1" style="--n:1" d="M 86 24 Q 90 74 85 124"/><path pathLength="1" style="--n:2" d="M 120 27 Q 116 76 121 121"/><path pathLength="1" style="--n:3" d="M 154 25 Q 157 72 152 123"/><path pathLength="1" style="--n:4" d="M 30 104 Q 110 70 190 42"/></g></svg></div>
<figcaption>Red chalk tallies, round by round. By the time they turned into numbers, things were getting serious.</figcaption>
</figure>

Zombies felt like something you weren't supposed to find. It was dark, properly difficult and nothing like the rest of the game, and before long it was all anybody talked about. We played it most nights until somebody's parents pulled the plug and sent them to bed. Fair enough, there was school in the morning.

## Infinite ammo

I was already fairly handy on a computer for a kid. But I think this was the first time I looked at a game and wondered if I could make it do stuff it wasn't meant to.

The first tutorial I followed wasn't anything fancy. It was infinite ammo for the campaign.

Copy the save onto a USB stick, plug it into the family laptop, pull `savegame.svg` out, open it in a hex editor. Then go hunting for this:

```text
player_sustainAmmo 1
```

That was it. A single `1`. Turn it on and the game stops taking bullets off you when you shoot. Infinite ammo was just a switch.

Put the save back on the Xbox, load it up, and... it worked. I couldn't believe it.

And `player_sustainAmmo` wasn't something I'd added. It was one of the game's own developer variables (dvars), a setting Treyarch had built in for themselves, and the save just had it written down. I wasn't really hacking in infinite ammo. I'd found a switch that was already there and left it on.

<figure class="fig" data-reveal>
<div class="split"><div class="box"><span>savegame.svg</span>player_sustainAmmo <em>1</em></div><i>→</i><div class="box"><span>every time you fire</span>if <em>sustainAmmo</em> is off:<br/>&nbsp;&nbsp;take a bullet</div></div>
<figcaption>The behaviour was already in the game. The save just told it which way the switch was set.</figcaption>
</figure>

<aside class="lesson">That's a <strong>feature flag</strong>: behaviour that ships switched off and gets turned on by configuration, not by changing the code. Teams use exactly this to roll things out gradually, test in production, or hide unfinished work. Treyarch's switches just happened to be sitting in a file a 12-year-old could open.</aside>

## What I was actually looking at

I had no clue what I was doing at the time, which is kind of the best part looking back.

A hex editor shows you a file as raw bytes. Hex on one side, and whatever those bytes are as text on the other. So `player_sustainAmmo 1` looked something like this:

<figure class="fig" data-reveal>
<div class="bytes"><span class="byte" style="--n:0"><b>70</b><i>p</i></span><span class="byte" style="--n:1"><b>6C</b><i>l</i></span><span class="byte" style="--n:2"><b>61</b><i>a</i></span><span class="byte" style="--n:3"><b>79</b><i>y</i></span><span class="byte" style="--n:4"><b>65</b><i>e</i></span><span class="byte" style="--n:5"><b>72</b><i>r</i></span><span class="byte" style="--n:6"><b>5F</b><i>_</i></span><span class="byte" style="--n:7"><b>73</b><i>s</i></span><span class="byte" style="--n:8"><b>75</b><i>u</i></span><span class="byte" style="--n:9"><b>73</b><i>s</i></span><span class="byte" style="--n:10"><b>74</b><i>t</i></span><span class="byte" style="--n:11"><b>61</b><i>a</i></span><span class="byte" style="--n:12"><b>69</b><i>i</i></span><span class="byte" style="--n:13"><b>6E</b><i>n</i></span><span class="byte" style="--n:14"><b>41</b><i>A</i></span><span class="byte" style="--n:15"><b>6D</b><i>m</i></span><span class="byte" style="--n:16"><b>6D</b><i>m</i></span><span class="byte" style="--n:17"><b>6F</b><i>o</i></span><span class="byte" style="--n:18"><b>20</b><i>␠</i></span><span class="byte hot" style="--n:19"><b>31</b><i>1</i></span></div>
<figcaption>Twenty bytes. The one I cared about is the last: 31, the character "1".</figcaption>
</figure>

Every pair is one byte. `70` is a `p`, `31` is the character `1`. The editor was showing me the same data twice and I thought the numbers on the left were some sort of secret code.

<aside class="lesson">This is <strong>character encoding</strong>. Text is just numbers plus an agreed lookup table (ASCII then, UTF-8 now). The <code>31</code> there isn't the number thirty-one, or even the number one. It's the <em>character</em> "1", which the game later reads and turns into a number. Value versus representation: it's everywhere in software, and it's behind a good half of the bugs I've ever chased. (Later, my ammo of 31 would turn out to be the byte <code>1F</code>. Same-looking number, completely different byte, because one's text and one isn't.)</aside>

## Doing it myself

The infinite ammo switch was somebody else's find, though. I'd just followed the steps. What I really wanted was to find something on my own.

So I went after the ammo counter itself. Note how much ammo I had, say 36, save, and pull the save onto the laptop. 36 in hex is `24` (I had to look that up online), so I'd search the file for `24`. Loads of hits, obviously. Then back on the Xbox, fire a few rounds until I was down to 31, save again, look up 31 (`1F`) and search the new file for that.

Most of those `24`s were still sitting exactly where they'd been. One of them had turned into a `1F`. That was it. That was the ammo.

<figure class="fig" data-reveal>
<div class="hunt"><div class="hunt-row"><span>save 1 · ammo 36</span><span class="cells"><b>00</b><b>1A</b><b class="match">24</b><b>00</b><b>07</b><b class="it">24</b><b>3C</b><b>00</b><b class="match">24</b></span></div><div class="hunt-row"><span>save 2 · ammo 31</span><span class="cells"><b>00</b><b>1A</b><b class="match">24</b><b>00</b><b>07</b><b class="it">1F</b><b>3C</b><b>00</b><b class="match">24</b></span></div><div class="hunt-key"><span><b>24</b> = 36</span><span><b class="it">1F</b> = 31</span><span>every 24 that didn't change was a red herring</span></div></div>
<figcaption>Two saves, one difference. The bytes around it are made up; the method is the real bit.</figcaption>
</figure>

From there it should have been easy. Look up what 999999 is in hex (`0F 42 3F`), type it over the top, put the save back on the Xbox and load in.

It wasn't. The number I typed in wasn't always the number I got back, and every so often the save just wouldn't load at all. It took a fair few goes, and a fair few trips back to the laptop, before an edit actually stuck. But when it did, I hadn't followed anyone's tutorial. I'd found it myself, and I've been chasing that feeling ever since.

<aside class="lesson">That's <strong>search and elimination</strong>: take a snapshot, change one thing, take another, keep only what changed with it. It's how memory scanners find values, how <code>git bisect</code> finds the commit that broke something, and honestly how most debugging works. Shrink the haystack until there's one needle left.</aside>

<aside class="lesson">Numbers coming back wrong is usually down to <strong>data types</strong>, which I didn't know were a thing. A number gets a fixed amount of room, and it's either <em>unsigned</em> (zero and up) or <em>signed</em> (it can go negative). The engine World at War is built on keeps your ammo in memory as a signed 4-byte integer, which can hold over two billion. But it sends that number from the game's server side to the bit that draws your HUD as a signed 2-byte <em>short</em>, which tops out at 32,767 (<a href="https://github.com/callofduty4x/CoD4x_Server/blob/master/src/player.h">the struct</a>, <a href="https://github.com/callofduty4x/CoD4x_Server/blob/master/src/msg.c">the network code</a>, from Call of Duty 4's engine). Anything bigger gets chopped down to its last two bytes on the way, and if what's left is over 32,767 it wraps round to negative. Write more bytes into a file than a value has room for, and the extras land on whatever was stored next to it, which is a very good way to crash a game.</aside>

<figure class="fig" data-reveal>
<table class="then-now"><thead><tr><th>where the ammo is</th><th>room</th><th>biggest value</th></tr></thead><tbody><tr><td>in memory, as an <code>int</code></td><td>4 bytes</td><td>2,147,483,647 · <code>7F FF FF FF</code></td></tr><tr><td>sent to your HUD, as a <code>short</code></td><td>2 bytes</td><td>32,767 · <code>7F FF</code></td></tr></tbody></table>
<div class="hunt" style="margin-top: 1.25rem"><div class="hunt-row"><span>2 bytes, maxed out</span><span class="cells"><b class="match">7F</b><b class="match">FF</b><em>= 32,767</em></span></div><div class="hunt-row"><span>add one</span><span class="cells"><b class="it">80</b><b class="it">00</b><em class="bad">= −32,768</em></span></div><div class="hunt-row"><span>999999, squeezed</span><span class="cells"><b>0F</b><b class="match">42</b><b class="match">3F</b><em>= 16,959</em></span></div></div>
<figcaption>Stored as an int, sent as a short. Go one past 32,767 and it wraps to the most negative number there is; send 999999 and only its last two bytes make it.</figcaption>
</figure>

<aside class="lesson">There's one more wrinkle I only learned about later: <strong>byte order</strong>. A number that takes up several bytes can be written biggest byte first (big-endian) or smallest byte first (little-endian), and it's the file format that decides. The 360's own processor was big-endian, but a format doesn't have to follow the hardware. Get the order wrong and you don't get your number back to front. You get a completely different number.</aside>

<figure class="fig" data-reveal>
<div class="hunt"><div class="hunt-row"><span>big-endian</span><span class="cells"><b class="match">0F</b><b class="match">42</b><b class="match">3F</b><em>= 999,999</em></span></div><div class="hunt-row"><span>little-endian</span><span class="cells"><b class="match">3F</b><b class="match">42</b><b class="match">0F</b><em>= 999,999</em></span></div><div class="hunt-row"><span>read the wrong way</span><span class="cells"><b class="it">3F</b><b class="it">42</b><b class="it">0F</b><em class="bad">= 4,145,679</em></span></div></div>
<figcaption>The same number, two byte orders. Mix them up and 999,999 quietly becomes four million.</figcaption>
</figure>

None of that was going through my head at 12. What I did pick up without knowing it was the loop. Change something in the file, load the game, see what happens, go back and change it again when it inevitably didn't work.

<figure class="fig" data-reveal>
<div class="loop-art"><span style="top: 5%; left: 50%; animation-delay: 0s">edit the save</span><span style="top: 50%; left: 95%; animation-delay: 1.5s">load the game</span><span style="top: 95%; left: 50%; animation-delay: 3s">see what happens</span><span style="top: 50%; left: 5%; animation-delay: 4.5s">work out why</span><i class="loop-dot" aria-hidden="true"></i></div>
<figcaption>Round and round. Every bug I've fixed since has gone through this loop.</figcaption>
</figure>

I'd call that debugging now. At the time I just thought I was getting one over on Treyarch.

<aside class="lesson">This is the <strong>feedback loop</strong> at the heart of programming: guess, change, run, look. Tests, hot reload, REPLs and fast CI all exist for one reason, which is to make that loop go round quicker. Mine had a USB stick in the middle of it, so it was not quick.</aside>

## Modded saves

Then I found the forums, and people posting their own custom _World at War_ saves. These were a different level. Someone had worked out that the game would read config strings out of a save and run commands it already knew about, so you could chain them together:

```text
bind BUTTON_BACK vstr mod2
mod2 = god; give ray_gun
```

`bind` hooks a button up to something, `vstr` runs another string by name, and the game already had commands for god mode, noclip and giving yourself guns. Stack enough of those together and you've got a mod menu. No new code anywhere, just a load of strings and button bindings set up carefully enough to feel like a menu.

<figure class="fig" data-reveal>
<div class="chain long"><span class="key" style="--n:0">Back button</span><i style="--n:1">→</i><span style="--n:2">bind</span><i style="--n:3">→</i><span style="--n:4">vstr mod2</span><i style="--n:5">→</i><span style="--n:6">god; give ray_gun</span><i style="--n:7">→</i><span class="out" style="--n:8">god mode</span><span class="out" style="--n:9">Ray Gun</span></div>
<figcaption>One button press, two levels of indirection, two of the game's own commands.</figcaption>
</figure>

<aside class="lesson">That's <strong>indirection</strong>: a name that points at code you'll run later. Chain enough of them and you've built a small program out of the game's own commands. That's the same idea behind shell scripts, the command pattern and, if you squint, every plugin system I've used since. The people writing those saves were programming. They just weren't calling it that.</aside>

Didn't matter to me how it worked. It was magic. I'd jump into a zombies lobby, play it straight for a few rounds, then hit Back and pull out a Ray Gun or walk through a wall. Keeping a straight face was the hard part.

The one annoying part was getting the Xbox to actually load the thing. Once you'd edited a save the console knew something was up, so it had to be rehashed and resigned before it'd take it. I did it in the same order every single time:

<figure class="fig" data-reveal>
<div class="chain"><span class="key" style="--n:0">edit</span><i style="--n:1">→</i><span style="--n:2">rehash</span><i style="--n:3">→</i><span style="--n:4">resign</span><i style="--n:5">→</i><span style="--n:6">USB</span><i style="--n:7">→</i><span class="out" style="--n:8">Xbox</span></div>
<div class="chain-notes"><p><b>rehash</b>The save carried fingerprints of its own contents. Change a byte and they no longer match, so they had to be worked out again.</p><p><b>resign</b>Then the package's signature had to be redone, so the console would accept it as genuine.</p><p><b>forget one</b>The console refuses the save, and it's back to the laptop.</p></div>
<figcaption>Same order, every single time.</figcaption>
</figure>

<aside class="lesson"><strong>Hashes catch changes; signatures prove who made something.</strong> I learned the order before I learned the words. HTTPS, signed git commits, package managers and app stores all rest on those two ideas. The console was doing what every secure system does: refusing to trust a file it couldn't check.</aside>

## Looking back

It's funny how much was going on in there when I think about it now. The save had a proper format, values lived at fixed spots you could hunt down, and the console wouldn't touch a file unless its hash and signature checked out. I was doing all of it to get infinite ammo in a game I shouldn't have been playing in the first place.

<figure class="fig" data-reveal>
<table class="then-now"><thead><tr><th>What I was doing</th><th>What it's called now</th></tr></thead><tbody><tr><td>The towel trick</td><td>A hotfix that hides the root cause</td></tr><tr><td>Flipping <code>player_sustainAmmo</code></td><td>Feature flags and config</td></tr><tr><td>Reading hex next to text</td><td>Binary formats and character encoding</td></tr><tr><td>Hunting the ammo with two saves</td><td>Search, diffing and elimination</td></tr><tr><td>Numbers that came back wrong</td><td>Data types, overflow and byte order</td></tr><tr><td>Edit, load, look, repeat</td><td>The debugging feedback loop</td></tr><tr><td>Chaining <code>bind</code> and <code>vstr</code></td><td>Indirection and scripting</td></tr><tr><td>Rehash, then resign</td><td>Hashing and digital signatures</td></tr></tbody></table>
<figcaption>Nobody told me any of these words. I just kept poking at a save file.</figcaption>
</figure>

Nearly twenty years later I ended up hex editing a _Borderlands_ save, and it felt exactly the same. The tools are better and I know a bit more about what I'm looking at, but it's the same buzz.

I write software for a living now, and I love pretty much anything to do with it. I can't say for sure this is where that started. But whenever I think about how I got into it all, a USB stick, the family laptop and a hex editor are some of the first things that come to mind.

Good memory, anyway. Sorry, Treyarch.
