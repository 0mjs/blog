---
{
  "title": "Hex Editing",
  "subtitle": "Simple beginnings.",
  "date": "2026-08-19",
  "read_time": 3,
  "draft": true,
  "tags": ["console-hacking"],
}
---

I've been thinking recently about _Call of Duty: World at War_ on the Xbox 360, and specifically the ridiculous modded Zombies saves we used to pass around.

The process felt unbelievably technical to me at the time. Put a save onto a USB stick, open it in Horizon, extract `savegame.svg`, change things in a hex editor, then rehash and resign it before putting it back onto the console. Most importantly, follow a Se7enSins tutorial written by somebody with a name like `xXModzKiller420Xx` and hope none of the MediaFire links were dead.

I still remember the first time it worked.

I'd join a four-player Zombies game and play normally for the first few rounds. I wouldn't let on that anything was different. Then, at the right moment, I'd turn on god mode, pull out some impossible weapon or disappear through a wall. To everyone else in the lobby it must have looked completely absurd. To me it felt like I'd discovered actual magic.

Looking back, I obviously hadn't engineered any of it. Some genuinely clever person had figured out that WaW trusted configuration strings stored inside its save file. The game had commands for binding controller buttons, changing variables and invoking developer functionality. With `bind`, `vstr` and a few HUD settings, people built entire scrolling "menus" out of text and controller inputs.

There wasn't really a custom UI hiding in the hex. It was closer to a tiny state machine assembled from bits the game already provided. Pressing down executed another string, printed another set of options and rebound the buttons. The save wasn't just progress either. It was structured data which the game loaded, interpreted and, crucially, trusted a little too much.

That's such a brilliant collection of ideas to accidentally encounter as a kid: file formats, serialization, command injection, hashes, signatures, indirection and interfaces made from completely unintended primitives. I didn't know the names for any of those things. I just knew that changing some bytes on a computer could make a supposedly sealed games console behave differently.

I recently considered trying to recreate the whole thing from scratch. I could make my own clean save, diff it, map the format and work out how the old menus were constructed. It would be a fun project, but probably also several evenings spent fighting obsolete Windows software and twenty-year-old forum archaeology.

I don't really need to rebuild it. Finally understanding the magic trick is enough.

I may only have been a kid following somebody else's tutorial, but that little glimpse behind the curtain stuck with me. It was probably one of the main reasons I became obsessed with computers, eventually made software my career and never stopped loving games.

Not bad for a USB stick and a file called `savegame.svg`.
