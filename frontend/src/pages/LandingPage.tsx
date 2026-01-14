import { useState, useCallback, useEffect } from 'react'
import { generateQuestions, generateOutputs, ApiError } from '@/lib/api'
import type { Question, GeneratedFile, Answer, ExperienceLevel, HookPreset } from '@/lib/api'
import * as storage from '@/lib/storage'
import type { Phase, SessionState } from '@/lib/storage'
import { ProjectInput } from '@/components/ProjectInput'
import { QuestionFlow } from '@/components/QuestionFlow'
import { OutputEditor } from '@/components/OutputEditor'
import { LoadingState } from '@/components/shared/LoadingState'
import { ErrorMessage } from '@/components/shared/ErrorMessage'
import { RateLimitCountdown } from '@/components/shared/RateLimitCountdown'
import { ExperienceLevelSelector } from '@/components/ExperienceLevelSelector'
import { HookPresetSelector } from '@/components/HookPresetSelector'
import { RestorePrompt } from '@/components/shared/RestorePrompt'
import { SuccessCelebration } from '@/components/shared/SuccessCelebration'
import { WelcomeGuide } from '@/components/WelcomeGuide'

interface LandingPageProps {
  onPhaseChange?: (phase: Phase) => void
  onViewInGallery?: (generationId: string) => void
}

type FailedOperation = 'questions' | 'outputs' | null

interface LandingPageState {
  phase: Phase
  experienceLevel: ExperienceLevel | null
  projectIdea: string
  hookPreset: HookPreset
  questions: Question[]
  answers: Map<number, string>
  currentQuestionIndex: number
  generatedFiles: GeneratedFile[]
  editedFiles: Map<string, string>
  error: string | null
  errorCode: string | null
  retryAfter: number | null
  showRestorePrompt: boolean
  pendingRestore: SessionState | null
  showCelebration: boolean
  canRetry: boolean
  failedOperation: FailedOperation
  loadingStartTime: number | null // Track when loading started for progress display
  generationId: string | null // ID of stored generation for gallery link
}

const EXAMPLE_IDEAS = [
  'A todo app with categories and due dates',
  'A REST API for a blog with authentication',
  'A CLI tool for managing dotfiles',
  'A real-time chat application',
  'A personal finance tracker',
]

// Helper to get user-friendly error message
function getErrorMessage(err: unknown): { message: string; retryAfter: number | null; code: string; canRetry: boolean } {
  if (err instanceof ApiError) {
    // Rate limit error - can retry after waiting
    if (err.status === 429) {
      return {
        message: 'Too many requests. Please wait before trying again.',
        retryAfter: err.retryAfter ?? null,
        code: 'RATE_LIMITED',
        canRetry: false, // Must wait for rate limit
      }
    }
    // Timeout error - can retry
    if (err.status === 504 || err.message.toLowerCase().includes('timeout')) {
      return {
        message: 'Request timed out. Your progress is saved.',
        retryAfter: null,
        code: 'TIMEOUT',
        canRetry: true,
      }
    }
    // Service unavailable - can retry
    if (err.status === 503) {
      return {
        message: 'Service temporarily unavailable. Please try again.',
        retryAfter: err.retryAfter ?? null,
        code: 'SERVICE_UNAVAILABLE',
        canRetry: true,
      }
    }
    // Server error - can retry
    if (err.status >= 500) {
      return {
        message: 'Something went wrong. Please try again.',
        retryAfter: null,
        code: 'SERVER_ERROR',
        canRetry: true,
      }
    }
    // Client error (bad input, etc.) - cannot retry without changes
    return {
      message: err.message,
      retryAfter: err.retryAfter ?? null,
      code: 'VALIDATION_ERROR',
      canRetry: false,
    }
  }
  
  // Network error (fetch failed) - can retry
  if (err instanceof TypeError && err.message.includes('fetch')) {
    return {
      message: 'Unable to connect. Check your connection and try again.',
      retryAfter: null,
      code: 'NETWORK_ERROR',
      canRetry: true,
    }
  }
  
  // Unknown error - can retry
  return {
    message: 'An unexpected error occurred. Please try again.',
    retryAfter: null,
    code: 'UNKNOWN_ERROR',
    canRetry: true,
  }
}

