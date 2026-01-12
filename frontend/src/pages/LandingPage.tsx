import { useState, useCallback } from 'react'
import { generateQuestions, generateOutputs, ApiError } from '@/lib/api'
import type { Question, GeneratedFile, Answer, ExperienceLevel } from '@/lib/api'
import { ProjectInput } from '@/components/ProjectInput'
import { QuestionFlow } from '@/components/QuestionFlow'
import { OutputEditor } from '@/components/OutputEditor'
import { LoadingState } from '@/components/shared/LoadingState'
import { ErrorMessage } from '@/components/shared/ErrorMessage'
import { RateLimitCountdown } from '@/components/shared/RateLimitCountdown'
import { ExperienceLevelSelector } from '@/components/ExperienceLevelSelector'

type Phase = 'level-select' | 'input' | 'questions' | 'generating' | 'output' | 'error'

interface LandingPageState {
  phase: Phase
  experienceLevel: ExperienceLevel | null
  projectIdea: string
  questions: Question[]
  answers: Map<number, string>
  currentQuestionIndex: number
  generatedFiles: GeneratedFile[]
  editedFiles: Map<string, string>
  error: string | null
  retryAfter: number | null
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

export function LandingPage() {
  const [state, setState] = useState<LandingPageState>({
    phase: 'level-select',
    experienceLevel: null,
    projectIdea: '',
    questions: [],
    answers: new Map(),
    currentQuestionIndex: 0,
    generatedFiles: [],
    editedFiles: new Map(),
    error: null,
    retryAfter: null,
  })

  const handleExperienceLevelSelect = useCallback((level: ExperienceLevel) => {
    setState(prev => ({ ...prev, experienceLevel: level, phase: 'input' }))
  }, [])

  const handleProjectSubmit = useCallback(async (idea: string) => {
    setState(prev => ({ ...prev, projectIdea: idea, phase: 'generating', error: null }))
    
    try {
      const response = await generateQuestions(idea, state.experienceLevel!)
      setState(prev => ({
        ...prev,
        phase: 'questions',
        questions: response.questions,
        currentQuestionIndex: 0,
        answers: new Map(),
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
  }, [state.experienceLevel])

  const handleAnswer = useCallback((questionId: number, answer: string) => {
    setState(prev => {
      const newAnswers = new Map(prev.answers)
      newAnswers.set(questionId, answer)
      return { ...prev, answers: newAnswers }
    })
  }, [])

  const handleBack = useCallback((questionId: number) => {
    setState(prev => {
      const questionIndex = prev.questions.findIndex(q => q.id === questionId)
      if (questionIndex >= 0) {
        return { ...prev, currentQuestionIndex: questionIndex }
      }
      return prev
    })
  }, [])

  const handleNext = useCallback(() => {
    setState(prev => ({
      ...prev,
      currentQuestionIndex: prev.currentQuestionIndex + 1,
    }))
  }, [])

  const handleQuestionsComplete = useCallback(async () => {
    setState(prev => ({ ...prev, phase: 'generating' }))
    
    try {
      const answers: Answer[] = Array.from(state.answers.entries()).map(([questionId, answer]) => ({
        questionId,
        answer,
      }))
      
      const response = await generateOutputs(state.projectIdea, answers, state.experienceLevel!)
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
  }, [state.answers, state.projectIdea, state.experienceLevel])

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

  return (
    <div className="max-w-3xl mx-auto">
      {state.phase === 'level-select' && (
        <ExperienceLevelSelector
          onSelect={handleExperienceLevelSelect}
          selected={state.experienceLevel ?? undefined}
        />
      )}

      {state.phase === 'input' && (
        <ProjectInput
          onSubmit={handleProjectSubmit}
          loading={false}
          examples={EXAMPLE_IDEAS}
        />
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
