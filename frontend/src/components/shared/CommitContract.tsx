export function CommitContract() {
  return (
    <div className="rounded-lg border border-border bg-muted p-4">
      <h3 className="mb-2 text-sm font-medium">Commit Message Contract</h3>
      <ul className="space-y-1 text-sm text-muted-foreground">
        <li>• Atomic: one concern per commit</li>
        <li>• Prefixed: feat:, fix:, docs:, chore:</li>
        <li>• One-sentence summary</li>
      </ul>
    </div>
  )
}
