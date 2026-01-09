interface StepIndicatorProps {
  currentStep: number
  totalSteps: number
  labels?: string[]
}

export function StepIndicator({ currentStep, totalSteps, labels }: StepIndicatorProps) {
  const currentLabel = labels?.[currentStep - 1] ?? `Step ${currentStep}`

  return (
    <div className="flex items-center gap-2" role="progressbar" aria-valuenow={currentStep} aria-valuemin={1} aria-valuemax={totalSteps} aria-label={`Step ${currentStep} of ${totalSteps}: ${currentLabel}`}>
      <span className="sr-only" aria-live="polite">Step {currentStep} of {totalSteps}: {currentLabel}</span>
      {Array.from({ length: totalSteps }, (_, i) => {
        const step = i + 1
        const isCompleted = step < currentStep
        const isCurrent = step === currentStep

        return (
          <div key={step} className="flex items-center gap-2">
            <div
              className={`flex h-8 w-8 items-center justify-center rounded-full text-sm font-medium ${
                isCompleted
                  ? 'bg-primary text-primary-foreground'
                  : isCurrent
                    ? 'border-2 border-primary text-primary'
                    : 'border border-muted text-muted-foreground'
              }`}
              aria-hidden="true"
            >
              {isCompleted ? '✓' : step}
            </div>
            {step < totalSteps && (
              <div className={`h-0.5 w-8 ${isCompleted ? 'bg-primary' : 'bg-muted'}`} aria-hidden="true" />
            )}
          </div>
        )
      })}
    </div>
  )
}
