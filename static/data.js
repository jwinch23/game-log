export const TODAY = new Date('2026-05-10');

export function computeStatus(rel) {
  if (!rel) return 'released';
  const d = new Date(rel);
  if (Number.isNaN(d.valueOf())) return 'released';
  const days = (d - TODAY) / 86400000;
  if (days > 30) return 'upcoming';
  if (days > 0) return 'soon';
  return -days <= 30 ? 'new' : 'released';
}

export function normalizeGame(game) {
  return {
    ...game,
    status: computeStatus(game.rel),
    released: !game.rel || new Date(game.rel) <= TODAY,
  };
}

export const RAW = [
  {id:1,yr:2020,t:"Doom Eternal",d:"Mar 20",p:"multi"},
  {id:2,yr:2020,t:"The Last of Us Part II",d:"Jun 19",p:"ps"},
  {id:3,yr:2020,t:"Ghost of Tsushima",d:"Jul 17",p:"ps"},
  {id:4,yr:2020,t:"Demon's Souls",d:"Nov 12",p:"ps"},
  {id:5,yr:2020,t:"Marvel's Spider-Man: Miles Morales",d:"Nov 12",p:"ps"},
  {id:6,yr:2020,t:"Call of Duty: Black Ops Cold War",d:"Nov 12",p:"multi"},
  {id:7,yr:2020,t:"Cyberpunk 2077",d:"Dec 10",p:"multi"},
  {id:8,yr:2021,t:"Returnal",d:"Apr 30",p:"ps"},
  {id:9,yr:2021,t:"Ratchet & Clank: Rift Apart",d:"Jun 11",p:"ps"},
  {id:10,yr:2021,t:"Scarlet Nexus",d:"Jun 25",p:"multi"},
  {id:11,yr:2021,t:"Far Cry 6",d:"Oct 07",p:"multi"},
  {id:12,yr:2021,t:"Call of Duty: Vanguard",d:"Nov 05",p:"multi"},
  {id:13,yr:2021,t:"Battlefield 2042",d:"Nov 19",p:"multi"},
  {id:14,yr:2021,t:"Halo Infinite",d:"Dec 08",p:"multi"},
  {id:15,yr:2022,t:"Dying Light 2",d:"Feb 04",p:"multi"},
  {id:16,yr:2022,t:"Horizon Forbidden West",d:"Feb 18",p:"ps"},
  {id:17,yr:2022,t:"Elden Ring",d:"Feb 25",p:"multi"},
  {id:18,yr:2022,t:"Gran Turismo 7",d:"Mar 04",p:"ps"},
  {id:19,yr:2022,t:"Call of Duty: Modern Warfare II",d:"Oct 28",p:"multi"},
  {id:20,yr:2022,t:"God of War Ragnarök",d:"Nov 09",p:"ps"},
  {id:21,yr:2023,t:"Dead Space",d:"Jan 27",p:"multi"},
  {id:22,yr:2023,t:"Hogwarts Legacy",d:"Feb 10",p:"multi"},
  {id:23,yr:2023,t:"Company of Heroes 3",d:"Feb 23",p:"multi"},
  {id:24,yr:2023,t:"Resident Evil 4 Remake",d:"Mar 24",p:"multi"},
  {id:25,yr:2023,t:"HFW: Burning Shores",d:"Apr 19",p:"ps"},
  {id:26,yr:2023,t:"Star Wars Jedi: Survivor",d:"Apr 28",p:"multi"},
  {id:27,yr:2023,t:"Diablo IV",d:"Jun 06",p:"multi"},
  {id:28,yr:2023,t:"Final Fantasy XVI",d:"Jun 22",p:"ps"},
  {id:29,yr:2023,t:"Baldur's Gate 3",d:"Aug 03",p:"multi"},
  {id:30,yr:2023,t:"Starfield",d:"Sep 06",p:"multi"},
  {id:31,yr:2023,t:"Marvel's Spider-Man 2",d:"Oct 20",p:"ps"},
  {id:32,yr:2023,t:"Alan Wake II",d:"Oct 27",p:"multi"},
  {id:33,yr:2023,t:"Call of Duty: Modern Warfare III",d:"Nov 10",p:"multi"},
  {id:34,yr:2024,t:"Persona 3 Reload",d:"Feb 02",p:"multi"},
  {id:35,yr:2024,t:"Helldivers 2",d:"Feb 08",p:"multi"},
  {id:36,yr:2024,t:"Rise of the Ronin",d:"Mar 22",p:"ps"},
  {id:37,yr:2024,t:"Stellar Blade",d:"Apr 26",p:"ps"},
  {id:38,yr:2024,t:"Destiny 2: The Final Shape",d:"Jun 04",p:"multi"},
  {id:39,yr:2024,t:"Elden Ring: Shadow of the Erdtree",d:"Jun 21",p:"multi"},
  {id:40,yr:2024,t:"Luigi's Mansion 2",d:"Jun 27",p:"nintendo"},
  {id:41,yr:2024,t:"Astro Bot",d:"Sep 06",p:"ps"},
  {id:42,yr:2024,t:"Space Marine 2",d:"Sep 09",p:"multi"},
  {id:43,yr:2024,t:"Diablo IV: Vessel of Hatred",d:"Oct 08",p:"multi"},
  {id:44,yr:2024,t:"Call of Duty: Black Ops 6",d:"Oct 25",p:"multi"},
  {id:45,yr:2024,t:"Indiana Jones and the Great Circle",d:"Dec 08",p:"multi"},
  {id:46,yr:2025,t:"Monster Hunter Wilds",d:"Feb 28",p:"multi"},
  {id:47,yr:2025,t:"Atomfall",d:"Mar 27",p:"multi"},
  {id:48,yr:2025,t:"TES IV: Oblivion Remastered",d:"Apr 22",p:"multi"},
  {id:49,yr:2025,t:"Clair Obscur: Expedition 33",d:"Apr 24",p:"multi"},
  {id:50,yr:2025,t:"Doom: The Dark Ages",d:"May 15",p:"multi"},
  {id:51,yr:2025,t:"Mario Kart World",d:"Jun 05",p:"nintendo"},
  {id:52,yr:2025,t:"Death Stranding 2: On the Beach",d:"Jun 26",p:"ps"},
  {id:53,yr:2025,t:"Donkey Kong Bonanza",d:"Jul 17",p:"nintendo"},
  {id:54,yr:2025,t:"Mafia: The Old Country",d:"Aug 08",p:"multi"},
  {id:55,yr:2025,t:"Metal Gear Solid: Snake Eater",d:"Aug 28",p:"multi"},
  {id:56,yr:2025,t:"Ghost of Yotei",d:"Oct 02",p:"ps"},
  {id:57,yr:2025,t:"Super Mario Galaxy",d:"Oct 02",p:"nintendo"},
  {id:58,yr:2025,t:"Battlefield 6",d:"Oct 10",p:"multi"},
  {id:59,yr:2025,t:"Pokémon Legends: Z-A",d:"Oct 16",p:"nintendo"},
  {id:60,yr:2025,t:"Metroid Prime 4",d:"Dec 04",p:"nintendo"},
  {id:61,yr:2026,t:"Resident Evil Requiem",d:"Feb 27",p:"multi",rel:"2026-02-27"},
  {id:62,yr:2026,t:"Toxic Commando",d:"Mar 12",p:"multi",rel:"2026-03-12"},
  {id:63,yr:2026,t:"Crimson Desert",d:"Mar 19",p:"multi",rel:"2026-03-19"},
  {id:64,yr:2026,t:"Mouse: P.I. for Hire",d:"Apr 16",p:"multi",rel:"2026-04-16"},
  {id:65,yr:2026,t:"Diablo IV: Lord of Hatred",d:"Apr 28",p:"multi",rel:"2026-04-28"},
  {id:66,yr:2026,t:"Mixtape",d:"May 07",p:"multi",rel:"2026-05-07"},
  {id:67,yr:2026,t:"Directive 8020",d:"May 12",p:"ps",rel:"2026-05-12"},
  {id:68,yr:2026,t:"Forza Horizon 6",d:"May 15",p:"multi",rel:"2026-05-15"},
  {id:69,yr:2026,t:"Lego Batman",d:"May 22",p:"multi",rel:"2026-05-22"},
  {id:70,yr:2026,t:"007 First Light",d:"May 27",p:"multi",rel:"2026-05-27"},
  {id:71,yr:2026,t:"Star Fox",d:"Jun 25",p:"nintendo",rel:"2026-06-25"},
  {id:72,yr:2026,t:"Assassin's Creed: Black Flag",d:"Jul 09",p:"multi",rel:"2026-07-09"},
  {id:73,yr:2026,t:"Splatoon Raiders",d:"Jul 23",p:"nintendo",rel:"2026-07-23"},
  {id:74,yr:2026,t:"Beast of Reincarnation",d:"Aug 04",p:"multi",rel:"2026-08-04"},
  {id:75,yr:2026,t:"Phantom Blade Zero",d:"Sep 09",p:"multi",rel:"2026-09-09"},
  {id:76,yr:2026,t:"Marvel's Wolverine",d:"Sep 15",p:"ps",rel:"2026-09-15"},
  {id:77,yr:2026,t:"Grand Theft Auto VI",d:"Nov 19",p:"multi",rel:"2026-11-19"},
];

export const GAMES = RAW.map(normalizeGame);

export const PL_LABEL = { multi: 'PC+', ps: 'PS', nintendo: 'NS' };

export const PLATFORM_OPTIONS = [
  { value: 'ps', label: 'PlayStation' },
  { value: 'nintendo', label: 'Nintendo' },
  { value: 'multi', label: 'Multi-platform' },
];

export const VALID_PLATFORMS = new Set(['ps', 'nintendo', 'multi']);
