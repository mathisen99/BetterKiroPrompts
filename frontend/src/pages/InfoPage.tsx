import { ArrowLeft, Home, Sparkles, FileText, Webhook, Shield, Server, ArrowRight, ImageIcon, Package, Settings2, Zap, BookOpen, Github, Mail } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'

interface InfoPageProps {
  onNavigateHome: () => void
  onNavigateGallery: () => void
  onNavigateScan?: () => void
}

const features = [
  {
    icon: Sparkles,
    title: 'Kickoff Prompts',
    description: 'Forces answer-first, no-coding-first thinking. Generate comprehensive prompts that make you think through architecture, data models, and edge cases before writing a single line of code.',
  },
  {
    icon: FileText,
    title: 'Steering Files',
    description: 'Creates .kiro/steering/ files with proper frontmatter. These files guide AI assistants with project-specific context, coding standards, and architectural decisions.',
  },
  {
    icon: Webhook,
    title: 'Hooks Generation',
    description: 'Creates .kiro/hooks/ files from presets. Automate common workflows like running tests on save, updating translations, or spell-checking documentation.',
  },
  {
    icon: Shield,
    title: 'Security Scanning',
    description: 'Scans repositories for vulnerabilities using local tools (Trivy, Semgrep, TruffleHog) with AI-powered remediation guidance for flagged issues.',
  },
]

