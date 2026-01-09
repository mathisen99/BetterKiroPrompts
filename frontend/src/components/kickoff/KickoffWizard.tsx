import { useState } from 'react'
import { StepIndicator } from '../shared/StepIndicator'
import { QuestionStep } from './QuestionStep'
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
  const [answers, setAnswers] = useState<KickoffAnswers>(initialAnswers)

  const updateAnswer = <K extends keyof KickoffAnswers>(key: K, value: KickoffAnswers[K]) => {
    setAnswers((prev) => ({ ...prev, [key]: value }))
  }

  const updateLifecycle = (field: keyof typeof answers.dataLifecycle, value: string) => {
    setAnswers((prev) => ({ ...prev, dataLifecycle: { ...prev.dataLifecycle, [field]: value } }))
  }

  const updateRisks = (field: keyof typeof answers.risksAndTradeoffs, value: string) => {
    const arr = value.split('\n').map((s) => s.trim()).filter(Boolean)
    setAnswers((prev) => ({ ...prev, risksAndTradeoffs: { ...prev.risksAndTradeoffs, [field]: arr } }))
  }

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
        {step === 1 && (
          <QuestionStep
            question="Restate your project in one sentence"
            description="What is this project about?"
            value={answers.projectIdentity}
            onChange={(v) => updateAnswer('projectIdentity', v)}
            placeholder="e.g., A task management app for small teams"
          />
        )}
        {step === 2 && (
          <QuestionStep
            question="What does 'done' mean?"
            description="Define the success criteria for this project"
            value={answers.successCriteria}
            onChange={(v) => updateAnswer('successCriteria', v)}
            placeholder="e.g., Users can create, assign, and complete tasks"
          />
        )}
        {step === 3 && (
          <QuestionStep
            question="Who uses this?"
            description="List user types: anonymous, authenticated, admin, etc."
            value={answers.usersAndRoles}
            onChange={(v) => updateAnswer('usersAndRoles', v)}
            placeholder="e.g., Anonymous visitors, registered users, team admins"
          />
        )}
        {step === 4 && (
          <div className="space-y-6">
            <QuestionStep
              question="What data is stored?"
              description="Label sensitive data explicitly (PII, credentials, financial, etc.)"
              value={answers.dataSensitivity}
              onChange={(v) => updateAnswer('dataSensitivity', v)}
              placeholder="e.g., User emails (PII), hashed passwords, task content"
            />
            <div className="space-y-4 border-t border-border pt-4">
              <p className="text-sm font-medium">Data Lifecycle</p>
              <div className="grid gap-3">
                <label className="block">
                  <span className="text-sm text-muted-foreground">Retention policy</span>
                  <input type="text" value={answers.dataLifecycle.retention} onChange={(e) => updateLifecycle('retention', e.target.value)} placeholder="e.g., 2 years after account deletion" className="mt-1 w-full rounded-md border border-input bg-background px-3 py-2 text-sm" />
                </label>
                <label className="block">
                  <span className="text-sm text-muted-foreground">Deletion process</span>
                  <input type="text" value={answers.dataLifecycle.deletion} onChange={(e) => updateLifecycle('deletion', e.target.value)} placeholder="e.g., Soft delete, hard delete after 30 days" className="mt-1 w-full rounded-md border border-input bg-background px-3 py-2 text-sm" />
                </label>
                <label className="block">
                  <span className="text-sm text-muted-foreground">Export capability</span>
                  <input type="text" value={answers.dataLifecycle.export} onChange={(e) => updateLifecycle('export', e.target.value)} placeholder="e.g., JSON export of user data on request" className="mt-1 w-full rounded-md border border-input bg-background px-3 py-2 text-sm" />
                </label>
                <label className="block">
                  <span className="text-sm text-muted-foreground">Audit logging</span>
                  <input type="text" value={answers.dataLifecycle.auditLogging} onChange={(e) => updateLifecycle('auditLogging', e.target.value)} placeholder="e.g., Log all data access and modifications" className="mt-1 w-full rounded-md border border-input bg-background px-3 py-2 text-sm" />
                </label>
                <label className="block">
                  <span className="text-sm text-muted-foreground">Backup strategy</span>
                  <input type="text" value={answers.dataLifecycle.backups} onChange={(e) => updateLifecycle('backups', e.target.value)} placeholder="e.g., Daily backups, 7-day retention" className="mt-1 w-full rounded-md border border-input bg-background px-3 py-2 text-sm" />
                </label>
              </div>
            </div>
          </div>
        )}
        {step === 5 && (
          <div className="space-y-3">
            <label htmlFor="auth-model" className="block text-lg font-medium">
              Authentication Model
            </label>
            <p className="text-sm text-muted-foreground">How will users authenticate?</p>
            <select
              id="auth-model"
              value={answers.authModel}
              onChange={(e) => updateAnswer('authModel', e.target.value as 'none' | 'basic' | 'external')}
              className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
            >
              <option value="none">None — No authentication required</option>
              <option value="basic">Basic — Username/password or email/password</option>
              <option value="external">External — OAuth, SSO, or third-party provider</option>
            </select>
          </div>
        )}
        {step === 6 && (
          <QuestionStep
            question="Concurrency Expectations"
            description="Multi-user access? Background jobs? Shared state?"
            value={answers.concurrency}
            onChange={(v) => updateAnswer('concurrency', v)}
            placeholder="e.g., Multiple users editing same document, background email jobs"
          />
        )}
        {step === 7 && (
          <div className="space-y-4">
            <div className="space-y-2">
              <label htmlFor="top-risks" className="block text-sm font-medium">Top 3 Risks</label>
              <textarea id="top-risks" rows={3} value={answers.risksAndTradeoffs.topRisks.join('\n')} onChange={(e) => updateRisks('topRisks', e.target.value)} placeholder="One risk per line" className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm" />
            </div>
            <div className="space-y-2">
              <label htmlFor="mitigations" className="block text-sm font-medium">Simplest Mitigations</label>
              <textarea id="mitigations" rows={3} value={answers.risksAndTradeoffs.mitigations.join('\n')} onChange={(e) => updateRisks('mitigations', e.target.value)} placeholder="One mitigation per line" className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm" />
            </div>
            <div className="space-y-2">
              <label htmlFor="not-handled" className="block text-sm font-medium">Explicitly Not Handled</label>
              <textarea id="not-handled" rows={3} value={answers.risksAndTradeoffs.notHandled.join('\n')} onChange={(e) => updateRisks('notHandled', e.target.value)} placeholder="One item per line" className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm" />
            </div>
          </div>
        )}
        {step === 8 && (
          <div className="space-y-4">
            <QuestionStep
              question="Public vs Private Boundaries"
              description="What data/actions are public vs private?"
              value={answers.boundaries}
              onChange={(v) => updateAnswer('boundaries', v)}
              placeholder="e.g., Task titles are public, task details are private to team members"
            />
            <div className="space-y-2">
              <label htmlFor="boundary-examples" className="block text-sm font-medium">Concrete Access Examples (2-3)</label>
              <p className="text-sm text-muted-foreground">Who can read/write what?</p>
              <textarea id="boundary-examples" rows={4} value={answers.boundaryExamples.join('\n')} onChange={(e) => updateAnswer('boundaryExamples', e.target.value.split('\n').map(s => s.trim()).filter(Boolean))} placeholder="One example per line, e.g.:\nAnonymous users can view public projects\nTeam members can edit their own tasks\nAdmins can delete any task" className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm" />
            </div>
          </div>
        )}
        {step === 9 && (
          <QuestionStep
            question="Non-Goals"
            description="What will NOT be built in this project?"
            value={answers.nonGoals}
            onChange={(v) => updateAnswer('nonGoals', v)}
            placeholder="e.g., Mobile app, real-time collaboration, third-party integrations"
          />
        )}
        {step === 10 && (
          <QuestionStep
            question="Constraints"
            description="Time limits, simplicity requirements, tech restrictions?"
            value={answers.constraints}
            onChange={(v) => updateAnswer('constraints', v)}
            placeholder="e.g., Must ship in 2 weeks, no external dependencies, PostgreSQL only"
          />
        )}
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
