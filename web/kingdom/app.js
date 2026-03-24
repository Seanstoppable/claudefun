// ── Tiny Kingdom Simulator – app.js ──

const STAT_CONFIG = [
  { key: 'treasury',   label: '💰 Treasury',   emoji: '💰', isCurrency: true },
  { key: 'population', label: '👥 Population', emoji: '👥', isCurrency: false },
  { key: 'happiness',  label: '😊 Happiness',  emoji: '😊' },
  { key: 'military',   label: '⚔️ Military',   emoji: '⚔️' },
  { key: 'culture',    label: '🎭 Culture',    emoji: '🎭' },
  { key: 'food',       label: '🌾 Food',       emoji: '🌾' },
  { key: 'reputation', label: '🌟 Reputation', emoji: '🌟' },
];

function escHtml(str) {
  if (!str) return '';
  return String(str).replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;');
}

const FACTION_EMOJI = {
  Farmers: '🌾', Merchants: '💰', Nobles: '👑', Scholars: '📚', Jesters: '🃏'
};

let pendingCallback = null;

// ── WASM loading ──
async function loadWasm() {
  const go = new Go();
  const result = await WebAssembly.instantiateStreaming(
    fetch('kingdom.wasm'), go.importObject
  );
  go.run(result.instance);
}

function initGame() {
  const raw = window.newGame('');
  const data = JSON.parse(raw);
  showWelcome(data);
}

// ── UI Helpers ──
function statColor(value, key) {
  if (key === 'treasury') {
    if (value > 100) return 'var(--stat-green)';
    if (value > 0) return 'var(--stat-yellow)';
    return 'var(--stat-red)';
  }
  if (key === 'population') {
    if (value > 80) return 'var(--stat-green)';
    if (value > 30) return 'var(--stat-yellow)';
    return 'var(--stat-red)';
  }
  if (value >= 60) return 'var(--stat-green)';
  if (value >= 30) return 'var(--stat-yellow)';
  return 'var(--stat-red)';
}

function statPercent(value, key) {
  if (key === 'treasury') return Math.max(0, Math.min(100, (value + 500) / 10));
  if (key === 'population') return Math.min(100, value / 2);
  return value;
}

function buildStatsGrid(kingdom) {
  const grid = document.getElementById('stats-grid');
  grid.innerHTML = '';

  for (const cfg of STAT_CONFIG) {
    const val = kingdom[cfg.key];
    const pct = statPercent(val, cfg.key);
    const color = statColor(val, cfg.key);
    const displayVal = cfg.isCurrency ? `${val}g` : val;

    const bar = document.createElement('div');
    bar.className = 'stat-bar';
    bar.innerHTML = `
      <div class="stat-fill" style="width: ${pct}%; background: ${escHtml(color)}"></div>
      <div class="stat-content">
        <span class="stat-label">${escHtml(cfg.label)}</span>
        <span class="stat-value" style="color: ${escHtml(color)}">${escHtml(String(displayVal))}</span>
      </div>
    `;
    grid.appendChild(bar);
  }
}

function buildFactions(factions) {
  const panel = document.getElementById('factions-panel');
  panel.innerHTML = '';

  for (const [name, mood] of Object.entries(factions)) {
    const chip = document.createElement('div');
    chip.className = 'faction-chip' + (mood > 5 ? ' positive' : mood < -5 ? ' negative' : '');
    const sign = mood > 0 ? '+' : '';
    chip.innerHTML = `
      <span>${escHtml(FACTION_EMOJI[name] || '❓')} ${escHtml(name)}</span>
      <span class="faction-mood" style="color: ${mood > 5 ? 'var(--stat-green)' : mood < -5 ? 'var(--stat-red)' : 'var(--text-dim)'}">${escHtml(sign + String(mood))}</span>
    `;
    panel.appendChild(chip);
  }
}

function updateHeader(kingdom) {
  document.getElementById('ruler-info').textContent =
    `${kingdom.name} — Turn ${kingdom.turn} — "${kingdom.rulerTitle}"`;
}

function showPolicy(policy) {
  document.getElementById('policy-question').textContent = policy.question;
  document.getElementById('btn-a').textContent = policy.optionA;
  document.getElementById('btn-b').textContent = policy.optionB;
  document.getElementById('policy-section').style.display = '';
}

