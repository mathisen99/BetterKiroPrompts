import { useState, useCallback } from 'react'
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
  retryAfter: number | null
  showRestorePrompt: boolean
  pendingRestore: SessionState | null
}

const EXAMPLE_IDEAS = [
  'A todo app with categories and due dates',
  'A REST API for a blog with authentication',
  'A CLI tool for managing dotfiles',
  'A real-time chat application',
  'A personal finance tracker',
]

// Helper to get user-friendly error message
function getErrorMessage(err: unknown): { message: string; retryAfter: number | null } {
  if (err instanceof ApiError) {
    // Rate limit error
    if (err.status === 429) {
      return {
        message: 'Too many requests. Please wait before trying again.',
        retryAfter: err.retryAfter ?? null,
      }
    }
    // Timeout error
    if (err.status === 504 || err.message.toLowerCase().includes('timeout')) {
      return {
        message: 'Request timed out. Please refresh and start over.',
        retryAfter: null,
      }
    }
    // Server error
    if (err.status >= 500) {
      return {
        message: 'Generation failed. Please refresh and try again.',
        retryAfter: null,
      }
    }
    // Client error (bad input, etc.)
    return {
      message: err.message,
      retryAfter: err.retryAfter ?? null,
    }
  }
  
  // Network error (fetch failed)
  if (err instanceof TypeError && err.message.includes('fetch')) {
    return {
      message: 'Unable to connect. Please check your connection and refresh.',
      retryAfter: null,
    }
  }
  
  // Unknown error
  return {
    message: 'An unexpected error occurred. Please refresh and try again.',
    retryAfter: null,
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
  const hasRestorableState = savedState !== null && savedState.phase !== 'level-select'
  
  return {
    phase: 'level-select',
    experienceLevel: null,
    projectIdea: '',
    hookPreset: 'default',
    questions: [],
    answers: new Map(),
    currentQuestionIndex: 0,
    generatedFiles: [],
    editedFiles: new Map(),
    error: null,
    retryAfter: null,
    showRestorePrompt: hasRestorableState,
    pendingRestore: hasRestorableState ? savedState : null,
  }
}

export function LandingPage() {
  const [state, setState] = useState<LandingPageState>(createInitialState)

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
    
    setState(prev => ({ ...prev, projectIdea: idea, phase: 'generating', error: null }))
    
    try {
      const response = await generateQuestions(idea, currentExperienceLevel!)
      setState(prev => {
        const newState = {
          ...prev,
          phase: 'questions' as Phase,
          questions: response.questions,
          currentQuestionIndex: 0,
          answers: new Map(),
        }
        saveState(newState)
        return newState
      })
    } catch (err) {
      const { message, retryAfter } = getErrorMessage(err)
      setState(prev => ({
        ...prev,
        phase: 'error',
        error: message,
        retryAfter,
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
    setState(prev => ({ ...prev, phase: 'generating' }))
    
    try {
      const answers: Answer[] = Array.from(state.answers.entries()).map(([questionId, answer]) => ({
        questionId,
        answer,
      }))
      
      const response = await generateOutputs(state.projectIdea, answers, state.experienceLevel!, state.hookPreset)
      
      // Clear saved state on successful generation
      storage.clear()
      
      setState(prev => ({
        ...prev,
        phase: 'output',
        generatedFiles: response.files,
        editedFiles: new Map(),
      }))
    } catch (err) {
      const { message, retryAfter } = getErrorMessage(err)
      setState(prev => ({
        ...prev,
        phase: 'error',
        error: message,
        retryAfter,
      }))
    }
  }, [state.answers, state.projectIdea, state.experienceLevel, state.hookPreset])

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
      {state.phase === 'level-select' && (
        <ExperienceLevelSelector
          onSelect={handleExperienceLevelSelect}
          selected={state.experienceLevel ?? undefined}
        />
      )}

      {state.phase === 'input' && (
        <div className="space-y-8">
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
        <QuestionFlow
          questions={state.questions}
          answers={state.answers}
          currentIndex={state.currentQuestionIndex}
          onAnswer={handleAnswer}
          onBack={handleBack}
          onNext={handleNext}
          onComplete={handleQuestionsComplete}
        />
      )}

      {state.phase === 'generating' && (
        <LoadingState
          message={
            state.questions.length === 0
              ? 'Generating questions for your project...'
              : 'Generating your Kiro files...'
          }
          estimatedTime="up to 60 seconds"
        />
      )}

      {state.phase === 'output' && (
        <OutputEditor
          files={state.generatedFiles}
          editedFiles={state.editedFiles}
          onEdit={handleEdit}
          onReset={handleReset}
          getFileContent={getFileContent}
        />
      )}

      {state.phase === 'error' && (
        <div className="py-8 space-y-4">
          <ErrorMessage message={state.error ?? 'An unexpected error occurred'} />
          {state.retryAfter ? (
            <RateLimitCountdown retryAfterSeconds={state.retryAfter} />
          ) : (
            <p className="text-sm text-muted-foreground text-center">
              Please refresh the page to start over.
            </p>
          )}
        </div>
      )}
    </div>
  )
}
