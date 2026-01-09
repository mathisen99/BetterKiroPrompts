interface PresetCardProps {
  value: string
  label: string
  description: string
  selected: boolean
  onSelect: () => void
}

export function PresetCard({ value, label, description, selected, onSelect }: PresetCardProps) {
  return (
    <label className="flex items-start gap-3 rounded border border-input p-3 cursor-pointer hover:bg-muted/50">
      <input
        type="radio"
        name="preset"
        value={value}
        checked={selected}
        onChange={onSelect}
        className="mt-1"
      />
      <div>
        <span className="font-medium">{label}</span>
        <p className="text-sm text-muted-foreground">{description}</p>
      </div>
    </label>
  )
}
