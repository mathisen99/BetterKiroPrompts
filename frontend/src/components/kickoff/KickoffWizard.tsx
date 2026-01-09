import { useState } from 'react'
import { StepIndicator } from '../shared/StepIndicator'
import type { KickoffAnswers } from '../../lib/api'

const STEP_LABELS = [
  'Project Identity',
  'Success Criteria',
  'Users & Roles',
  'Data Sensitivity',
  'Auth Model',
  'Concurrency',
  'Risks & Tradeoffs',
  'Boundaries',
  'Non-Goals',
  'Constraints',
]

const initialAnswers: KickoffAnswers = {
  projectIdentity: '',
  successCriteria: '',
  usersAndRoles: '',
  dataSensitivity: '',
  dataLifecycle: { retention: '', deletion: '', export: '', auditLogging: '', backups: '' },
  authModel: 'none',
  concurrency: '',
  risksAndTradeoffs: { topRisks: [], mitigations: [], notHandled: [] },
  boundaries: '',
  boundaryExamples: [],
  nonGoals: '',
  constraints: '',
}

export function KickoffWizard() {
  const [step, setStep] = useState(1)
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  const [answers, setAnswers] = useState<KickoffAnswers>(initialAnswers)

  const canProceed = (): boolean => {
    // Validation will be implemented per-step in tasks 14-19
    return true
  }

  const handleNext = () => {
    if (canProceed() && step < STEP_LABELS.length) {
      setStep(step + 1)
    }
  }

  const handlePrev = () => {
    if (step > 1) {
      setStep(step - 1)
    }
  }

  return (
    <div className="space-y-6">
      <StepIndicator currentStep={step} totalSteps={STEP_LABELS.length} labels={STEP_LABELS} />

      <div className="min-h-[200px] rounded-lg border border-border bg-card p-6">
        <h2 className="mb-4 text-xl font-semibold">{STEP_LABELS[step - 1]}</h2>
        {/* QuestionStep components will be added in tasks 13-19 */}
        <p className="text-muted-foreground">Step content coming soon...</p>
      </div>

      <div className="flex justify-between">
        <button
          onClick={handlePrev}
          disabled={step === 1}
          className="rounded px-4 py-2 text-sm bg-secondary text-secondary-foreground hover:bg-secondary/80 disabled:opacity-50"
          aria-label="Previous step"
        >
          Previous
        </button>
        <button
          onClick={handleNext}
          disabled={step === STEP_LABELS.length}
          className="rounded px-4 py-2 text-sm bg-primary text-primary-foreground hover:bg-primary/80 disabled:opacity-50"
          aria-label="Next step"
        >
          Next
        </button>
      </div>
    </div>
  )
}
