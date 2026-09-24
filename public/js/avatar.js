// Hex portrait: the headshot decodes out of a hex dump of its own pixels,
// and falls back into one on hover. Without JS (or with reduced motion) the plain <img> stays.
(() => {
  const wrap = document.querySelector(".avatar");
  const img = wrap && wrap.querySelector("img");
  if (!img || matchMedia("(prefers-reduced-motion: reduce)").matches) return;

  const N = 16; // cells per side: 256 cells, one byte each
  const HEX = "0123456789ABCDEF";
  const canvas = document.createElement("canvas");
  const ctx = canvas.getContext("2d");
  canvas.setAttribute("aria-hidden", "true");

  let size, cell, photo, cells, raf;
  let progress = 0; // 0 = hex dump, 1 = photo
  let target = 1;
  let last = 0;

  // Sample the image down to N×N: each cell keeps its colour and a luminance byte.
  function sample() {
    const small = document.createElement("canvas");
    small.width = small.height = N;
    const sctx = small.getContext("2d");
    sctx.imageSmoothingQuality = "high";
    sctx.drawImage(img, 0, 0, N, N);
    const data = sctx.getImageData(0, 0, N, N).data;
    cells = [];
    for (let i = 0; i < N * N; i++) {
      const [r, g, b] = data.subarray(i * 4, i * 4 + 3);
      const x = i % N, y = Math.floor(i / N);
      const lum = Math.round(0.2126 * r + 0.7152 * g + 0.0722 * b);
      // Decode from the centre outwards, with some noise so it doesn't look like a wipe.
      const dist = Math.hypot(x - (N - 1) / 2, y - (N - 1) / 2) / (N / Math.SQRT2);
      cells.push({
        x, y,
        byte: HEX[lum >> 4] + HEX[lum & 15],
        ink: `rgb(${Math.max(r, 58)} ${Math.max(g, 56)} ${Math.max(b, 53)})`,
        tint: `rgb(${r} ${g} ${b} / .1)`,
        at: 0.12 + dist * 0.55 + Math.random() * 0.3,
      });
    }
  }

  // Pre-render the photo once at display resolution so each frame only copies cells.
  function resize() {
    const dpr = Math.min(devicePixelRatio || 1, 3);
    size = img.getBoundingClientRect().width;
    canvas.width = canvas.height = Math.round(size * dpr);
    cell = canvas.width / N;
    photo = document.createElement("canvas");
    photo.width = photo.height = canvas.width;
    photo.getContext("2d").drawImage(img, 0, 0, photo.width, photo.height);
    draw(true);
  }

  function draw(settled) {
    ctx.clearRect(0, 0, canvas.width, canvas.height);
    ctx.font = `500 ${cell * 0.5}px "IBM Plex Mono", ui-monospace, monospace`;
    ctx.textAlign = "center";
    ctx.textBaseline = "middle";
    const brand = getComputedStyle(wrap).getPropertyValue("--brand").trim();
    for (const c of cells) {
      const dx = c.x * cell, dy = c.y * cell;
      if (progress >= c.at) {
        ctx.drawImage(photo, dx, dy, cell, cell, dx, dy, cell, cell);
        continue;
      }
      ctx.fillStyle = c.tint;
      ctx.fillRect(dx, dy, cell, cell);
      // Cells about to flip scramble briefly, like a value being rewritten.
      const scrambling = !settled && progress > c.at - 0.14;
      ctx.fillStyle = scrambling ? brand : c.ink;
      const text = scrambling ? HEX[(Math.random() * 16) | 0] + HEX[(Math.random() * 16) | 0] : c.byte;
      ctx.fillText(text, dx + cell / 2, dy + cell / 2 + cell * 0.04);
    }
  }

  function tick(now) {
    const dt = Math.min(now - (last || now), 50);
    last = now;
    const step = dt / 1100;
    progress = target > progress ? Math.min(progress + step, target) : Math.max(progress - step * 1.4, target);
    const done = progress === target;
    draw(done);
    wrap.classList.toggle("is-hex", progress < 0.5);
    raf = done ? 0 : requestAnimationFrame(tick);
  }

  function go(to) {
    target = to;
    if (!raf) {
      last = 0;
      raf = requestAnimationFrame(tick);
    }
  }

  async function start() {
    await Promise.all([img.decode(), document.fonts.ready]);
    sample();
    // Decode once per visit; coming back to the homepage shows the photo straight away.
    let seen = false;
    try {
      seen = sessionStorage.getItem("avatar-decoded") === "1";
      sessionStorage.setItem("avatar-decoded", "1");
    } catch {}
    progress = seen ? 1 : 0;
    wrap.append(canvas);
    wrap.classList.add("is-live");
    resize();
    if (!seen) setTimeout(() => go(1), 350);
    wrap.addEventListener("pointerenter", (e) => e.pointerType === "mouse" && go(0));
    wrap.addEventListener("pointerleave", (e) => e.pointerType === "mouse" && go(1));
    wrap.addEventListener("click", () => go(target ? 0 : 1));
    let pending;
    addEventListener("resize", () => {
      clearTimeout(pending);
      pending = setTimeout(resize, 150);
    });
  }

  start().catch(() => {});
})();
