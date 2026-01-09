interface SteeringOptionsProps {
  includeConditional: boolean
  onChange: (value: boolean) => void
}

export function SteeringOptions({ includeConditional, onChange }: SteeringOptionsProps) {
  return (
    <div className="space-y-2">
      <label className="flex items-center gap-2">
        <input
          type="checkbox"
          checked={includeConditional}
          onChange={(e) => onChange(e.target.checked)}
          className="rounded border-input"
        />
        <span className="text-sm font-medium">Include conditional steering files</span>
      </label>
      <p className="text-xs text-muted-foreground ml-6">
        Adds security and quality rules for Go (*.go) and web (*.ts, *.tsx) files
      </p>
    </div>
  )
}
