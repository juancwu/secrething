import { ModeToggle } from "@/components/ui/mode-toggle";
import { Button } from "@/components/ui/button";
import { GitHubLogo } from "@/components/github-logo";
import { Sparkles, Package, Target, Shield, Rocket } from "lucide-react";
import {
	Card,
	CardTitle,
	CardDescription,
	CardHeader,
	CardContent,
} from "@/components/ui/card";

function Landing() {
	return (
		<div className="min-h-screen bg-background text-foreground flex flex-col">
			<header className="container mx-auto py-6 px-4 flex justify-between items-center">
				<div className="flex items-center gap-2">
					<span className="text-2xl font-bold flex items-center gap-2">
						<Package className="h-6 w-6" /> Konbini
					</span>
				</div>
				<div className="flex items-center gap-4">
					<nav className="hidden md:flex gap-6">
						<a
							href="#features"
							className="hover:text-primary transition-colors"
						>
							Features
						</a>
						<a
							href="#what-is-bento"
							className="hover:text-primary transition-colors"
						>
							What is a Bento?
						</a>
						<a
							href="#security"
							className="hover:text-primary transition-colors"
						>
							Security
						</a>
						<a
							href="#getting-started"
							className="hover:text-primary transition-colors"
						>
							Getting Started
						</a>
					</nav>
					<div className="flex items-center gap-4">
						<Button asChild variant="outline" size="icon">
							<a
								href="https://github.com/juancwu/konbini"
								target="_blank"
								rel="noreferrer"
							>
								<GitHubLogo className="w-4 h-4" />
								<span className="sr-only">GitHub</span>
							</a>
						</Button>
						<ModeToggle />
					</div>
				</div>
			</header>

			<main>
				{/* Hero Section */}
				<section className="py-20 md:py-32 border-b">
					<div className="container mx-auto text-center px-4">
						<h1 className="text-4xl md:text-6xl font-bold mb-6">
							Secure Secret Management Made Simple
						</h1>
						<p className="text-xl md:text-2xl text-muted-foreground max-w-3xl mx-auto mb-10">
							Konbini (Japanese for "convenience store") is your go-to solution
							for securely storing, managing, and sharing sensitive information
							within your organization. Like a well-organized bento box, Konbini
							keeps your secrets neatly compartmentalized and protected.
						</p>
						<div className="flex flex-col sm:flex-row gap-4 justify-center">
							<Button size="lg" asChild>
								<a href="#getting-started">Get Started</a>
							</Button>
							<Button size="lg" variant="outline" asChild>
								<a
									href="https://github.com/juancwu/konbini"
									target="_blank"
									rel="noreferrer"
								>
									View on GitHub
								</a>
							</Button>
						</div>
					</div>
				</section>

				{/* Features Section */}
				<section id="features" className="py-20 border-b">
					<div className="container mx-auto px-4">
						<h2 className="text-3xl md:text-4xl font-bold mb-12 text-center">
							<span className="flex items-center justify-center gap-2">
								<Sparkles className="w-8 h-8 text-yellow-500 dark:text-yellow-400" />
								Features
							</span>
						</h2>
						<div className="grid md:grid-cols-2 lg:grid-cols-3 gap-8">
							<FeatureCard
								title="End-to-End Encryption"
								description="Secrets are encrypted on the client side - Konbini never sees plaintext data"
							/>
							<FeatureCard
								title="Team Sharing"
								description="Securely share credentials with team members through the groups system"
							/>
							<FeatureCard
								title="Fine-grained Permissions"
								description="Control who can access, view, and modify your secrets"
							/>
							<FeatureCard
								title="Two-Factor Authentication"
								description="Enhanced security with TOTP (Time-based One-Time Password)"
							/>
							<FeatureCard
								title="Intuitive CLI"
								description="Command-line interface with TUI support for easy management"
							/>
							<FeatureCard
								title="API Access"
								description="RESTful API for integration with your existing tools"
							/>
							<FeatureCard
								title="Audit Logs"
								description="Track who accessed what and when"
							/>
						</div>
					</div>
				</section>

				{/* What is a Bento Section */}
				<section id="what-is-bento" className="py-20 border-b">
					<div className="container mx-auto px-4">
						<h2 className="text-3xl md:text-4xl font-bold mb-12 text-center">
							<span className="flex items-center justify-center gap-2">
								<Target className="h-8 w-8 text-red-500 dark:test-red-400" />
								What is a Bento?
							</span>
						</h2>
						<div className="max-w-3xl mx-auto">
							<p className="text-lg mb-6">
								In Konbini, a "bento" is a container for your secrets:
							</p>
							<ul className="space-y-4 list-disc pl-6 text-lg">
								<li>
									Each bento has a unique name and can contain multiple
									"ingredients" (key-value pairs)
								</li>
								<li>Bentos can be shared with other users through groups</li>
								<li>Permissions control who can view or modify each bento</li>
								<li>
									All bento contents are encrypted on the client side before
									being sent to the server
								</li>
							</ul>
							<Card className="mt-10">
								<CardHeader>
									<CardTitle>Example CLI Usage:</CardTitle>
								</CardHeader>
								<CardContent>
									<pre className="bg-black text-white p-4 rounded overflow-x-auto">
										<code>{`# Create a new bento
konbini-cli bento new my-api-keys

# Add a secret to a bento
konbini-cli bento add my-api-keys AWS_SECRET_KEY=abcdefg

# List all bentos
konbini-cli bento list

# Share a bento with a group
konbini-cli group invite DevTeam john@example.com`}</code>
									</pre>
								</CardContent>
							</Card>
						</div>
					</div>
				</section>

				{/* Security Section */}
				<section id="security" className="py-20 border-b">
					<div className="container mx-auto px-4">
						<h2 className="text-3xl md:text-4xl font-bold mb-12 text-center">
							<span className="flex items-center justify-center gap-2">
								<Shield className="h-8 w-8 text-blue-500 dark:text-blue-400" />
								Security
							</span>
						</h2>
						<div className="max-w-3xl mx-auto">
							<p className="text-lg mb-6">
								Konbini is designed with security at its core:
							</p>
							<ul className="space-y-4 list-disc pl-6 text-lg">
								<li>
									Client-side encryption ensures your secrets never leave your
									machine in plaintext
								</li>
								<li>Two-factor authentication (TOTP) protects your account</li>
								<li>
									Fine-grained permission system prevents unauthorized access
								</li>
								<li>No plaintext storage of sensitive data</li>
								<li>Email verification for new accounts</li>
							</ul>
						</div>
					</div>
				</section>

				{/* Getting Started Section */}
				<section id="getting-started" className="py-20">
					<div className="container mx-auto px-4">
						<h2 className="text-3xl md:text-4xl font-bold mb-12 text-center">
							<span className="flex items-center justify-center gap-2">
								<Rocket className="h-8 w-8 text-orange-500 dark:text-orange-400" />
								Getting Started
							</span>
						</h2>
						<div className="max-w-3xl mx-auto">
							<h3 className="text-xl font-semibold mb-4">Prerequisites</h3>
							<ul className="space-y-2 list-disc pl-6 mb-8">
								<li>Go 1.21+</li>
								<li>SQLite database (or Turso for production)</li>
								<li>Resend.com account (for email verification)</li>
								<li>
									<a
										href="https://github.com/pressly/goose"
										className="text-primary hover:underline"
									>
										Goose
									</a>{" "}
									for database migrations
								</li>
							</ul>

							<h3 className="text-xl font-semibold mb-4">Installation</h3>
							<div className="space-y-6">
								<div>
									<p className="mb-2">1. Clone the repository</p>
									<pre className="bg-black text-white p-4 rounded overflow-x-auto">
										<code>{`git clone https://github.com/juancwu/konbini.git
cd konbini`}</code>
									</pre>
								</div>
								<div>
									<p className="mb-2">2. Install dependencies</p>
									<pre className="bg-black text-white p-4 rounded overflow-x-auto">
										<code>go mod download</code>
									</pre>
								</div>
								<div>
									<p className="mb-2">
										3. Create a .env file in the project root with the following
										variables:
									</p>
									<pre className="bg-black text-white p-4 rounded overflow-x-auto">
										<code>{`PORT=8080
DB_URL=file:konbini.db
JWT_SECRET=your-secret-key
RESEND_API_KEY=your-resend-api-key
APP_URL=http://localhost:8080`}</code>
									</pre>
								</div>
								<div>
									<p className="mb-2">4. Build the project</p>
									<pre className="bg-black text-white p-4 rounded overflow-x-auto">
										<code>{`# Build the server
go build -o bin/konbini cmd/server/main.go

# Build the CLI
go build -o bin/konbini-cli cmd/cli/main.go`}</code>
									</pre>
								</div>
							</div>
						</div>
					</div>
				</section>
			</main>

			<footer className="mt-auto border-t py-10">
				<div className="container mx-auto px-4">
					<div className="flex flex-col md:flex-row justify-between items-center">
						<div className="mb-4 md:mb-0">
							<span className="text-xl font-bold flex items-center gap-2">
								<Package className="h-5 w-5" /> Konbini
							</span>
							<p className="text-muted-foreground">
								Secure Secret Management Made Simple
							</p>
						</div>
						<div className="flex flex-col md:flex-row items-center gap-4">
							<a
								href="https://github.com/juancwu/konbini"
								className="text-muted-foreground hover:text-foreground transition-colors"
							>
								GitHub
							</a>
							<a
								href="#features"
								className="text-muted-foreground hover:text-foreground transition-colors"
							>
								Features
							</a>
							<a
								href="#what-is-bento"
								className="text-muted-foreground hover:text-foreground transition-colors"
							>
								What is a Bento?
							</a>
							<a
								href="#security"
								className="text-muted-foreground hover:text-foreground transition-colors"
							>
								Security
							</a>
						</div>
					</div>
					<div className="mt-6 text-center text-muted-foreground">
						<p>Licensed under the MIT License</p>
					</div>
				</div>
			</footer>
		</div>
	);
}

function FeatureCard({
	title,
	description,
}: { title: string; description: string }) {
	return (
		<Card>
			<CardHeader>
				<CardTitle>{title}</CardTitle>
				<CardDescription>{description}</CardDescription>
			</CardHeader>
		</Card>
	);
}

export default Landing;
