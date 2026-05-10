import { useState, useEffect, useCallback, useMemo, createElement } from 'https://esm.sh/react@18';
import { createRoot } from 'https://esm.sh/react-dom@18/client';
import htm from 'https://esm.sh/htm@3';
import { GAMES, normalizeGame, PLATFORM_OPTIONS } from './data.js';
import { api } from './api.js';
import { StatCard, FilterGroup, Stars, GameRow, YearSection, AddGameForm } from './components.js';

const html = htm.bind(createElement);

function App() {
  const [gameData, setGameData] = useState({});
  const [customGames, setCustomGames] = useState([]);
  const [loading, setLoading] = useState(true);
  const [initError, setInitError] = useState(null);
  const [saveErr, setSaveErr] = useState(false);
  const [search, setSearch] = useState('');
  const [platform, setPlatform] = useState('all');
  const [statusFilter, setStatusFilter] = useState('all');
  const [collapsed, setCollapsed] = useState({});

  useEffect(() => {
    api.state()
      .then((data) => {
        setGameData(data.state || {});
        setCustomGames((data.games || []).map(normalizeGame));
        setLoading(false);
      })
      .catch((err) => {
        setInitError(err.message);
        setLoading(false);
      });
  }, []);

  const mergedGames = useMemo(
    () => [...GAMES, ...customGames.map(normalizeGame)],
    [customGames],
  );

  const filtered = useMemo(
    () =>
      mergedGames.filter((g) => {
        if (search && !g.t.toLowerCase().includes(search.toLowerCase())) return false;
        if (platform !== 'all' && g.p !== platform) return false;
        if (statusFilter === 'released' && !g.released) return false;
        if (statusFilter === 'upcoming' && g.released) return false;
        return true;
      }),
    [search, platform, statusFilter, mergedGames],
  );

  const years = useMemo(
    () => [...new Set(filtered.map((g) => g.yr))].sort((a, b) => a - b),
    [filtered],
  );

  const totalPlayed = useMemo(
    () => Object.values(gameData).filter((d) => d.played).length,
    [gameData],
  );

  const totalRel = useMemo(
    () => mergedGames.filter((g) => g.released).length,
    [mergedGames],
  );

  const totalUp = useMemo(
    () => mergedGames.filter((g) => !g.released).length,
    [mergedGames],
  );

  const flashSaveErr = useCallback(() => {
    setSaveErr(true);
    setTimeout(() => setSaveErr(false), 2500);
  }, []);

  const togglePlayed = useCallback(
    async (id) => {
      const key = String(id);
      const current = gameData[key] || {};
      const next = { ...current, played: !current.played };
      setGameData((prev) => ({ ...prev, [key]: next }));
      try {
        await api.update(id, next);
      } catch {
        setGameData((prev) => ({ ...prev, [key]: current }));
        flashSaveErr();
      }
    },
    [gameData, flashSaveErr],
  );

  const setRating = useCallback(
    async (id, rating) => {
      const key = String(id);
      const current = gameData[key] || {};
      const next = { ...current, rating };
      setGameData((prev) => ({ ...prev, [key]: next }));
      try {
        await api.update(id, next);
      } catch {
        setGameData((prev) => ({ ...prev, [key]: current }));
        flashSaveErr();
      }
    },
    [gameData, flashSaveErr],
  );

  const addGame = useCallback(async (game) => {
    const created = await api.createGame(game);
    setCustomGames((prev) => [...prev, normalizeGame(created)]);
    return created;
  }, []);

  const toggleCollapse = useCallback((year) => {
    setCollapsed((prev) => ({ ...prev, [year]: !prev[year] }));
  }, []);

  if (loading) {
    return html`
      <div class="loading">
        <span class="loading-text mono">Loading…</span>
      </div>
    `;
  }

  if (initError) {
    return html`
      <div class="error-state">
        <h2>Cannot connect to API</h2>
        <p>Make sure the Go server is running:</p>
        <p><code>go run .</code></p>
        <p class="error-detail">${initError}</p>
      </div>
    `;
  }

  return html`
    <div class="app">
      <header>
        <h1>Game log</h1>
        <p class="header-sub mono">2020 – 2026 · PlayStation · Nintendo · Multi-platform</p>
      </header>

      <${AddGameForm} platforms=${PLATFORM_OPTIONS} onAdd=${addGame} />

      <div class="stats">
        <${StatCard} label="Total" value=${mergedGames.length} />
        <${StatCard} label="Released" value=${totalRel} />
        <${StatCard} label="Upcoming" value=${totalUp} />
        <${StatCard} label="Played" value=${totalPlayed} accent />
      </div>

      <div class="controls">
        <div class="search-wrap">
          <i class="ti ti-search search-icon" aria-hidden="true" />
          <input
            type="text"
            placeholder="Search titles…"
            value=${search}
            onInput=${(e) => setSearch(e.target.value)}
          />
        </div>
        <${FilterGroup}
          options=${[
            { value: 'all', label: 'All' },
            { value: 'ps', label: 'PlayStation' },
            { value: 'nintendo', label: 'Nintendo' },
            { value: 'multi', label: 'Multi-platform' },
          ]}
          value=${platform}
          accentVar="--acc"
          onChange=${setPlatform}
        />
        <${FilterGroup}
          options=${[
            { value: 'all', label: 'All' },
            { value: 'released', label: 'Released' },
            { value: 'upcoming', label: 'Upcoming' },
          ]}
          value=${statusFilter}
          accentVar="--soon-fg"
          onChange=${setStatusFilter}
        />
      </div>

      <main>
        ${filtered.length === 0
          ? html`<div class="no-results">No results</div>`
          : years.map((yr) => html`
              <${YearSection}
                key=${yr}
                year=${yr}
                games=${filtered.filter((g) => g.yr === yr)}
                gameData=${gameData}
                onToggle=${togglePlayed}
                onRate=${setRating}
                isCollapsed=${!!collapsed[yr]}
                onCollapse=${toggleCollapse}
              />
            `)}
      </main>

      <footer class="app-footer mono">
        <span>Click a title to mark played · stars to rate · year to collapse</span>
        <span>${totalPlayed}/${mergedGames.length} played</span>
      </footer>

      <div class=${'save-err' + (saveErr ? ' visible' : '')}>
        Save failed — check server
      </div>
    </div>
  `;
}

createRoot(document.getElementById('root')).render(html`<${App} />`);
