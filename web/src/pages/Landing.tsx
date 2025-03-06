import { useState } from "react";
import { Button } from "@/components/ui/button";
import {
	Sparkles,
	Package,
	Target,
	Shield,
	Copy,
	Terminal,
	CheckCircle2,
} from "lucide-react";
import {
	Card,
	CardTitle,
	CardDescription,
	CardHeader,
	CardContent,
	CardFooter,
} from "@/components/ui/card";

function Landing() {
	return (
		<div className="min-h-screen bg-background text-foreground flex flex-col">
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
								<a href="#getting-started">
									<Terminal />
									Get Started
								</a>
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
							<CLIExampleCard />
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

function CLIExampleCard() {
	const [copied, setCopied] = useState<string | null>(null);

	const commands = [
		{
			id: "create-new-bento",
			title: "Create a new bento",
			command: "konbini-cli bento new my-api-keys",
			description: "Quickly create a new bento to store your API keys",
		},
		{
			id: "add-secret-to-bento",
			title: "Add a secret to a bento",
			command: "konbini-cli bento add my-api-keys SECRET_KEY",
			description:
				"The CLI will ask for the secret value and not echo it to the console",
		},
		{
			id: "list-all-bentos",
			title: "List all bentos",
			command: "konbini-cli bento list",
			description: "Get a list of all of your bentos",
		},
		{
			id: "share-bento-with-group",
			title: "Share a bento with a group",
			command: "konbini-cli group invite DevTeam john@example.com",
			description:
				"Invite john to join the group DevTeam to share the group's bentos",
		},
	];

	const copyToClipboard = (text: string, id: string) => {
		navigator.clipboard.writeText(text);
		setCopied(id);
		setTimeout(() => setCopied(null), 2000);
	};

	return (
		<Card className="mt-10 w-full max-w-2xl border-border shadow-md">
			<CardHeader className="bg-card border-b border-border">
				<div className="flex items-center gap-2">
					<Terminal className="h-5 w-5 text-primary" />
					<CardTitle>Example CLI Usage</CardTitle>
				</div>
				<CardDescription className="mb-6">
					Learn how to use the command line interface with these examples
				</CardDescription>
			</CardHeader>
			<CardContent className="p-6 pt-0 pb-0 space-y-6">
				{commands.map((cmd) => (
					<div key={cmd.id} className="space-y-2">
						<h3 className="text-sm font-medium">{cmd.title}</h3>
						<div className="bg-muted rounded-md p-4 relative">
							<div className="flex items-start">
								<div className="flex-1 font-mono text-sm overflow-x-auto">
									<span className="text-muted-foreground">$</span> {cmd.command}
								</div>
								<Button
									variant="ghost"
									size="icon"
									className="h-8 w-8 absolute right-2 top-2"
									onClick={() => copyToClipboard(cmd.command, cmd.id)}
								>
									{copied === cmd.id ? (
										<CheckCircle2 className="h-4 w-4 text-green-500" />
									) : (
										<Copy className="h-4 w-4" />
									)}
									<span className="sr-only">Copy command</span>
								</Button>
							</div>
							{cmd.id === "init" && (
								<div className="mt-2 text-sm text-muted-foreground font-mono">
									<div className="text-green-500">
										✓ Created directory my-project
									</div>
									<div className="text-green-500">
										✓ Initialized configuration
									</div>
									<div className="text-green-500">✓ Installed dependencies</div>
									<div>Project ready! Run 'cd my-project' to get started</div>
								</div>
							)}
							{cmd.id === "build" && (
								<div className="mt-2 text-sm text-muted-foreground font-mono">
									<div>Building project...</div>
									<div className="text-green-500">
										✓ Compiled successfully in 2.34s
									</div>
									<div className="text-green-500">
										✓ Output written to ./dist
									</div>
									<div>Build complete! 🚀</div>
								</div>
							)}
						</div>
						<p className="text-sm text-muted-foreground">{cmd.description}</p>
					</div>
				))}
			</CardContent>
			<CardFooter className="bg-muted/50 p-4 border-t border-border">
				<div className="text-sm text-muted-foreground">
					<span className="font-medium">Pro tip:</span> Use the{" "}
					<code className="bg-muted px-1 py-0.5 rounded text-xs">--help</code>{" "}
					flag with any command to see available options.
				</div>
			</CardFooter>
		</Card>
	);
}

export default Landing;
