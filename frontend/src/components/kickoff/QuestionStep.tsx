interface QuestionStepProps {
  question: string
  description?: string
  value: string
  onChange: (value: string) => void
  placeholder?: string
}

export function QuestionStep({ question, description, value, onChange, placeholder }: QuestionStepProps) {
  const id = question.toLowerCase().replace(/\s+/g, '-')
  const descId = `${id}-desc`

  return (
    <div className="space-y-3">
      <label htmlFor={id} className="block text-lg font-medium">
        {question}
      </label>
      {description && (
        <p id={descId} className="text-sm text-muted-foreground">{description}</p>
      )}
      <textarea
        id={id}
        value={value}
        onChange={(e) => onChange(e.target.value)}
        placeholder={placeholder}
        rows={4}
        className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring"
        aria-describedby={description ? descId : undefined}
      />
    </div>
  )
}
