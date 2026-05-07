<script setup lang="ts">
import { onMounted, onUnmounted, ref, watch } from 'vue'
import * as THREE from 'three'
import { OrbitControls } from 'three/examples/jsm/controls/OrbitControls.js'
import { easeOutBack, easeOutCubic } from '~/lib/easing'
import { MOODS, moodForColor } from '~/lib/journal'

const props = withDefaults(
  defineProps<{
    daysInMonth: number
    entriesByDay: Record<number, string>
    todayDay: number | null
    monthKey: string
    centerColor: string
    preview?: boolean
  }>(),
  { preview: false },
)

const emit = defineEmits<{
  'petal-click': [day: number]
}>()

const container = ref<HTMLDivElement | null>(null)

interface Placement {
  angle: number
  radius: number
  elevation: number
  tilt: number
  length: number
}

interface PetalEntry {
  day: number
  anchor: THREE.Group
  placement: THREE.Group
  bloom: THREE.Group
  mesh: THREE.Mesh
  outline: THREE.Mesh
  material: THREE.MeshToonMaterial
  geometry: THREE.BufferGeometry
  edgeGeometry: THREE.BufferGeometry
  color: string
  current: Placement
  painted: boolean
}

interface Animation {
  start: number
  duration: number
  update: (t: number) => void
  onComplete?: () => void
}

interface Particle {
  mesh: THREE.Mesh
  velocity: THREE.Vector3
  life: number
  duration: number
  geometry: THREE.BufferGeometry
  material: THREE.MeshBasicMaterial
}

let renderer: THREE.WebGLRenderer
let scene: THREE.Scene
let camera: THREE.PerspectiveCamera
let controls: OrbitControls
let flowerGroup: THREE.Group
let petalsGroup: THREE.Group
let pollenGroup: THREE.Group
let centerOrb: THREE.Mesh
let centerHalo: THREE.Mesh
let centerMaterial: THREE.MeshToonMaterial
let centerHaloMaterial: THREE.MeshBasicMaterial
let waterMesh: THREE.Mesh
let waterGeom: THREE.PlaneGeometry
let waterBaseZ: Float32Array
let ripples: { mesh: THREE.Mesh; mat: THREE.MeshBasicMaterial; phase: number }[] = []
let smallPads: { mesh: THREE.Mesh; baseY: number; phase: number; freq: number }[] = []
let outlineMaterial: THREE.MeshBasicMaterial
let hoverOutlineMaterial: THREE.MeshBasicMaterial
let todayOutlineMaterial: THREE.MeshBasicMaterial
let edgeMaterial: THREE.LineBasicMaterial
let petalTexture: THREE.CanvasTexture
let hovered: PetalEntry | null = null
const raycaster = new THREE.Raycaster()
const pointerNDC = new THREE.Vector2()
let pointerActive = false
let rafId = 0
let lastTime = 0
let baseY = 0
let elapsed = 0
let mounted = false
let transitioning = false
let pointerDownNDC = new THREE.Vector2()
let pointerDownTime = 0

const PETAL_LENGTH = 1.0

// Capacity per ring, outer → inner. Total >= 31 covers any month.
const RING_SLOTS = [12, 10, 8, 4]
const MAX_CAPACITY = RING_SLOTS.reduce((a, b) => a + b, 0)

function distribute(count: number): number[] {
  const result: number[] = []
  let rem = Math.min(count, MAX_CAPACITY)
  for (const cap of RING_SLOTS) {
    if (rem <= 0) break
    result.push(Math.min(cap, rem))
    rem -= result[result.length - 1]
  }
  return result
}

const petals: PetalEntry[] = []
const animations: Animation[] = []
const particles: Particle[] = []

