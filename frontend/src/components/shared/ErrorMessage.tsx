interface ErrorMessageProps {
  message: string
}

export function ErrorMessage({ message }: ErrorMessageProps) {
  return (
    <div className="rounded-md border border-destructive/50 bg-destructive/10 p-4" role="alert">
      <p className="text-sm text-destructive">{message}</p>
    </div>
  )
}
