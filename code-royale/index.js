const Width = 1920
const Height = 1000
const Queen = -1
const QueenRadius = 30
const QueenMovement = 60
const NoSite = -1
const Mine = 0
const Tower = 1
const Barracks = 2
const NoStructure = -1
const Friendly = 0
const Enemy = 1
const Knight = 0
const KnightCost = 80
const KnightBatch = 4
const Archer = 1
const ArcherCost = 100
const ArcherBatch = 2
const Giant = 2
const GiantCost = 140
const GiantBatch = 1
const MinimumEvade = 300
const MinimumRepair = 10
/** @returns {number[]} */
let line = () => readline().split(' ').map(Number)

let debug = (...args) => printErr(...args)
let head = a => a[0]
let tail = a => a.slice(1)
let assign = Object.assign.bind(Object)
let distance = (a, b) => Math.hypot(b.x - a.x, b.y - a.y)
let withDistance = (a, list) =>
  list.map(b => assign(b, { distance: distance(a, b) }))
let toString = obj =>
  `{${Object.keys(obj).map(k => `${k}: ${obj[k]}`).join(', ')}}`

let distanceAsc = (a, b) => a.distance > b.distance
let byOwner = owner => x => x.owner === owner
let byType = type => x => x.type === type
let bySpeciality = uType => x => x.param2 === uType
let int = n => Math.floor(n)
// end helperfunctions
function evade (player, sites, units) {
  const enemyKnights = units.filter(by(byOwner(Enemy), byType(Knight)))
  const nearest = head(enemyKnights)
  if (nearest) {
    if (nearest.distance < MinimumEvade) {
      let dx = player.x - nearest.x
      let dy = player.y - nearest.y
      return `MOVE ${int(dx)} ${int(dy)} HALP!`
    }
  }
  return build(player, sites, units)
}
function build (player, sites, units) {
  const { touchedSite } = player
  if (!touchedSite) return move(player, sites, units)

  const mine = ByType(Mine)
  if (mine(touchedSite)) {
    if (touchedSite.param1 < touchedSite.maxMineSize) {
      return `BUILD MINE`
    }
    return move(player, sites, units)
  }
  const tower = byType(Tower)
  if (tower(touchedSite)) {
    if (touchedSite.param1 < MinimumRepair) {
      return `BUILD TOWER`
    }
    return move(player, sites, units)
  }
  if (touchedSite.type != NoStructure) return move(player, sites, units)
  if (mySites.length === 0) return `BUILD MINE`

  const friendly = byOwner(Friendly)
  const barracks = byType(Barracks)

  const mySites = sites.filter(friendly)
  const myMines = mySites.filter(mine)
  const myBarracks = mySites.filter(barracks)
  const myTowers = mySites.filter(tower)

  const knight = bySpeciality(Knight)
  const archer = bySpeciality(Archer)
  const giant = bySpeciality(Giant)
  return 'WAIT'
}
function move (player, sites, units) {}
function turn (player, sites, units) {
  return evade(player, sites, units)
}
function train (player, sites, units) {
  return 'TRAIN'
}
// init
let sites = new Map()
let [numSites] = line()

for (let i = 0; i < numSites; i++) {
  let [id, x, y, radius] = line()
  sites.set(id, { id, x, y, radius })
}

// game loop
while (true) {
  let [gold, touchedSite] = line()
  let player = { gold, touchedSite }
  for (let i = 0; i < numSites; i++) {
    let [id, goldRemaining, maxMineSize, type, owner, param1, param2] = line()
    sites.set(
      id,
      assign(sites.get(id), {
        goldRemaining,
        maxMineSize,
        type,
        owner,
        param1,
        param2
      })
    )
  }
  let [numUnits] = line()
  let units = []
  for (let i = 0; i < numUnits; i++) {
    let [x, y, owner, type, health] = line()
    if (owner == Friendly && type == Queen) {
      player = assign(player, {
        x,
        y,
        owner,
        type,
        health,
        touchedSite: sites.get(player.touchedSite)
      })
    }
    units.push({ x, y, owner, type, health })
  }
  withDistance(player, units)
  withDistance(player, Array.from(sites.values()))
  units.sort(distanceAsc)
  debug(units.map(toString).join('\n'))
  let sitesList = Array.from(sites.values())
  sitesList.sort(distanceAsc)
  debug(sites.get(-1))
  print(turn(player, sitesList, units))
  print(train(player, sitesList, units))
}