function makePetalTexture(): THREE.CanvasTexture {
  const size = 256
  const c = document.createElement('canvas')
  c.width = size
  c.height = size
  const ctx = c.getContext('2d')!

  ctx.fillStyle = '#ffffff'
  ctx.fillRect(0, 0, size, size)

  const grad = ctx.createLinearGradient(0, size, 0, 0)
  grad.addColorStop(0, 'rgba(255,255,255,1)')
  grad.addColorStop(0.55, 'rgba(248,240,242,1)')
  grad.addColorStop(1, 'rgba(220,200,208,1)')
  ctx.fillStyle = grad
  ctx.fillRect(0, 0, size, size)

  const r = ctx.createRadialGradient(size / 2, size * 0.85, 4, size / 2, size * 0.85, size * 0.55)
  r.addColorStop(0, 'rgba(255,255,255,0.55)')
  r.addColorStop(1, 'rgba(255,255,255,0)')
  ctx.fillStyle = r
  ctx.fillRect(0, 0, size, size)

  ctx.strokeStyle = 'rgba(150,130,140,0.3)'
  ctx.lineWidth = 1.5
  for (let i = -2; i <= 2; i++) {
    const x0 = size / 2
    const x1 = size / 2 + i * 22
    ctx.beginPath()
    ctx.moveTo(x0, size - 8)
    ctx.bezierCurveTo(x0 + i * 6, size * 0.7, x1, size * 0.35, x1, 16)
    ctx.stroke()
  }

  const img = ctx.getImageData(0, 0, size, size)
  for (let i = 0; i < img.data.length; i += 4) {
    const n = (Math.random() - 0.5) * 10
    img.data[i] = Math.max(0, Math.min(255, img.data[i] + n))
    img.data[i + 1] = Math.max(0, Math.min(255, img.data[i + 1] + n))
    img.data[i + 2] = Math.max(0, Math.min(255, img.data[i + 2] + n))
  }
  ctx.putImageData(img, 0, 0)

  const tex = new THREE.CanvasTexture(c)
  tex.colorSpace = THREE.SRGBColorSpace
  tex.anisotropy = 4
  return tex
}

function makePetalGeometry(): THREE.BufferGeometry {
  const shape = new THREE.Shape()
  shape.moveTo(0, 0)
  shape.bezierCurveTo(0.16, 0.1, 0.32, 0.32, 0.3, 0.55)
  shape.bezierCurveTo(0.28, 0.82, 0.14, 1.02, 0, 1.08)
  shape.bezierCurveTo(-0.14, 1.02, -0.28, 0.82, -0.3, 0.55)
  shape.bezierCurveTo(-0.32, 0.32, -0.16, 0.1, 0, 0)

  const geom = new THREE.ExtrudeGeometry(shape, {
    depth: 0.05,
    bevelEnabled: true,
    bevelThickness: 0.03,
    bevelSize: 0.025,
    bevelSegments: 3,
    curveSegments: 28,
  })
  geom.translate(0, 0, -0.025)
  geom.scale(1, PETAL_LENGTH, 1)
  const pos = geom.attributes.position
  const v = new THREE.Vector3()
  for (let i = 0; i < pos.count; i++) {
    v.fromBufferAttribute(pos, i)
    const t = THREE.MathUtils.clamp(v.y / PETAL_LENGTH, 0, 1)
    v.z += -0.28 * t * t
    pos.setXYZ(i, v.x, v.y, v.z)
  }
  pos.needsUpdate = true
  geom.rotateZ(-Math.PI / 2)
  geom.rotateX(Math.PI / 2)
  geom.computeVertexNormals()
  return geom
}

function makeShadowTexture(): THREE.CanvasTexture {
  const size = 256
  const c = document.createElement('canvas')
  c.width = size
  c.height = size
  const ctx = c.getContext('2d')!
  const g = ctx.createRadialGradient(size / 2, size / 2, 8, size / 2, size / 2, size / 2)
  g.addColorStop(0, 'rgba(40,30,40,0.85)')
  g.addColorStop(0.4, 'rgba(40,30,40,0.5)')
  g.addColorStop(1, 'rgba(40,30,40,0)')
  ctx.fillStyle = g
  ctx.fillRect(0, 0, size, size)
  const tex = new THREE.CanvasTexture(c)
  tex.colorSpace = THREE.SRGBColorSpace
  return tex
}

