import { useState, type FormEvent } from 'react'
import { Button } from '@/components/ui/button'
import { Textarea } from '@/components/ui/textarea'
import type { Question } from '@/lib/api'

interface QuestionFlowProps {
  questions: Question[]
  answers: Map<number, string>
  currentIndex: number
  onAnswer: (questionId: number, answer: string) => void
  onBack: (questionId: number) => void
  onNext: () => void
  onComplete: () => void
}

export function QuestionFlow({
  questions,
  answers,
  currentIndex,
  onAnswer,
  onBack,
  onNext,
  onComplete,
}: QuestionFlowProps) {
  const currentQuestion = questions[currentIndex]
  const [currentAnswer, setCurrentAnswer] = useState(
    answers.get(currentQuestion?.id) ?? ''
  )

  const previousQuestions = questions.slice(0, currentIndex)
  const isLastQuestion = currentIndex === questions.length - 1
  const canGoBack = currentIndex > 0

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault()
    const trimmed = currentAnswer.trim()
    if (!trimmed) return

    onAnswer(currentQuestion.id, trimmed)

    if (isLastQuestion) {
      onComplete()
    } else {
      onNext()
      // Pre-fill next answer if it exists
      const nextQuestion = questions[currentIndex + 1]
      setCurrentAnswer(answers.get(nextQuestion?.id) ?? '')
    }
  }

  const handleBack = () => {
    // Save current answer before going back
    if (currentAnswer.trim()) {
      onAnswer(currentQuestion.id, currentAnswer.trim())
    }
    const prevQuestion = questions[currentIndex - 1]
    onBack(prevQuestion.id)
    setCurrentAnswer(answers.get(prevQuestion.id) ?? '')
  }

  const handlePreviousClick = (questionId: number) => {
    // Save current answer before jumping
    if (currentAnswer.trim()) {
      onAnswer(currentQuestion.id, currentAnswer.trim())
    }
    onBack(questionId)
    setCurrentAnswer(answers.get(questionId) ?? '')
  }

  if (!currentQuestion) {
    return null
  }

  return (
    <div className="py-8">
      {/* Progress indicator */}
      <div className="mb-6">
        <div className="flex justify-between text-sm text-muted-foreground mb-2">
          <span>Question {currentIndex + 1} of {questions.length}</span>
          <span>{Math.round(((currentIndex + 1) / questions.length) * 100)}%</span>
        </div>
        <div className="h-2 bg-secondary rounded-full overflow-hidden">
          <div
            className="h-full bg-primary transition-all duration-300"
            style={{ width: `${((currentIndex + 1) / questions.length) * 100}%` }}
          />
        </div>
      </div>

      {/* Previous Q&A */}
      {previousQuestions.length > 0 && (
        <div className="mb-8 space-y-3">
          <p className="text-sm text-muted-foreground">Previous answers (click to edit):</p>
          {previousQuestions.map((q) => (
            <button
              key={q.id}
              onClick={() => handlePreviousClick(q.id)}
              className="w-full text-left p-3 rounded-lg border border-border bg-card/50 hover:bg-card transition-colors"
            >
              <p className="text-sm text-muted-foreground mb-1">{q.text}</p>
              <p className="text-sm">{answers.get(q.id) || '(no answer)'}</p>
            </button>
          ))}
        </div>
      )}

      {/* Current question */}
      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label htmlFor="answer" className="block text-lg font-medium mb-2">
            {currentQuestion.text}
          </label>
          {currentQuestion.hint && (
            <p className="text-sm text-muted-foreground mb-3">{currentQuestion.hint}</p>
          )}
          <Textarea
            id="answer"
            value={currentAnswer}
            onChange={(e) => setCurrentAnswer(e.target.value)}
            placeholder="Type your answer..."
            rows={4}
            className="resize-none"
            aria-describedby={currentQuestion.hint ? 'hint' : undefined}
          />
        </div>

        <div className="flex gap-3">
          {canGoBack && (
            <Button
              type="button"
              variant="outline"
              onClick={handleBack}
            >
              Back
            </Button>
          )}
          <Button
            type="submit"
            disabled={!currentAnswer.trim()}
            className="flex-1"
          >
            {isLastQuestion ? 'Generate Files' : 'Next'}
          </Button>
        </div>
      </form>
    </div>
  )
}
