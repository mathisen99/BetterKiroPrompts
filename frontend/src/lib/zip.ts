import JSZip from 'jszip'
import type { GeneratedFile } from './api'

export async function downloadAsZip(files: GeneratedFile[], zipName: string): Promise<void> {
  const zip = new JSZip()
  for (const file of files) {
    zip.file(file.path, file.content)
  }
  const blob = await zip.generateAsync({ type: 'blob' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = zipName
  a.click()
  URL.revokeObjectURL(url)
}