// Convert Map to Record for storage
function answersMapToRecord(answers: Map<number, string>): Record<number, string> {
  const record: Record<number, string> = {}
  answers.forEach((value, key) => {
    record[key] = value
  })
  return record
}

// Convert Record to Map for state
function answersRecordToMap(answers: Record<number, string>): Map<number, string> {
  return new Map(Object.entries(answers).map(([k, v]) => [Number(k), v]))
}

// Save current state to localStorage
function saveState(state: LandingPageState): void {
  // Only save meaningful phases (not generating, output, or error)
  if (state.phase === 'generating' || state.phase === 'output' || state.phase === 'error') {
    return
  }
  
  storage.save({
    phase: state.phase,
    experienceLevel: state.experienceLevel,
    projectIdea: state.projectIdea,
    hookPreset: state.hookPreset,
    questions: state.questions,
    answers: answersMapToRecord(state.answers),
    currentQuestionIndex: state.currentQuestionIndex,
  })
}

// Create initial state, checking for restorable session
function createInitialState(): LandingPageState {
  const savedState = storage.load()
  const hasRestorableState = savedState !== null && savedState.phase !== 'welcome' && savedState.phase !== 'level-select'
  
  // Skip welcome screen if user has seen it before
  const skipWelcome = storage.hasSeenWelcome()
  
  return {
    phase: skipWelcome ? 'level-select' : 'welcome',
    experienceLevel: null,
    projectIdea: '',
    hookPreset: 'default',
    questions: [],
    answers: new Map(),
    currentQuestionIndex: 0,
    generatedFiles: [],
    editedFiles: new Map(),
    error: null,
    errorCode: null,
    retryAfter: null,
    showRestorePrompt: hasRestorableState,
    pendingRestore: hasRestorableState ? savedState : null,
    showCelebration: false,
    canRetry: false,
    failedOperation: null,
    loadingStartTime: null,
    generationId: null,
  }
}