function makeFlowerBase() {
  flowerGroup = new THREE.Group()

  const orbGeom = new THREE.SphereGeometry(0.34, 32, 24)
  orbGeom.scale(1, 0.85, 1)
  centerMaterial = new THREE.MeshToonMaterial({ color: new THREE.Color(props.centerColor) })
  centerOrb = new THREE.Mesh(orbGeom, centerMaterial)
  centerOrb.position.y = 0.22
  flowerGroup.add(centerOrb)

  const orbOutline = new THREE.Mesh(
    orbGeom,
    new THREE.MeshBasicMaterial({ color: 0x3a2230, side: THREE.BackSide }),
  )
  orbOutline.scale.setScalar(1.06)
  centerOrb.add(orbOutline)

  const haloGeom = new THREE.SphereGeometry(0.48, 24, 18)
  centerHaloMaterial = new THREE.MeshBasicMaterial({
    color: new THREE.Color(props.centerColor),
    transparent: true,
    opacity: 0.18,
    depthWrite: false,
  })
  centerHalo = new THREE.Mesh(haloGeom, centerHaloMaterial)
  centerHalo.position.y = 0.22
  flowerGroup.add(centerHalo)

  if (!props.preview) {
    const padGeom = new THREE.CylinderGeometry(1.5, 1.5, 0.1, 48)
    const padMat = new THREE.MeshToonMaterial({ color: 0x4f9156 })
    const pad = new THREE.Mesh(padGeom, padMat)
    pad.position.y = -0.2
    flowerGroup.add(pad)
  } else {
    const cupGeom = new THREE.SphereGeometry(0.95, 32, 16, 0, Math.PI * 2, 0, Math.PI / 2)
    cupGeom.scale(1, 0.45, 1)
    const cupMat = new THREE.MeshToonMaterial({
      color: new THREE.Color(props.centerColor).lerp(new THREE.Color(0x000000), 0.15),
    })
    const cup = new THREE.Mesh(cupGeom, cupMat)
    cup.position.y = -0.08
    flowerGroup.add(cup)

    const cupOutline = new THREE.Mesh(
      cupGeom,
      new THREE.MeshBasicMaterial({ color: 0x3a2230, side: THREE.BackSide }),
    )
    cupOutline.scale.setScalar(1.04)
    cup.add(cupOutline)

    const shadowGeom = new THREE.PlaneGeometry(2.6, 2.6)
    shadowGeom.rotateX(-Math.PI / 2)
    const shadowMat = new THREE.MeshBasicMaterial({
      map: makeShadowTexture(),
      transparent: true,
      opacity: 0.45,
      depthWrite: false,
    })
    const shadow = new THREE.Mesh(shadowGeom, shadowMat)
    shadow.position.y = -0.6
    flowerGroup.add(shadow)
  }

  petalsGroup = new THREE.Group()
  flowerGroup.add(petalsGroup)

  baseY = 0.05
  flowerGroup.position.y = baseY
  scene.add(flowerGroup)

  pollenGroup = new THREE.Group()
  scene.add(pollenGroup)
}

function makeWater() {
  waterGeom = new THREE.PlaneGeometry(18, 18, 64, 64)
  waterGeom.rotateX(-Math.PI / 2)
  const pos = waterGeom.attributes.position
  waterBaseZ = new Float32Array(pos.count)
  for (let i = 0; i < pos.count; i++) waterBaseZ[i] = pos.getY(i)

  const waterMat = new THREE.MeshToonMaterial({
    color: 0x6fa8c5,
    transparent: true,
    opacity: 0.92,
  })
  waterMesh = new THREE.Mesh(waterGeom, waterMat)
  waterMesh.position.y = -0.3
  scene.add(waterMesh)

  // darker depth ring under lily pad
  const depthGeom = new THREE.CircleGeometry(2.2, 48)
  depthGeom.rotateX(-Math.PI / 2)
  const depthMat = new THREE.MeshBasicMaterial({
    color: 0x3d6a85,
    transparent: true,
    opacity: 0.35,
  })
  const depth = new THREE.Mesh(depthGeom, depthMat)
  depth.position.y = -0.29
  scene.add(depth)

  // scattered small lily pads
  const padShades = [0x4f9156, 0x5fa066, 0x447a4b, 0x6fb074]
  const placements = [
    { r: 2.6, a: 0.3, scale: 0.55 },
    { r: 3.2, a: 1.7, scale: 0.4 },
    { r: 2.9, a: 2.9, scale: 0.6 },
    { r: 3.6, a: 4.1, scale: 0.45 },
    { r: 2.4, a: 5.2, scale: 0.5 },
    { r: 4.1, a: 0.9, scale: 0.35 },
    { r: 3.7, a: 3.5, scale: 0.5 },
  ]
  for (let i = 0; i < placements.length; i++) {
    const p = placements[i]
    const shape = new THREE.Shape()
    const segs = 40
    const notch = 0.45
    for (let s = 0; s <= segs; s++) {
      const ang = notch / 2 + (s / segs) * (Math.PI * 2 - notch)
      const x = Math.cos(ang)
      const y = Math.sin(ang)
      if (s === 0) shape.moveTo(x, y)
      else shape.lineTo(x, y)
    }
    shape.lineTo(0, 0)
    const padG = new THREE.ExtrudeGeometry(shape, {
      depth: 0.06,
      bevelEnabled: true,
      bevelThickness: 0.02,
      bevelSize: 0.03,
      bevelSegments: 2,
    })
    padG.rotateX(-Math.PI / 2)
    padG.scale(p.scale, 1, p.scale)
    const padM = new THREE.MeshToonMaterial({ color: padShades[i % padShades.length] })
    const m = new THREE.Mesh(padG, padM)
    m.position.set(Math.cos(p.a) * p.r, -0.22, Math.sin(p.a) * p.r)
    m.rotation.y = Math.random() * Math.PI * 2
    scene.add(m)
    smallPads.push({
      mesh: m,
      baseY: -0.22,
      phase: Math.random() * Math.PI * 2,
      freq: 0.7 + Math.random() * 0.5,
    })
  }

  // expanding ripple rings around pad
  const ringGeom = new THREE.RingGeometry(0.98, 1.0, 64)
  ringGeom.rotateX(-Math.PI / 2)
  for (let i = 0; i < 3; i++) {
    const mat = new THREE.MeshBasicMaterial({
      color: 0xcfe4ee,
      transparent: true,
      opacity: 0,
      side: THREE.DoubleSide,
      depthWrite: false,
    })
    const mesh = new THREE.Mesh(ringGeom, mat)
    mesh.position.y = -0.28
    scene.add(mesh)
    ripples.push({ mesh, mat, phase: i / 3 })
  }
}