function updateStability(kingdom) {
  const warn = document.getElementById('instability-warning');
  if (!kingdom.isStable) {
    warn.classList.add('show');
  } else {
    warn.classList.remove('show');
  }
}

// ── Screens ──
function showWelcome(data) {
  document.getElementById('loading').style.display = 'none';
  document.getElementById('welcome').style.display = 'block';
  document.getElementById('welcome-text').textContent = data.narration;
  updateHeader(data.kingdom);

  document.getElementById('begin-btn').onclick = () => {
    document.getElementById('welcome').style.display = 'none';
    document.getElementById('game').style.display = 'block';
    buildStatsGrid(data.kingdom);
    buildFactions(data.kingdom.factions);
    updateStability(data.kingdom);
    showPolicy(data.policy);
    document.getElementById('turn-narration').textContent = '';
    document.getElementById('bard-narration-panel').style.display = 'none';
  };
}

function showTurnResult(data) {
  const kingdom = data.kingdom;
  updateHeader(kingdom);
  buildStatsGrid(kingdom);
  buildFactions(kingdom.factions);
  updateStability(kingdom);

  // Show bard narration of the policy result
  document.getElementById('bard-text').textContent = data.narration;
  document.getElementById('bard-narration-panel').style.display = '';

  // Turn start narration
  if (data.turnStart) {
    document.getElementById('turn-narration').textContent = data.turnStart;
  }

  // Handle event
  if (data.event) {
    showEvent(data.event, () => {
      if (data.gameOver) {
        showGameOver(data.gameOver, kingdom);
      } else {
        showPolicy(data.policy);
      }
    });
    document.getElementById('policy-section').style.display = 'none';
  } else if (data.gameOver) {
    document.getElementById('policy-section').style.display = 'none';
    showGameOver(data.gameOver, kingdom);
  } else {
    showPolicy(data.policy);
  }
}

function showEvent(eventData, callback) {
  document.getElementById('event-name').textContent = `⚡ ${eventData.name}`;
  document.getElementById('event-desc').textContent = eventData.description;
  document.getElementById('event-narration').textContent = eventData.narration;
  document.getElementById('event-overlay').classList.add('show');
  pendingCallback = callback;
}

function dismissEvent() {
  document.getElementById('event-overlay').classList.remove('show');
  if (pendingCallback) {
    const cb = pendingCallback;
    pendingCallback = null;
    cb();
  }
}

function showGameOver(goData, kingdom) {
  const card = document.getElementById('gameover-card');
  card.className = 'gameover-card ' + (goData.victory ? 'victory' : 'defeat');

  document.getElementById('gameover-title').textContent =
    goData.victory ? '🏆 Victory!' : '💀 Defeat';
  document.getElementById('gameover-message').textContent = goData.message;
  document.getElementById('gameover-narration').textContent = goData.narration;
  document.getElementById('gameover-stats').textContent =
    `Final stats — Turn ${kingdom.turn} | Treasury: ${kingdom.treasury}g | Pop: ${kingdom.population} | ` +
    `😊${kingdom.happiness} ⚔️${kingdom.military} 🎭${kingdom.culture} 🌾${kingdom.food} 🌟${kingdom.reputation}`;

  document.getElementById('gameover').classList.add('show');
}

function restartGame() {
  document.getElementById('gameover').classList.remove('show');
  document.getElementById('game').style.display = 'none';
  document.getElementById('bard-narration-panel').style.display = 'none';
  document.getElementById('turn-narration').textContent = '';
  initGame();
}

// ── Policy choice handler ──
function handleChoice(chooseA) {
  const raw = window.choosePolicy(chooseA);
  const data = JSON.parse(raw);
  showTurnResult(data);
}

// ── Event listeners ──
document.getElementById('btn-a').addEventListener('click', () => handleChoice(true));
document.getElementById('btn-b').addEventListener('click', () => handleChoice(false));
document.getElementById('event-dismiss').addEventListener('click', dismissEvent);
document.getElementById('restart-btn').addEventListener('click', restartGame);

// ── Boot ──
loadWasm().then(() => {
  initGame();
}).catch(err => {
  document.getElementById('loading').innerHTML =
    `<div style="color: var(--stat-red)">Failed to load: ${escHtml(err.message)}</div>`;
  console.error(err);
});
