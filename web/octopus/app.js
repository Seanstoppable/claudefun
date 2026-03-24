// Mood Octopus — Web App
(function () {
  'use strict';

  // ── State ──
  let wasmReady = false;
  let currentEmotion = 'curiosity';

  // ── Arm path definitions per emotion ──
  const armPaths = {
    idle: {
      l1: 'M120,200 Q80,230 60,270 Q50,290 65,300',
      l2: 'M130,205 Q95,240 80,280 Q75,300 85,310',
      l3: 'M140,210 Q115,250 105,290 Q100,310 110,315',
      l4: 'M148,215 Q135,255 135,295 Q135,315 140,320',
      r1: 'M200,200 Q240,230 260,270 Q270,290 255,300',
      r2: 'M190,205 Q225,240 240,280 Q245,300 235,310',
      r3: 'M180,210 Q205,250 215,290 Q220,310 210,315',
      r4: 'M172,215 Q185,255 185,295 Q185,315 180,320',
    },
    joy: {
      l1: 'M120,200 Q60,210 40,240 Q25,255 35,265',
      l2: 'M130,205 Q75,220 55,260 Q45,280 60,285',
      l3: 'M140,210 Q100,235 85,275 Q80,295 95,298',
      l4: 'M148,215 Q125,250 120,290 Q118,310 128,315',
      r1: 'M200,200 Q260,210 280,240 Q295,255 285,265',
      r2: 'M190,205 Q245,220 265,260 Q275,280 260,285',
      r3: 'M180,210 Q220,235 235,275 Q240,295 225,298',
      r4: 'M172,215 Q195,250 200,290 Q202,310 192,315',
    },
    sadness: {
      l1: 'M120,200 Q110,240 115,280 Q118,305 120,320',
      l2: 'M130,205 Q125,245 128,285 Q130,308 132,325',
      l3: 'M140,210 Q138,250 140,290 Q141,310 142,328',
      l4: 'M148,215 Q147,255 148,295 Q149,315 150,330',
      r1: 'M200,200 Q210,240 205,280 Q202,305 200,320',
      r2: 'M190,205 Q195,245 192,285 Q190,308 188,325',
      r3: 'M180,210 Q182,250 180,290 Q179,310 178,328',
      r4: 'M172,215 Q173,255 172,295 Q171,315 170,330',
    },
    anger: {
      l1: 'M120,200 Q65,195 40,210 Q20,220 15,235',
      l2: 'M130,205 Q80,208 55,225 Q35,240 30,250',
      l3: 'M140,210 Q100,220 80,245 Q65,265 60,275',
      l4: 'M148,215 Q120,235 108,265 Q100,285 95,295',
      r1: 'M200,200 Q255,195 280,210 Q300,220 305,235',
      r2: 'M190,205 Q240,208 265,225 Q285,240 290,250',
      r3: 'M180,210 Q220,220 240,245 Q255,265 260,275',
      r4: 'M172,215 Q200,235 212,265 Q220,285 225,295',
    },
    fear: {
      l1: 'M120,200 Q115,210 125,230 Q135,245 130,255',
      l2: 'M130,205 Q125,220 132,240 Q140,255 135,268',
      l3: 'M140,210 Q138,228 142,250 Q148,268 145,280',
      l4: 'M148,215 Q147,235 150,258 Q155,275 152,290',
      r1: 'M200,200 Q205,210 195,230 Q185,245 190,255',
      r2: 'M190,205 Q195,220 188,240 Q180,255 185,268',
      r3: 'M180,210 Q182,228 178,250 Q172,268 175,280',
      r4: 'M172,215 Q173,235 170,258 Q165,275 168,290',
    },
    curiosity: {
      l1: 'M120,200 Q80,230 60,270 Q50,290 65,300',
      l2: 'M130,205 Q95,240 80,280 Q75,300 85,310',
      l3: 'M140,210 Q115,250 105,290 Q100,310 110,315',
      l4: 'M148,215 Q135,255 135,295 Q135,315 140,320',
      r1: 'M200,200 Q240,210 255,180 Q265,160 275,140',
      r2: 'M190,205 Q225,240 240,280 Q245,300 235,310',
      r3: 'M180,210 Q205,250 215,290 Q220,310 210,315',
      r4: 'M172,215 Q185,255 185,295 Q185,315 180,320',
    },
    sleepy: {
      l1: 'M120,200 Q110,240 115,285 Q118,310 120,330',
      l2: 'M130,205 Q125,245 128,290 Q130,315 132,335',
      l3: 'M140,210 Q138,252 140,295 Q141,318 142,338',
      l4: 'M148,215 Q147,258 148,300 Q149,320 150,340',
      r1: 'M200,200 Q210,240 205,285 Q202,310 200,330',
      r2: 'M190,205 Q195,245 192,290 Q190,315 188,335',
      r3: 'M180,210 Q182,252 180,295 Q179,318 178,338',
      r4: 'M172,215 Q173,258 172,300 Q171,320 170,340',
    },
    silly: {
      l1: 'M120,200 Q70,220 90,260 Q110,230 80,280',
      l2: 'M130,205 Q85,235 105,270 Q120,250 95,295',
      l3: 'M140,210 Q110,245 130,280 Q145,260 120,305',
      l4: 'M148,215 Q130,250 145,285 Q155,270 140,310',
      r1: 'M200,200 Q250,220 230,260 Q210,230 240,280',
      r2: 'M190,205 Q235,235 215,270 Q200,250 225,295',
      r3: 'M180,210 Q210,245 190,280 Q175,260 200,305',
      r4: 'M172,215 Q190,250 175,285 Q165,270 180,310',
    },
    love: {
      l1: 'M120,200 Q80,210 60,230 Q50,250 75,260 Q100,250 120,230',
      l2: 'M130,205 Q95,225 80,255 Q75,275 90,280',
      l3: 'M140,210 Q115,250 105,290 Q100,310 110,315',
      l4: 'M148,215 Q135,255 135,295 Q135,315 140,320',
      r1: 'M200,200 Q240,210 260,230 Q270,250 245,260 Q220,250 200,230',
      r2: 'M190,205 Q225,225 240,255 Q245,275 230,280',
      r3: 'M180,210 Q205,250 215,290 Q220,310 210,315',
      r4: 'M172,215 Q185,255 185,295 Q185,315 180,320',
    },
  };

  // Mouth shapes per emotion
  const mouthPaths = {
    joy: 'M142,165 Q160,182 178,165',
    sadness: 'M145,172 Q160,162 175,172',
    anger: 'M145,170 L175,170',
    fear: 'M150,168 Q160,178 170,168',
    curiosity: 'M152,168 Q160,173 168,168',
    sleepy: 'M148,168 Q160,172 172,168',
    silly: 'M142,165 Q160,185 178,165',
    love: 'M145,165 Q160,180 175,165',
  };

  // Pupil adjustments per emotion
  const pupils = {
    joy:       { lx: 137, ly: 140, rx: 187, ry: 140, lr: 8, rr: 8 },
    sadness:   { lx: 135, ly: 143, rx: 185, ry: 143, lr: 9, rr: 9 },
    anger:     { lx: 138, ly: 141, rx: 188, ry: 141, lr: 10, rr: 10 },
    fear:      { lx: 135, ly: 139, rx: 185, ry: 139, lr: 11, rr: 11 },
    curiosity: { lx: 137, ly: 140, rx: 189, ry: 138, lr: 9, rr: 8 },
    sleepy:    { lx: 136, ly: 142, rx: 186, ry: 142, lr: 7, rr: 7 },
    silly:     { lx: 133, ly: 138, rx: 189, ry: 140, lr: 8, rr: 9 },
    love:      { lx: 137, ly: 140, rx: 187, ry: 140, lr: 9, rr: 9 },
  };

  function escHtml(str) {
    if (!str) return '';
    return String(str).replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;');
  }

  // ── DOM refs ──
  const $ = (sel) => document.querySelector(sel);
  const app = $('#app');
  const input = $('#mood-input');
  const moodDisplay = $('#mood-display');
  const moodEmoji = $('#mood-emoji');
  const moodName = $('#mood-name');
  const confidenceFill = $('#confidence-fill');
  const moodBreakdown = $('#mood-breakdown');
  const moodHistory = $('#mood-history');
  const adviceBubble = $('#advice-bubble');
  const loading = $('#loading');
  const welcome = $('#welcome');
  const octopusSvg = $('#octopus-svg');

  // ── WASM init ──
  async function init() {
    try {
      await loadWasm('octopus.wasm');
      wasmReady = true;
      loading.style.display = 'none';
      app.style.display = '';
      checkWelcome();
      loadHistory();
    } catch (e) {
      loading.textContent = 'Failed to load WASM: ' + String(e.message);
    }
  }

  // ── Welcome ──
  function checkWelcome() {
    if (!localStorage.getItem('octopus-visited')) {
      welcome.style.display = '';
    }
  }

  window.dismissWelcome = function () {
    welcome.style.display = 'none';
    localStorage.setItem('octopus-visited', '1');
  };

  // ── History ──
  function getHistory() {
    try {
      return JSON.parse(localStorage.getItem('octopus-history') || '[]');
    } catch { return []; }
  }

  function addToHistory(emoji) {
    const h = getHistory();
    h.push(emoji);
    if (h.length > 30) h.splice(0, h.length - 30);
    localStorage.setItem('octopus-history', JSON.stringify(h));
    renderHistory();
  }

  function loadHistory() {
    renderHistory();
  }

  function renderHistory() {
    const h = getHistory();
    if (h.length === 0) {
      moodHistory.innerHTML = '';
      return;
    }
    moodHistory.innerHTML =
      '<span class="mood-history-label">History:</span>' +
      h.map((e) => `<span class="history-emoji">${escHtml(e)}</span>`).join('');
  }

  // ── Set emotion on octopus ──
  function setOctopusEmotion(emotion) {
    const lower = emotion.toLowerCase();
    if (lower === currentEmotion) return;

    // Remove old emotion class, add new
    octopusSvg.classList.remove(currentEmotion);
    octopusSvg.classList.add(lower);
    currentEmotion = lower;

    // Update arms
    const paths = armPaths[lower] || armPaths.idle;
    Object.entries(paths).forEach(([key, d]) => {
      const arm = $(`#arm-${key}`);
      if (arm) arm.setAttribute('d', d);
    });

    // Update mouth
    const mouth = $('#octo-mouth');
    if (mouthPaths[lower]) {
      mouth.setAttribute('d', mouthPaths[lower]);
    }

    // Update pupils
    const p = pupils[lower] || pupils.curiosity;
    const pl = $('#pupil-l');
    const pr = $('#pupil-r');
    pl.setAttribute('cx', p.lx);
    pl.setAttribute('cy', p.ly);
    pl.setAttribute('r', p.lr);
    pr.setAttribute('cx', p.rx);
    pr.setAttribute('cy', p.ry);
    pr.setAttribute('r', p.rr);
  }

  // ── Update background tint ──
  function setBackgroundTint(color) {
    document.documentElement.style.setProperty('--bg-tint', color);
    document.documentElement.style.setProperty('--current-emotion', color);
  }

  // ── Show advice ──
  function showAdvice(text) {
    if (!text) {
      adviceBubble.classList.remove('visible');
      return;
    }
    adviceBubble.textContent = `🐙 "${text}"`;
    adviceBubble.classList.add('visible');
    setTimeout(() => adviceBubble.classList.remove('visible'), 6000);
  }

  // ── Render mood breakdown ──
  function renderBreakdown(allMoods) {
    if (!allMoods || allMoods.length <= 1) {
      moodBreakdown.innerHTML = '';
      return;
    }
    moodBreakdown.innerHTML = allMoods
      .map(
        (m) => `
        <span class="mood-chip">
          ${escHtml(m.emoji)}
          <span class="chip-bar"><span class="chip-fill" style="width:${Math.round(m.confidence * 100)}%;background:${escHtml(m.color)}"></span></span>
          ${Math.round(m.confidence * 100)}%
        </span>`
      )
      .join('');
  }

  // ── Handle input ──
  let debounceTimer;
  input.addEventListener('input', () => {
    clearTimeout(debounceTimer);
    debounceTimer = setTimeout(processInput, 300);
  });

  input.addEventListener('keydown', (e) => {
    if (e.key === 'Enter') {
      clearTimeout(debounceTimer);
      processInput();
    }
  });

  function processInput() {
    if (!wasmReady) return;
    const text = input.value.trim();
    if (!text) return;

    try {
      const raw = window.analyzeMood(text);
      const result = JSON.parse(raw);
      const dom = result.dominant;

      // Update display
      moodDisplay.style.display = '';
      moodEmoji.textContent = dom.emoji;
      moodName.textContent = dom.name;
      confidenceFill.style.width = Math.round(dom.confidence * 100) + '%';

      // Update octopus
      setOctopusEmotion(dom.name);
      setBackgroundTint(dom.color);

      // Breakdown
      renderBreakdown(result.all);

      // History
      addToHistory(dom.emoji);

      // Advice
      if (result.advice) {
        showAdvice(result.advice);
      }

      // Dismiss welcome on first interaction
      if (welcome.style.display !== 'none') {
        dismissWelcome();
      }
    } catch (e) {
      console.error('analyzeMood error:', e);
    }
  }

  // ── Start ──
  init();
})();
