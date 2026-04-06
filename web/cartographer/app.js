// Imaginary Cartographer — app logic
(async function () {
    const loadingEl = document.getElementById('loading');
    const appEl = document.getElementById('app');
    const inputEl = document.getElementById('textInput');
    const resultEl = document.getElementById('result');
    const emptyStateEl = document.getElementById('emptyState');
    const svgContainer = document.getElementById('svgContainer');
    const placeNameEl = document.getElementById('placeName');
    const landmarkListEl = document.getElementById('landmarkList');
    const loreMotto = document.getElementById('loreMotto');
    const loreMyth = document.getElementById('loreMyth');
    const loreLegends = document.getElementById('loreLegends');
    const loreWarning = document.getElementById('loreWarning');
    const downloadBtn = document.getElementById('downloadBtn');
    const shareBtn = document.getElementById('shareBtn');
    const toastEl = document.getElementById('toast');
    const lastMsgEl = document.getElementById('lastMessage');

    let currentSVG = '';
    let currentPlace = '';

    function sanitizeSVG(svgString) {
        const parser = new DOMParser();
        const doc = parser.parseFromString(svgString, 'image/svg+xml');
        const errorNode = doc.querySelector('parsererror');
        if (errorNode) return '';

        doc.querySelectorAll('script').forEach(el => el.remove());

        const allEls = doc.querySelectorAll('*');
        for (const el of allEls) {
            for (const attr of Array.from(el.attributes)) {
                if (attr.name.startsWith('on')) {
                    el.removeAttribute(attr.name);
                }
            }
            if (el.hasAttributeNS('http://www.w3.org/1999/xlink', 'href')) {
                const href = el.getAttributeNS('http://www.w3.org/1999/xlink', 'href');
                if (href && href.trim().toLowerCase().startsWith('javascript:')) {
                    el.removeAttributeNS('http://www.w3.org/1999/xlink', 'href');
                }
            }
            if (el.hasAttribute('href')) {
                const href = el.getAttribute('href');
                if (href && href.trim().toLowerCase().startsWith('javascript:')) {
                    el.removeAttribute('href');
                }
            }
        }

        return doc.documentElement.outerHTML;
    }

    function escapeHTML(str) {
        const div = document.createElement('div');
        div.textContent = str;
        return div.innerHTML;
    }

    // Load WASM
    try {
        await loadWasm('cartographer.wasm');
    } catch (err) {
        loadingEl.textContent = '❌ Failed to load: ' + err.message;
        return;
    }

    loadingEl.style.display = 'none';
    appEl.style.display = 'block';

    // Check for ?place= query param
    const params = new URLSearchParams(window.location.search);
    const initialPlace = params.get('place');
    if (initialPlace) {
        inputEl.value = initialPlace;
        generate(initialPlace);
    }

    // Submit on Enter only
    inputEl.addEventListener('keydown', function (e) {
        if (e.key === 'Enter') {
            const text = inputEl.value.trim();
            if (text) {
                generate(text);
                lastMsgEl.textContent = '🗺️ "' + text + '"';
                lastMsgEl.style.display = 'block';
                inputEl.value = '';
            }
        }
    });

    inputEl.focus();

    function generate(placeName) {
        if (typeof generateMap !== 'function') {
            return;
        }

        try {
            const json = generateMap(placeName);
            const data = JSON.parse(json);

            currentSVG = data.svg;
            currentPlace = data.placeName;

            svgContainer.innerHTML = sanitizeSVG(data.svg);
            placeNameEl.textContent = data.placeName;

            // Landmarks
            landmarkListEl.innerHTML = '';
            if (data.landmarks) {
                data.landmarks.forEach(function (lm) {
                    const item = document.createElement('div');
                    item.className = 'landmark-item';
                    item.innerHTML =
                        '<span class="landmark-symbol">' + escapeHTML(lm.symbol) + '</span>' +
                        '<span class="landmark-name">' + escapeHTML(lm.name) + '</span>' +
                        '<span class="landmark-type">' + escapeHTML(lm.type) + '</span>';
                    landmarkListEl.appendChild(item);
                });
            }

            // Lore
            if (data.lore) {
                loreMotto.textContent = '— ' + data.lore.motto + ' —';
                loreMyth.textContent = data.lore.creationMyth;
                loreWarning.textContent = data.lore.warning;

                loreLegends.innerHTML = '';
                if (data.lore.legends) {
                    data.lore.legends.forEach(function (leg) {
                        const entry = document.createElement('div');
                        entry.className = 'legend-entry';
                        entry.innerHTML =
                            '<div class="legend-title">' + escapeHTML(leg.title) + '</div>' +
                            '<span class="legend-category">' + escapeHTML(leg.category) + '</span>' +
                            '<div class="legend-story">' + escapeHTML(leg.story) + '</div>';
                        loreLegends.appendChild(entry);
                    });
                }
            }

            emptyStateEl.style.display = 'none';
            resultEl.style.display = 'block';
        } catch (err) {
            console.error('Generation error:', err);
        }
    }

    // Download SVG
    downloadBtn.addEventListener('click', function () {
        if (!currentSVG) return;
        const blob = new Blob([currentSVG], { type: 'image/svg+xml' });
        const url = URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = (currentPlace || 'map') + '.svg';
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
        URL.revokeObjectURL(url);
        showToast('SVG downloaded!');
    });

    // Share URL
    shareBtn.addEventListener('click', function () {
        const place = currentPlace;
        if (!place) return;
        const url = window.location.origin + window.location.pathname + '?place=' + encodeURIComponent(place);
        navigator.clipboard.writeText(url).then(function () {
            showToast('Link copied to clipboard!');
        }).catch(function () {
            prompt('Copy this link:', url);
        });
    });

    function showToast(message) {
        toastEl.textContent = message;
        toastEl.classList.add('show');
        setTimeout(function () {
            toastEl.classList.remove('show');
        }, 2500);
    }
})();