export function LandingPage({ onPhaseChange, onViewInGallery }: LandingPageProps) {
  const [state, setState] = useState<LandingPageState>(createInitialState)

  // Notify parent of phase changes
  useEffect(() => {
    onPhaseChange?.(state.phase)
  }, [state.phase, onPhaseChange])

  const handleRestoreAccept = useCallback(() => {
    setState(prev => {
      if (!prev.pendingRestore) return prev
      
      const restored = prev.pendingRestore
      return {
        ...prev,
        phase: restored.phase,
        experienceLevel: restored.experienceLevel,
        projectIdea: restored.projectIdea,
        hookPreset: restored.hookPreset,
        questions: restored.questions,
        answers: answersRecordToMap(restored.answers),
        currentQuestionIndex: restored.currentQuestionIndex,
        showRestorePrompt: false,
        pendingRestore: null,
      }
    })
  }, [])

  const handleRestoreDecline = useCallback(() => {
    storage.clear()
    setState(prev => ({
      ...prev,
      showRestorePrompt: false,
      pendingRestore: null,
    }))
  }, [])

  const handleWelcomeContinue = useCallback(() => {
    storage.markWelcomeSeen()
    setState(prev => ({ ...prev, phase: 'level-select' }))
  }, [])

  const handleExperienceLevelSelect = useCallback((level: ExperienceLevel) => {
    setState(prev => {
      const newState = { ...prev, experienceLevel: level, phase: 'input' as Phase }
      saveState(newState)
      return newState
    })
  }, [])

  const handleHookPresetSelect = useCallback((preset: HookPreset) => {
    setState(prev => {
      const newState = { ...prev, hookPreset: preset }
      saveState(newState)
      return newState
    })
  }, [])

  const handleProjectSubmit = useCallback(async (idea: string) => {
    const currentExperienceLevel = state.experienceLevel
    
    // Save state before making API call
    const preApiState = { ...state, projectIdea: idea, phase: 'input' as Phase }
    saveState(preApiState)
    
    setState(prev => ({ ...prev, projectIdea: idea, phase: 'generating', error: null, errorCode: null, loadingStartTime: Date.now() }))
    
    try {
      const response = await generateQuestions(idea, currentExperienceLevel!)
      setState(prev => {
        const newState = {
          ...prev,
          phase: 'questions' as Phase,
          questions: response.questions,
          currentQuestionIndex: 0,
          answers: new Map(),
          canRetry: false,
          failedOperation: null as FailedOperation,
          loadingStartTime: null,
        }
        saveState(newState)
        return newState
      })
    } catch (err) {
      const { message, retryAfter, code, canRetry } = getErrorMessage(err)
      setState(prev => ({
        ...prev,
        phase: 'error',
        error: message,
        errorCode: code,
        retryAfter,
        canRetry,
        failedOperation: 'questions' as FailedOperation,
        loadingStartTime: null,
      }))
    }
  }, [state])

  const handleAnswer = useCallback((questionId: number, answer: string) => {
    setState(prev => {
      const newAnswers = new Map(prev.answers)
      newAnswers.set(questionId, answer)
      const newState = { ...prev, answers: newAnswers }
      saveState(newState)
      return newState
    })
  }, [])

  const handleBack = useCallback((questionId: number) => {
    setState(prev => {
      const questionIndex = prev.questions.findIndex(q => q.id === questionId)
      if (questionIndex >= 0) {
        const newState = { ...prev, currentQuestionIndex: questionIndex }
        saveState(newState)
        return newState
      }
      return prev
    })
  }, [])

  const handleNext = useCallback(() => {
    setState(prev => {
      const newState = {
        ...prev,
        currentQuestionIndex: prev.currentQuestionIndex + 1,
      }
      saveState(newState)
      return newState
    })
  }, [])

  const handleQuestionsComplete = useCallback(async () => {
    setState(prev => ({ ...prev, phase: 'generating', error: null, errorCode: null, loadingStartTime: Date.now() }))
    
    try {
      const answers: Answer[] = Array.from(state.answers.entries()).map(([questionId, answer]) => ({
        questionId,
        answer,
      }))
      
      const response = await generateOutputs(state.projectIdea, answers, state.experienceLevel!, state.hookPreset)
      
      // Clear saved state on successful generation
      storage.clear()
      
      // Show celebration first, then transition to output
      setState(prev => ({
        ...prev,
        generatedFiles: response.files,
        editedFiles: new Map(),
        showCelebration: true,
        canRetry: false,
        failedOperation: null,
        loadingStartTime: null,
        generationId: response.generationId ?? null,
      }))
    } catch (err) {
      const { message, retryAfter, code, canRetry } = getErrorMessage(err)
      setState(prev => ({
        ...prev,
        phase: 'error',
        error: message,
        errorCode: code,
        retryAfter,
        canRetry,
        failedOperation: 'outputs' as FailedOperation,
        loadingStartTime: null,
      }))
    }
  }, [state.answers, state.projectIdea, state.experienceLevel, state.hookPreset])

  const handleCelebrationComplete = useCallback(() => {
    setState(prev => ({
      ...prev,
      phase: 'output',
      showCelebration: false,
    }))
  }, [])

  const handleRetry = useCallback(async () => {
    if (!state.canRetry) return
    
    if (state.failedOperation === 'questions') {
      // Retry question generation
      await handleProjectSubmit(state.projectIdea)
    } else if (state.failedOperation === 'outputs') {
      // Retry output generation
      await handleQuestionsComplete()
    }
  }, [state.canRetry, state.failedOperation, state.projectIdea, handleProjectSubmit, handleQuestionsComplete])

  const handleStartOver = useCallback(() => {
    storage.clear()
    setState(createInitialState())
  }, [])

  const handleEdit = useCallback((path: string, content: string) => {
    setState(prev => {
      const newEditedFiles = new Map(prev.editedFiles)
      newEditedFiles.set(path, content)
      return { ...prev, editedFiles: newEditedFiles }
    })
  }, [])

  const handleReset = useCallback((path: string) => {
    setState(prev => {
      const newEditedFiles = new Map(prev.editedFiles)
      newEditedFiles.delete(path)
      return { ...prev, editedFiles: newEditedFiles }
    })
  }, [])

  const getFileContent = useCallback((path: string): string => {
    if (state.editedFiles.has(path)) {
      return state.editedFiles.get(path)!
    }
    const file = state.generatedFiles.find(f => f.path === path)
    return file?.content ?? ''
  }, [state.editedFiles, state.generatedFiles])

  // Show restore prompt if there's saved state
  if (state.showRestorePrompt && state.pendingRestore) {
    return (
      <div className="max-w-3xl mx-auto">
        <RestorePrompt
          projectIdea={state.pendingRestore.projectIdea}
          onRestore={handleRestoreAccept}
          onStartFresh={handleRestoreDecline}
        />
      </div>
    )
  }

  return (
    <div className="max-w-3xl mx-auto">
      {/* Success celebration overlay */}
      {state.showCelebration && (
        <SuccessCelebration onComplete={handleCelebrationComplete} />
      )}

      {state.phase === 'welcome' && (
        <div className="animate-phase-enter">
          <WelcomeGuide onContinue={handleWelcomeContinue} />
        </div>
      )}

      {state.phase === 'level-select' && (
        <div className="animate-phase-enter">
          <ExperienceLevelSelector
            onSelect={handleExperienceLevelSelect}
            selected={state.experienceLevel ?? undefined}
          />
        </div>
      )}

      {state.phase === 'input' && (
        <div className="animate-phase-enter space-y-8">
          <ProjectInput
            onSubmit={handleProjectSubmit}
            loading={false}
            examples={EXAMPLE_IDEAS}
          />
          <HookPresetSelector
            onSelect={handleHookPresetSelect}
            selected={state.hookPreset}
          />
        </div>
      )}

      {state.phase === 'questions' && (
        <div className="animate-phase-enter">
          <QuestionFlow
            questions={state.questions}
            answers={state.answers}
            currentIndex={state.currentQuestionIndex}
            onAnswer={handleAnswer}
            onBack={handleBack}
            onNext={handleNext}
            onComplete={handleQuestionsComplete}
          />
        </div>
      )}

      {state.phase === 'generating' && (
        <div className="animate-phase-enter">
          <LoadingState
            type={state.questions.length === 0 ? 'questions' : 'outputs'}
            startTime={state.loadingStartTime ?? undefined}
          />
        </div>
      )}

      {state.phase === 'output' && (
        <div className="animate-phase-enter">
          <OutputEditor
            files={state.generatedFiles}
            editedFiles={state.editedFiles}
            onEdit={handleEdit}
            onReset={handleReset}
            getFileContent={getFileContent}
            generationId={state.generationId}
            onViewInGallery={onViewInGallery}
          />
        </div>
      )}

      {state.phase === 'error' && (
        <div className="animate-phase-enter py-8 space-y-4">
          <ErrorMessage 
            message={state.error ?? 'An unexpected error occurred'} 
            canRetry={state.canRetry}
            onRetry={handleRetry}
            onStartOver={handleStartOver}
          />
          {state.retryAfter && (
            <RateLimitCountdown retryAfterSeconds={state.retryAfter} />
          )}
        </div>
      )}
    </div>
  )
}
