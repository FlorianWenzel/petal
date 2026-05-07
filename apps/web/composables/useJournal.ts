import type { Mood, MonthEntries } from '~/lib/journal'

export function useJournal() {
  const entries = useState<MonthEntries>('journal.entries', () => ({}))
  const config = useRuntimeConfig()
  const base = config.public.apiBase

  async function load(year: number, month: number) {
    const m = String(month + 1).padStart(2, '0')
    entries.value = await $fetch<MonthEntries>(`${base}/entries/${year}/${m}`, {
      credentials: 'include',
    })
  }

  async function save(
    year: number,
    month: number,
    day: number,
    mood: Mood,
    text?: string,
  ) {
    const day2 = String(day).padStart(2, '0')
    const m = String(month + 1).padStart(2, '0')
    entries.value = await $fetch<MonthEntries>(`${base}/entries/${year}/${m}/${day2}`, {
      method: 'PUT',
      body: { mood, text },
      credentials: 'include',
    })
  }

  return { entries, load, save }
}
