---
{
  "title": "Simple Beginnings",
  "subtitle": "Super Ultra Hackerman 3000.",
  "date": "2026-08-19",
  "read_time": 6,
  "draft": false,
  "tags": ["console-hacking"],
}
---

Random memories from when I was a kid keep popping into my head lately. Must be an age thing.

This one's about messing with Xbox 360 saves. I'd have been 11 or 12.

I loved that console. By the end of that generation my Microsoft account reckoned I'd owned six of them, which sounds mad until you remember the Red Ring of Death was a thing. Me and my dad did the towel trick more than once. If you've never heard of it, you wrapped a dead Xbox in a towel, let it basically cook itself, and hoped whatever had come loose inside would stick back down long enough to get a few more days out of it. Stupid idea. Worked sometimes, though.

That was the generation that got me properly into games.

I remember bringing _Call of Duty 4: Modern Warfare_ home after begging my parents in Abbeycentre. I definitely shouldn't have been playing it at that age, but parents were a bit more lenient back then (or just had no idea what was in it), and I turned out fine.

Mostly.

Then _Call of Duty: World at War_ came out in 2008 and that was it for me. I was obsessed with anything WWII at the time. _Saving Private Ryan_, documentaries, the lot. A Call of Duty set in WWII was more or less made with me in mind.

At school a few of us used to compare how far we'd got in the campaign. Then one day the guy furthest ahead said he'd finished it and there was a zombies mode at the end.

Nobody believed him. Turns out he was telling the truth.

Zombies felt like something you weren't supposed to find. It was dark, properly difficult and nothing like the rest of the game, and before long it was all anybody talked about. We played it most nights until somebody's parents pulled the plug and sent them to bed. Fair enough, there was school in the morning.

## Infinite ammo

I was already fairly handy on a computer for a kid. But I think this was the first time I looked at a game and wondered if I could make it do stuff it wasn't meant to.

The first tutorial I followed wasn't anything fancy. It was infinite ammo for the campaign.

Copy the save onto a USB stick, plug it into the family laptop, pull `savegame.svg` out, open it in a hex editor. Then go hunting for this:

```text
player_sustainAmmo 1
```

That was it. A single `1`. Turn it on and the game stops taking bullets off you when you shoot. Other tutorials had huge numbers like `999999` floating around for health and all sorts, but infinite ammo was just a switch.

Put the save back on the Xbox, load it up, and... it worked. I couldn't believe it.

## What I was actually looking at

I had no clue what I was doing at the time, which is kind of the best part looking back.

A hex editor shows you a file as raw bytes. Hex on one side, and whatever those bytes are as text on the other. So `player_sustainAmmo 1` looked something like this:

```hex
70 6C 61 79 65 72 5F 73 75 73 74 61 69 6E 41 6D 6D 6F 20 31
 p  l  a  y  e  r  _  s  u  s  t  a  i  n  A  m  m  o     1
```

Every pair is one byte. `70` is a `p`, `31` is the character `1`. The editor was showing me the same data twice and I thought the numbers on the left were some sort of secret code.

And `player_sustainAmmo` wasn't something I'd added. It was one of the game's own developer variables (dvars), a setting Treyarch had built in for themselves, and the save just had it written down. I wasn't really hacking in infinite ammo. I'd found a switch that was already there and left it on.

Numbers were a bit weirder. Later on I'd come across values in saves like `3F 42 0F 00`, which is `999999`, stored backwards. The game reads it as a little-endian integer, meaning the smallest byte goes first. If you didn't know that (I definitely didn't), you'd type your big number in the wrong order and end up with something completely different to what you wanted.

None of that was going through my head at 12. What I did pick up without knowing it was the loop. Change something in the file, load the game, see what happens, go back and change it again when it inevitably didn't work.

I'd call that debugging now. At the time I just thought I was getting one over on Treyarch.

## Modded saves

Then I found the forums, and people posting their own custom _World at War_ saves. These were a different level. Someone had worked out that the game would read config strings out of a save and run commands it already knew about, so you could chain them together:

```text
bind BUTTON_BACK vstr mod2
mod2 = god; give ray_gun
```

`bind` hooks a button up to something, `vstr` runs another string by name, and the game already had commands for god mode, noclip and giving yourself guns. Stack enough of those together and you've got a mod menu. No new code anywhere, just a load of strings and button bindings set up carefully enough to feel like a menu.

Didn't matter to me how it worked. It was magic. I'd jump into a zombies lobby, play it straight for a few rounds, then hit Back and pull out a Ray Gun or walk through a wall. Keeping a straight face was the hard part.

The one annoying part was getting the Xbox to actually load the thing. Once you'd edited a save the console knew something was up, so it had to be rehashed and resigned before it'd take it. I did it in the same order every single time:

```text
edit → rehash → resign → USB → Xbox
```

Forget one and it's back to the laptop.

## Looking back

It's funny how much was going on in there when I think about it now. The save had a proper format, numbers were stored backwards, and the console wouldn't touch a file unless its hash and signature checked out. File formats, endianness, integrity checks. I was doing all of it to get infinite ammo in a game I shouldn't have been playing in the first place.

Nearly twenty years later I ended up hex editing a _Borderlands_ save, and it felt exactly the same. The tools are better and I know a bit more about what I'm looking at, but it's the same buzz.

I write software for a living now, and I love pretty much anything to do with it. I can't say for sure this is where that started. But whenever I think about how I got into it all, a USB stick, the family laptop and a hex editor are some of the first things that come to mind.

Good memory, anyway. Sorry, Treyarch.