const RING_PROFILES = [
  { radius: 0.58, elevation: 0.0,  tilt: THREE.MathUtils.degToRad(4),  length: 1.4  },
  { radius: 0.46, elevation: 0.12, tilt: THREE.MathUtils.degToRad(22), length: 1.2  },
  { radius: 0.34, elevation: 0.24, tilt: THREE.MathUtils.degToRad(42), length: 1.0  },
  { radius: 0.22, elevation: 0.36, tilt: THREE.MathUtils.degToRad(60), length: 0.85 },
]

function placementsFor(count: number): Placement[] {
  const rings = distribute(count)
  const out: Placement[] = []
  rings.forEach((slots, idx) => {
    const profile = RING_PROFILES[Math.min(idx, RING_PROFILES.length - 1)]
    const offset = idx % 2 === 1 ? Math.PI / slots : 0
    for (let s = 0; s < slots; s++) {
      const angle = ((s + 0.5) / slots) * Math.PI * 2 - Math.PI / 2 + offset
      out.push({
        angle,
        radius: profile.radius,
        elevation: profile.elevation,
        tilt: profile.tilt,
        length: profile.length,
      })
    }
  })
  return out
}

function spawnPollen(originWorld: THREE.Vector3, outwardWorld: THREE.Vector3) {
  const count = 12 + Math.floor(Math.random() * 4)
  const baseGeom = new THREE.SphereGeometry(0.045, 8, 6)
  for (let i = 0; i < count; i++) {
    const mat = new THREE.MeshBasicMaterial({
      color: 0xffe066,
      transparent: true,
      opacity: 1,
    })
    const mesh = new THREE.Mesh(baseGeom, mat)
    mesh.position.copy(originWorld)
    pollenGroup.add(mesh)

    const out = outwardWorld.clone().normalize()
    const jitter = new THREE.Vector3(
      (Math.random() - 0.5) * 0.6,
      Math.random() * 0.6 + 0.4,
      (Math.random() - 0.5) * 0.6,
    )
    const speed = 1.2 + Math.random() * 0.8
    const velocity = out.multiplyScalar(speed).add(jitter)

    particles.push({
      mesh,
      velocity,
      life: 0,
      duration: 0.9 + Math.random() * 0.3,
      geometry: baseGeom,
      material: mat,
    })
  }
}

function flowerBounce(now: number) {
  const start = now
  const dur = 0.45
  animations.push({
    start,
    duration: dur,
    update: (t) => {
      const s = t < 0.4 ? 1 + (0.05 * (t / 0.4)) : 1.05 - (0.05 * ((t - 0.4) / 0.6))
      flowerGroup.scale.setScalar(s)
    },
    onComplete: () => flowerGroup.scale.setScalar(1),
  })
}

const PETAL_WHITE = 0xfaf4ee

