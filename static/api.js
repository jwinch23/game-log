function ensureOk(response) {
  if (!response.ok) {
    throw new Error(`${response.status} ${response.statusText}`);
  }
  return response;
}

export const api = {
  state: () =>
    fetch('/api/state')
      .then(ensureOk)
      .then((response) => response.json()),

  update: (id, gs) =>
    fetch(`/api/games/${id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(gs),
    })
      .then(ensureOk)
      .then((response) => response.json()),

  createGame: (game) =>
    fetch('/api/games', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(game),
    })
      .then(ensureOk)
      .then((response) => response.json()),
};
