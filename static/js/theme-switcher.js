class ThemeSwitcher {
  constructor() {
    this.theme = this.getStoredTheme() || this.getSystemTheme();
    this.init();
  }

  getSystemTheme() {
    return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
  }

  getStoredTheme() {
    return localStorage.getItem('theme');
  }

  setStoredTheme(theme) {
    localStorage.setItem('theme', theme);
  }

  init() {
    this.applyTheme(this.theme);
    this.createSwitcher();
    this.bindEvents();
    this.updateIcon();
  }

  applyTheme(theme) {
    document.documentElement.setAttribute('data-theme', theme);
    this.setStoredTheme(theme);
  }

  createSwitcher() {
    const header = document.querySelector('header');
    if (!header) return;

    const switcher = document.createElement('div');
    switcher.className = 'theme-switcher';
    switcher.innerHTML = `
      <button class="theme-toggle" aria-label="Toggle theme">
        <span class="theme-icon">🌞</span>
      </button>
    `;

    const nav = header.querySelector('nav');
    nav.parentNode.insertBefore(switcher, nav.nextSibling);
  }

  bindEvents() {
    const toggle = document.querySelector('.theme-toggle');
    if (!toggle) return;

    toggle.addEventListener('click', () => {
      this.toggleTheme();
    });
  }

  toggleTheme() {
    this.theme = this.theme === 'light' ? 'dark' : 'light';
    this.applyTheme(this.theme);
    this.updateIcon();
  }

  updateIcon() {
    const icon = document.querySelector('.theme-icon');
    if (!icon) return;

    icon.textContent = this.theme === 'light' ? '🌝' : '🌞';
  }
}

document.addEventListener('DOMContentLoaded', () => new ThemeSwitcher());