export function InfoPage({ onNavigateHome, onNavigateGallery, onNavigateScan }: InfoPageProps) {
  return (
    <div className="min-h-screen">
      <div className="container mx-auto px-4 py-8">
        {/* Header */}
        <div className="mb-8 flex items-center justify-between">
          <div className="flex items-center gap-4">
            <Button variant="ghost" size="icon" onClick={onNavigateHome} aria-label="Go back">
              <ArrowLeft className="h-5 w-5" />
            </Button>
            <div>
              <h1 className="text-2xl font-bold">About BetterKiroPrompts</h1>
              <p className="text-muted-foreground">
                Think before you code
              </p>
            </div>
          </div>
          <div className="flex items-center gap-2">
            <Button
              variant="ghost"
              size="sm"
              onClick={onNavigateGallery}
              className="gap-1.5"
            >
              <ImageIcon className="h-4 w-4" />
              Gallery
            </Button>
            <Button
              variant="outline"
              onClick={onNavigateHome}
              className="gap-2"
            >
              <Home className="h-4 w-4" />
              Back to Home
            </Button>
          </div>
        </div>

        {/* Hero Section */}
        <section className="mb-12 text-center">
          <div className="max-w-3xl mx-auto">
            <h2 className="text-4xl font-bold mb-4 bg-linear-to-r from-blue-400 to-indigo-400 bg-clip-text text-transparent">
              Stop Vibe-Coding. Start Thinking.
            </h2>
            <p className="text-xl text-muted-foreground mb-6">
              BetterKiroPrompts helps developers think through their projects before writing code.
              Generate better prompts, steering documents, and hooks that force you to consider
              architecture, security, data models, and edge cases upfront.
            </p>
          </div>
        </section>

        {/* What You Get Section */}
        <section className="mb-12">
          <h3 className="text-2xl font-bold mb-6 text-center">What You'll Download</h3>
          
          {/* ZIP structure */}
          <Card className="bg-card/50 backdrop-blur mb-6">
            <CardHeader>
              <div className="flex items-center gap-3">
                <div className="p-2 rounded-lg bg-primary/10">
                  <Package className="h-5 w-5 text-primary" />
                </div>
                <CardTitle>ZIP File Contents</CardTitle>
              </div>
            </CardHeader>
            <CardContent>
              <div className="font-mono text-sm bg-muted/50 rounded-md p-4 space-y-1">
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

          {/* File explanations */}
          <div className="grid md:grid-cols-2 gap-4">
            <Card className="bg-card/50 backdrop-blur">
              <CardHeader className="pb-2">
                <div className="flex items-center gap-3">
                  <div className="p-2 rounded-lg bg-blue-500/10">
                    <FileText className="h-4 w-4 text-blue-500" />
                  </div>
                  <CardTitle className="text-base">kickoff-prompt.md</CardTitle>
                </div>
              </CardHeader>
              <CardContent>
                <CardDescription>
                  A structured document that forces "thinking before coding". Defines your project's identity, success criteria, users, data handling, and boundaries. Paste this into Kiro chat to start your session with full context.
                </CardDescription>
              </CardContent>
            </Card>

            <Card className="bg-card/50 backdrop-blur">
              <CardHeader className="pb-2">
                <div className="flex items-center gap-3">
                  <div className="p-2 rounded-lg bg-purple-500/10">
                    <BookOpen className="h-4 w-4 text-purple-500" />
                  </div>
                  <CardTitle className="text-base">AGENTS.md</CardTitle>
                </div>
              </CardHeader>
              <CardContent>
                <CardDescription>
                  An{' '}
                  <a href="https://agents.md" target="_blank" rel="noopener noreferrer" className="text-primary hover:underline">
                    open standard
                  </a>
                  {' '}used by 60k+ repos. Think of it as a README for AI agents — defines commit standards, coding principles, and what to do when stuck. Works with Kiro, Cursor, Copilot, and more.
                </CardDescription>
              </CardContent>
            </Card>

            <Card className="bg-card/50 backdrop-blur">
              <CardHeader className="pb-2">
                <div className="flex items-center gap-3">
                  <div className="p-2 rounded-lg bg-green-500/10">
                    <Settings2 className="h-4 w-4 text-green-500" />
                  </div>
                  <CardTitle className="text-base">Steering Files</CardTitle>
                </div>
              </CardHeader>
              <CardContent>
                <CardDescription>
                  Markdown files in <code className="text-xs bg-muted px-1 rounded">.kiro/steering/</code> that give Kiro persistent knowledge. Core files (product, tech, structure) load every session. Language-specific files (security-go.md, quality-web.md) load only when editing matching files.
                </CardDescription>
              </CardContent>
            </Card>

            <Card className="bg-card/50 backdrop-blur">
              <CardHeader className="pb-2">
                <div className="flex items-center gap-3">
                  <div className="p-2 rounded-lg bg-yellow-500/10">
                    <Zap className="h-4 w-4 text-yellow-500" />
                  </div>
                  <CardTitle className="text-base">Hooks</CardTitle>
                </div>
              </CardHeader>
              <CardContent>
                <CardDescription>
                  Event-driven automation in <code className="text-xs bg-muted px-1 rounded">.kiro/hooks/</code>. You choose a preset: Light (just formatters), Basic (+ linters), Default (+ secret scanning), or Strict (+ static analysis). Hooks run on agent stop, prompt submit, or file changes.
                </CardDescription>
              </CardContent>
            </Card>
          </div>
        </section>

        {/* How to Use */}
        <section className="mb-12">
          <h3 className="text-2xl font-bold mb-6 text-center">How to Use the Generated Files</h3>
          <Card className="bg-card/50 backdrop-blur border-primary/20">
            <CardContent className="pt-6">
              <div className="space-y-4">
                <div className="flex gap-4">
                  <div className="shrink-0 h-8 w-8 rounded-full bg-primary/20 flex items-center justify-center text-primary font-semibold">1</div>
                  <div>
                    <p className="font-medium">Create an empty project folder</p>
                    <p className="text-sm text-muted-foreground">This will be your new project's home</p>
                  </div>
                </div>
                <div className="flex gap-4">
                  <div className="shrink-0 h-8 w-8 rounded-full bg-primary/20 flex items-center justify-center text-primary font-semibold">2</div>
                  <div>
                    <p className="font-medium">Extract the ZIP into your project folder</p>
                    <p className="text-sm text-muted-foreground">Make sure hidden files are included (the <code className="bg-muted px-1 rounded">.kiro/</code> folder)</p>
                  </div>
                </div>
                <div className="flex gap-4">
                  <div className="shrink-0 h-8 w-8 rounded-full bg-primary/20 flex items-center justify-center text-primary font-semibold">3</div>
                  <div>
                    <p className="font-medium">Open the folder in Kiro</p>
                    <p className="text-sm text-muted-foreground">Steering files load automatically — Kiro already knows your project</p>
                  </div>
                </div>
                <div className="flex gap-4">
                  <div className="shrink-0 h-8 w-8 rounded-full bg-primary/20 flex items-center justify-center text-primary font-semibold">4</div>
                  <div>
                    <p className="font-medium">Paste the kickoff prompt into Kiro chat</p>
                    <p className="text-sm text-muted-foreground">This gives Kiro the full context to start building your project</p>
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>
        </section>

        {/* Problem Statement */}
        <section className="mb-12">
          <Card className="bg-card/50 backdrop-blur border-amber-500/20">
            <CardHeader>
              <CardTitle className="text-amber-400">The Problem</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-muted-foreground">
                Many developers, especially beginners, jump straight into coding without thinking through
                their project's architecture. They "vibe-code" — writing code based on intuition without
                considering security implications, data structures, concurrency issues, or edge cases.
                This leads to technical debt, security vulnerabilities, and projects that need to be
                rewritten from scratch.
              </p>
            </CardContent>
          </Card>
        </section>

        {/* Who It's For */}
        <section className="mb-12">
          <div className="grid md:grid-cols-2 gap-6">
            <Card className="bg-card/50 backdrop-blur border-blue-500/20">
              <CardHeader>
                <CardTitle className="text-blue-400">For Beginners</CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-muted-foreground">
                  If you're new to programming, this tool helps you avoid the trap of bad initial prompts.
                  By answering guided questions about your project, you'll naturally think through important
                  aspects like data persistence, error handling, and user authentication before you start coding.
                </p>
              </CardContent>
            </Card>
            <Card className="bg-card/50 backdrop-blur border-indigo-500/20">
              <CardHeader>
                <CardTitle className="text-indigo-400">For Experienced Developers</CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-muted-foreground">
                  Even experienced developers benefit from structured thinking. Use the expert-level questions
                  to generate comprehensive steering files and hooks that capture your architectural decisions
                  and coding standards for AI assistants to follow.
                </p>
              </CardContent>
            </Card>
          </div>
        </section>

        {/* Features */}
        <section className="mb-12">
          <h3 className="text-2xl font-bold mb-6 text-center">Features</h3>
          <div className="grid md:grid-cols-2 gap-6">
            {features.map((feature) => (
              <Card key={feature.title} className="bg-card/50 backdrop-blur">
                <CardHeader>
                  <div className="flex items-center gap-3">
                    <div className="p-2 rounded-lg bg-primary/10">
                      <feature.icon className="h-5 w-5 text-primary" />
                    </div>
                    <CardTitle className="text-lg">{feature.title}</CardTitle>
                  </div>
                </CardHeader>
                <CardContent>
                  <CardDescription className="text-sm">
                    {feature.description}
                  </CardDescription>
                </CardContent>
              </Card>
            ))}
          </div>
        </section>

        {/* Self-Hosting */}
        <section className="mb-12">
          <Card className="bg-card/50 backdrop-blur border-green-500/20">
            <CardHeader>
              <div className="flex items-center gap-3">
                <div className="p-2 rounded-lg bg-green-500/10">
                  <Server className="h-5 w-5 text-green-400" />
                </div>
                <CardTitle className="text-green-400">Open Source & Self-Hostable</CardTitle>
              </div>
            </CardHeader>
            <CardContent>
              <p className="text-muted-foreground mb-4">
                BetterKiroPrompts is an open-source project that you can self-host with your own API keys.
                This means you have full control over your data and can customize the tool to fit your
                team's specific needs.
              </p>
              <ul className="list-disc list-inside text-muted-foreground space-y-1 text-sm">
                <li>Use your own OpenAI API key for AI-powered features</li>
                <li>Configure GitHub tokens for private repository scanning</li>
                <li>Deploy with Docker Compose for easy setup</li>
                <li>Customize prompts and questions for your organization</li>
              </ul>
            </CardContent>
          </Card>
        </section>

        {/* CTA */}
        <section className="text-center mb-12">
          <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
            <Button
              size="lg"
              onClick={onNavigateHome}
              className="gap-2 shadow-lg shadow-primary/25"
            >
              Get Started
              <ArrowRight className="h-4 w-4" />
            </Button>
            {onNavigateScan && (
              <Button
                size="lg"
                variant="outline"
                onClick={onNavigateScan}
                className="gap-2"
              >
                <Shield className="h-4 w-4" />
                Try Security Scan
              </Button>
            )}
          </div>
        </section>

        {/* Disclaimer */}
        <section className="text-center text-sm text-muted-foreground border-t border-border/50 pt-6">
          <p>
            This is a community project, not affiliated with or endorsed by AWS or Kiro.
          </p>
          <p className="mt-1">
            <a 
              href="https://kiro.dev" 
              target="_blank" 
              rel="noopener noreferrer"
              className="inline-flex items-center gap-1 text-primary hover:underline"
            >
              Visit the official Kiro website
              <ArrowRight className="h-3 w-3" />
            </a>
          </p>
          <div className="mt-4 flex items-center justify-center gap-4">
            <a
              href="https://github.com/mathisen99/BetterKiroPrompts"
              target="_blank"
              rel="noopener noreferrer"
              className="inline-flex items-center gap-1.5 text-muted-foreground hover:text-foreground transition-colors"
            >
              <Github className="h-4 w-4" />
              GitHub
            </a>
            <a
              href="mailto:tommy.mathisen@aland.net"
              className="inline-flex items-center gap-1.5 text-muted-foreground hover:text-foreground transition-colors"
            >
              <Mail className="h-4 w-4" />
              Contact
            </a>
          </div>
          <p className="mt-2 text-xs text-muted-foreground/70">
            Created by Tommy Mathisen
          </p>
        </section>
      </div>
    </div>
  )
}