function createPetal(target: Placement, day: number, index: number): PetalEntry {
  const anchor = new THREE.Group()
  const placement = new THREE.Group()
  const bloom = new THREE.Group()
  const geometry = makePetalGeometry()
  const material = new THREE.MeshToonMaterial({ color: PETAL_WHITE, map: petalTexture })
  const mesh = new THREE.Mesh(geometry, material)

  const outline = new THREE.Mesh(geometry, outlineMaterial)
  outline.scale.setScalar(1.04)
  outline.renderOrder = -1
  mesh.add(outline)

  const edgeGeom = new THREE.BufferGeometry()

  bloom.add(mesh)
  placement.add(bloom)
  anchor.add(placement)
  petalsGroup.add(anchor)

  anchor.rotation.y = target.angle
  placement.position.set(target.radius, target.elevation, 0)
  placement.rotation.z = target.tilt
  mesh.scale.x = target.length
  bloom.scale.setScalar(1)

  mesh.userData.petalIndex = index
  outline.userData.petalIndex = index

  return {
    day,
    anchor,
    placement,
    bloom,
    mesh,
    outline,
    material,
    geometry,
    edgeGeometry: edgeGeom,
    color: '#ffffff',
    current: { ...target },
    painted: false,
  }
}

function buildPetals(count: number) {
  for (const p of petals) {
    petalsGroup.remove(p.anchor)
    p.geometry.dispose()
    p.edgeGeometry.dispose()
    p.material.dispose()
  }
  petals.length = 0

  const targets = placementsFor(count)
  const built = Math.min(count, targets.length)
  for (let i = 0; i < built; i++) {
    petals.push(createPetal(targets[i], i + 1, i))
  }
}

function moodTreatedColor(rawHex: string): THREE.Color {
  const c = new THREE.Color(rawHex)
  const mood = moodForColor(rawHex)
  if (!mood) return c
  const profile = MOODS.find((m) => m.id === mood)
  const sat = profile?.saturation ?? 1
  if (sat !== 1) {
    const hsl = { h: 0, s: 0, l: 0 }
    c.getHSL(hsl)
    hsl.s = Math.max(0, Math.min(1, hsl.s * sat))
    c.setHSL(hsl.h, hsl.s, hsl.l)
  }
  return c
}

function moodGeometry(mood: ReturnType<typeof moodForColor>) {
  const profile = mood ? MOODS.find((m) => m.id === mood) : null
  return {
    droop: profile?.droopRad ?? 0,
    lenScale: profile?.lengthScale ?? 1,
  }
}

function paintPetal(p: PetalEntry, color: string, animate: boolean) {
  p.painted = true
  p.color = color
  const endColor = moodTreatedColor(color)
  const mood = moodForColor(color)
  const { droop, lenScale } = moodGeometry(mood)
  const targetTilt = p.current.tilt + droop
  const targetLen = p.current.length * lenScale

  if (!animate) {
    p.material.color.copy(endColor)
    p.placement.rotation.z = targetTilt
    p.mesh.scale.x = targetLen
    return
  }

  const now = performance.now() / 1000
  const startColor = p.material.color.clone()
  const startTilt = p.placement.rotation.z
  const startLen = p.mesh.scale.x

  animations.push({
    start: now,
    duration: 0.55,
    update: (t) => {
      const e = easeOutCubic(t)
      p.material.color.lerpColors(startColor, endColor, e)
      p.placement.rotation.z = startTilt + (targetTilt - startTilt) * e
      p.mesh.scale.x = startLen + (targetLen - startLen) * e
    },
  })

  animations.push({
    start: now,
    duration: 0.7,
    update: (t) => {
      const pulse = Math.sin(t * Math.PI) * 0.18
      p.bloom.scale.setScalar(1 + pulse * easeOutBack(Math.min(t * 2, 1)))
    },
    onComplete: () => p.bloom.scale.setScalar(1),
  })

  const tipLocal = new THREE.Vector3(PETAL_LENGTH, 0, 0)
  const tipWorld = p.mesh.localToWorld(tipLocal.clone())
  const baseWorld = p.mesh.localToWorld(new THREE.Vector3(0, 0, 0))
  const outward = tipWorld.clone().sub(baseWorld)
  spawnPollen(tipWorld, outward)
  flowerBounce(now)
}

function unpaintPetal(p: PetalEntry) {
  p.painted = false
  p.color = '#ffffff'
  p.material.color.setHex(PETAL_WHITE)
  p.bloom.scale.setScalar(1)
  p.placement.rotation.z = p.current.tilt
  p.mesh.scale.x = p.current.length
}

