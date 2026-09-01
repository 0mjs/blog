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

As I get older, memories from being a youngster seem to show up at random.

This one was about messing with Xbox 360 saves. I must've been 11 or 12.

I loved that console. By the end of the generation, my Microsoft account reckoned I'd owned six of them, which feels less excessive if you remember the Red Ring of Death. My dad and I did the towel trick more than once: wrap a dead Xbox in a towel, make it far too hot, then hope the dodgy connection inside woke up long enough for another few days of play. It was a terrible idea. We did it anyway.

It was the generation that got me properly into games.

I remember bringing _Call of Duty 4: Modern Warfare_ home after pleading with my parents in Abbeycentre. Looking back, I probably shouldn't have been playing it. They were a bit more lenient and blissfully unaware back then. I got to experience it all and it didn't do me a bit of harm.

...or so I think.

Then came _Call of Duty: World at War_. I was obsessed with WWII at the time — _Saving Private Ryan_, documentaries, all of it — so it was basically made for me.

At school, a few of us would spend breaks comparing campaign progress. One day, the friend furthest through it said he'd finished it and there was a Zombies mode at the end.

We assumed he was making it up. He wasn't.

Zombies had this strange, hidden-away feel to it. It was dark, difficult and nothing like the campaign. Before long it was all anybody talked about, and we played it most evenings until somebody's parents put a stop to it (and sent them rightfully to bed, for school the next day).

I was already fairly handy with a laptop or PC for a kid that age. But this was probably the first time I'd looked at a game and been properly intrigued by the idea of breaking its rules.

The first tutorial I followed was not one of the massive Zombies mod menus. It was for infinite ammo in the campaign.

I copied a save to a USB stick, opened it on the family laptop, extracted `savegame.svg`, then opened it in a hex editor. The actual setting was almost disappointingly small:

```text
player_sustainAmmo 1
```

That `1` told the game not to take ammo away when I fired. I remember huge values like `999999` floating around in other tutorials, probably for health or jump height, but infinite ammo was just a switch.

The clever bit is that the hex editor was not showing some secret ammo language. It was showing the same text as bytes. In ASCII, the start of `player_sustainAmmo 1` looks like this:

```text
70 6C 61 79 65 72 5F 73 75 73 74 61 69 6E 41 6D 6D 6F 20 31
 p  l  a  y  e  r  _  s  u  s  t  a  i  n  A  m  m  o     1
```

Each pair is one byte. `70` is the letter `p`; `31` is the character `1`. The editor showed both columns, but I had no idea at the time that they were just two ways of looking at the same data.

`player_sustainAmmo` was a setting the game already knew about — one of its developer variables, or dvars. The save file contained a name and a value; when the game loaded it, it treated `1` as enabled. I wasn't adding infinite ammo to _World at War_. I was finding a switch that Treyarch had already built and leaving it on.

That tiny discovery connected a lot of dots without me realising. A save was not just a mysterious blob: it had a format. Hex was not a weird programming language: it was a way to display bytes. And those bytes could be text, numbers, or something the game used as a setting.

Later I would find actual numbers in saves too. If a four-byte field contains `3F 42 0F 00`, and the game reads it as a little-endian integer, that is `999999` in decimal. Little-endian just means the least-significant byte comes first. But the infinite-ammo tutorial taught me the more useful lesson first: before changing bytes, you need to know what the game thinks they are.

Back in those days I had no idea how to say any of that. I just knew I could change a thing in a file, put the save back on the Xbox, and suddenly never run out of bullets in the campaign.

That was enough to get its hooks into me.

I started to understand the rhythm of it, even if I had none of the words for it. Change the file. Load the game. See what happened. If it did not work, go back and change something else.

That is debugging, really. I just thought I was getting away with something.

The more advanced stuff came later. I found forums full of custom _World at War_ saves, made by people who had worked out that the game would load configuration strings from the save and run commands it already understood.

Something as simple as this could be enough to start a chain of nonsense:

```text
bind BUTTON_BACK vstr mod2
mod2 = god; give ray_gun
```

`bind` connected a controller button to an action. `vstr` ran another named string. The game already had commands for god mode, noclip and giving you ammo. Somebody had realised that a save file could wire those pieces together.

That is how the old mod menus worked. They were not new software bolted onto the game. They were a pile of strings and button bindings, arranged carefully enough to feel like a menu.

To me, it was still magic. I would join a Zombies lobby, play normally for a few rounds, then hit a button and pull out a Ray Gun or walk straight through a wall. It was difficult to keep a straight face.

There was one last annoying part. Editing the contents made the Xbox save container unhappy, so it had to be rehashed and resigned before it would load again. I treated that as a sacred ritual:

```text
edit → rehash → resign → USB → Xbox
```

Miss a step and it was back to the laptop.

Looking back now, it is funny how much was hiding in that little routine. The save had a format. The bytes had an order. The edit had to pass an integrity check. Without realising it, I was doing a small bit of reverse engineering.

I was not consciously learning any of it. I was a kid trying to get infinite ammo in a game I probably should not have been playing.

Nearly twenty years later, I ended up hex editing a _Borderlands_ save. The tools are better now and I understand a bit more of what I am looking at, but the feeling was exactly the same: there is some state in a file, the game has rules for reading it, and perhaps those rules can be nudged.

I did not become a software engineer because I changed a few bytes in an Xbox 360 save. But it was one of the first times a computer stopped feeling like a sealed box.

A USB stick and a hex editor were a pretty good introduction, especially when the game started behaving in ways it absolutely should not.
