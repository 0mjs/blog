---
{
  "title": "Simple Beginnings",
  "subtitle": "Super Ultra Hackerman 3000.",
  "date": "2026-08-19",
  "read_time": 7,
  "draft": false,
  "tags": ["console-hacking"],
}
---

As I get older, memories from when I was a youngster seem to randomly pop into my head at any given moment, on any given day.

Recently I had one about "hacking" console games. I must've been 11 or 12.

I loved my Xbox 360. By the end of that generation, my Microsoft account said I'd owned six of the things. I don't think I was quite early enough for the 512MB Arcade model, so the first was probably a 20GB Pro / Premium, or whatever Microsoft was calling it that week.

That was the console generation that really got me into games.

I remember bringing _Call of Duty 4: Modern Warfare_ home after pleading with my parents in Abbeycentre. Looking back, I probably shouldn't have been playing it at that age, but thank the high Heavens parents were a bit more lenient and blissfully unaware back then. I got to experience it all and it didn't do me a bit of harm.

...or so I think.

A year later I was allowed _Call of Duty: World at War_. It landed at exactly the right time. I was obsessed with WWII: _Saving Private Ryan_, documentaries, the whole lot. Getting to play through the Pacific and Eastern Front in a Call of Duty game was basically peak childhood.

At school, a few of us would spend breaks comparing campaign progress. Then the friend furthest through it dropped this:

> "I finished it last night. There's Nazi Zombies at the end."

Calling his bluff would be an understatement. None of us believed him, not really. But it sounded just plausible enough that we all had to find out.

And somehow he wasn't talking absolute shite.

There really was a Zombies mode hidden after the campaign. In 2008, that was basically crack for a group of kids. Zombies were everywhere around then anyway — I'd already played _Dead Rising_, and _Left 4 Dead_ came out the same year — but this weird, dark survival mode bolted onto _World at War_ felt special.

Within a few months it was the thing everyone talked about and played every evening, usually far later than our parents appreciated.

I was already fairly handy with a laptop or PC for a kid that age. But this was probably the first time I'd looked at a game and been properly intrigued by the idea of breaking its rules.

## Changing numbers

The first bit that stuck with me was much more basic than downloading somebody else's ridiculous custom save. It was opening a save in a hex editor and changing a number.

You'd make a little controlled change in-game — spend a few rounds, fire a few shots, pick up some ammo — then compare the before and after files. Most of the file still looked like total nonsense. Somewhere in it, though, a small run of bytes might have moved in a way that lined up with the number you had changed.

Say an ammo counter was 36 and you found this nearby:

```text
24 00
```

If that field is an unsigned 16-bit integer stored little-endian, those two bytes mean `0x0024`: 36 in decimal. Change the counter by three in-game and you might see:

```text
27 00
```

That is `0x0027`, or 39.

The important bit is the **if**. Finding a value that looks right once proves almost nothing. It could be a copy, a checksum input, a timestamp coincidence, or just the wrong data entirely. But repeat the experiment, change only one thing at a time, and patterns start to become evidence.

That was the bit I found addictive. A meaningless blob was slowly becoming a map.

At the time, hexadecimal may as well have been hieroglyphics. I didn't know that hex was just a compact way to write bytes, or that `0x24` is `2 × 16 + 4`. I definitely didn't know what little-endian meant. I just knew that a number in the game had changed, a number in a file had changed too, and I could poke it and see what happened.

Looking back, that is a surprisingly good introduction to reverse engineering. Create two saves, make one deliberate change, binary-diff them, then form a theory about an offset. Test it. Be wrong. Try again.

It also explains why hex editing sounds more mystical than it is. A save is game state serialised into bytes. Those bytes only mean anything because the game has rules for reading them back. Work out the rules, change the bytes correctly, and you change the state. The hard bit is working out the rules.

There are plenty of ways for that tidy story to fall apart, of course. Fields can be packed, compressed, encrypted, duplicated or checksummed. The byte you want might not even exist in plain form. But none of that changes the basic loop: change something, compare, learn a little more.

## Then came the modded saves

Only after that did I find the properly daft _World at War_ saves floating around forums. This was the USB stick, Horizon, `savegame.svg`, rehash-and-resign era. And, critically, the Se7enSins tutorial by someone whose name probably contained `xXx`, `Modz`, or both.

The ritual went roughly like this:

1. Copy the save to a USB stick.
2. Open the Xbox container in Horizon and extract `savegame.svg`.
3. Edit the save's configuration strings.
4. Put it back, then rehash and resign the container.
5. Copy it to the Xbox and hope you'd not missed a step.

That felt unbelievably technical. The first time one worked, I'd join a Zombies lobby and play normally for a few rounds without saying anything. Then I'd press a button, turn on god mode, pull out a Ray Gun, or disappear through a wall. To everyone else it must've looked completely absurd. To me it felt like actual magic.

The clever part of those saves wasn't really hex editing. People had worked out that _World at War_ would load configuration strings from the save and execute commands the game already understood. A very stripped-down version looked like this:

```text
activeaction = bind BUTTON_START vstr mod1; bind BUTTON_BACK vstr mod2
mod1         = noclip; give ammo; vstr activeaction
mod2         = god; give ray_gun; vstr activeaction
```

`bind` attached an action to a controller button; `vstr` evaluated another named string. The game already knew commands such as `god`, `noclip`, `give ammo`, `give all`, and `give ray_gun`. The save just wired them together.

From there, people made entire mod menus. Press down, execute a different string, show a different set of options, rebind the buttons, then press A to run the current choice. No new UI system had been added. It was controller bindings, text and existing console commands stitched into something menu-shaped.

Looking at it now, it is basically a tiny state machine. At the time I didn't know what a state machine was either.

But I was accidentally meeting a whole pile of ideas:

- binary data and hexadecimal
- offsets, integers and endianness
- serialisation and file formats
- command parsing and indirection
- input mapping and state machines
- checksums, hashes and signatures

Even the annoying rehash-and-resign step was teaching me something. You couldn't casually edit a file inside an Xbox save container and put it straight back. The integrity data no longer matched the contents, and the container was associated with an Xbox profile. Horizon rebuilt that metadata so the console would accept it again.

Obviously, I didn't know any of that. As far as I was concerned, the sacred ritual was:

```text
edit file → rehash → resign → USB → Xbox
```

Miss one step and it didn't work.

## Nearly twenty years later

The reason this all came flooding back is that I ended up hex editing a _Borderlands_ save again, nearly twenty years later.

It was bizarrely familiar. The tools are nicer now and I understand a bit more of what I'm looking at, but the feeling is the same: change something, load the game, see whether you were right. The save has state. It gets serialised. The layout is the puzzle.

I wasn't building software back then, and I definitely wasn't the person who figured out those _World at War_ menus. Most of the genuinely clever work had already been done by somebody on a dodgy modding forum.

But I was getting a glimpse behind the curtain. Computers stopped feeling quite as magical. There were files, formats, numbers and rules. If you understood enough of the rules, you could make the machine do things it definitely wasn't intended to do.

I didn't become a software engineer because I changed a few bytes in an Xbox 360 save. But it was one of the first times a computer stopped feeling like a sealed box.

Turns out a USB stick, a hex editor and a game behaving in ways it absolutely shouldn't was a pretty decent introduction.
