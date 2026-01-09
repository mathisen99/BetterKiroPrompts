import { KickoffWizard } from '../components/kickoff/KickoffWizard'

export function KickoffPage() {
  return (
    <main className="container mx-auto max-w-3xl px-4 py-8">
      <h1 className="mb-8 text-3xl font-bold">Kickoff Prompt Generator</h1>
      <KickoffWizard />
    </main>
  )
}
