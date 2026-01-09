import { Toaster as Sonner } from "sonner"

export function Toaster() {
  return (
    <Sonner
      theme="dark"
      className="toaster group"
      toastOptions={{
        classNames: {
          toast: "group toast bg-background text-foreground border-border shadow-lg",
          success: "bg-background text-foreground border-border",
        },
      }}
    />
  )
}
