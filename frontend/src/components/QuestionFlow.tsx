import { useState, type FormEvent } from 'react'
import { Button } from '@/components/ui/button'
import { Textarea } from '@/components/ui/textarea'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import type { Question } from '@/lib/api'
import { ArrowLeft, ArrowRight, CheckCircle2, MessageSquare, Send } from 'lucide-react'

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

  const progress = ((currentIndex + 1) / questions.length) * 100

  return (
    <div className="py-8 space-y-6">
      {/* Progress indicator */}
      <Card className="border-border/30 bg-card/30">
        <CardContent className="py-4">
          <div className="flex items-center justify-between mb-3">
            <div className="flex items-center gap-2">
              <MessageSquare className="h-4 w-4 text-primary" />
              <span className="text-sm font-medium">
                Question {currentIndex + 1} of {questions.length}
              </span>
            </div>
            <span className="text-sm text-muted-foreground font-medium">
              {Math.round(progress)}% complete
            </span>
          </div>
          <div className="h-2 bg-secondary/50 rounded-full overflow-hidden">
            <div
              className="h-full bg-linear-to-r from-primary to-primary/80 transition-all duration-500 ease-out"
              style={{ width: `${progress}%` }}
            />
          </div>
          {/* Step indicators */}
          <div className="flex justify-between mt-3">
            {questions.map((_, idx) => (
              <div
                key={idx}
                className={`h-1.5 w-1.5 rounded-full transition-colors ${
                  idx <= currentIndex ? 'bg-primary' : 'bg-secondary'
                }`}
              />
            ))}
          </div>
        </CardContent>
      </Card>

      {/* Previous Q&A */}
      {previousQuestions.length > 0 && (
        <Card className="border-border/30 bg-card/30">
          <CardHeader className="pb-3">
            <CardTitle className="flex items-center gap-2 text-base font-medium">
              <CheckCircle2 className="h-4 w-4 text-primary/70" />
              Previous Answers
            </CardTitle>
            <CardDescription>Click to edit any previous answer</CardDescription>
          </CardHeader>
          <CardContent className="space-y-2">
            {previousQuestions.map((q, idx) => (
              <button
                key={q.id}
                onClick={() => handlePreviousClick(q.id)}
                className="w-full text-left p-4 rounded-lg border border-border/30 bg-background/30 hover:bg-background/50 hover:border-primary/30 transition-all group"
              >
                <div className="flex items-start gap-3">
                  <span className="shrink-0 w-6 h-6 rounded-full bg-primary/10 text-primary text-xs font-medium flex items-center justify-center">
                    {idx + 1}
                  </span>
                  <div className="flex-1 min-w-0">
                    <p className="text-sm text-muted-foreground mb-1 line-clamp-1">{q.text}</p>
                    <p className="text-sm line-clamp-2 group-hover:text-primary transition-colors">
                      {answers.get(q.id) || '(no answer)'}
                    </p>
                  </div>
                </div>
              </button>
            ))}
          </CardContent>
        </Card>
      )}

      {/* Current question */}
      <Card className="border-border/50 bg-card/50 backdrop-blur">
        <CardHeader>
          <div className="flex items-start gap-3">
            <span className="shrink-0 w-8 h-8 rounded-full bg-primary text-primary-foreground text-sm font-semibold flex items-center justify-center">
              {currentIndex + 1}
            </span>
            <div className="flex-1">
              <CardTitle className="text-xl leading-relaxed">
                {currentQuestion.text}
              </CardTitle>
              {currentQuestion.hint && (
                <CardDescription className="mt-2 text-base">
                  {currentQuestion.hint}
                </CardDescription>
              )}
            </div>
          </div>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-6">
            <Textarea
              id="answer"
              value={currentAnswer}
              onChange={(e) => setCurrentAnswer(e.target.value)}
              placeholder="Type your answer..."
              rows={5}
              className="resize-none text-base bg-background/50 border-border/50 focus:border-primary/50"
              aria-describedby={currentQuestion.hint ? 'hint' : undefined}
            />

            <div className="flex gap-3">
              {canGoBack && (
                <Button
                  type="button"
                  variant="outline"
                  onClick={handleBack}
                  className="gap-2"
                >
                  <ArrowLeft className="h-4 w-4" />
                  Back
                </Button>
              )}
              <Button
                type="submit"
                disabled={!currentAnswer.trim()}
                className="flex-1 h-11 gap-2"
              >
                {isLastQuestion ? (
                  <>
                    <Send className="h-4 w-4" />
                    Generate Files
                  </>
                ) : (
                  <>
                    Next
                    <ArrowRight className="h-4 w-4" />
                  </>
                )}
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  )
}