function syncEntries(entries: Record<number, string>, animateNew: boolean) {
  for (const p of petals) {
    const c = entries[p.day]
    if (c) {
      const incoming = new THREE.Color(c).getHex()
      const current = p.material.color.getHex()
      if (!p.painted) {
        paintPetal(p, c, animateNew)
      } else if (incoming !== current) {
        paintPetal(p, c, animateNew)
      }
    } else if (p.painted) {
      unpaintPetal(p)
    }
  }
}

function startMonthTransition() {
  // cancel any in-flight transition anims to avoid stacking
  animations.length = 0
  transitioning = true
  const now = performance.now() / 1000
  const exitDur = 0.45
  for (const p of petals) {
    p.material.transparent = true
  }
  animations.push({
    start: now,
    duration: exitDur,
    update: (t) => {
      const e = easeOutCubic(t)
      petalsGroup.scale.setScalar(Math.max(0.001, 1 - e))
      petalsGroup.rotation.y = e * Math.PI * 0.55
      for (const p of petals) {
        p.material.opacity = Math.max(0, 1 - e * 1.1)
      }
    },
    onComplete: () => {
      buildPetals(props.daysInMonth)
      syncEntries(props.entriesByDay, false)
      applyTodayOutline()

      petalsGroup.scale.setScalar(0.001)
      petalsGroup.rotation.y = -Math.PI * 0.55
      for (const p of petals) {
        p.material.transparent = true
        p.material.opacity = 0
      }

      const now2 = performance.now() / 1000
      animations.push({
        start: now2,
        duration: 0.65,
        update: (t) => {
          const eb = easeOutBack(t)
          petalsGroup.scale.setScalar(Math.max(0.001, eb))
          petalsGroup.rotation.y = -Math.PI * 0.55 * (1 - easeOutCubic(t))
          const op = Math.min(1, t * 1.6)
          for (const p of petals) {
            p.material.opacity = op
          }
        },
        onComplete: () => {
          petalsGroup.scale.setScalar(1)
          petalsGroup.rotation.y = 0
          for (const p of petals) {
            p.material.opacity = 1
            p.material.transparent = false
          }
          transitioning = false
        },
      })
    },
  })
}

function applyTodayOutline() {
  for (const p of petals) {
    const isToday = p.day === props.todayDay
    if (p === hovered) continue
    p.outline.material = isToday ? todayOutlineMaterial : outlineMaterial
    p.outline.scale.setScalar(isToday ? 1.08 : 1.04)
  }
}

