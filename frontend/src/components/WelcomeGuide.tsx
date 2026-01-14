import { useState } from 'react'
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { 
  Lightbulb, 
  FileText, 
  Settings2, 
  Zap, 
  ChevronDown, 
  ChevronUp,
  FolderOpen,
  Download,
  Terminal,
  BookOpen,
  Package
} from 'lucide-react'

interface WelcomeGuideProps {
  onContinue: () => void
}

export function WelcomeGuide({ onContinue }: WelcomeGuideProps) {
  const [expanded, setExpanded] = useState(false)

  return (
    <div className="py-8 space-y-6">
      <div className="text-center space-y-3">
        <h2 className="text-3xl font-bold tracking-tight">
          Welcome to BetterKiroPrompts
        </h2>
        <p className="text-muted-foreground text-lg max-w-2xl mx-auto">
          Generate tailored configuration files for{' '}
          <a 
            href="https://kiro.dev" 
            target="_blank" 
            rel="noopener noreferrer"
            className="text-primary hover:underline"
          >
            Kiro
          </a>
          , AWS's spec-driven AI coding assistant.
        </p>
      </div>

      {/* How it works - always visible */}
      <Card className="border-primary/20 bg-primary/5">
        <CardHeader className="pb-3">
          <CardTitle className="text-lg flex items-center gap-2">
            <Lightbulb className="h-5 w-5 text-primary" />
            How it works
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="grid gap-4 sm:grid-cols-3">
            <div className="flex gap-3">
              <div className="shrink-0 h-8 w-8 rounded-full bg-primary/20 flex items-center justify-center text-primary font-semibold">
                1
              </div>
              <div>
                <p className="font-medium">Describe your project</p>
                <p className="text-sm text-muted-foreground">Tell us what you want to build</p>
              </div>
            </div>
            <div className="flex gap-3">
              <div className="shrink-0 h-8 w-8 rounded-full bg-primary/20 flex items-center justify-center text-primary font-semibold">
                2
              </div>
              <div>
                <p className="font-medium">Answer questions</p>
                <p className="text-sm text-muted-foreground">AI asks about your needs</p>
              </div>
            </div>
            <div className="flex gap-3">
              <div className="shrink-0 h-8 w-8 rounded-full bg-primary/20 flex items-center justify-center text-primary font-semibold">
                3
              </div>
              <div>
                <p className="font-medium">Download ZIP</p>
                <p className="text-sm text-muted-foreground">Extract and start coding</p>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Expandable section for more details */}
      <div className="space-y-4">
        <button
          onClick={() => setExpanded(!expanded)}
          className="w-full flex items-center justify-center gap-2 text-sm text-muted-foreground hover:text-foreground transition-colors py-2"
        >
          {expanded ? (
            <>
              <ChevronUp className="h-4 w-4" />
              Hide details
            </>
          ) : (
            <>
              <ChevronDown className="h-4 w-4" />
              What's in the ZIP file?
            </>
          )}
        </button>

        {expanded && (
          <div className="space-y-4 animate-in fade-in slide-in-from-top-2 duration-200">
            {/* ZIP contents */}
            <Card className="border-dashed">
              <CardHeader className="pb-2">
                <CardTitle className="text-base flex items-center gap-2">
                  <Package className="h-4 w-4 text-primary" />
                  What you'll download
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="font-mono text-sm bg-muted/50 rounded-md p-3 space-y-1">
                  <p className="text-muted-foreground">your-project.zip</p>
                  <p className="pl-4">├── kickoff-prompt.md</p>
                  <p className="pl-4">├── AGENTS.md</p>
                  <p className="pl-4">└── .kiro/</p>
                  <p className="pl-8">├── steering/</p>
                  <p className="pl-12 text-muted-foreground">├── product.md <span className="text-xs">(always loaded)</span></p>
                  <p className="pl-12 text-muted-foreground">├── tech.md <span className="text-xs">(always loaded)</span></p>
                  <p className="pl-12 text-muted-foreground">├── structure.md <span className="text-xs">(always loaded)</span></p>
                  <p className="pl-12 text-muted-foreground">├── security-*.md <span className="text-xs">(per language)</span></p>
                  <p className="pl-12 text-muted-foreground">└── quality-*.md <span className="text-xs">(per language)</span></p>
                  <p className="pl-8">└── hooks/</p>
                  <p className="pl-12 text-muted-foreground">├── format-on-stop.kiro.hook</p>
                  <p className="pl-12 text-muted-foreground">├── lint-on-stop.kiro.hook</p>
                  <p className="pl-12 text-muted-foreground">├── secret-scan.kiro.hook</p>
                  <p className="pl-12 text-muted-foreground">└── ... <span className="text-xs">(based on preset)</span></p>
                </div>
              </CardContent>
            </Card>

            {/* Generated files explanation */}
            <div className="grid gap-4 sm:grid-cols-2">
              <Card>
                <CardHeader className="pb-2">
                  <CardTitle className="text-base flex items-center gap-2">
                    <FileText className="h-4 w-4 text-blue-500" />
                    kickoff-prompt.md
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <CardDescription>
                    A structured document that forces "thinking before coding". Defines your project's identity, success criteria, users, data handling, and boundaries. <strong>Paste this into Kiro chat</strong> to start your session with full context.
                  </CardDescription>
                </CardContent>
              </Card>

              <Card>
                <CardHeader className="pb-2">
                  <CardTitle className="text-base flex items-center gap-2">
                    <BookOpen className="h-4 w-4 text-purple-500" />
                    AGENTS.md
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <CardDescription>
                    An{' '}
                    <a 
                      href="https://agents.md" 
                      target="_blank" 
                      rel="noopener noreferrer"
                      className="text-primary hover:underline"
                    >
                      open standard
                    </a>
                    {' '}used by 60k+ repos. Think of it as a README for AI agents — defines commit standards, coding principles, and what to do when stuck. Works with Kiro, Cursor, Copilot, and more.
                  </CardDescription>
                </CardContent>
              </Card>

              <Card>
                <CardHeader className="pb-2">
                  <CardTitle className="text-base flex items-center gap-2">
                    <Settings2 className="h-4 w-4 text-green-500" />
                    Steering Files
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <CardDescription>
                    Markdown files in <code className="text-xs bg-muted px-1 rounded">.kiro/steering/</code> that give Kiro persistent knowledge. Core files (product, tech, structure) load <strong>every session</strong>. Language-specific files (security-go.md, quality-web.md) load <strong>only when editing matching files</strong>.
                  </CardDescription>
                </CardContent>
              </Card>

              <Card>
                <CardHeader className="pb-2">
                  <CardTitle className="text-base flex items-center gap-2">
                    <Zap className="h-4 w-4 text-yellow-500" />
                    Hooks
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <CardDescription>
                    Event-driven automation in <code className="text-xs bg-muted px-1 rounded">.kiro/hooks/</code>. You choose a preset: <strong>Light</strong> (just formatters), <strong>Basic</strong> (+ linters), <strong>Default</strong> (+ secret scanning), or <strong>Strict</strong> (+ static analysis). Hooks run on agent stop, prompt submit, or file changes.
                  </CardDescription>
                </CardContent>
              </Card>
            </div>

            {/* How to use */}
            <Card className="border-primary/20 bg-primary/5">
              <CardHeader className="pb-2">
                <CardTitle className="text-base">How to use the generated files</CardTitle>
              </CardHeader>
              <CardContent className="space-y-3 text-sm">
                <div className="flex gap-3">
                  <div className="shrink-0 h-6 w-6 rounded-full bg-primary/20 flex items-center justify-center text-primary text-xs font-semibold">
                    1
                  </div>
                  <div>
                    <p className="font-medium">Create an empty project folder</p>
                    <p className="text-muted-foreground">This will be your new project's home</p>
                  </div>
                </div>
                <div className="flex gap-3">
                  <div className="shrink-0 h-6 w-6 rounded-full bg-primary/20 flex items-center justify-center text-primary text-xs font-semibold">
                    2
                  </div>
                  <div className="flex items-start gap-2">
                    <Download className="h-4 w-4 shrink-0 mt-0.5 text-muted-foreground" />
                    <div>
                      <p className="font-medium">Extract the ZIP into your project folder</p>
                      <p className="text-muted-foreground">Make sure hidden files are included (the <code className="bg-muted px-1 rounded">.kiro/</code> folder)</p>
                    </div>
                  </div>
                </div>
                <div className="flex gap-3">
                  <div className="shrink-0 h-6 w-6 rounded-full bg-primary/20 flex items-center justify-center text-primary text-xs font-semibold">
                    3
                  </div>
                  <div className="flex items-start gap-2">
                    <FolderOpen className="h-4 w-4 shrink-0 mt-0.5 text-muted-foreground" />
                    <div>
                      <p className="font-medium">Open the folder in Kiro</p>
                      <p className="text-muted-foreground">Steering files load automatically — Kiro already knows your project</p>
                    </div>
                  </div>
                </div>
                <div className="flex gap-3">
                  <div className="shrink-0 h-6 w-6 rounded-full bg-primary/20 flex items-center justify-center text-primary text-xs font-semibold">
                    4
                  </div>
                  <div className="flex items-start gap-2">
                    <Terminal className="h-4 w-4 shrink-0 mt-0.5 text-muted-foreground" />
                    <div>
                      <p className="font-medium">Paste the kickoff prompt into Kiro chat</p>
                      <p className="text-muted-foreground">This gives Kiro the full context to start building your project</p>
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* Why AGENTS.md */}
            <Card>
              <CardHeader className="pb-2">
                <CardTitle className="text-sm text-muted-foreground">Why include AGENTS.md?</CardTitle>
              </CardHeader>
              <CardContent>
                <CardDescription>
                  AGENTS.md is an open standard that works across many AI tools — not just Kiro. If you later use Cursor, Copilot, Windsurf, or other AI assistants, they'll all read the same file. It's like a README, but for AI agents: commit message formats, coding standards, and project conventions in one place.
                </CardDescription>
              </CardContent>
            </Card>
          </div>
        )}
      </div>

      {/* Continue button */}
      <div className="flex justify-center pt-4">
        <Button size="lg" onClick={onContinue} className="px-8">
          Get Started
        </Button>
      </div>
    </div>
  )
}
