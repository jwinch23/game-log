import { createElement, useState } from 'https://esm.sh/react@18';
import htm from 'https://esm.sh/htm@3';
import { PL_LABEL } from './data.js';
const html = htm.bind(createElement);

export function StatCard({ label, value, accent }) {
  return html`
    <div class="stat">
      <div class="stat-label">${label}</div>
      <div class="stat-value mono ${accent ? 'accent' : ''}">${value ?? '—'}</div>
    </div>
  `;
}

export function FilterGroup({ options, value, accentVar, onChange }) {
  return html`
    <div class="filter-group">
      ${options.map((opt) => html`
        <button
          key=${opt.value}
          type="button"
          onClick=${() => onChange(opt.value)}
          style=${
            value === opt.value
              ? { borderColor: `var(${accentVar})`, color: `var(${accentVar})` }
              : {}
          }
        >${opt.label}</button>
      `)}
    </div>
  `;
}

export function Stars({ id, rating, onRate }) {
  return html`
    <span class="stars" onClick=${(e) => e.stopPropagation()}>
      ${[1, 2, 3, 4, 5].map((s) => html`
        <span
          key=${s}
          class=${'star' + (s <= rating ? ' lit' : '')}
          onClick=${() => onRate(id, s === rating ? 0 : s)}
        >★</span>
      `)}
    </span>
  `;
}

export function GameRow({ game, data = {}, onToggle, onRate }) {
  const cls = [
    'game-row',
    data.played ? 'played' : '',
    game.status === 'soon' ? 'soon' : !game.released ? 'dim' : '',
  ]
    .filter(Boolean)
    .join(' ');

  return html`
    <div class=${cls} onClick=${() => onToggle(game.id)}>
      <div class="dot" />
      <span class=${'badge badge-' + game.p}>${PL_LABEL[game.p]}</span>
      <span class="game-title">${game.t}</span>
      ${game.status === 'new' && html`<span class="pill pill-new">new</span>`}
      ${game.status === 'soon' && html`<span class="pill pill-soon">soon</span>`}
      ${data.played && html`<${Stars} id=${game.id} rating=${data.rating || 0} onRate=${onRate} />`}
      <span class=${'game-date mono' + (game.status === 'soon' ? ' soon-date' : '')}>${game.d}</span>
    </div>
  `;
}

export function YearSection({ year, games, gameData, onToggle, onRate, isCollapsed, onCollapse }) {
  const played = games.filter((g) => gameData[g.id]?.played).length;
  const pct = games.length > 0 ? Math.round((played / games.length) * 100) : 0;
  const meta = `${games.length} titles${played > 0 ? ` · ${played}/${games.length} played` : ''}`;

  return html`
    <div>
      <div class="year-header" onClick=${() => onCollapse(year)}>
        <span class="year-num mono">${year}</span>
        <span class="year-meta mono">${meta}</span>
        <div class="year-line" />
        <span class="year-toggle">${isCollapsed ? '▸' : '▾'}</span>
      </div>

      ${played > 0 && !isCollapsed && html`
        <div class="year-bar">
          <div class="year-bar-fill" style=${{ width: pct + '%' }} />
        </div>
      `}

      ${!isCollapsed && html`
        <div class="game-grid">
          ${games.map((g) => html`
            <${GameRow}
              key=${g.id}
              game=${g}
              data=${gameData[String(g.id)]}
              onToggle=${onToggle}
              onRate=${onRate}
            />
          `)}
        </div>
      `}
    </div>
  `;
}

export function AddGameForm({ platforms, onAdd }) {
  const [title, setTitle] = useState('');
  const [year, setYear] = useState('');
  const [platform, setPlatform] = useState(platforms[0]?.value || 'ps');
  const [date, setDate] = useState('');
  const [rel, setRel] = useState('');
  const [error, setError] = useState('');
  const [saving, setSaving] = useState(false);

  const handleSubmit = async (event) => {
    event.preventDefault();
    setError('');

    const yearNumber = Number(year);
    if (!title.trim()) {
      setError('Title is required');
      return;
    }
    if (!Number.isInteger(yearNumber) || yearNumber < 1970 || yearNumber > 2100) {
      setError('Year must be valid');
      return;
    }
    if (!date.trim()) {
      setError('Display date is required');
      return;
    }

    setSaving(true);
    try {
      await onAdd({
        title: title.trim(),
        year: yearNumber,
        platform,
        date: date.trim(),
        rel,
      });
      setTitle('');
      setYear('');
      setPlatform(platforms[0]?.value || 'ps');
      setDate('');
      setRel('');
    } catch (err) {
      setError(err?.message || 'Failed to save game');
    } finally {
      setSaving(false);
    }
  };

  return html`
    <section class="add-game">
      <div class="add-game-head">
        <div>
          <strong>Add a custom game</strong>
          <div class="add-game-note">Saved to disk and available after refresh.</div>
        </div>
      </div>

      <form class="add-game-form" onSubmit=${handleSubmit}>
        <div class="field">
          <label>
            Title
            <input
              type="text"
              value=${title}
              onInput=${(e) => setTitle(e.target.value)}
              placeholder="Game title"
            />
          </label>
        </div>

        <div class="add-game-row">
          <div class="field">
            <label>
              Year
              <input
                type="number"
                value=${year}
                onInput=${(e) => setYear(e.target.value)}
                placeholder="2026"
              />
            </label>
          </div>
          <div class="field">
            <label>
              Platform
              <select value=${platform} onChange=${(e) => setPlatform(e.target.value)}>
                ${platforms.map((option) => html`
                  <option key=${option.value} value=${option.value}>${option.label}</option>
                `)}
              </select>
            </label>
          </div>
        </div>

        <div class="field">
          <label>
            Display date
            <input
              type="text"
              value=${date}
              onInput=${(e) => setDate(e.target.value)}
              placeholder="May 12"
            />
          </label>
        </div>

        <div class="field">
          <label>
            Release date
            <input
              type="date"
              value=${rel}
              onInput=${(e) => setRel(e.target.value)}
            />
          </label>
        </div>

        <div class="form-actions">
          <button class="form-button" type="submit" disabled=${saving}>
            ${saving ? 'Saving…' : 'Add game'}
          </button>
          ${error && html`<span class="form-error">${error}</span>`}
        </div>
      </form>
    </section>
  `;
}
