// Elevator Music Composer — app logic
(async function () {
    const loadingEl = document.getElementById('loading');
    const appEl = document.getElementById('app');
    const inputEl = document.getElementById('moodInput');
    const resultEl = document.getElementById('result');
    const emptyStateEl = document.getElementById('emptyState');
    const lastMsgEl = document.getElementById('lastMessage');
    const coverArtEl = document.getElementById('coverArt');
    const albumTitleEl = document.getElementById('albumTitle');
    const albumArtistEl = document.getElementById('albumArtist');
    const albumGenreEl = document.getElementById('albumGenre');
    const albumYearEl = document.getElementById('albumYear');
    const albumLabelEl = document.getElementById('albumLabel');
    const albumBioEl = document.getElementById('albumBio');
    const avgRatingEl = document.getElementById('avgRating');
    const totalPlaysEl = document.getElementById('totalPlays');
    const monthlyListenersEl = document.getElementById('monthlyListeners');
    const tracklistBodyEl = document.getElementById('tracklistBody');
    const reviewsContainerEl = document.getElementById('reviewsContainer');
    const similarArtistsEl = document.getElementById('similarArtists');
    const topPlaylistEl = document.getElementById('topPlaylist');
    const anotherBtn = document.getElementById('anotherBtn');
    const shareBtn = document.getElementById('shareBtn');
    const toastEl = document.getElementById('toast');

    const randomMoods = [
        'melancholy', 'corporate despair', 'anxious', 'smooth jazz',
        'bored in a meeting', 'cheerful', 'existential dread', 'mellow',
        'peaceful', 'blue', 'nervous', 'happy', 'calm', 'sad',
        'waiting for the bus', 'fluorescent lighting', 'monday morning',
        'lukewarm coffee', 'gentle rain', 'suburban ennui'
    ];

    // Load WASM
    try {
        await loadWasm('elevator.wasm');
    } catch (err) {
        loadingEl.textContent = '❌ Failed to load: ' + err.message;
        return;
    }

    loadingEl.style.display = 'none';
    appEl.style.display = 'block';

    // Check for ?mood= query param
    var params = new URLSearchParams(window.location.search);
    var initialMood = params.get('mood');
    if (initialMood) {
        inputEl.value = initialMood;
        generate(initialMood);
    }

    // Submit on Enter
    inputEl.addEventListener('keydown', function (e) {
        if (e.key === 'Enter') {
            var mood = inputEl.value.trim();
            if (mood) {
                generate(mood);
                lastMsgEl.textContent = '🎵 Mood: "' + mood + '"';
                lastMsgEl.style.display = 'block';
                inputEl.value = '';
            }
        }
    });

    // Mood suggestion clicks
    document.querySelectorAll('.mood-suggestions span[data-mood]').forEach(function (el) {
        el.addEventListener('click', function () {
            var mood = el.getAttribute('data-mood');
            generate(mood);
            lastMsgEl.textContent = '🎵 Mood: "' + mood + '"';
            lastMsgEl.style.display = 'block';
            inputEl.value = '';
        });
    });

    // Generate Another button
    anotherBtn.addEventListener('click', function () {
        var mood = randomMoods[Math.floor(Math.random() * randomMoods.length)];
        generate(mood);
        lastMsgEl.textContent = '🎵 Mood: "' + mood + '"';
        lastMsgEl.style.display = 'block';
        inputEl.value = '';
    });

    // Share URL
    shareBtn.addEventListener('click', function () {
        if (!currentMood) return;
        var url = window.location.origin + window.location.pathname + '?mood=' + encodeURIComponent(currentMood);
        navigator.clipboard.writeText(url).then(function () {
            showToast('Link copied to clipboard!');
        }).catch(function () {
            prompt('Copy this link:', url);
        });
    });

    inputEl.focus();

    var currentMood = '';

    function generate(mood) {
        if (typeof generateAlbum !== 'function') return;

        currentMood = mood;

        try {
            var json = generateAlbum(mood, 8);
            var data = JSON.parse(json);

            // Album header
            coverArtEl.textContent = data.coverArt;
            albumTitleEl.textContent = data.title;
            albumArtistEl.textContent = data.artist;
            albumGenreEl.textContent = data.genre;
            albumYearEl.textContent = data.year;
            albumLabelEl.textContent = data.label;
            albumBioEl.textContent = '"' + data.bio + '"';

            // Stats
            avgRatingEl.textContent = data.avgRatingStars;
            totalPlaysEl.textContent = data.totalPlays;
            monthlyListenersEl.textContent = data.monthlyListeners;

            // Track listing
            tracklistBodyEl.innerHTML = '';
            data.tracks.forEach(function (t) {
                var row = document.createElement('tr');
                var featHtml = t.featured ? '<span class="track-featured">' + escapeHtml(t.featured) + '</span>' : '';
                row.innerHTML =
                    '<td>' + t.number + '</td>' +
                    '<td class="track-title">' + escapeHtml(t.title) + featHtml + '</td>' +
                    '<td>' + escapeHtml(t.duration) + '</td>' +
                    '<td>' + t.bpm + '</td>' +
                    '<td class="track-key">' + escapeHtml(t.key) + '</td>';
                tracklistBodyEl.appendChild(row);
            });

            // Reviews
            reviewsContainerEl.innerHTML = '';
            data.reviews.forEach(function (r) {
                var div = document.createElement('div');
                div.className = 'review';
                div.innerHTML =
                    '<div class="review-header">' +
                        '<span class="reviewer-name">@' + escapeHtml(r.reviewer) + '</span>' +
                        '<span class="review-stars">' + escapeHtml(r.stars) + '</span>' +
                    '</div>' +
                    '<p class="review-text">"' + escapeHtml(r.text) + '"</p>' +
                    '<div class="review-footer">' +
                        '<span>' + escapeHtml(r.playCount) + ' plays</span>' +
                        '<span>👍 ' + r.helpful + ' found this helpful</span>' +
                    '</div>';
                reviewsContainerEl.appendChild(div);
            });

            // Similar artists
            similarArtistsEl.innerHTML = '';
            data.similarTo.forEach(function (name) {
                var li = document.createElement('li');
                li.textContent = name;
                similarArtistsEl.appendChild(li);
            });

            // Playlist
            topPlaylistEl.textContent = data.topPlaylist;

            emptyStateEl.style.display = 'none';
            resultEl.style.display = 'block';
        } catch (err) {
            console.error('Generation error:', err);
        }
    }

    function escapeHtml(str) {
        var div = document.createElement('div');
        div.appendChild(document.createTextNode(str));
        return div.innerHTML;
    }

    function showToast(message) {
        toastEl.textContent = message;
        toastEl.classList.add('show');
        setTimeout(function () {
            toastEl.classList.remove('show');
        }, 2500);
    }
})();