function tick() {
  const now = performance.now() / 1000
  const dt = lastTime === 0 ? 0 : Math.min(now - lastTime, 0.1)
  lastTime = now
  elapsed += dt

  flowerGroup.rotation.y += 0.06 * dt
  flowerGroup.position.y = baseY + Math.sin(elapsed * 1.0) * 0.03

  if (waterMesh) {
    const pos = waterGeom.attributes.position
    const t = elapsed
    for (let i = 0; i < pos.count; i++) {
      const x = pos.getX(i)
      const z = pos.getZ(i)
      const w =
        Math.sin(x * 0.7 + t * 1.1) * 0.06 +
        Math.sin(z * 0.9 - t * 0.8) * 0.05 +
        Math.sin((x + z) * 0.4 + t * 0.5) * 0.04
      pos.setY(i, w)
    }
    pos.needsUpdate = true

    for (const sp of smallPads) {
      sp.mesh.position.y = sp.baseY + Math.sin(elapsed * sp.freq + sp.phase) * 0.04
      sp.mesh.rotation.z = Math.sin(elapsed * sp.freq * 0.6 + sp.phase) * 0.04
    }

    for (const r of ripples) {
      const cycle = 4.0
      const local = ((elapsed + r.phase * cycle) % cycle) / cycle
      const s = 1 + local * 1.6
      r.mesh.scale.setScalar(s)
      r.mat.opacity = (1 - local) * 0.35
    }
  }

  if (centerOrb) {
    const breath = 1 + Math.sin(elapsed * 1.6) * 0.04
    centerOrb.scale.setScalar(breath)
    centerHalo.scale.setScalar(1 + Math.sin(elapsed * 1.6 + 0.4) * 0.08)
    centerHaloMaterial.opacity = 0.16 + Math.sin(elapsed * 1.6) * 0.06
  }

  if (pointerActive && !transitioning) {
    raycaster.setFromCamera(pointerNDC, camera)
    const meshes = petals.map((p) => p.mesh)
    const hits = raycaster.intersectObjects(meshes, false)
    const next = hits.length > 0 ? petals[hits[0].object.userData.petalIndex] : null
    if (next !== hovered) {
      if (hovered) {
        const wasToday = hovered.day === props.todayDay
        hovered.outline.material = wasToday ? todayOutlineMaterial : outlineMaterial
        hovered.outline.scale.setScalar(wasToday ? 1.08 : 1.04)
      }
      hovered = next
      if (hovered) {
        hovered.outline.material = hoverOutlineMaterial
        hovered.outline.scale.setScalar(1.1)
      }
      const dom = renderer.domElement
      dom.style.cursor = hovered ? 'pointer' : ''
    }
  } else if (hovered) {
    const wasToday = hovered.day === props.todayDay
    hovered.outline.material = wasToday ? todayOutlineMaterial : outlineMaterial
    hovered.outline.scale.setScalar(wasToday ? 1.08 : 1.04)
    hovered = null
    renderer.domElement.style.cursor = ''
  }

  for (let i = animations.length - 1; i >= 0; i--) {
    const a = animations[i]
    const t = (now - a.start) / a.duration
    if (t >= 1) {
      a.update(1)
      a.onComplete?.()
      animations.splice(i, 1)
    } else {
      a.update(Math.max(0, t))
    }
  }

  for (let i = particles.length - 1; i >= 0; i--) {
    const p = particles[i]
    p.life += dt
    const t = p.life / p.duration
    if (t >= 1) {
      pollenGroup.remove(p.mesh)
      p.material.dispose()
      particles.splice(i, 1)
      continue
    }
    p.mesh.position.addScaledVector(p.velocity, dt)
    p.velocity.y -= 0.6 * dt
    p.velocity.multiplyScalar(1 - 0.8 * dt)
    p.material.opacity = easeOutCubic(1 - t)
    const s = 1 - t * 0.4
    p.mesh.scale.setScalar(s)
  }

  controls.update()
  renderer.render(scene, camera)
  rafId = requestAnimationFrame(tick)
}

function setNDC(ev: PointerEvent, target: HTMLElement, out: THREE.Vector2) {
  const rect = target.getBoundingClientRect()
  out.x = ((ev.clientX - rect.left) / rect.width) * 2 - 1
  out.y = -((ev.clientY - rect.top) / rect.height) * 2 + 1
}

function onPointerMove(ev: PointerEvent) {
  setNDC(ev, ev.currentTarget as HTMLElement, pointerNDC)
  pointerActive = true
}

function onPointerLeave() {
  pointerActive = false
}

function onPointerDown(ev: PointerEvent) {
  setNDC(ev, ev.currentTarget as HTMLElement, pointerDownNDC)
  pointerDownTime = performance.now()
}

function onPointerUp(ev: PointerEvent) {
  if (transitioning) return
  const upNDC = new THREE.Vector2()
  setNDC(ev, ev.currentTarget as HTMLElement, upNDC)
  const dx = upNDC.x - pointerDownNDC.x
  const dy = upNDC.y - pointerDownNDC.y
  const dragDist = Math.sqrt(dx * dx + dy * dy)
  const elapsedMs = performance.now() - pointerDownTime
  if (dragDist > 0.02 || elapsedMs > 400) return

  raycaster.setFromCamera(upNDC, camera)
  const meshes = petals.map((p) => p.mesh)
  const hits = raycaster.intersectObjects(meshes, false)
  if (hits.length > 0) {
    const idx = hits[0].object.userData.petalIndex as number
    const p = petals[idx]
    if (p) emit('petal-click', p.day)
  }
}

function onResize() {
  if (!container.value) return
  const w = container.value.clientWidth
  const h = container.value.clientHeight
  camera.aspect = w / h
  camera.updateProjectionMatrix()
  renderer.setSize(w, h)
}

watch(
  () => [props.monthKey, props.entriesByDay] as const,
  ([newKey], oldVals) => {
    if (!mounted) return
    const oldKey = oldVals?.[0]
    if (newKey !== oldKey) {
      startMonthTransition()
    } else if (!transitioning) {
      syncEntries(props.entriesByDay, true)
    }
  },
  { deep: true },
)

watch(
  () => props.todayDay,
  () => {
    if (!mounted) return
    applyTodayOutline()
  },
)

