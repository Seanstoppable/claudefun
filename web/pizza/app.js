(function () {
  'use strict';

  const generateBtn = document.getElementById('generate-btn');
  const menuToggle = document.getElementById('menu-toggle');
  const menuControls = document.getElementById('menu-controls');
  const courseCountSelect = document.getElementById('course-count');
  const output = document.getElementById('output');
  const emptyState = document.getElementById('empty-state');
  const loading = document.getElementById('loading');

  let menuMode = false;
  let wasmReady = false;

  // Check URL for seed param
  const urlParams = new URLSearchParams(window.location.search);
  const urlSeed = urlParams.get('seed');

  // Load WASM
  loadWasm('pizza.wasm').then(() => {
    wasmReady = true;
    generateBtn.disabled = false;
    generateBtn.textContent = '🍕 Generate Pizza';

    // Auto-generate if seed in URL
    if (urlSeed) {
      generateSingle(parseInt(urlSeed, 10));
    }
  }).catch(err => {
    output.innerHTML = `<div class="card"><p style="color:var(--pizza-red)">Failed to load WASM: ${escHtml(err.message)}</p></div>`;
  });

  // Event Listeners
  generateBtn.addEventListener('click', () => {
    if (!wasmReady) return;
    spinButton();
    if (menuMode) {
      generateTastingMenu();
    } else {
      generateSingle();
    }
  });

  menuToggle.addEventListener('click', () => {
    menuMode = !menuMode;
    menuToggle.classList.toggle('active', menuMode);
    menuControls.style.display = menuMode ? 'flex' : 'none';
    generateBtn.innerHTML = menuMode
      ? '<span class="pizza-icon">📋</span> Generate Tasting Menu'
      : '<span class="pizza-icon">🍕</span> Generate Pizza';
  });

  function spinButton() {
    generateBtn.classList.add('spinning');
    setTimeout(() => generateBtn.classList.remove('spinning'), 600);
  }

  function generateSingle(seed) {
    showLoading();
    setTimeout(() => {
      try {
        const jsonStr = seed !== undefined
          ? window.generatePizza(seed)
          : window.generatePizza();
        const pizza = JSON.parse(jsonStr);
        hideLoading();
        output.innerHTML = renderRecipeCard(pizza) + renderSeedBar(pizza.seed);
        animatePretensionBars();
      } catch (e) {
        hideLoading();
        output.innerHTML = `<div class="card"><p style="color:var(--pizza-red)">Error: ${escHtml(e.message)}</p></div>`;
      }
    }, 100);
  }

  function generateTastingMenu() {
    showLoading();
    setTimeout(() => {
      try {
        const count = parseInt(courseCountSelect.value, 10);
        const jsonStr = window.generateMenu(count);
        const data = JSON.parse(jsonStr);
        hideLoading();
        let html = '';
        data.pizzas.forEach((pizza, i) => {
          html += renderRecipeCard(pizza, i + 1, data.pizzas.length);
        });
        html += renderSeedBar(data.seed);
        output.innerHTML = html;
        animatePretensionBars();
      } catch (e) {
        hideLoading();
        output.innerHTML = `<div class="card"><p style="color:var(--pizza-red)">Error: ${escHtml(e.message)}</p></div>`;
      }
    }, 100);
  }

  function showLoading() {
    emptyState.style.display = 'none';
    loading.style.display = 'block';
    output.style.display = 'none';
  }

  function hideLoading() {
    loading.style.display = 'none';
    output.style.display = 'block';
  }

  function renderRecipeCard(pizza, courseNum, totalCourses) {
    const isCourse = courseNum !== undefined;

    let html = `<div class="recipe-card" style="animation-delay: ${(courseNum || 1) * 0.1}s">`;

    if (isCourse) {
      html += `<div class="course-label">· Course ${courseNum} of ${totalCourses} ·</div>`;
    }

    html += `<div class="pizza-name">"${escHtml(pizza.name)}"</div>`;
    html += '<div class="recipe-card-body">';

    // Ingredients
    html += renderIngredient('base', pizza.base);
    html += renderIngredient('sauce', pizza.sauce);
    pizza.cheeses.forEach((c, i) => html += renderIngredient(i === 0 ? 'cheese' : '', c));
    pizza.proteins.forEach((p, i) => html += renderIngredient(i === 0 ? 'protein' : '', p));
    pizza.toppings.forEach((t, i) => html += renderIngredient(i === 0 ? 'topping' : '', t));
    html += renderIngredient('garnish', pizza.garnish);
    html += renderIngredient('drizzle', pizza.drizzle);

    // Tasting Notes
    html += '<hr class="section-divider">';
    html += '<div class="tasting-section">';
    html += '<div class="section-title">Tasting Notes</div>';
    html += `<div class="tasting-note">${escHtml(pizza.tastingNote)}</div>`;
    html += '</div>';

    // Pairing
    html += '<div class="tasting-section">';
    html += '<div class="section-title">Pairs With</div>';
    html += `<div class="pairing"><span class="pairing-icon">🍷</span>${escHtml(pizza.pairing)}</div>`;
    html += '</div>';

    // Chef Quote
    html += '<div class="tasting-section">';
    html += '<div class="section-title">Chef\'s Word</div>';
    html += `<div class="chef-quote"><span class="chef-quote-icon">👨‍🍳</span>${escHtml(pizza.chefQuote)}</div>`;
    html += '</div>';

    // Pretension Meter
    html += '<hr class="section-divider">';
    html += '<div class="pretension-section">';
    html += '<div class="section-title">Pretension Level</div>';
    const pct = Math.min((pizza.pretension / 5) * 100, 100);
    html += '<div class="pretension-bar-container">';
    html += `<div class="pretension-bar" data-width="${pct}" style="width: 0%"></div>`;
    html += '</div>';
    html += `<div class="pretension-rating-text">${escHtml(pizza.pretensionRating)}</div>`;
    html += '</div>';

    html += '</div></div>';
    return html;
  }

  function renderIngredient(category, ing) {
    const labelClass = category ? `label-${category}` : '';
    const labelText = category ? category.toUpperCase() : '';
    let html = '<div class="ingredient-row">';
    html += `<span class="ingredient-label ${labelClass}">${labelText}</span>`;
    html += '<div class="ingredient-info">';
    html += `<span class="ingredient-name">${escHtml(ing.name)}</span>`;
    html += `<span class="ingredient-flavor">${escHtml(ing.flavor)}</span>`;
    html += `<div class="ingredient-desc">${escHtml(ing.description)}</div>`;
    html += '</div></div>';
    return html;
  }

  function renderSeedBar(seed) {
    if (!seed) return '';
    return `
      <div class="seed-bar">
        <span>Seed:</span>
        <span class="seed-value">${seed}</span>
        <button class="btn-share" onclick="shareSeed(${seed})">📋 Copy Link</button>
      </div>`;
  }

  function animatePretensionBars() {
    requestAnimationFrame(() => {
      document.querySelectorAll('.pretension-bar[data-width]').forEach(bar => {
        bar.style.width = bar.dataset.width + '%';
      });
    });
  }

  // Share
  window.shareSeed = function (seed) {
    const url = new URL(window.location.href);
    url.searchParams.set('seed', seed);
    navigator.clipboard.writeText(url.toString()).then(() => {
      const btn = document.querySelector('.btn-share');
      if (btn) {
        btn.classList.add('copied');
        btn.textContent = '✓ Copied!';
        setTimeout(() => {
          btn.classList.remove('copied');
          btn.textContent = '📋 Copy Link';
        }, 2000);
      }
    }).catch(() => {
      prompt('Copy this link:', url.toString());
    });
  };

  function escHtml(str) {
    if (!str) return '';
    return str.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;');
  }
})();
