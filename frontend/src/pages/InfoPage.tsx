import { ArrowLeft, Home, Sparkles, FileText, Webhook, Shield, Server, ArrowRight, ImageIcon } from 'lucide-react'
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
        <section className="text-center">
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
      </div>
    </div>
  )
}