watch(
  () => props.centerColor,
  (c) => {
    if (!mounted) return
    const now = performance.now() / 1000
    const startA = centerMaterial.color.clone()
    const startB = centerHaloMaterial.color.clone()
    const target = new THREE.Color(c)
    animations.push({
      start: now,
      duration: 0.6,
      update: (t) => {
        const e = easeOutCubic(t)
        centerMaterial.color.lerpColors(startA, target, e)
        centerHaloMaterial.color.lerpColors(startB, target, e)
      },
    })
  },
)

onMounted(() => {
  if (!container.value) return
  const w = container.value.clientWidth
  const h = container.value.clientHeight

  scene = new THREE.Scene()
  scene.background = null

  camera = new THREE.PerspectiveCamera(40, w / h, 0.1, 100)
  if (props.preview) {
    camera.position.set(0, 5.7, 4.0)
  } else {
    camera.position.set(0, 2.4, 5.0)
  }
  camera.lookAt(0, 0.3, 0)

  renderer = new THREE.WebGLRenderer({ antialias: true, alpha: true })
  renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2))
  renderer.setSize(w, h)
  renderer.domElement.style.display = 'block'
  renderer.outputColorSpace = THREE.SRGBColorSpace
  container.value.appendChild(renderer.domElement)

  const ambient = new THREE.AmbientLight(0xfff2e0, 0.6)
  scene.add(ambient)

  const dir = new THREE.DirectionalLight(0xffffff, 0.85)
  dir.position.set(2.5, 4, 3)
  scene.add(dir)

  const fill = new THREE.DirectionalLight(0xffd9ec, 0.25)
  fill.position.set(-3, 1, -2)
  scene.add(fill)

  petalTexture = makePetalTexture()

  outlineMaterial = new THREE.MeshBasicMaterial({
    color: 0x3a2230,
    side: THREE.BackSide,
  })
  hoverOutlineMaterial = new THREE.MeshBasicMaterial({
    color: 0xfff2a8,
    side: THREE.BackSide,
  })
  todayOutlineMaterial = new THREE.MeshBasicMaterial({
    color: 0x6cd2c8,
    side: THREE.BackSide,
  })
  edgeMaterial = new THREE.LineBasicMaterial({
    color: 0x3a2230,
    transparent: true,
    opacity: 0.85,
  })

  controls = new OrbitControls(camera, renderer.domElement)
  controls.enableDamping = true
  controls.dampingFactor = 0.08
  controls.target.set(0, 0.3, 0)
  controls.minDistance = 3
  controls.maxDistance = 8
  controls.maxPolarAngle = Math.PI / 2 - 0.1
  controls.minPolarAngle = 0.2
  controls.enablePan = false
  if (props.preview) {
    controls.autoRotate = true
    controls.autoRotateSpeed = 0.5
  }

  makeFlowerBase()
  if (!props.preview) makeWater()
  buildPetals(props.daysInMonth)
  syncEntries(props.entriesByDay, false)
  applyTodayOutline()
  mounted = true

  const dom = renderer.domElement
  dom.addEventListener('pointermove', onPointerMove)
  dom.addEventListener('pointerleave', onPointerLeave)
  dom.addEventListener('pointerdown', onPointerDown)
  dom.addEventListener('pointerup', onPointerUp)

  window.addEventListener('resize', onResize)
  rafId = requestAnimationFrame(tick)
})

onUnmounted(() => {
  cancelAnimationFrame(rafId)
  window.removeEventListener('resize', onResize)

  for (const p of petals) {
    petalsGroup.remove(p.anchor)
    p.geometry.dispose()
    p.edgeGeometry.dispose()
    p.material.dispose()
  }
  petals.length = 0
  for (const part of particles) {
    pollenGroup.remove(part.mesh)
    part.material.dispose()
  }
  particles.length = 0

  scene.traverse((obj) => {
    const mesh = obj as THREE.Mesh
    if (mesh.geometry) mesh.geometry.dispose()
    const m = mesh.material as THREE.Material | THREE.Material[] | undefined
    if (Array.isArray(m)) m.forEach((mm) => mm.dispose())
    else if (m) m.dispose()
  })

  outlineMaterial.dispose()
  hoverOutlineMaterial.dispose()
  todayOutlineMaterial.dispose()
  edgeMaterial.dispose()
  petalTexture.dispose()

  renderer.dispose()
  if (renderer.domElement.parentElement) {
    renderer.domElement.parentElement.removeChild(renderer.domElement)
  }
})
</script>

<template>
  <div ref="container" class="w-full h-full"></div>
</template>
