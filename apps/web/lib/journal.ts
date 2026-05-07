export type Mood = 'bright' | 'tender' | 'calm' | 'heavy' | 'stormy'

export interface MoodMeta {
  id: Mood
  label: string
  color: string
  blurb: string
  /** extra tilt in radians applied to petal placement; positive lifts up, negative droops down */
  droopRad?: number
  /** scalar applied to petal length (1.0 = unchanged) */
  lengthScale?: number
  /** scalar applied to color saturation (1.0 = unchanged, <1 = washed out) */
  saturation?: number
}

export const MOODS: MoodMeta[] = [
  {
    id: 'bright',  label: 'bright',  color: '#FFD93D', blurb: 'lit up, awake, glad',
    droopRad: 0.08, lengthScale: 1.05,
  },
  {
    id: 'tender',  label: 'tender',  color: '#FF9EBB', blurb: 'warm, loving, soft on the edges',
    lengthScale: 1.0,
  },
  {
    id: 'calm',    label: 'calm',    color: '#B5C7F7', blurb: 'settled, slow, content',
    lengthScale: 1.0,
  },
  {
    id: 'heavy',   label: 'heavy',   color: '#4D6F9E', blurb: 'tired, low, blue',
    droopRad: -0.28, lengthScale: 0.88, saturation: 0.55,
  },
  {
    id: 'stormy',  label: 'stormy',  color: '#FFA372', blurb: 'tense, restless, on edge',
    droopRad: 0.05, lengthScale: 1.06,
  },
]

export function colorFor(mood: Mood): string {
  return MOODS.find((m) => m.id === mood)?.color ?? '#FFD93D'
}

export function moodForColor(hex: string): Mood | null {
  const h = hex.toUpperCase()
  return MOODS.find((m) => m.color.toUpperCase() === h)?.id ?? null
}

export interface DayEntry {
  mood: Mood
  text?: string
}

export type MonthEntries = Record<string, DayEntry>

export function daysInMonth(year: number, month: number): number {
  return new Date(year, month + 1, 0).getDate()
}

const EMPTY_CENTER = '#f4ece2'

function hexToRgb(hex: string): [number, number, number] {
  const h = hex.replace('#', '')
  const n = parseInt(h, 16)
  return [(n >> 16) & 255, (n >> 8) & 255, n & 255]
}

function rgbToHex(r: number, g: number, b: number): string {
  const c = (v: number) => Math.max(0, Math.min(255, Math.round(v))).toString(16).padStart(2, '0')
  return `#${c(r)}${c(g)}${c(b)}`
}

export function averageColor(entries: MonthEntries): string {
  const vals = Object.values(entries)
  if (vals.length === 0) return EMPTY_CENTER
  let r = 0, g = 0, b = 0
  for (const e of vals) {
    const [cr, cg, cb] = hexToRgb(colorFor(e.mood))
    r += cr; g += cg; b += cb
  }
  return rgbToHex(r / vals.length, g / vals.length, b / vals.length)
}
