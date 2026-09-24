// Byte field: the card's own text sits behind it as bytes, and the pointer
// uncovers a soft patch of it. Both layers are drawn once; CSS masks do the reveal.
(() => {
  const card = document.querySelector(".home-shell");
  if (!card || !matchMedia("(hover: hover) and (pointer: fine)").matches) return;

  const HEX = "0123456789ABCDEF";
  const CELL_W = 22, CELL_H = 18;
  const bytes = new TextEncoder().encode(card.innerText.replace(/\s+/g, " ").trim());
  const [field, hot] = ["byte-field", "byte-field byte-field-hot"].map((className) => {
    const canvas = document.createElement("canvas");
    canvas.className = className;
    canvas.setAttribute("aria-hidden", "true");
    card.prepend(canvas);
    return canvas;
  });

  function paint(canvas, colour) {
    const dpr = Math.min(devicePixelRatio || 1, 2);
    const w = card.clientWidth, h = card.clientHeight;
    const cols = Math.ceil(w / CELL_W), rows = Math.ceil(h / CELL_H);
    canvas.width = Math.round(w * dpr);
    canvas.height = Math.round(h * dpr);
    const ctx = canvas.getContext("2d");
    ctx.scale(dpr, dpr);
    ctx.font = '10px "IBM Plex Mono", ui-monospace, monospace';
    ctx.textAlign = "center";
    ctx.textBaseline = "middle";
    ctx.fillStyle = colour;
    for (let y = 0; y < rows; y++) {
      for (let x = 0; x < cols; x++) {
        const b = bytes[(y * cols + x) % bytes.length];
        ctx.fillText(HEX[b >> 4] + HEX[b & 15], x * CELL_W + CELL_W / 2, y * CELL_H + CELL_H / 2);
      }
    }
  }

  function draw() {
    const style = getComputedStyle(card);
    paint(field, style.getPropertyValue("--faint").trim());
    paint(hot, style.getPropertyValue("--brand").trim());
  }

  let frame = 0, x = 0, y = 0;
  card.addEventListener("pointermove", (e) => {
    const rect = card.getBoundingClientRect();
    x = e.clientX - rect.left;
    y = e.clientY - rect.top;
    frame ||= requestAnimationFrame(() => {
      frame = 0;
      card.style.setProperty("--x", x + "px");
      card.style.setProperty("--y", y + "px");
    });
  });
  card.addEventListener("pointerenter", () => card.classList.add("is-probing"));
  card.addEventListener("pointerleave", () => card.classList.remove("is-probing"));

  document.fonts.ready.then(() => {
    draw();
    new ResizeObserver(draw).observe(card);
    // Repaint when the accent changes (god mode).
    new MutationObserver(draw).observe(document.documentElement, { attributes: true, attributeFilter: ["class"] });
  });
})();
