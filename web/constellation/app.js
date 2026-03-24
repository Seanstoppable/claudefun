// Constellation Composer — app logic
(async function () {
    const loadingEl = document.getElementById('loading');
    const appEl = document.getElementById('app');
    const inputEl = document.getElementById('textInput');
    const resultEl = document.getElementById('result');
    const emptyStateEl = document.getElementById('emptyState');
    const svgContainer = document.getElementById('svgContainer');
    const nameEl = document.getElementById('constellationName');
    const starCountEl = document.getElementById('starCount');
    const edgeCountEl = document.getElementById('edgeCount');
    const mythCulture = document.getElementById('mythCulture');
    const mythStory = document.getElementById('mythStory');
    const mythMoral = document.getElementById('mythMoral');
    const mythViewing = document.getElementById('mythViewing');
    const downloadBtn = document.getElementById('downloadBtn');
    const shareBtn = document.getElementById('shareBtn');
    const toastEl = document.getElementById('toast');

    let currentSVG = '';
    let currentName = '';

    // Sanitize SVG: parse and strip scripts/event handlers
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

        const svgEl = doc.documentElement;
        return svgEl.outerHTML;
    }

    // Load WASM
    try {
        await loadWasm('constellation.wasm');
    } catch (err) {
        loadingEl.textContent = '❌ Failed to load: ' + err.message;
        return;
    }

    loadingEl.style.display = 'none';
    appEl.style.display = 'block';

    // Check for ?text= query param
    const params = new URLSearchParams(window.location.search);
    const initialText = params.get('text');
    if (initialText) {
        inputEl.value = initialText;
        generate(initialText);
    }

    // Submit on Enter only
    inputEl.addEventListener('keydown', function (e) {
        if (e.key === 'Enter') {
            const text = inputEl.value.trim();
            if (text) {
                generate(text);
            }
        }
    });

    inputEl.focus();

    function generate(text) {
        if (typeof generateConstellation !== 'function') {
            return;
        }

        try {
            const json = generateConstellation(text, 800, 600);
            const data = JSON.parse(json);

            currentSVG = data.svg;
            currentName = data.name;

            svgContainer.innerHTML = sanitizeSVG(data.svg);
            nameEl.textContent = data.name;
            starCountEl.textContent = data.stars;
            edgeCountEl.textContent = data.edges;

            mythCulture.textContent = data.mythology.culture;
            mythStory.textContent = data.mythology.story;
            mythMoral.textContent = data.mythology.moral;
            mythViewing.textContent = data.mythology.bestViewing;

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
        a.download = (currentName || 'constellation') + '.svg';
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
        URL.revokeObjectURL(url);
        showToast('SVG downloaded!');
    });

    // Share URL
    shareBtn.addEventListener('click', function () {
        const text = inputEl.value.trim();
        if (!text) return;
        const url = window.location.origin + window.location.pathname + '?text=' + encodeURIComponent(text);
        navigator.clipboard.writeText(url).then(function () {
            showToast('Link copied to clipboard!');
        }).catch(function () {
            // Fallback
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
