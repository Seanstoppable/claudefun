(async function () {
    const loadingEl = document.getElementById('loading');
    const appEl = document.getElementById('app');
    const codeInput = document.getElementById('code-input');
    const languageSelect = document.getElementById('language');
    const growBtn = document.getElementById('grow-btn');
    const gardenDisplay = document.getElementById('garden-display');
    const emptyState = document.getElementById('empty-state');
    const gardenContent = document.getElementById('garden-content');
    const statsPanel = document.getElementById('stats-panel');
    const healthScore = document.getElementById('health-score');
    const healthFill = document.getElementById('health-fill');
    const elementCounts = document.getElementById('element-counts');
    const codeStats = document.getElementById('code-stats');
    const reportText = document.getElementById('report-text');

    // Load WASM
    try {
        await loadWasm('garden.wasm');
        loadingEl.style.display = 'none';
        appEl.style.display = 'block';
    } catch (err) {
        loadingEl.textContent = '❌ Failed to sprout: ' + err.message;
        return;
    }

    function growGarden() {
        const code = codeInput.value.trim();
        if (!code) return;

        if (typeof analyzeCode !== 'function') {
            console.error('analyzeCode not available');
            return;
        }

        try {
            const lang = languageSelect.value;
            const json = analyzeCode(code, lang);
            const data = JSON.parse(json);

            if (data.error) {
                console.error('Analysis error:', data.error);
                return;
            }

            renderGarden(data);
            renderStats(data);
        } catch (err) {
            console.error('Garden error:', err);
        }
    }

    function renderGarden(data) {
        let html = '';

        html += '<div class="garden-title">🌱 ' + escapeHtml(data.title) + ' 🌱</div>';
        html += '<div class="garden-meta">' + escapeHtml(data.season) + ' · ' + escapeHtml(data.weather) + '</div>';
        html += '<div class="garden-grid">';

        for (let y = 0; y < data.grid.length; y++) {
            html += '<div class="garden-row">';
            const row = data.grid[y];
            for (let x = 0; x < row.length; x++) {
                const cell = row[x];
                const title = cell.label ? cell.label : '';
                html += '<span class="garden-cell" style="color:' + cell.color + '"'
                    + (title ? ' title="' + escapeHtml(title) + '"' : '')
                    + '>' + escapeHtml(cell.symbol) + '</span>';
            }
            html += '</div>';
        }

        html += '</div>';

        emptyState.style.display = 'none';
        gardenContent.innerHTML = html;
        gardenContent.style.display = 'block';
    }

    function renderStats(data) {
        // Health bar
        const score = Math.round(data.healthScore);
        healthScore.textContent = score + '/100';

        let barColor = '#4CAF50';
        if (score < 40) barColor = '#F44336';
        else if (score < 70) barColor = '#FF9800';

        healthFill.style.width = score + '%';
        healthFill.style.background = barColor;

        // Element counts
        const stats = data.stats;
        const elements = [
            { icon: '🌳', label: 'Trees', count: stats.trees, color: '#006400' },
            { icon: '✿', label: 'Flowers', count: stats.flowers, color: '#FF69B4' },
            { icon: '┃', label: 'Fences', count: stats.fences, color: '#8B4513' },
            { icon: '🐛', label: 'Bugs', count: stats.bugs, color: '#FF0000' },
            { icon: '🦋', label: 'Butterflies', count: stats.butterflies, color: '#87CEEB' },
            { icon: '●', label: 'Rocks', count: stats.rocks, color: '#808080' },
            { icon: '⌇', label: 'Weeds', count: stats.weeds, color: '#DAA520' },
        ];

        let countsHtml = '';
        for (const el of elements) {
            countsHtml += '<div class="element-count">'
                + '<span style="color:' + el.color + '">' + el.icon + '</span> '
                + '<span class="num">' + el.count + '</span> ' + el.label
                + '</div>';
        }
        elementCounts.innerHTML = countsHtml;

        // Code stats
        const code = data.code;
        codeStats.textContent = '📊 ' + code.totalLines + ' lines · '
            + code.totalFuncs + ' functions · '
            + code.totalTypes + ' types · '
            + code.totalTests + ' tests · '
            + code.totalComments + ' comments · '
            + code.totalTODOs + ' TODOs';

        // Report
        reportText.textContent = '"' + data.report + '"';

        statsPanel.style.display = 'block';
    }

    function escapeHtml(str) {
        const div = document.createElement('div');
        div.appendChild(document.createTextNode(str));
        return div.innerHTML;
    }

    // Event listeners
    growBtn.addEventListener('click', growGarden);

    codeInput.addEventListener('keydown', function (e) {
        if (e.ctrlKey && e.key === 'Enter') {
            e.preventDefault();
            growGarden();
        }
        // Allow Tab to insert spaces
        if (e.key === 'Tab') {
            e.preventDefault();
            const start = this.selectionStart;
            const end = this.selectionEnd;
            this.value = this.value.substring(0, start) + '    ' + this.value.substring(end);
            this.selectionStart = this.selectionEnd = start + 4;
        }
    });
})();
